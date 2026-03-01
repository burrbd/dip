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

// specs lists all integration test scenarios. DATC tests are numbered by
// section (e.g. "DATC 6.A.11") and run in DATC document order. Tests that
// require unimplemented features are commented out; the comment states which
// open task unblocks them. To focus on one scenario during development, set
// focus: true (and unset before committing).
//
// DATC implementation order:
//   NOW (army engine, no extra features needed):
//     6.A.8, 6.A.11, 6.A.12, 6.C.1–6.C.3, 6.D.1–6.D.5, 6.D.9, 6.D.15, 6.D.21
//   NEXT – country-aware resolution (self-dislodgement):
//     6.D.10–6.D.14, 6.D.16, 6.D.20
//   AFTER fleet movement + coasts (Task 4):
//     6.A.1–6.A.3, 6.A.9–6.A.10, all of 6.B, fleet variants in 6.D, all of 6.E
//   AFTER convoy (Task 5):
//     6.C.4–6.C.7, 6.D.6, 6.D.8, 6.D.27, all of 6.F, all of 6.G
//   AFTER retreat phase:
//     all of 6.H
//   AFTER build phase:
//     all of 6.I
var specs = []spec{
	// ─── Non-DATC regression tests ────────────────────────────────────────────

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
		description: "given units attack in circular chain without support, then all attacking units bounce back",
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
		description: "given a counterattack, and another attacks one counterattack party," +
			"then all units bounce",
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
		description: "given a unit holds and another unit supports holding unit," +
			"then both units remain in position",
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
			"then support is still cut and attack fails (DATC 6.D.9)",
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
			"then support is still cut and attack fails (DATC 6.D.9)",
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
			"then support is still cut and attack on supported territory fails (DATC 6.D.9)",
		orders: []*result{
			{order: "A Gal-Vie", position: "gal"},
			{order: "A Boh S A Gal-Vie", position: "boh"},
			{order: "A Vie H", position: "vie"},
			{order: "A Mun-Boh", position: "mun"},
		},
	},

	// ─── DATC 6.A: Basic Checks ───────────────────────────────────────────────

	// DATC 6.A.1: Moving to non-adjacent area
	// England: F NTH → PIC (Picardy; non-adjacent). Order fails; F NTH stays.
	// Requires: sea territories + fleet movement (Task 4)
	// {
	// 	description: "DATC 6.A.1: moving to non-adjacent area",
	// 	orders: []*result{
	// 		// NTH (North Sea) is a sea territory, not yet in the territory map.
	// 		{order: "F Nth H", position: "nth"},
	// 	},
	// },

	// DATC 6.A.2: Move army to sea
	// England: A LVP → IRI (Irish Sea). Order fails; A LVP stays.
	// Requires: sea territories + unit-type validation (Task 4)
	// {
	// 	description: "DATC 6.A.2: move army to sea",
	// 	orders: []*result{
	// 		// IRI (Irish Sea) is not in the territory map.
	// 		{order: "A Lvp H", position: "lvp"},
	// 	},
	// },

	// DATC 6.A.3: Move fleet to land
	// Germany: F KIE → MUN. Order fails; F KIE stays.
	// Requires: fleet/army unit-type movement validation (Task 4)
	// {
	// 	description: "DATC 6.A.3: move fleet to land",
	// 	orders: []*result{
	// 		{order: "F Kie H", position: "kie"},
	// 	},
	// },

	// DATC 6.A.4: Move to own sector
	// Germany: F KIE → KIE. Order fails; F KIE stays.
	// Requires: fleet movement (Task 4); self-loop detection
	// {
	// 	description: "DATC 6.A.4: move to own sector",
	// 	orders: []*result{
	// 		{order: "F Kie H", position: "kie"},
	// 	},
	// },

	// DATC 6.A.5: Move to own sector with convoy
	// Complex: A YOR → YOR via convoy. Self-targeting convoy is void.
	// Requires: convoy + fleet movement (Task 5)

	// DATC 6.A.6: Ordering a unit of another country
	// Not applicable to current single-country engine.
	// Requires: multi-country order submission

	// DATC 6.A.7: Only armies can be convoyed
	// Requires: convoy validation (Task 5)

	// DATC 6.A.8: Support to hold yourself is not possible
	// Italy: A VEN → TRI (supported by A TYR); Austria: A TRI holds and
	// gives self-support (void). A TRI is dislodged.
	// Army-only adaptation (original uses F TRI).
	{
		description: "DATC 6.A.8: support to hold yourself is not possible (army adaptation)",
		orders: []*result{
			{order: "A Ven-Tri", position: "tri"},
			{order: "A Tyr S A Ven-Tri", position: "tyr"},
			// Self-support "A Tri S A Tri" is void; A Tri holds with strength 0.
			{order: "A Tri S A Tri", position: "tri", defeated: true},
		},
	},

	// DATC 6.A.9: Fleets must follow coast if not on sea
	// Italy: F ROM → VEN (must follow coast). Order fails; F ROM stays.
	// Requires: fleet + coast movement (Task 4)
	// {
	// 	description: "DATC 6.A.9: fleets must follow coast if not on sea",
	// 	orders: []*result{
	// 		{order: "F Rom H", position: "rom"},
	// 	},
	// },

	// DATC 6.A.10: Support on unreachable destination not possible
	// Italy: F ROM S A APU → VEN (ROM cannot reach VEN via fleet route). Void.
	// Requires: fleet + coast movement (Task 4)
	// {
	// 	description: "DATC 6.A.10: support on unreachable destination not possible",
	// 	orders: []*result{
	// 		// F ROM cannot support A APU-VEN because ROM cannot reach VEN by fleet.
	// 		{order: "F Rom H", position: "rom"},
	// 		{order: "A Apu-Ven", position: "apu"},
	// 		{order: "A Ven H", position: "ven"},
	// 	},
	// },

	// DATC 6.A.11: Simple bounce
	// Austria: A VIE → TYR; Italy: A VEN → TYR. Both bounce.
	{
		description: "DATC 6.A.11: simple bounce",
		orders: []*result{
			{order: "A Vie-Tyr", position: "vie"},
			{order: "A Ven-Tyr", position: "ven"},
		},
	},

	// DATC 6.A.12: Bounce of three units
	// Austria: A VIE → TYR; Germany: A MUN → TYR; Italy: A VEN → TYR. All bounce.
	{
		description: "DATC 6.A.12: bounce of three units",
		orders: []*result{
			{order: "A Vie-Tyr", position: "vie"},
			{order: "A Mun-Tyr", position: "mun"},
			{order: "A Ven-Tyr", position: "ven"},
		},
	},

	// ─── DATC 6.B: Coastal Issues ─────────────────────────────────────────────
	// All 6.B tests require fleet movement + coast designators (Task 4).
	// Territories SPA/NC, SPA/SC, GAS, MAR, POR, BUL/SC, BUL/EC, STP/NC, STP/SC
	// plus sea routes (MAO, LYO, WES, NTH, IRI, NAO) are not yet modelled.

	// DATC 6.B.1: Moving with unspecified coast when coast is necessary
	// DATC 6.B.2: Moving with unspecified coast when coast is not necessary
	// DATC 6.B.3: Moving with wrong coast when coast is not necessary
	// DATC 6.B.4: Support to unreachable coast allowed
	// DATC 6.B.5: Support from unreachable coast not allowed
	// DATC 6.B.6: Support can be cut with other coast
	// DATC 6.B.7: Supporting with unspecified coast
	// DATC 6.B.8: Supporting with unspecified coast when only one coast is possible
	// DATC 6.B.9: Supporting with wrong coast
	// DATC 6.B.10: Unit ordered with wrong coast
	// DATC 6.B.11: Coast cannot be ordered to change
	// DATC 6.B.12: Army movement with coastal specification
	// DATC 6.B.13: Coastal crawl not allowed
	// DATC 6.B.14: Building with unspecified coast (build phase)
	//
	// Uncomment after implementing Task 4 (fleet movement + coasts).

	// ─── DATC 6.C: Circular Movement ──────────────────────────────────────────

	// DATC 6.C.1: Three army circular movement
	// Turkey: A ANK → CON, A CON → SMY, A SMY → ANK. All succeed.
	// (Army adaptation; original uses F ANK.)
	{
		description: "DATC 6.C.1: three army circular movement (army adaptation)",
		orders: []*result{
			{order: "A Ank-Con", position: "con"},
			{order: "A Con-Smy", position: "smy"},
			{order: "A Smy-Ank", position: "ank"},
		},
	},

	// DATC 6.C.2: Three army circular movement with support
	// Turkey: A ANK → CON (supported by A BUL), A CON → SMY, A SMY → ANK.
	// All succeed; support has no negative effect on the rotation.
	// (Army adaptation; original uses F ANK.)
	{
		description: "DATC 6.C.2: three army circular movement with support (army adaptation)",
		orders: []*result{
			{order: "A Ank-Con", position: "con"},
			{order: "A Con-Smy", position: "smy"},
			{order: "A Smy-Ank", position: "ank"},
			{order: "A Bul S A Ank-Con", position: "bul"},
		},
	},

	// DATC 6.C.3: A disrupted three army circular movement
	// Turkey: A ANK → CON, A CON → SMY, A SMY → ANK; A BUL → CON (disrupts).
	// All bounce; the interloper at CON locks everything.
	// (Army adaptation; original uses F ANK.)
	{
		description: "DATC 6.C.3: disrupted three army circular movement (army adaptation)",
		orders: []*result{
			{order: "A Ank-Con", position: "ank"},
			{order: "A Con-Smy", position: "con"},
			{order: "A Smy-Ank", position: "smy"},
			{order: "A Bul-Con", position: "bul"},
		},
	},

	// DATC 6.C.4: A circular movement with attacked convoy
	// Austria/Turkey circular via convoy; Italy attacks a convoying fleet.
	// The convoy survives; circular movement succeeds.
	// Requires: convoy (Task 5)
	// {
	// 	description: "DATC 6.C.4: circular movement with attacked convoy",
	// 	orders: []*result{
	// 		// Austria: A TRI-SER, A SER-BUL
	// 		// Turkey: A BUL-TRI (via convoy through ADR+ION+AEG), fleets convoy
	// 		// Italy: F NAP-ION (attacks convoy fleet, fails)
	// 		// Sea territories ADR, ION, AEG, NAP not yet in map.
	// 	},
	// },

	// DATC 6.C.5: A disrupted circular movement due to dislodged convoy
	// Same as 6.C.4 but Italy (with support) dislodges convoying fleet.
	// Circular movement fails; A BUL has no convoy route.
	// Requires: convoy (Task 5)

	// DATC 6.C.6: Two armies with two convoys
	// England: A LON → BEL (convoyed by F NTH); France: A BEL → LON (by F ENG).
	// Both succeed; armies swap.
	// Requires: convoy + sea territories (Task 5)

	// DATC 6.C.7: Disrupted unit swap
	// Same as 6.C.6 but France also has A BUR → BEL which causes a bounce.
	// Requires: convoy + sea territories (Task 5)

	// ─── DATC 6.D: Supports ───────────────────────────────────────────────────

	// DATC 6.D.1: Supported hold can prevent dislodgement
	// Austria: A TRI → VEN (supported by A APU); Italy: A VEN H (hold-supported
	// by A TYR). Both sides at strength 1; A TRI bounces.
	// Army adaptation: APU replaces fleet ADR.
	{
		description: "DATC 6.D.1: supported hold can prevent dislodgement (army adaptation)",
		orders: []*result{
			{order: "A Tri-Ven", position: "tri"},
			{order: "A Apu S A Tri-Ven", position: "apu"},
			{order: "A Ven H", position: "ven"},
			{order: "A Tyr S A Ven", position: "tyr"},
		},
	},

	// DATC 6.D.2: A move cuts support on hold
	// Austria: A TRI → VEN (supp. A APU), A VIE → TYR;
	// Italy: A VEN H (hold-supp. A TYR), A TYR support cut by A VIE.
	// Result: A VEN dislodged (unsupported hold vs. supported attack), A VIE bounces.
	// Army adaptation: APU replaces fleet ADR.
	{
		description: "DATC 6.D.2: a move cuts support on hold (army adaptation)",
		orders: []*result{
			{order: "A Tri-Ven", position: "ven"},
			{order: "A Apu S A Tri-Ven", position: "apu"},
			{order: "A Vie-Tyr", position: "vie"},
			{order: "A Ven H", position: "ven", defeated: true},
			{order: "A Tyr S A Ven", position: "tyr"},
		},
	},

	// DATC 6.D.3: A move cuts support on move
	// Austria: A TRI → VEN (supp. A APU); Italy: A VEN H, A NAP → APU (cuts support).
	// Result: A APU's support cut; A TRI has strength 0; A TRI bounces. A NAP bounces.
	// Army adaptation: NAP replaces fleet ION, APU replaces fleet ADR.
	{
		description: "DATC 6.D.3: a move cuts support on move (army adaptation)",
		orders: []*result{
			{order: "A Tri-Ven", position: "tri"},
			{order: "A Apu S A Tri-Ven", position: "apu"},
			{order: "A Ven H", position: "ven"},
			{order: "A Nap-Apu", position: "nap"},
		},
	},

	// DATC 6.D.4: Support to hold on unit supporting a hold is allowed
	// Germany: A BER S A KIE (hold), A KIE S A BER (hold).
	// Russia: A PRU → BER (supp. A SIL). A PRU has str 1; A BER has str 1 (KIE supp.).
	// A BER's support for KIE IS cut (PRU attacks BER), but KIE's support for BER is not.
	// Result: A PRU bounces.
	// Army adaptation: SIL replaces sea fleet BAL.
	{
		description: "DATC 6.D.4: support to hold on unit supporting a hold is allowed (army adaptation)",
		orders: []*result{
			{order: "A Ber S A Kie", position: "ber"},
			{order: "A Kie S A Ber", position: "kie"},
			{order: "A Sil S A Pru-Ber", position: "sil"},
			{order: "A Pru-Ber", position: "pru"},
		},
	},

	// DATC 6.D.5: Support to hold on unit supporting a move is allowed
	// Germany: A BER S A MUN → BOH (move supp., cut by PRU attack), A KIE S A BER (hold supp.),
	//          A MUN → BOH (succeeds despite support cut).
	// Russia: A PRU → BER (supp. A SIL). A BER holds (str 1 from KIE). A PRU bounces.
	// Army adaptation: Mun-Boh instead of Mun-Sil; SIL replaces sea fleet BAL.
	{
		description: "DATC 6.D.5: support to hold on unit supporting a move is allowed (army adaptation)",
		orders: []*result{
			{order: "A Ber S A Mun-Boh", position: "ber"},
			{order: "A Kie S A Ber", position: "kie"},
			{order: "A Mun-Boh", position: "boh"},
			{order: "A Sil S A Pru-Ber", position: "sil"},
			{order: "A Pru-Ber", position: "pru"},
		},
	},

	// DATC 6.D.6: Support to hold on convoying unit allowed
	// Germany: A BER → SWE (convoyed by F BAL, hold-supported by F PRU).
	// Russia: F LVN → BAL (supp. F BOT). F PRU's support keeps F BAL from being dislodged.
	// Requires: convoy + sea territories (Task 5)
	// {
	// 	description: "DATC 6.D.6: support to hold on convoying unit allowed",
	// 	orders: []*result{
	// 		// Sea territories BAL, BOT, LVN (sea route); SWE reachable only by convoy.
	// 	},
	// },

	// DATC 6.D.7: Support to hold on moving unit not allowed
	// Germany: F BAL → SWE, F PRU S F BAL (hold support on moving unit = void).
	// Russia: F LVN → BAL (supp. F BOT), A FIN → SWE.
	// F BAL is dislodged; A FIN bounces at SWE.
	// Requires: fleet + sea territories (Task 4)
	// {
	// 	description: "DATC 6.D.7: support to hold on moving unit not allowed",
	// 	orders: []*result{
	// 		// Sea territories BAL, BOT, PRU (sea fleet), LVN, SWE not in army map.
	// 	},
	// },

	// DATC 6.D.8: Failed convoy cannot receive hold support
	// Austria: F ION H, A SER S A ALB → GRE, A ALB → GRE.
	// Turkey: A GRE → NAP (no convoy route), A BUL S A GRE (support void; GRE moves away).
	// A GRE dislodged by A ALB.
	// Requires: convoy (Task 5)

	// DATC 6.D.9: Support to move on holding unit not allowed
	// Italy: A VEN → TRI (supp. A TYR); Austria: A ALB S A TRI → SER (void; TRI holds),
	//        A TRI H.
	// A VEN dislodges A TRI; A ALB's support is void.
	{
		description: "DATC 6.D.9: support to move on holding unit not allowed",
		orders: []*result{
			{order: "A Ven-Tri", position: "tri"},
			{order: "A Tyr S A Ven-Tri", position: "tyr"},
			{order: "A Alb S A Tri-Ser", position: "alb"},
			{order: "A Tri H", position: "tri", defeated: true},
		},
	},

	// DATC 6.D.10: Self dislodgement prohibited
	// Germany: A BER H, F KIE → BER (own unit), A MUN S F KIE → BER. All void.
	// Requires: country-aware self-dislodgement prohibition
	// {
	// 	description: "DATC 6.D.10: self dislodgement prohibited",
	// 	orders: []*result{
	// 		{order: "A Ber H", position: "ber"},
	// 		{order: "F Kie-Ber", position: "kie"}, // same country — order void
	// 		{order: "A Mun S F Kie-Ber", position: "mun"},
	// 	},
	// },

	// DATC 6.D.11: No self dislodgement of returning unit
	// Germany: A BER → PRU, F KIE → BER, A MUN S F KIE → BER.
	// Russia: A WAR → PRU. A BER bounces at PRU; F KIE cannot take BER (own unit returns).
	// Requires: country-aware self-dislodgement prohibition
	// {
	// 	description: "DATC 6.D.11: no self dislodgement of returning unit",
	// 	orders: []*result{
	// 		{order: "A Ber-Pru", position: "ber"},
	// 		{order: "F Kie-Ber", position: "kie"},
	// 		{order: "A Mun S F Kie-Ber", position: "mun"},
	// 		{order: "A War-Pru", position: "war"},
	// 	},
	// },

	// DATC 6.D.12: Supporting a foreign unit to dislodge own unit prohibited
	// Austria: F TRI H, A VIE S A VEN → TRI. Italy: A VEN → TRI.
	// Austria cannot support Italy to dislodge own F TRI; support is void.
	// Requires: country-aware self-dislodgement prohibition
	// {
	// 	description: "DATC 6.D.12: supporting a foreign unit to dislodge own unit prohibited",
	// 	orders: []*result{
	// 		{order: "A Tri H", position: "tri"},
	// 		{order: "A Vie S A Ven-Tri", position: "vie"},
	// 		{order: "A Ven-Tri", position: "ven"},
	// 	},
	// },

	// DATC 6.D.13: Supporting a foreign unit to dislodge a returning own unit prohibited
	// Austria: F TRI → ADR, A VIE S A VEN → TRI. Italy: A VEN → TRI, F APU → ADR.
	// F TRI and F APU bounce at ADR. Austria's support for A VEN → TRI is void.
	// Requires: country-aware self-dislodgement prohibition + sea territory ADR (Task 4)

	// DATC 6.D.14: Supporting a foreign unit is not enough to prevent dislodgement
	// Austria: F TRI H, A VIE S A VEN → TRI. Italy: A VEN → TRI (2 extra supports).
	// Austria's own support is void; Italy dislodges with net strength 3.
	// Requires: country-aware self-dislodgement prohibition
	// {
	// 	description: "DATC 6.D.14: supporting a foreign unit is not enough to prevent dislodgement",
	// 	orders: []*result{
	// 		{order: "A Tri H", position: "tri", defeated: true},
	// 		{order: "A Vie S A Ven-Tri", position: "vie"},  // void — own unit
	// 		{order: "A Ven-Tri", position: "tri"},
	// 		{order: "A Tyr S A Ven-Tri", position: "tyr"},
	// 		{order: "A Apu S A Ven-Tri", position: "apu"},
	// 	},
	// },

	// DATC 6.D.15: Defender cannot cut support for attack on itself
	// Russia: A CON S A SMY → ANK, A SMY → ANK.
	// Turkey: A ANK → CON (attacks the supporter, but cannot cut support aimed at itself).
	// Result: A SMY dislodges A ANK; A CON's support is not cut.
	// Army adaptation (original uses sea fleets in Black Sea).
	{
		description: "DATC 6.D.15: defender cannot cut support for attack on itself (army adaptation)",
		orders: []*result{
			{order: "A Con S A Smy-Ank", position: "con"},
			{order: "A Smy-Ank", position: "ank"},
			{order: "A Ank-Con", position: "ank", defeated: true},
		},
	},

	// DATC 6.D.16: Convoying a unit dislodging a unit of same power is allowed
	// England: A LON H, F NTH C A BEL → LON. France: F ENG S A BEL → LON, A BEL → LON.
	// A BEL dislodges A LON (same power's convoy is allowed to carry the dislodging army).
	// Requires: convoy + country-aware rules (Task 5)

	// DATC 6.D.17: Dislodgement cuts supports
	// Russia: F CON S F BLA → ANK, F BLA → ANK. Turkey: F ANK → CON, A SMY S F ANK → CON,
	//         A ARM → ANK.
	// F CON is dislodged; its support is cut. F BLA bounces. F ANK succeeds.
	// Requires: sea territories BLA, sea fleet movement (Task 4)

	// DATC 6.D.18: A surviving unit will sustain support
	// Same as 6.D.17 but Russia also has A BUL S F CON (hold support).
	// F CON survives; F BLA dislodges F ANK.
	// Requires: sea territories (Task 4)

	// DATC 6.D.19: Even when surviving is in alternative way
	// Russia: F CON S F BLA → ANK, F BLA → ANK, A SMY S F ANK → CON (supports own attacker).
	// Turkey: F ANK → CON. A SMY's support is void (supporting attack on CON while CON
	// supports attack on ANK). F BLA dislodges F ANK.
	// Requires: sea territories (Task 4)

	// DATC 6.D.20: Unit cannot cut support of its own country
	// England: F LON S F NTH → ENG, F NTH → ENG, A YOR → LON (own unit, void cut).
	// France: F ENG H. F ENG dislodged.
	// Requires: country-aware support cutting + sea territories (Task 4)

	// DATC 6.D.21: Dislodging does not cancel a support cut
	// Austria: A TRI H. Italy: A VEN → TRI, A TYR S A VEN → TRI.
	// Germany: A MUN → TYR (cuts Italy's support). Russia: A SIL → MUN (supp. A BER).
	// A MUN cuts A TYR's support then is dislodged by A SIL. Cut stands; A VEN bounces.
	{
		description: "DATC 6.D.21: dislodging does not cancel a support cut",
		orders: []*result{
			{order: "A Tri H", position: "tri"},
			{order: "A Ven-Tri", position: "ven"},
			{order: "A Tyr S A Ven-Tri", position: "tyr"},
			{order: "A Mun-Tyr", position: "mun", defeated: true},
			{order: "A Sil-Mun", position: "mun"},
			{order: "A Ber S A Sil-Mun", position: "ber"},
		},
	},

	// DATC 6.D.22: Impossible fleet move cannot be supported
	// Germany: F KIE → MUN (invalid fleet move), A BUR S F KIE → MUN (void support).
	// Russia: A MUN → KIE (supp. A BER). F KIE is dislodged.
	// Requires: fleet movement validation (Task 4)
	// {
	// 	description: "DATC 6.D.22: impossible fleet move cannot be supported",
	// 	orders: []*result{
	// 		{order: "F Kie H", position: "kie", defeated: true}, // illegal fleet move
	// 		{order: "A Bur S F Kie-Mun", position: "bur"},       // void support
	// 		{order: "A Mun-Kie", position: "kie"},
	// 		{order: "A Ber S A Mun-Kie", position: "ber"},
	// 	},
	// },

	// DATC 6.D.23: Impossible coast move cannot be supported
	// Italy: F LYO → SPA/SC (supported); France: F SPA/NC → LYO (impossible coast). Void.
	// Requires: coast designators + fleet movement (Task 4)

	// DATC 6.D.24: Impossible army move cannot be supported
	// France: A MAR → LYO (army to sea; void), F SPA/SC S A MAR → LYO (void support).
	// Turkey: F WES → LYO (supp. F TYS). F LYO dislodged.
	// Requires: sea territories + fleet movement (Task 4)

	// DATC 6.D.25: Failing hold support can be supported
	// Germany: A BER S A PRU (hold, but BER is attacked → support cut), F KIE S A BER.
	// Russia: F BAL S A PRU → BER, A PRU → BER. A PRU bounces (A BER holds with KIE supp.).
	// Requires: sea fleet BAL (Task 4)
	// {
	// 	description: "DATC 6.D.25: failing hold support can be supported",
	// 	orders: []*result{
	// 		// BAL is a sea territory. Replace with army territory when Task 4 done.
	// 	},
	// },

	// DATC 6.D.26: Failing move support can be supported
	// Germany: A BER S A PRU → SIL (support cut; BER attacked), F KIE S A BER.
	// Russia: F BAL S A PRU → BER, A PRU → BER. A PRU bounces.
	// Requires: sea fleet BAL (Task 4)

	// DATC 6.D.27: Failing convoy can be supported
	// England: F SWE → BAL (supp. F DEN). Germany: A BER H.
	// Russia: F BAL C A BER → LVN (convoy), F PRU S F BAL. F SWE bounces; convoy survives.
	// Requires: convoy + sea territories (Task 5)

	// DATC 6.D.28: Impossible move and support
	// Austria: A BUD S F RUM. Russia: F RUM → HOL (impossible; landlocked route).
	// Turkey: F BLA → RUM (supp. A BUL). F BLA bounces (A BUD's support for F RUM valid).
	// Requires: fleet movement + sea territory BLA (Task 4)

	// DATC 6.D.29: Move to impossible coast and support
	// Similar to 6.D.28 with F RUM → BUL/SC (impossible coast).
	// Requires: coast designators + sea territory BLA (Task 4)

	// DATC 6.D.30: Move without coast and support
	// Italy: F AEG S F CON. Russia: F CON → BUL (no coast; ambiguous).
	// Turkey: F BLA → CON (supp. A BUL). F BLA bounces.
	// Requires: sea territories AEG, BLA (Task 4)

	// DATC 6.D.31: A tricky impossible support
	// Austria: A RUM → ARM (land move impossible; needs convoy).
	// Turkey: F BLA S A RUM → ARM. Support void (move requires convoy).
	// Requires: sea territory BLA + convoy (Task 5)

	// DATC 6.D.32–6.D.34: (content needed from DATC document)
	// TODO: add these tests once the full DATC document is available.

	// ─── DATC 6.E: Head to Head Battles and Beleaguered Garrison ──────────────
	// All 6.E tests involve the head-to-head rule: when A→B and B→A, support
	// "from behind" the defender (i.e. from A's territory) does not count.
	// Many tests use fleet routes. All are commented pending the fixed-point
	// resolver being verified correct for these edge cases (see open tasks).
	// TODO: add exact orders from DATC document for 6.E.1–6.E.14.

	// DATC 6.E.1: Prevented head-to-head battle
	// DATC 6.E.2: Swapping places with unequal strength
	// DATC 6.E.3: Swapping places with supported hold
	// DATC 6.E.4: Support on attack on itself
	// DATC 6.E.5: Three units in a beleaguered garrison
	// DATC 6.E.6: Beleaguered garrison with no move
	// DATC 6.E.7: No self dislodgement in beleaguered garrison
	// DATC 6.E.8: No help in dislodging own unit
	// DATC 6.E.9: Beleaguered garrison with almost unresolvable bounce
	// DATC 6.E.10: Beleaguered garrison with almost circular movement
	// DATC 6.E.11: A move from convoy can be seen as going over land
	// DATC 6.E.12: Convoyed unit can bounce without dislodging
	// DATC 6.E.13: Dislodgement of multi route convoy
	// DATC 6.E.14: Contested for both coasts
	// Requires: fleet movement (Task 4); some require convoy (Task 5)

	// ─── DATC 6.F: Convoys ────────────────────────────────────────────────────
	// All 6.F tests require convoy resolution (Task 5) and sea territories.
	// TODO: add exact orders from DATC document for 6.F.1–6.F.24.

	// DATC 6.F.1:  No convoy in land area
	// DATC 6.F.2:  Army being convoyed can bounce as normal
	// DATC 6.F.3:  Army can still be on hold support
	// DATC 6.F.4:  Convoy with disrupted fleet
	// DATC 6.F.5:  Disrupted multi-route convoy
	// DATC 6.F.6:  Two armies with two convoys
	// DATC 6.F.7:  Convoy disrupted by dislodged fleet
	// DATC 6.F.8:  Non-disrupted multi-route convoy
	// DATC 6.F.9:  Dislodgement is prevented by double convoy
	// DATC 6.F.10: Convoying to adjacent place with other route
	// DATC 6.F.11: Swapping with intent to convoy
	// DATC 6.F.12: Support on attack on itself via convoy
	// DATC 6.F.13: Missing fleet convoyed to adjacent place
	// DATC 6.F.14: Army can bounce while convoying
	// DATC 6.F.15: Bounce by convoy with disrupted convoy
	// DATC 6.F.16: Convoy alone is not sufficient for swapping
	// DATC 6.F.17: Convoy can pass through not related fleet
	// DATC 6.F.18: Multi-convoys with disruption
	// DATC 6.F.19: Multi-routes with disruption
	// DATC 6.F.20: Impossible fleet movement resolved before convoy
	// DATC 6.F.21: Possible fleet movement resolved before convoy
	// DATC 6.F.22: Self-dislodgement supported by convoy not allowed
	// DATC 6.F.23: No self-dislodgement with beleaguered garrison and convoy
	// DATC 6.F.24: Dislodgement of multi-route convoy

	// ─── DATC 6.G: Convoying to Adjacent Places ───────────────────────────────
	// All 6.G tests require convoy (Task 5).
	// TODO: add exact orders from DATC document for 6.G.1–6.G.8.

	// DATC 6.G.1:  Two units can swap places by convoy
	// DATC 6.G.2:  Kidnapping an army
	// DATC 6.G.3:  Swapping with unintended intent
	// DATC 6.G.4:  Support on attack on own unit
	// DATC 6.G.5:  Swapping with illegal move not ignored
	// DATC 6.G.6:  Illegal support for convoy swap
	// DATC 6.G.7:  No convoy to adjacent place with attack
	// DATC 6.G.8:  No self-dislodgement with beleaguered garrison

	// ─── DATC 6.H: Retreating ─────────────────────────────────────────────────
	// All 6.H tests require a retreat phase, which is not yet implemented.
	// TODO: add exact orders from DATC document for 6.H.1–6.H.12.

	// DATC 6.H.1:  No support during retreat
	// DATC 6.H.2:  No throwing of support during retreat
	// DATC 6.H.3:  No retreat with another army
	// DATC 6.H.4:  No retreat to area left from
	// DATC 6.H.5:  Retreat when dislodged
	// DATC 6.H.6:  Unit cannot retreat to contested area
	// DATC 6.H.7:  Multiple retreat to same area
	// DATC 6.H.8:  Three-unit circular retreat
	// DATC 6.H.9:  Optional retreat to mother country
	// DATC 6.H.10: Unit retreats to country of another power
	// DATC 6.H.11: Retreat to contested area after failure
	// DATC 6.H.12: Retreat to area that was a starting point

	// ─── DATC 6.I: Building ───────────────────────────────────────────────────
	// All 6.I tests require a build/adjustment phase, which is not yet implemented.
	// TODO: add exact orders from DATC document for 6.I.1–6.I.7.

	// DATC 6.I.1:  Building with too many units
	// DATC 6.I.2:  Building too many units
	// DATC 6.I.3:  Three builds in a home supply centre
	// DATC 6.I.4:  Building in occupied supply centre
	// DATC 6.I.5:  Civil disorder
	// DATC 6.I.6:  Removing the wrong unit
	// DATC 6.I.7:  Failing to remove own unit
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
