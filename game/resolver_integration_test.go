package game_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/burrbd/dip/game"
	"github.com/burrbd/dip/game/order"
	"github.com/burrbd/dip/game/order/board"
	"github.com/cheekybits/is"
)

// specs contains integration test cases in two groups:
//
//  1. Non-DATC baseline tests that cover core mechanics with descriptive names.
//  2. DATC tests in canonical order (6.A through 6.G).
//
// Unimplemented DATC tests are commented out with a note of what feature is
// needed. To implement a test: uncomment it, run `go test -v ./game/...`,
// watch it fail, implement the feature, watch it pass. See CLAUDE.md for the
// recommended implementation order.
//
// Sea-territory abbreviations used in commented-out fleet tests:
//
//	adr=Adriatic Sea, aeg=Aegean Sea, bal=Baltic Sea, bla=Black Sea,
//	ech=English Channel, eme=Eastern Mediterranean, gob=Gulf of Bothnia,
//	gol=Gulf of Lyon, hel=Helgoland Bight, ion=Ionian Sea, iri=Irish Sea,
//	mao=Mid-Atlantic Ocean, nao=North Atlantic Ocean, nth=North Sea,
//	nwg=Norwegian Sea, ska=Skagerrak, tys=Tyrrhenian Sea, wme=Western Mediterranean
var specs = []spec{

	// ===== NON-DATC BASELINE TESTS =====

	{
		description: "given a unit moves unchallenged, then unit changes territory",
		orders: []*result{
			{order: "A Bud-Vie", position: "vie"},
		},
	},
	{
		description: "given two units attack same territory without support, then neither unit wins territory",
		orders: []*result{
			{order: "A Bud-Vie", position: "bud"},
			{order: "A Gal-Vie", position: "gal"},
		},
	},
	{
		description: "given units attack in circular chain with a hold, then all units stay",
		orders: []*result{
			{order: "A Bud-Gal", position: "bud"},
			{order: "A Gal-Vie", position: "gal"},
			{order: "A Vie H", position: "vie"},
		},
	},
	{
		description: "given two units attack an empty territory, then supported attack wins",
		orders: []*result{
			{order: "A Gal-Vie", position: "vie"},
			{order: "A Boh S A Gal-Vie", position: "boh"},
			{order: "A Bud-Vie", position: "bud"},
		},
	},
	{
		description: "given two units attack an empty territory, then unit with greatest support wins",
		orders: []*result{
			{order: "A Gal-Vie", position: "vie"},
			{order: "A Boh S A Gal-Vie", position: "boh"},
			{order: "A Tri S A Gal-Vie", position: "tri"},
			{order: "A Bud-Vie", position: "bud"},
			{order: "A Tyr S A Bud-Vie", position: "tyr"},
		},
	},
	{
		description: "given unit holds territory, then unit remains on territory",
		orders: []*result{
			{order: "A Vie H", position: "vie"},
		},
	},
	{
		// This also covers DATC 6.D.15: the defender (Vie-Boh) attacks the
		// supporter (Boh), but moveSupportCut exempts it because the cutter
		// moves FROM the attacked territory.
		description: "given unit attacks territory and defending territory attacks support, " +
			"then attacking unit still wins",
		orders: []*result{
			{order: "A Gal-Vie", position: "vie"},
			{order: "A Boh S A Gal-Vie", position: "boh"},
			{order: "A Vie-Boh", position: "vie", defeated: true},
		},
	},
	{
		description: "given two units attack each other (counterattack), then both units bounce",
		orders: []*result{
			{order: "A Vie-Bud", position: "vie"},
			{order: "A Bud-Vie", position: "bud"},
		},
	},
	{
		description: "given a counterattack, and another attacks one counterattack party, then all units bounce",
		orders: []*result{
			{order: "A Vie-Bud", position: "vie"},
			{order: "A Bud-Vie", position: "bud"},
			{order: "A Boh-Vie", position: "boh"},
		},
	},
	{
		description: "given a counterattack, and another unit attacks one counterattack party with support, " +
			"then supported unit wins",
		orders: []*result{
			{order: "A Vie-Bud", position: "vie", defeated: true},
			{order: "A Bud-Vie", position: "bud"},
			{order: "A Boh-Vie", position: "vie"},
			{order: "A Tyr S A Boh-Vie", position: "tyr"},
		},
	},
	{
		description: "given a counterattack and a supported second attack, where one counterattack party has support, " +
			"then all units bounce",
		orders: []*result{
			{order: "A Vie-Bud", position: "vie"},
			{order: "A Bud-Vie", position: "bud"},
			{order: "A Sil S A Bud-Vie", position: "sil"},
			{order: "A Boh-Vie", position: "boh"},
			{order: "A Tyr S A Boh-Vie", position: "tyr"},
		},
	},
	{
		description: "given a unit holds and another unit supports holding unit, then both units remain in position",
		orders: []*result{
			{order: "A Vie H", position: "vie"},
			{order: "A Bud S A Vie", position: "bud"},
		},
	},
	{
		description: "given a unit moves to a non-contiguous territory, then the move will be invalid",
		orders: []*result{
			{order: "A Vie-Lon", position: "vie"},
		},
	},
	{
		description: "given a supported unit holds and is attacked by unit with equal strength, " +
			"then attacking unit bounces",
		orders: []*result{
			{order: "A Vie H", position: "vie"},
			{order: "A Bud S A Vie", position: "bud"},
			{order: "A Boh-Vie", position: "boh"},
			{order: "A Tyr S A Boh-Vie", position: "tyr"},
		},
	},
	{
		description: "given a supported attack where the support cutter is itself bounced, " +
			"then support is still cut and attack fails",
		orders: []*result{
			{order: "A Boh-Vie", position: "boh"},
			{order: "A Gal S A Boh-Vie", position: "gal"},
			{order: "A Vie H", position: "vie"},
			{order: "A Bud-Gal", position: "bud"},
			{order: "A Sil-Gal", position: "sil"},
		},
	},
	{
		description: "given a supported attack where both support cutters tie at the supporter's territory, " +
			"then support is still cut and attack fails",
		orders: []*result{
			{order: "A Boh-Gal", position: "boh"},
			{order: "A Vie S A Boh-Gal", position: "vie"},
			{order: "A Gal H", position: "gal"},
			{order: "A Bud-Vie", position: "bud"},
			{order: "A Tri-Vie", position: "tri"},
		},
	},
	{
		description: "given a lone attack on a supporter that itself bounces, " +
			"then support is still cut and attack on supported territory fails",
		orders: []*result{
			{order: "A Gal-Vie", position: "gal"},
			{order: "A Boh S A Gal-Vie", position: "boh"},
			{order: "A Vie H", position: "vie"},
			{order: "A Mun-Boh", position: "mun"},
		},
	},

	// ===== DATC 6.A. BASIC CHECKS =====
	//
	// 6.A.1 through 6.A.10 require the fleet model (sea territories, fleet
	// adjacency graph) or army-to-sea validation. Uncomment when implementing
	// Phase 4 (fleet model). See CLAUDE.md.

	/*
		// DATC 6.A.1. MOVING TO AN AREA THAT IS NOT A NEIGHBOUR
		// Needs: fleet model + sea territories (nth, pic adjacency check)
		// England: F North Sea - Picardy → order should fail; fleet stays in nth
		{
			description: "DATC 6.A.1 - moving to an area that is not a neighbour",
			orders: []*result{
				{order: "F Nth-Pic", position: "nth"},
			},
		},
	*/

	/*
		// DATC 6.A.2. MOVE ARMY TO SEA
		// Needs: sea territories and army-to-sea validation
		// England: A Liverpool - Irish Sea → order should fail
		{
			description: "DATC 6.A.2 - move army to sea",
			orders: []*result{
				{order: "A Lvp-Iri", position: "lvp"},
			},
		},
	*/

	/*
		// DATC 6.A.3. MOVE FLEET TO LAND
		// Needs: fleet model; fleet cannot enter land territory
		// Germany: F Kiel - Munich → order should fail
		{
			description: "DATC 6.A.3 - move fleet to land",
			orders: []*result{
				{order: "F Kie-Mun", position: "kie"},
			},
		},
	*/

	/*
		// DATC 6.A.4. MOVE TO OWN SECTOR
		// Needs: fleet model; moving to current territory is illegal
		// Germany: F Kiel - Kiel → program should not crash
		{
			description: "DATC 6.A.4 - move to own sector",
			orders: []*result{
				{order: "F Kie-Kie", position: "kie"},
			},
		},
	*/

	/*
		// DATC 6.A.5. MOVE TO OWN SECTOR WITH CONVOY
		// Needs: fleet model + convoy model
		// England: F North Sea Convoys A Yorkshire - Yorkshire (illegal)
		//          A Yorkshire - Yorkshire, A Liverpool Supports A Yorkshire - Yorkshire
		// Germany: F London - Yorkshire, A Wales Supports F London - Yorkshire
		// Result: Yorkshire's self-move is illegal; German dislodges Yorkshire army.
		{
			description: "DATC 6.A.5 - move to own sector with convoy",
			orders: []*result{
				{order: "F Nth C A Yor-Yor", position: "nth"},
				{order: "A Yor-Yor", position: "yor", defeated: true},
				{order: "A Lvp S A Yor-Yor", position: "lvp"},
				{order: "F Lon-Yor", position: "yor"},
				{order: "A Wal S F Lon-Yor", position: "wal"},
			},
		},
	*/

	/*
		// DATC 6.A.6. ORDERING A UNIT OF ANOTHER COUNTRY
		// Needs: multi-country validation; Germany cannot order England's fleet
		// England has F London. Germany: F London - North Sea → order should fail
		{
			description: "DATC 6.A.6 - ordering a unit of another country",
			orders: []*result{
				{order: "F Lon-Nth", position: "lon"}, // Germany's illegal order; England's fleet stays
			},
		},
	*/

	/*
		// DATC 6.A.7. ONLY ARMIES CAN BE CONVOYED
		// Needs: fleet model + convoy validation
		// England: F London - Belgium (fleet cannot be convoyed), F North Sea Convoys A London - Belgium
		// Move from London to Belgium should fail.
		{
			description: "DATC 6.A.7 - only armies can be convoyed",
			orders: []*result{
				{order: "F Lon-Bel", position: "lon"},
				{order: "F Nth C A Lon-Bel", position: "nth"},
			},
		},
	*/

	/*
		// DATC 6.A.8. SUPPORT TO HOLD YOURSELF IS NOT POSSIBLE
		// Needs: self-support validation (a unit cannot support itself holding)
		// Italy: A Venice - Trieste, A Tyrolia Supports A Venice - Trieste
		// Austria: F Trieste Supports F Trieste (illegal self-support)
		// Result: Trieste is dislodged (self-support adds 0 strength).
		// Note: with army-only engine and no self-support check, this test passes
		// already — Trieste's "support" of itself is simply ignored by the engine.
		// Uncomment once self-support order validation is wired in.
		{
			description: "DATC 6.A.8 - support to hold yourself is not possible",
			orders: []*result{
				{order: "A Ven-Tri", position: "tri"},
				{order: "A Tyr S A Ven-Tri", position: "tyr"},
				{order: "A Tri S A Tri", position: "tri", defeated: true},
			},
		},
	*/

	/*
		// DATC 6.A.9. FLEETS MUST FOLLOW COAST IF NOT ON SEA
		// Needs: fleet model; fleet adjacency differs from army adjacency
		// Italy: F Rome - Venice → move fails (fleet cannot go Rome→Venice by coast)
		{
			description: "DATC 6.A.9 - fleets must follow coast if not on sea",
			orders: []*result{
				{order: "F Rom-Ven", position: "rom"},
			},
		},
	*/

	/*
		// DATC 6.A.10. SUPPORT ON UNREACHABLE DESTINATION NOT POSSIBLE
		// Needs: fleet model; fleet adjacency validation on support orders
		// Austria: A Venice Hold
		// Italy: F Rome Supports A Apulia - Venice (Rome can't reach Venice by fleet), A Apulia - Venice
		// Result: support of Rome is illegal; Venice is not dislodged.
		{
			description: "DATC 6.A.10 - support on unreachable destination not possible",
			orders: []*result{
				{order: "A Ven H", position: "ven"},
				{order: "F Rom S A Apu-Ven", position: "rom"},
				{order: "A Apu-Ven", position: "apu"},
			},
		},
	*/

	// DATC 6.A.11. SIMPLE BOUNCE
	// Two armies move to the same empty territory; both bounce.
	{
		description: "DATC 6.A.11 - simple bounce",
		orders: []*result{
			{order: "A Vie-Tyr", position: "vie"},
			{order: "A Ven-Tyr", position: "ven"},
		},
	},

	// DATC 6.A.12. BOUNCE OF THREE UNITS
	// Three armies all move to the same territory; all three bounce.
	{
		description: "DATC 6.A.12 - bounce of three units",
		orders: []*result{
			{order: "A Vie-Tyr", position: "vie"},
			{order: "A Mun-Tyr", position: "mun"},
			{order: "A Ven-Tyr", position: "ven"},
		},
	},

	// ===== DATC 6.B. COASTAL ISSUES =====
	// All 6.B tests require the fleet model with coast designators.
	// Uncomment when implementing Phase 4. See CLAUDE.md.

	/*
		// DATC 6.B.1. MOVING WITH UNSPECIFIED COAST WHEN COAST IS NECESSARY
		// Needs: fleet model + coast designators on Spain (nc/sc)
		// France: F Portugal - Spain → move should fail (coast required)
		{
			description: "DATC 6.B.1 - moving with unspecified coast when coast is necessary",
			orders: []*result{
				{order: "F Por-Spa", position: "por"},
			},
		},
	*/

	/*
		// DATC 6.B.2. MOVING WITH UNSPECIFIED COAST WHEN COAST IS NOT NECESSARY
		// Needs: fleet model + coast designators
		// France: F Gascony - Spain → north coast is the only reachable coast; move succeeds
		{
			description: "DATC 6.B.2 - moving with unspecified coast when coast is not necessary",
			orders: []*result{
				{order: "F Gas-Spa", position: "spa"},
			},
		},
	*/

	/*
		// DATC 6.B.3. MOVING WITH WRONG COAST WHEN COAST IS NOT NECESSARY
		// Needs: fleet model + coast designators
		// France: F Gascony - Spain(sc) → wrong coast; order is illegal; fleet holds
		{
			description: "DATC 6.B.3 - moving with wrong coast when coast is not necessary",
			orders: []*result{
				{order: "F Gas-Spa(sc)", position: "gas"},
			},
		},
	*/

	/*
		// DATC 6.B.4. SUPPORT TO UNREACHABLE COAST ALLOWED
		// Needs: fleet model + coast designators
		// France: F Gascony - Spain(nc), F Marseilles Supports F Gascony - Spain(nc)
		// Italy: F Western Mediterranean - Spain(sc)
		// Result: Gascony moves to Spain(nc); Italian fleet fails.
		{
			description: "DATC 6.B.4 - support to unreachable coast allowed",
			orders: []*result{
				{order: "F Gas-Spa(nc)", position: "spa"},
				{order: "F Mar S F Gas-Spa(nc)", position: "mar"},
				{order: "F Wme-Spa(sc)", position: "wme"},
			},
		},
	*/

	/*
		// DATC 6.B.5. SUPPORT FROM UNREACHABLE COAST NOT ALLOWED
		// Needs: fleet model + coast designators
		// France: F Marseilles - Gulf of Lyon, F Spain(nc) Supports F Marseilles - Gulf of Lyon
		// Italy: F Gulf of Lyon Hold
		// Result: Spain(nc) cannot reach Gulf of Lyon; support illegal; Gulf of Lyon not dislodged.
		{
			description: "DATC 6.B.5 - support from unreachable coast not allowed",
			orders: []*result{
				{order: "F Mar-Gol", position: "mar"},
				{order: "F Spa(nc) S F Mar-Gol", position: "spa"},
				{order: "F Gol H", position: "gol"},
			},
		},
	*/

	/*
		// DATC 6.B.6. SUPPORT CAN BE CUT WITH OTHER COAST
		// Needs: fleet model + coast designators
		// England: F Irish Sea Supports F North Atlantic - Mid-Atlantic, F North Atlantic - Mid-Atlantic
		// France: F Spain(nc) Supports F Mid-Atlantic, F Mid-Atlantic Hold
		// Italy: F Gulf of Lyon - Spain(sc)
		// Result: Italian fleet cuts Spanish support; French Mid-Atlantic dislodged.
		{
			description: "DATC 6.B.6 - support can be cut with other coast",
			orders: []*result{
				{order: "F Iri S F Nao-Mao", position: "iri"},
				{order: "F Nao-Mao", position: "mao"},
				{order: "F Spa(nc) S F Mao", position: "spa"},
				{order: "F Mao H", position: "mao", defeated: true},
				{order: "F Gol-Spa(sc)", position: "spa"},
			},
		},
	*/

	/*
		// DATC 6.B.7-6.B.15: Various coastal edge cases.
		// All need fleet model + coast designators. Uncomment for Phase 4.
	*/

	// ===== DATC 6.C. CIRCULAR MOVEMENT =====

	// DATC 6.C.1. THREE ARMY CIRCULAR MOVEMENT
	// Three units rotate through a triangle; all three move.
	// Note: original DATC uses F Ankara; army adjacency is identical so this
	// tests the same resolution logic.
	{
		description: "DATC 6.C.1 - three army circular movement",
		orders: []*result{
			{order: "A Ank-Con", position: "con"},
			{order: "A Con-Smy", position: "smy"},
			{order: "A Smy-Ank", position: "ank"},
		},
	},

	// DATC 6.C.2. THREE ARMY CIRCULAR MOVEMENT WITH SUPPORT
	// One unit in the rotation receives support; all three still move.
	{
		description: "DATC 6.C.2 - three army circular movement with support",
		orders: []*result{
			{order: "A Ank-Con", position: "con"},
			{order: "A Con-Smy", position: "smy"},
			{order: "A Smy-Ank", position: "ank"},
			{order: "A Bul S A Ank-Con", position: "bul"},
		},
	},

	// DATC 6.C.3. A DISRUPTED THREE ARMY CIRCULAR MOVEMENT
	// A fourth unit attacks a destination in the rotation; none of the
	// rotating units move.
	{
		description: "DATC 6.C.3 - a disrupted three army circular movement",
		orders: []*result{
			{order: "A Ank-Con", position: "ank"},
			{order: "A Con-Smy", position: "con"},
			{order: "A Smy-Ank", position: "smy"},
			{order: "A Bul-Con", position: "bul"},
		},
	},

	/*
		// DATC 6.C.4. A CIRCULAR MOVEMENT WITH ATTACKED CONVOY
		// Needs: convoy model
		// Austria: A Trieste - Serbia, A Serbia - Bulgaria
		// Turkey: A Bulgaria - Trieste (via convoy), F Aegean/Ionian/Adriatic Convoys
		// Italy: F Naples - Ionian Sea (attacks convoy fleet but does not dislodge)
		// Result: circular movement succeeds; all three armies advance.
		{
			description: "DATC 6.C.4 - circular movement with attacked convoy",
			orders: []*result{
				{order: "A Tri-Ser", position: "ser"},
				{order: "A Ser-Bul", position: "bul"},
				{order: "A Bul-Tri", position: "tri"},
				{order: "F Aeg C A Bul-Tri", position: "aeg"},
				{order: "F Ion C A Bul-Tri", position: "ion"},
				{order: "F Adr C A Bul-Tri", position: "adr"},
				{order: "F Nap-Ion", position: "nap"},
			},
		},
	*/

	/*
		// DATC 6.C.5. A DISRUPTED CIRCULAR MOVEMENT DUE TO DISLODGED CONVOY
		// Needs: convoy model
		// Same as 6.C.4 but Italy adds F Tunis Supports F Naples - Ionian Sea.
		// The Ionian convoy fleet is dislodged; circular movement fails; all armies stay.
		{
			description: "DATC 6.C.5 - disrupted circular movement due to dislodged convoy",
			orders: []*result{
				{order: "A Tri-Ser", position: "tri"},
				{order: "A Ser-Bul", position: "ser"},
				{order: "A Bul-Tri", position: "bul"},
				{order: "F Aeg C A Bul-Tri", position: "aeg"},
				{order: "F Ion C A Bul-Tri", position: "ion", defeated: true},
				{order: "F Adr C A Bul-Tri", position: "adr"},
				{order: "F Nap-Ion", position: "ion"},
				{order: "F Tun S F Nap-Ion", position: "tun"},
			},
		},
	*/

	/*
		// DATC 6.C.6. TWO ARMIES WITH TWO CONVOYS
		// Needs: convoy model
		// England: F North Sea Convoys A London - Belgium, A London - Belgium
		// France: F English Channel Convoys A Belgium - London, A Belgium - London
		// Both convoys succeed; armies swap.
		{
			description: "DATC 6.C.6 - two armies with two convoys",
			orders: []*result{
				{order: "F Nth C A Lon-Bel", position: "nth"},
				{order: "A Lon-Bel", position: "bel"},
				{order: "F Ech C A Bel-Lon", position: "ech"},
				{order: "A Bel-Lon", position: "lon"},
			},
		},
	*/

	/*
		// DATC 6.C.7. DISRUPTED UNIT SWAP
		// Needs: convoy model
		// Same as 6.C.6 but France adds A Burgundy - Belgium; the swap is disrupted
		// and neither army moves.
		{
			description: "DATC 6.C.7 - disrupted unit swap",
			orders: []*result{
				{order: "F Nth C A Lon-Bel", position: "nth"},
				{order: "A Lon-Bel", position: "lon"},
				{order: "F Ech C A Bel-Lon", position: "ech"},
				{order: "A Bel-Lon", position: "bel"},
				{order: "A Bur-Bel", position: "bur"},
			},
		},
	*/

	/*
		// DATC 6.C.8. NO SELF DISLODGEMENT IN DISRUPTED CIRCULAR MOVEMENT
		// Needs: fleet model + country/self-dislodgement rules
		// Turkey: F Constantinople - Black Sea, A Bulgaria - Constantinople,
		//         A Smyrna Supports A Bulgaria - Constantinople
		// Russia: F Black Sea - Bulgaria(ec)
		// Austria: A Serbia - Bulgaria
		// Result: none of the units move (self-dislodgement blocked).
		{
			description: "DATC 6.C.8 - no self dislodgement in disrupted circular movement",
			orders: []*result{
				{order: "F Con-Bla", position: "con"},
				{order: "A Bul-Con", position: "bul"},
				{order: "A Smy S A Bul-Con", position: "smy"},
				{order: "F Bla-Bul(ec)", position: "bla"},
				{order: "A Ser-Bul", position: "ser"},
			},
		},
	*/

	/*
		// DATC 6.C.9. NO HELP IN DISLODGEMENT OF OWN UNIT IN DISRUPTED CIRCULAR MOVEMENT
		// Needs: fleet model + country/self-dislodgement rules
		// Turkey: F Constantinople - Black Sea, A Smyrna Supports A Bulgaria - Constantinople
		// Russia: F Black Sea - Bulgaria(ec)
		// Austria: A Serbia - Bulgaria, A Bulgaria - Constantinople
		// Result: none of the units move.
		{
			description: "DATC 6.C.9 - no help in dislodgement of own unit in disrupted circular movement",
			orders: []*result{
				{order: "F Con-Bla", position: "con"},
				{order: "A Smy S A Bul-Con", position: "smy"},
				{order: "F Bla-Bul(ec)", position: "bla"},
				{order: "A Ser-Bul", position: "ser"},
				{order: "A Bul-Con", position: "bul"},
			},
		},
	*/

	// ===== DATC 6.D. SUPPORTS AND DISLODGES =====

	// DATC 6.D.1. SUPPORTED HOLD CAN PREVENT DISLODGEMENT
	// Italy: A Apulia Supports A Trieste - Venice, A Trieste - Venice
	// Austria: A Venice Hold, A Tyrolia Supports A Venice
	// Note: original DATC uses F Adriatic Sea instead of A Apulia.
	// Result: supported hold (str 1) ties with supported attack (str 1); Trieste bounces.
	{
		description: "DATC 6.D.1 - supported hold can prevent dislodgement",
		orders: []*result{
			{order: "A Apu S A Tri-Ven", position: "apu"},
			{order: "A Tri-Ven", position: "tri"},
			{order: "A Ven H", position: "ven"},
			{order: "A Tyr S A Ven", position: "tyr"},
		},
	},

	// DATC 6.D.2. A MOVE CUTS SUPPORT ON HOLD
	// Italy: A Apulia Supports A Trieste - Venice, A Trieste - Venice, A Vienna - Tyrolia
	// Austria: A Venice Hold, A Tyrolia Supports A Venice
	// Note: original uses F Adriatic Sea instead of A Apulia.
	// Result: Vienna cuts Tyrolia's hold support; Venice dislodged.
	{
		description: "DATC 6.D.2 - a move cuts support on hold",
		orders: []*result{
			{order: "A Apu S A Tri-Ven", position: "apu"},
			{order: "A Tri-Ven", position: "ven"},
			{order: "A Vie-Tyr", position: "vie"},
			{order: "A Ven H", position: "ven", defeated: true},
			{order: "A Tyr S A Ven", position: "tyr"},
		},
	},

	// DATC 6.D.3. A MOVE CUTS SUPPORT ON MOVE
	// Italy: A Apulia Supports A Trieste - Venice, A Trieste - Venice
	// Austria: A Venice Hold, A Naples - Apulia (cuts A Apulia's move support)
	// Note: original uses F Adriatic Sea (supporter) and F Ionian Sea (cutter).
	// Result: Apulia's support is cut; Venice is not dislodged; Naples bounces.
	{
		description: "DATC 6.D.3 - a move cuts support on move",
		orders: []*result{
			{order: "A Apu S A Tri-Ven", position: "apu"},
			{order: "A Tri-Ven", position: "tri"},
			{order: "A Ven H", position: "ven"},
			{order: "A Nap-Apu", position: "nap"},
		},
	},

	// DATC 6.D.4. SUPPORT TO HOLD ON UNIT SUPPORTING A HOLD ALLOWED
	// Germany: A Berlin Supports A Kiel (hold), A Kiel Supports A Berlin (hold)
	// Russia: A Warsaw Supports A Prussia - Berlin, A Prussia - Berlin
	// Note: original uses F Kiel and F Baltic Sea.
	// Result: Berlin's hold support from Kiel is not cut; Russian attack bounces.
	{
		description: "DATC 6.D.4 - support to hold on unit supporting a hold allowed",
		orders: []*result{
			{order: "A Ber S A Kie", position: "ber"},
			{order: "A Kie S A Ber", position: "kie"},
			{order: "A War S A Pru-Ber", position: "war"},
			{order: "A Pru-Ber", position: "pru"},
		},
	},

	// DATC 6.D.5. SUPPORT TO HOLD ON UNIT SUPPORTING A MOVE ALLOWED
	// Germany: A Berlin Supports A Munich - Silesia, A Kiel Supports A Berlin, A Munich - Silesia
	// Russia: A Warsaw Supports A Prussia - Berlin, A Prussia - Berlin
	// Note: original uses F Kiel and F Baltic Sea.
	// Result: Berlin's hold support from Kiel is not cut; Russian attack bounces.
	{
		description: "DATC 6.D.5 - support to hold on unit supporting a move allowed",
		orders: []*result{
			{order: "A Ber S A Mun-Sil", position: "ber"},
			{order: "A Kie S A Ber", position: "kie"},
			{order: "A Mun-Sil", position: "sil"},
			{order: "A War S A Pru-Ber", position: "war"},
			{order: "A Pru-Ber", position: "pru"},
		},
	},

	/*
		// DATC 6.D.6. SUPPORT TO HOLD ON CONVOYING UNIT ALLOWED
		// Needs: convoy model
		// Germany: A Berlin - Sweden (via convoy), F Baltic Sea Convoys A Berlin - Sweden,
		//          F Prussia Supports F Baltic Sea
		// Russia: F Livonia - Baltic Sea, F Gulf of Bothnia Supports F Livonia - Baltic Sea
		// Result: Baltic Sea not dislodged; convoy succeeds.
		{
			description: "DATC 6.D.6 - support to hold on convoying unit allowed",
			orders: []*result{
				{order: "A Ber-Swe", position: "swe"},
				{order: "F Bal C A Ber-Swe", position: "bal"},
				{order: "F Pru S F Bal", position: "pru"},
				{order: "F Lvn-Bal", position: "lvn"},
				{order: "F Gob S F Lvn-Bal", position: "gob"},
			},
		},
	*/

	/*
		// DATC 6.D.7. SUPPORT TO HOLD ON MOVING UNIT NOT ALLOWED
		// Needs: fleet model (all units are fleets in original)
		// Germany: F Baltic Sea - Sweden, F Prussia Supports F Baltic Sea
		// Russia: F Livonia - Baltic Sea, F Gulf of Bothnia Supports F Livonia - Baltic Sea,
		//         A Finland - Sweden
		// Result: Prussia's support for the moving Baltic Sea is invalid; Baltic dislodged.
		{
			description: "DATC 6.D.7 - support to hold on moving unit not allowed",
			orders: []*result{
				{order: "F Bal-Swe", position: "bal", defeated: true},
				{order: "F Pru S F Bal", position: "pru"},
				{order: "F Lvn-Bal", position: "bal"},
				{order: "F Gob S F Lvn-Bal", position: "gob"},
				{order: "A Fin-Swe", position: "swe"},
			},
		},
	*/

	/*
		// DATC 6.D.8. FAILED CONVOY CANNOT RECEIVE HOLD SUPPORT
		// Needs: convoy model
		// Austria: F Ionian Sea Hold, A Serbia Supports A Albania - Greece, A Albania - Greece
		// Turkey: A Greece - Naples (convoy attempt), A Bulgaria Supports A Greece
		// Result: convoy fails; Greece cannot receive hold support; Albania dislodges Greece.
		{
			description: "DATC 6.D.8 - failed convoy cannot receive hold support",
			orders: []*result{
				{order: "F Ion H", position: "ion"},
				{order: "A Ser S A Alb-Gre", position: "ser"},
				{order: "A Alb-Gre", position: "gre"},
				{order: "A Gre-Nap", position: "gre", defeated: true},
				{order: "A Bul S A Gre", position: "bul"},
			},
		},
	*/

	// DATC 6.D.9. SUPPORT TO MOVE ON HOLDING UNIT NOT ALLOWED
	// Italy: A Venice - Trieste, A Tyrolia Supports A Venice - Trieste
	// Austria: A Albania Supports A Trieste - Serbia (invalid: Trieste is holding),
	//          A Trieste Hold
	// Result: Albanian support is ignored (no matching move); Venice dislodges Trieste.
	{
		description: "DATC 6.D.9 - support to move on holding unit not allowed",
		orders: []*result{
			{order: "A Ven-Tri", position: "tri"},
			{order: "A Tyr S A Ven-Tri", position: "tyr"},
			{order: "A Alb S A Tri-Ser", position: "alb"},
			{order: "A Tri H", position: "tri", defeated: true},
		},
	},

	/*
		// DATC 6.D.10. SELF DISLODGEMENT PROHIBITED
		// Needs: country/self-dislodgement rules
		// Germany: A Berlin Hold, F Kiel - Berlin, A Munich Supports F Kiel - Berlin
		// Result: Germany cannot dislodge its own unit; move fails.
		{
			description: "DATC 6.D.10 - self dislodgement prohibited",
			orders: []*result{
				{order: "A Ber H", position: "ber"},
				{order: "F Kie-Ber", position: "kie"},
				{order: "A Mun S F Kie-Ber", position: "mun"},
			},
		},
	*/

	/*
		// DATC 6.D.11. NO SELF DISLODGEMENT OF RETURNING UNIT
		// Needs: country/self-dislodgement rules
		// Germany: A Berlin - Prussia, F Kiel - Berlin, A Munich Supports F Kiel - Berlin
		// Russia: A Warsaw - Prussia
		// Result: Berlin bounces with Warsaw; Berlin is not dislodged by Kiel.
		{
			description: "DATC 6.D.11 - no self dislodgement of returning unit",
			orders: []*result{
				{order: "A Ber-Pru", position: "ber"},
				{order: "F Kie-Ber", position: "kie"},
				{order: "A Mun S F Kie-Ber", position: "mun"},
				{order: "A War-Pru", position: "war"},
			},
		},
	*/

	/*
		// DATC 6.D.12. SUPPORTING A FOREIGN UNIT TO DISLODGE OWN UNIT PROHIBITED
		// Needs: country/self-dislodgement rules
		// Austria: F Trieste Hold, A Vienna Supports A Venice - Trieste
		// Italy: A Venice - Trieste
		// Result: Austria cannot help Italy dislodge its own fleet; Trieste not dislodged.
		{
			description: "DATC 6.D.12 - supporting a foreign unit to dislodge own unit prohibited",
			orders: []*result{
				{order: "A Tri H", position: "tri"},
				{order: "A Vie S A Ven-Tri", position: "vie"},
				{order: "A Ven-Tri", position: "ven"},
			},
		},
	*/

	/*
		// DATC 6.D.13. SUPPORTING A FOREIGN UNIT TO DISLODGE A RETURNING OWN UNIT PROHIBITED
		// Needs: country/self-dislodgement rules + fleet model
		// Austria: F Trieste - Adriatic Sea, A Vienna Supports A Venice - Trieste
		// Italy: A Venice - Trieste, F Apulia - Adriatic Sea
		// Result: Trieste bounces with Apulia; not dislodged.
		{
			description: "DATC 6.D.13 - supporting foreign unit to dislodge returning own unit prohibited",
			orders: []*result{
				{order: "F Tri-Adr", position: "tri"},
				{order: "A Vie S A Ven-Tri", position: "vie"},
				{order: "A Ven-Tri", position: "ven"},
				{order: "F Apu-Adr", position: "apu"},
			},
		},
	*/

	// DATC 6.D.14. SUPPORTING A FOREIGN UNIT IS NOT ENOUGH TO PREVENT DISLODGEMENT
	// Austria: A Trieste Hold, A Vienna Supports A Venice - Trieste (Austria helps Italy)
	// Italy: A Venice - Trieste, A Tyrolia Supports A Venice - Trieste,
	//        A Albania Supports A Venice - Trieste
	// Note: original uses F Trieste and F Adriatic Sea; army adaptation uses A Albania.
	// Result: 3 supports > 0 hold strength; Trieste dislodged despite Austria's help.
	{
		description: "DATC 6.D.14 - supporting a foreign unit is not enough to prevent dislodgement",
		orders: []*result{
			{order: "A Tri H", position: "tri", defeated: true},
			{order: "A Vie S A Ven-Tri", position: "vie"},
			{order: "A Ven-Tri", position: "tri"},
			{order: "A Tyr S A Ven-Tri", position: "tyr"},
			{order: "A Alb S A Ven-Tri", position: "alb"},
		},
	},

	// DATC 6.D.15. DEFENDER CANNOT CUT SUPPORT FOR ATTACK ON ITSELF
	// Russia: A Galicia Supports A Vienna - Budapest (attacking Bud)
	// Turkey: A Budapest - Galicia (defender attacks the supporter)
	// Note: original uses Constantinople/Black Sea/Ankara (fleets). Army adaptation
	// with Gal/Vie/Bud gives the same result.
	// Result: Bud's attack on Gal does not cut Gal's support; Budapest dislodged.
	{
		description: "DATC 6.D.15 - defender cannot cut support for attack on itself",
		orders: []*result{
			{order: "A Gal S A Vie-Bud", position: "gal"},
			{order: "A Vie-Bud", position: "bud"},
			{order: "A Bud-Gal", position: "bud", defeated: true},
		},
	},

	/*
		// DATC 6.D.16. CONVOYING A UNIT DISLODGING A UNIT OF SAME POWER IS ALLOWED
		// Needs: convoy model
		// England: A London Hold, F North Sea Convoys A Belgium - London
		// France: F English Channel Supports A Belgium - London, A Belgium - London
		// Result: English army dislodged by French army via convoy.
		{
			description: "DATC 6.D.16 - convoying a unit dislodging a unit of same power is allowed",
			orders: []*result{
				{order: "A Lon H", position: "lon", defeated: true},
				{order: "F Nth C A Bel-Lon", position: "nth"},
				{order: "F Ech S A Bel-Lon", position: "ech"},
				{order: "A Bel-Lon", position: "lon"},
			},
		},
	*/

	/*
		// DATC 6.D.17. DISLODGEMENT CUTS SUPPORTS
		// Needs: fleet model + support recalculation after dislodgement
		// Russia: F Constantinople Supports F Black Sea - Ankara, F Black Sea - Ankara
		// Turkey: F Ankara - Constantinople, A Smyrna Supports F Ankara - Constantinople,
		//         A Armenia - Ankara
		// Result: Constantinople dislodged → support cut → Black Sea bounces with Armenia.
		{
			description: "DATC 6.D.17 - dislodgement cuts supports",
			orders: []*result{
				{order: "F Con S F Bla-Ank", position: "con", defeated: true},
				{order: "F Bla-Ank", position: "bla"},
				{order: "F Ank-Con", position: "con"},
				{order: "A Smy S F Ank-Con", position: "smy"},
				{order: "A Arm-Ank", position: "ank"},
			},
		},
	*/

	/*
		// DATC 6.D.18. A SURVIVING UNIT WILL SUSTAIN SUPPORT
		// Needs: fleet model + support recalculation
		// Same as 6.D.17 but Russia adds A Bulgaria Supports F Constantinople.
		// Result: Constantinople survives → support holds → Black Sea dislodges Ankara.
		{
			description: "DATC 6.D.18 - a surviving unit will sustain support",
			orders: []*result{
				{order: "F Con S F Bla-Ank", position: "con"},
				{order: "F Bla-Ank", position: "ank"},
				{order: "A Bul S F Con", position: "bul"},
				{order: "F Ank-Con", position: "ank", defeated: true},
				{order: "A Smy S F Ank-Con", position: "smy"},
				{order: "A Arm-Ank", position: "arm"},
			},
		},
	*/

	/*
		// DATC 6.D.19. EVEN WHEN SURVIVING IS IN ALTERNATIVE WAY
		// Needs: fleet model + support recalculation + country rules
		// Russia: F Constantinople Supports F Black Sea - Ankara, F Black Sea - Ankara,
		//         A Smyrna Supports F Ankara - Constantinople
		// Turkey: F Ankara - Constantinople
		// Result: Constantinople not dislodged (Russian unit supports Ankara-Con, blocked);
		//         Ankara dislodged.
		{
			description: "DATC 6.D.19 - even when surviving is in alternative way",
			orders: []*result{
				{order: "F Con S F Bla-Ank", position: "con"},
				{order: "F Bla-Ank", position: "ank"},
				{order: "A Smy S F Ank-Con", position: "smy"},
				{order: "F Ank-Con", position: "ank", defeated: true},
			},
		},
	*/

	/*
		// DATC 6.D.20. UNIT CANNOT CUT SUPPORT OF ITS OWN COUNTRY
		// Needs: country rules (same-power support-cut exemption) + fleet model
		// England: F London Supports F North Sea - English Channel,
		//          F North Sea - English Channel, A Yorkshire - London
		// France: F English Channel Hold
		// Result: Yorkshire (England) does not cut London (England)'s support;
		//         English Channel dislodged.
		{
			description: "DATC 6.D.20 - unit cannot cut support of its own country",
			orders: []*result{
				{order: "F Lon S F Nth-Ech", position: "lon"},
				{order: "F Nth-Ech", position: "ech"},
				{order: "A Yor-Lon", position: "yor"},
				{order: "F Ech H", position: "ech", defeated: true},
			},
		},
	*/

	// DATC 6.D.21. DISLODGING DOES NOT CANCEL A SUPPORT CUT
	// Austria: A Trieste Hold
	// Italy: A Venice - Trieste, A Tyrolia Supports A Venice - Trieste
	// Germany: A Munich - Tyrolia (cuts Italy's support even though Munich is dislodged)
	// Russia: A Silesia - Munich, A Berlin Supports A Silesia - Munich
	// Result: Munich cuts Tyrolia's support before being dislodged; Venice bounces; Trieste safe.
	{
		description: "DATC 6.D.21 - dislodging does not cancel a support cut",
		orders: []*result{
			{order: "A Tri H", position: "tri"},
			{order: "A Ven-Tri", position: "ven"},
			{order: "A Tyr S A Ven-Tri", position: "tyr"},
			{order: "A Mun-Tyr", position: "mun", defeated: true},
			{order: "A Sil-Mun", position: "mun"},
			{order: "A Ber S A Sil-Mun", position: "ber"},
		},
	},

	/*
		// DATC 6.D.22. IMPOSSIBLE FLEET MOVE CANNOT BE SUPPORTED
		// Needs: fleet model + illegal-move support invalidation
		// Germany: F Kiel - Munich (illegal), A Burgundy Supports F Kiel - Munich
		// Russia: A Munich - Kiel, A Berlin Supports A Munich - Kiel
		// Result: Kiel's illegal move makes Burgundy's support invalid; Munich dislodges Kiel.
		{
			description: "DATC 6.D.22 - impossible fleet move cannot be supported",
			orders: []*result{
				{order: "F Kie-Mun", position: "kie", defeated: true},
				{order: "A Bur S F Kie-Mun", position: "bur"},
				{order: "A Mun-Kie", position: "kie"},
				{order: "A Ber S A Mun-Kie", position: "ber"},
			},
		},
	*/

	/*
		// DATC 6.D.23. IMPOSSIBLE COAST MOVE CANNOT BE SUPPORTED
		// Needs: fleet model + coast designators
		// Italy: F Gulf of Lyon - Spain(sc), F Western Med Supports F Gulf - Spain(sc)
		// France: F Spain(nc) - Gulf of Lyon, F Marseilles Supports F Spain(nc) - Gulf of Lyon
		// Result: French move is illegal (wrong coast); Marseilles support fails; Spain dislodged.
		{
			description: "DATC 6.D.23 - impossible coast move cannot be supported",
			orders: []*result{
				{order: "F Gol-Spa(sc)", position: "spa"},
				{order: "F Wme S F Gol-Spa(sc)", position: "wme"},
				{order: "F Spa(nc)-Gol", position: "spa", defeated: true},
				{order: "F Mar S F Spa(nc)-Gol", position: "mar"},
			},
		},
	*/

	/*
		// DATC 6.D.24. IMPOSSIBLE ARMY MOVE CANNOT BE SUPPORTED
		// Needs: army-to-sea validation + fleet model
		// France: A Marseilles - Gulf of Lyon (illegal), F Spain(sc) Supports A Marseilles - Gulf of Lyon
		// Italy: F Gulf of Lyon Hold
		// Turkey: F Tyrrhenian Sea Supports F Western Mediterranean - Gulf of Lyon,
		//         F Western Mediterranean - Gulf of Lyon
		// Result: French move is illegal; Spain's support fails; Gulf dislodged by Turkey.
		{
			description: "DATC 6.D.24 - impossible army move cannot be supported",
			orders: []*result{
				{order: "A Mar-Gol", position: "mar"},
				{order: "F Spa(sc) S A Mar-Gol", position: "spa"},
				{order: "F Gol H", position: "gol", defeated: true},
				{order: "F Tys S F Wme-Gol", position: "tys"},
				{order: "F Wme-Gol", position: "gol"},
			},
		},
	*/

	// DATC 6.D.25. FAILING HOLD SUPPORT CAN BE SUPPORTED
	// Germany: A Berlin Supports A Prussia (hold support that fails — Prussia is moving),
	//          A Kiel Supports A Berlin (hold support for Berlin — valid)
	// Russia: A Warsaw Supports A Prussia - Berlin, A Prussia - Berlin
	// Note: original uses F Kiel and F Baltic Sea.
	// Result: Berlin's hold support from Kiel is still valid despite Berlin's mismatched
	//         support order; Russian attack bounces.
	{
		description: "DATC 6.D.25 - failing hold support can be supported",
		orders: []*result{
			{order: "A Ber S A Pru", position: "ber"},
			{order: "A Kie S A Ber", position: "kie"},
			{order: "A War S A Pru-Ber", position: "war"},
			{order: "A Pru-Ber", position: "pru"},
		},
	},

	// DATC 6.D.26. FAILING MOVE SUPPORT CAN BE SUPPORTED
	// Germany: A Berlin Supports A Prussia - Silesia (mismatched: Prussia moves to Berlin),
	//          A Kiel Supports A Berlin
	// Russia: A Warsaw Supports A Prussia - Berlin, A Prussia - Berlin
	// Note: original uses F Kiel and F Baltic Sea.
	// Result: Berlin's hold support from Kiel is valid; Russian attack bounces.
	{
		description: "DATC 6.D.26 - failing move support can be supported",
		orders: []*result{
			{order: "A Ber S A Pru-Sil", position: "ber"},
			{order: "A Kie S A Ber", position: "kie"},
			{order: "A War S A Pru-Ber", position: "war"},
			{order: "A Pru-Ber", position: "pru"},
		},
	},

	/*
		// DATC 6.D.27. FAILING CONVOY CAN BE SUPPORTED
		// Needs: convoy model
		// England: F Sweden - Baltic Sea, F Denmark Supports F Sweden - Baltic Sea
		// Germany: A Berlin Hold
		// Russia: F Baltic Sea Convoys A Berlin - Livonia (unmatched convoy),
		//         F Prussia Supports F Baltic Sea
		// Result: Baltic's convoy is unmatched but support of Prussia is still valid;
		//         Baltic not dislodged.
		{
			description: "DATC 6.D.27 - failing convoy can be supported",
			orders: []*result{
				{order: "F Swe-Bal", position: "swe"},
				{order: "F Den S F Swe-Bal", position: "den"},
				{order: "A Ber H", position: "ber"},
				{order: "F Bal C A Ber-Lvn", position: "bal"},
				{order: "F Pru S F Bal", position: "pru"},
			},
		},
	*/

	/*
		// DATC 6.D.28. IMPOSSIBLE MOVE AND SUPPORT
		// Needs: fleet model + impossible-move validation
		// Austria: A Budapest Supports F Rumania (hold support for Rumania)
		// Russia: F Rumania - Holland (illegal: too far)
		// Turkey: F Black Sea - Rumania, A Bulgaria Supports F Black Sea - Rumania
		// Result: Rumania's illegal order ignored; Rumania holds with support; not dislodged.
		{
			description: "DATC 6.D.28 - impossible move and support",
			orders: []*result{
				{order: "A Bud S F Rum", position: "bud"},
				{order: "F Rum-Hol", position: "rum"},
				{order: "F Bla-Rum", position: "bla"},
				{order: "A Bul S F Bla-Rum", position: "bul"},
			},
		},
	*/

	/*
		// DATC 6.D.29. MOVE TO IMPOSSIBLE COAST AND SUPPORT
		// Needs: fleet model + coast designators
		// Austria: A Budapest Supports F Rumania
		// Russia: F Rumania - Bulgaria(sc) (impossible coast from Rumania)
		// Turkey: F Black Sea - Rumania, A Bulgaria Supports F Black Sea - Rumania
		// Result: same as 6.D.28; Rumania not dislodged.
		{
			description: "DATC 6.D.29 - move to impossible coast and support",
			orders: []*result{
				{order: "A Bud S F Rum", position: "bud"},
				{order: "F Rum-Bul(sc)", position: "rum"},
				{order: "F Bla-Rum", position: "bla"},
				{order: "A Bul S F Bla-Rum", position: "bul"},
			},
		},
	*/

	/*
		// DATC 6.D.30. MOVE WITHOUT COAST AND SUPPORT
		// Needs: fleet model + coast designators
		// Italy: F Aegean Sea Supports F Constantinople
		// Russia: F Constantinople - Bulgaria (coast unspecified; illegal)
		// Turkey: F Black Sea - Constantinople, A Bulgaria Supports F Black Sea - Constantinople
		// Result: Constantinople not dislodged.
		{
			description: "DATC 6.D.30 - move without coast and support",
			orders: []*result{
				{order: "F Aeg S F Con", position: "aeg"},
				{order: "F Con-Bul", position: "con"},
				{order: "F Bla-Con", position: "bla"},
				{order: "A Bul S F Bla-Con", position: "bul"},
			},
		},
	*/

	/*
		// DATC 6.D.31. A TRICKY IMPOSSIBLE SUPPORT
		// Needs: fleet model + convoy-route awareness for support validation
		// Austria: A Rumania - Armenia
		// Turkey: F Black Sea Supports A Rumania - Armenia (impossible: only route is via Black Sea,
		//         which can't convoy and support simultaneously)
		// Result: support is illegal and ignored.
		{
			description: "DATC 6.D.31 - a tricky impossible support",
			orders: []*result{
				{order: "A Rum-Arm", position: "rum"},
				{order: "F Bla S A Rum-Arm", position: "bla"},
			},
		},
	*/

	/*
		// DATC 6.D.32. A MISSING FLEET
		// Needs: fleet model + convoy-route validation for move support
		// England: F Edinburgh Supports A Liverpool - Yorkshire, A Liverpool - Yorkshire
		// France: F London Supports A Yorkshire
		// Germany: A Yorkshire - Holland (requires North Sea fleet that isn't there)
		// Result: Yorkshire's illegal order ignored; French support holds; Yorkshire not dislodged.
		{
			description: "DATC 6.D.32 - a missing fleet",
			orders: []*result{
				{order: "F Edi S A Lvp-Yor", position: "edi"},
				{order: "A Lvp-Yor", position: "yor"},
				{order: "F Lon S A Yor", position: "lon"},
				{order: "A Yor-Hol", position: "yor"},
			},
		},
	*/

	// DATC 6.D.33. UNWANTED SUPPORT ALLOWED
	// Austria: A Serbia - Budapest, A Vienna - Budapest
	// Russia: A Galicia Supports A Serbia - Budapest (unwanted by Austria)
	// Turkey: A Bulgaria - Serbia
	// Result: Russia's support gives Serbia strength 2; Serbia wins Budapest;
	//         Turkey captures the now-empty Serbia.
	{
		description: "DATC 6.D.33 - unwanted support allowed",
		orders: []*result{
			{order: "A Ser-Bud", position: "bud"},
			{order: "A Vie-Bud", position: "vie"},
			{order: "A Gal S A Ser-Bud", position: "gal"},
			{order: "A Bul-Ser", position: "ser"},
		},
	},

	/*
		// DATC 6.D.34. SUPPORT TARGETING OWN AREA NOT ALLOWED
		// Needs: fleet model + self-targeting support validation
		// Germany: A Berlin - Prussia, A Silesia Supports A Berlin - Prussia,
		//          F Baltic Sea Supports A Berlin - Prussia
		// Italy: A Prussia Supports A Livonia - Prussia (illegal: Prussia can't support
		//        into its own area)
		// Russia: A Warsaw Supports A Livonia - Prussia, A Livonia - Prussia
		// Result: Italian order illegal; German attack succeeds.
		{
			description: "DATC 6.D.34 - support targeting own area not allowed",
			orders: []*result{
				{order: "A Ber-Pru", position: "pru"},
				{order: "A Sil S A Ber-Pru", position: "sil"},
				{order: "F Bal S A Ber-Pru", position: "bal"},
				{order: "A Pru S A Lvn-Pru", position: "pru", defeated: true},
				{order: "A War S A Lvn-Pru", position: "war"},
				{order: "A Lvn-Pru", position: "lvn"},
			},
		},
	*/

	// ===== DATC 6.E. HEAD-TO-HEAD BATTLES AND BELEAGUERED GARRISON =====
	//
	// All 6.E tests require the head-to-head battle algorithm (Phase 3).
	// In a head-to-head battle, a dislodged loser has no effect on the winner's
	// origin territory. The current simultaneous-pass algorithm does not model
	// this correctly. Uncomment when implementing Phase 3. See CLAUDE.md.

	/*
		// DATC 6.E.1. DISLODGED UNIT HAS NO EFFECT ON ATTACKER'S AREA
		// Germany: A Berlin - Prussia, A Kiel - Berlin, A Silesia Supports A Berlin - Prussia
		// Russia: A Prussia - Berlin
		// Note: original uses F Kiel; army adaptation gives same result.
		// Result: Berlin dislodges Prussia; Kiel follows into the vacated Berlin.
		{
			description: "DATC 6.E.1 - dislodged unit has no effect on attacker's area",
			orders: []*result{
				{order: "A Ber-Pru", position: "pru"},
				{order: "A Kie-Ber", position: "ber"},
				{order: "A Sil S A Ber-Pru", position: "sil"},
				{order: "A Pru-Ber", position: "pru", defeated: true},
			},
		},
	*/

	/*
		// DATC 6.E.2. NO SELF DISLODGEMENT IN HEAD-TO-HEAD BATTLE
		// Needs: country rules
		// Germany: A Berlin - Kiel, F Kiel - Berlin, A Munich Supports A Berlin - Kiel
		// Result: no unit moves (self-dislodgement blocked).
		{
			description: "DATC 6.E.2 - no self dislodgement in head-to-head battle",
			orders: []*result{
				{order: "A Ber-Kie", position: "ber"},
				{order: "F Kie-Ber", position: "kie"},
				{order: "A Mun S A Ber-Kie", position: "mun"},
			},
		},
	*/

	/*
		// DATC 6.E.3. NO HELP IN DISLODGING OWN UNIT
		// Needs: country rules
		// Germany: A Berlin - Kiel, A Munich Supports F Kiel - Berlin
		// England: F Kiel - Berlin
		// Result: Germany cannot help England dislodge its own unit; no move.
		{
			description: "DATC 6.E.3 - no help in dislodging own unit",
			orders: []*result{
				{order: "A Ber-Kie", position: "ber"},
				{order: "A Mun S F Kie-Ber", position: "mun"},
				{order: "F Kie-Ber", position: "kie"},
			},
		},
	*/

	/*
		// DATC 6.E.4. NON-DISLODGED LOSER STILL HAS EFFECT
		// Needs: head-to-head algorithm + fleet model
		// Complex beleaguered garrison scenario with fleets. See DATC.txt 6.E.4.
		{
			description: "DATC 6.E.4 - non-dislodged loser still has effect",
			orders: []*result{
				{order: "F Hol-Nth", position: "hol"},
				{order: "F Hel S F Hol-Nth", position: "hel"},
				{order: "F Ska S F Hol-Nth", position: "ska"},
				{order: "F Nth-Hol", position: "nth"},
				{order: "F Bel S F Nth-Hol", position: "bel"},
				{order: "F Edi S F Nwg-Nth", position: "edi"},
				{order: "F Yor S F Nwg-Nth", position: "yor"},
				{order: "F Nwg-Nth", position: "nwg"},
				{order: "A Kie S A Ruh-Hol", position: "kie"},
				{order: "A Ruh-Hol", position: "ruh"},
			},
		},
	*/

	/*
		// DATC 6.E.5–6.E.15: Further head-to-head and beleaguered garrison scenarios.
		// All need head-to-head algorithm (Phase 3) and most need fleet model (Phase 4).
		// See DATC.txt sections 6.E.5 through 6.E.15 for full scenario details.
	*/

	// ===== DATC 6.F. CONVOYS =====
	// All 6.F tests require the convoy model (Phase 5). See CLAUDE.md.

	/*
		// DATC 6.F.1. NO CONVOY IN COASTAL AREAS
		// Turkey: A Greece - Sevastopol (convoy), F Aegean Convoys, F Constantinople Convoys,
		//         F Black Sea Convoys
		// Result: Constantinople is coastal; convoy fails; army stays in Greece.
		{
			description: "DATC 6.F.1 - no convoy in coastal areas",
			orders: []*result{
				{order: "A Gre-Sev", position: "gre"},
				{order: "F Aeg C A Gre-Sev", position: "aeg"},
				{order: "F Con C A Gre-Sev", position: "con"},
				{order: "F Bla C A Gre-Sev", position: "bla"},
			},
		},
	*/

	/*
		// DATC 6.F.2. AN ARMY BEING CONVOYED CAN BOUNCE AS NORMAL
		// England: F English Channel Convoys A London - Brest, A London - Brest
		// France: A Paris - Brest
		// Result: London and Paris bounce at Brest.
		{
			description: "DATC 6.F.2 - an army being convoyed can bounce as normal",
			orders: []*result{
				{order: "F Ech C A Lon-Bre", position: "ech"},
				{order: "A Lon-Bre", position: "lon"},
				{order: "A Par-Bre", position: "par"},
			},
		},
	*/

	/*
		// DATC 6.F.3–6.F.25: Convoy mechanics, paradox resolutions, and multi-route convoys.
		// All need convoy model (Phase 5). See DATC.txt for full scenario details.
	*/

	// ===== DATC 6.G. CONVOYING TO ADJACENT PROVINCES =====
	// All 6.G tests require the convoy model (Phase 5). See CLAUDE.md.

	/*
		// DATC 6.G.1. TWO UNITS CAN SWAP PROVINCES BY CONVOY
		// England: A Norway - Sweden (via convoy), F Skagerrak Convoys A Norway - Sweden
		// Russia: A Sweden - Norway
		// Result: convoy intent given by own fleet; armies swap.
		{
			description: "DATC 6.G.1 - two units can swap provinces by convoy",
			orders: []*result{
				{order: "A Nwy-Swe", position: "swe"},
				{order: "F Ska C A Nwy-Swe", position: "ska"},
				{order: "A Swe-Nwy", position: "nwy"},
			},
		},
	*/

	/*
		// DATC 6.G.2–6.G.10: Adjacent-province convoy edge cases.
		// All need convoy model (Phase 5). See DATC.txt for full scenario details.
	*/
}

type result struct {
	order    string
	position string
	defeated bool
	unit     *board.Unit
}

type spec struct {
	description string
	orders      []*result
	focus       bool
}

func TestMainPhaseResolver_ResolveCases(t *testing.T) {
	country := "a_country"

	for _, spec := range filter(specs) {
		t.Run(spec.description, func(t *testing.T) {

			is := is.New(t)

			logTableHeading(t)
			positionManager := board.NewPositionManager()

			orders := order.Set{}

			for _, result := range spec.orders {
				o, err := order.Decode(result.order, country)
				is.NoErr(err)
				var terr board.Territory
				switch v := o.(type) {
				case order.Move:
					terr = v.From
					orders.AddMove(v)
				case order.Hold:
					terr = v.At
					orders.AddHold(v)
				case order.MoveSupport:
					terr = v.By
					orders.AddMoveSupport(v)
				case order.HoldSupport:
					terr = v.By
					orders.AddHoldSupport(v)
				case order.MoveConvoy:
					terr = v.By
					orders.AddMoveConvoy(v)
				}

				u := &board.Unit{Country: country}
				positionManager.AddUnit(u, board.LookupTerritory(terr.Abbr))
				result.unit = u
			}

			validator := order.NewValidator(board.CreateArmyGraph())
			orderHandler := game.OrderHandler{
				Validator: validator,
			}
			orderHandler.ApplyOrders(orders, positionManager)
			game.ResolveOrders(positionManager)

			for _, result := range spec.orders {
				logTableRow(t, *result)
				is.NotNil(result.unit)
				is.Equal(result.defeated, positionManager.Defeated(result.unit))
				is.Equal(result.position, positionManager.Position(result.unit).Territory.Abbr)
			}
			is.Equal(len(spec.orders), len(positionManager.Positions()))
		})
	}
}

func filter(specs []spec) []spec {
	focused := make([]spec, 0)
	for _, c := range specs {
		if c.focus {
			focused = append(focused, c)
		}
	}
	if len(focused) == 0 {
		return specs
	}
	return focused
}

func logTableHeading(t *testing.T) {
	t.Helper()
	t.Log("  | order             | result | defeated |")
	t.Log("  +---------------------------------------+")
}

func logTableRow(t *testing.T, o result) {
	t.Helper()
	t.Logf("  | %s%s| %s    | %t%s|",
		o.order,
		strings.Repeat(" ", 18-len(o.order)),
		o.position,
		o.defeated,
		strings.Repeat(" ", 9-len(fmt.Sprintf("%t", o.defeated))))
}
