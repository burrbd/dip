{ % This program is copyright (C) 1982 by D. E. Knuth; all rights are reserved.
% Unlimited copying and redistribution of this file are permitted as long
% as this file is not modified. Modifications are permitted, but only if
% the resulting file is not named tex.web. (The WEB system provides
% for alterations via an auxiliary file; the master file should stay intact.)
% See Appendix H of the WEB manual for hints on how to install this program.
% And see Appendix A of the TRIP manual for details about how to validate it.

% TeX is a trademark of the American Mathematical Society.
% METAFONT is a trademark of Addison-Wesley Publishing Company.

% Version 0 was released in September 1982 after it passed a variety of tests.
% Version 1 was released in November 1983 after thorough testing.
% Version 1.1 fixed ``disappearing font identifiers'' et alia (July 1984).
% Version 1.2 allowed `0' in response to an error, et alia (October 1984).
% Version 1.3 made memory allocation more flexible and local (November 1984).
% Version 1.4 fixed accents right after line breaks, et alia (April 1985).
% Version 1.5 fixed \the\toks after other expansion in \edefs (August 1985).
% Version 2.0 (almost identical to 1.5) corresponds to "Volume B" (April 1986).
% Version 2.1 corrected anomalies in discretionary breaks (January 1987).
% Version 2.2 corrected "(Please type...)" with null \endlinechar (April 1987).
% Version 2.3 avoided incomplete page in premature termination (August 1987).
% Version 2.4 fixed \noaligned rules in indented displays (August 1987).
% Version 2.5 saved cur_order when expanding tokens (September 1987).
% Version 2.6 added 10sp slop when shipping leaders (November 1987).
% Version 2.7 improved rounding of negative-width characters (November 1987).
% Version 2.8 fixed weird bug if no \patterns are used (December 1987).
% Version 2.9 made \csname\endcsname's "relax" local (December 1987).
% Version 2.91 fixed \outer\def\a0[]\a\a bug (April 1988).
% Version 2.92 fixed \patterns, also file names with complex macros (May 1988).
% Version 2.93 fixed negative halving in allocator when mem_min<0 (June 1988).
% Version 2.94 kept open_log_file from calling fatal_error (November 1988).
% Version 2.95 solved that problem a better way (December 1988).
% Version 2.96 corrected bug in "Infinite shrinkage" recovery (January 1989).
% Version 2.97 corrected blunder in creating 2.95 (February 1989).
% Version 2.98 omitted save_for_after at outer level (March 1989).
% Version 2.99 caught $$\begingroup\halign..$$ (June 1989).
% Version 2.991 caught .5\ifdim.6... (June 1989).
% Version 2.992 introduced major changes for 8-bit extensions (September 1989).
% Version 2.993 fixed a save_stack synchronization bug et alia (December 1989).
% Version 3.0 fixed unusual displays; was more \output robust (March 1990).
% Version 3.1 fixed nullfont, disabled \write[\the\prevgraf] (September 1990).
% Version 3.14 fixed unprintable font names and corrected typos (March 1991).
% Version 3.141 more of same; reconstituted ligatures better (March 1992).
% Version 3.1415 preserved nonexplicit kerns, tidied up (February 1993).
% Version 3.14159 allowed fontmemsize to change; bulletproofing (March 1995).
% Version 3.141592 fixed \xleaders, glueset, weird alignments (December 2002).
% Version 3.1415926 was a general cleanup with minor fixes (February 2008).
% Version 3.14159265 was similar (January 2014).
% Version 3.141592653 was similar but more extensive (January 2021).

% A reward of $327.68 will be paid to the first finder of any remaining bug.

% Although considerable effort has been expended to make the TeX program
% correct and reliable, no warranty is implied; the author disclaims any
% obligation or liability for damages, including but not limited to
% special, indirect, or consequential damages arising out of or in
% connection with the use or performance of this software. This work has
% been a ``labor of love'' and the author hopes that users enjoy it.

% Here is TeX material that gets inserted after \input webmac
\def\hang[\hangindent 3em\noindent\ignorespaces]
\def\hangg#1 [\hang\hbox[#1 ]]
\def\textindent#1[\hangindent2.5em\noindent\hbox to2.5em[\hss#1 ]\ignorespaces]
\font\ninerm=cmr9
\let\mc=\ninerm % medium caps for names like SAIL
\def\PASCAL[Pascal]
\def\ph[\hbox[Pascal-H]]
\def\pct![[\char`\%]] % percent sign in ordinary text
\font\logo=logo10 % font used for the METAFONT logo
\def\MF[[\logo META]\-[\logo FONT]]
\def\<#1>[$\langle#1\rangle$]
\def\section[\mathhexbox278]

\def\(#1)[] % this is used to make section names sort themselves better
\def\9#1[] % this is used for sort keys in the index via @:sort key][entry@>

\outer\def\N#1. \[#2]#3.[\MN#1.\vfil\eject % begin starred section
  \def\rhead[PART #2:\uppercase[#3]] % define running headline
  \message[*\modno] % progress report
  \edef\next[\write\cont[\Z[\?#2]#3][\modno][\the\pageno]]]\next
  \ifon\startsection[\bf\ignorespaces#3.\quad]\ignorespaces]
\let\?=\relax % we want to be able to \write a \?

\def\title[\TeX82]
\def\topofcontents[\hsize 5.5in
  \vglue 0pt plus 1fil minus 1.5in
  \def\?##1][\hbox[Changes to ##1.\ ]]
  ]
\let\maybe=\iffalse
\def\botofcontents[\vskip 0pt plus 1fil minus 1.5in]
\pageno=3
\def\glob[13] % this should be the section number of "<Global...>"
\def\gglob[20, 26] % this should be the next two sections of "<Global...>"

 }

{ 1. \[1] Introduction }

{tangle:pos tex.web:95:22: }

{ This is \TeX, a document compiler intended to produce typesetting of high
quality.
The \PASCAL\ program that follows is the definition of \TeX82, a standard
\xref[PASCAL][\PASCAL]
 \xref[TeX82][\TeX82]
version of \TeX\ that is designed to be highly portable so that identical output
will be obtainable on a great variety of computers.

The main purpose of the following program is to explain the algorithms of \TeX\
as clearly as possible. As a result, the program will not necessarily be very
efficient when a particular \PASCAL\ compiler has translated it into a
particular machine language. However, the program has been written so that it
can be tuned to run efficiently in a wide variety of operating environments
by making comparatively few changes. Such flexibility is possible because
the documentation that follows is written in the \.[WEB] language, which is
at a higher level than \PASCAL; the preprocessing step that converts \.[WEB]
to \PASCAL\ is able to introduce most of the necessary refinements.
Semi-automatic translation to other languages is also feasible, because the
program below does not make extensive use of features that are peculiar to
\PASCAL.

A large piece of software like \TeX\ has inherent complexity that cannot
be reduced below a certain level of difficulty, although each individual
part is fairly simple by itself. The \.[WEB] language is intended to make
the algorithms as readable as possible, by reflecting the way the
individual program pieces fit together and by providing the
cross-references that connect different parts. Detailed comments about
what is going on, and about why things were done in certain ways, have
been liberally sprinkled throughout the program.  These comments explain
features of the implementation, but they rarely attempt to explain the
\TeX\ language itself, since the reader is supposed to be familiar with
[\sl The \TeX book].
\xref[WEB]
\xref[TeXbook][\sl The \TeX book] }

{ 2. }

{tangle:pos tex.web:131:3: }

{ The present implementation has a long ancestry, beginning in the summer
of~1977, when Michael~F. Plass and Frank~M. Liang designed and coded
a prototype
\xref[Plass, Michael Frederick]
\xref[Liang, Franklin Mark]
\xref[Knuth, Donald Ervin]
based on some specifications that the author had made in May of that year.
This original proto\TeX\ included macro definitions and elementary
manipulations on boxes and glue, but it did not have line-breaking,
page-breaking, mathematical formulas, alignment routines, error recovery,
or the present semantic nest; furthermore,
it used character lists instead of token lists, so that a control sequence
like \.[\\halign] was represented by a list of seven characters. A
complete version of \TeX\ was designed and coded by the author in late
1977 and early 1978; that program, like its prototype, was written in the
[\mc SAIL] language, for which an excellent debugging system was
available. Preliminary plans to convert the [\mc SAIL] code into a form
somewhat like the present ``web'' were developed by Luis Trabb~Pardo and
\xref[Trabb Pardo, Luis Isidoro]
the author at the beginning of 1979, and a complete implementation was
created by Ignacio~A. Zabala in 1979 and 1980. The \TeX82 program, which
\xref[Zabala Salelles, Ignacio Andr\'es]
was written by the author during the latter part of 1981 and the early
part of 1982, also incorporates ideas from the 1979 implementation of
\xref[Guibas, Leonidas Ioannis]
\xref[Sedgewick, Robert]
\xref[Wyatt, Douglas Kirk]
\TeX\ in [\mc MESA] that was written by Leonidas Guibas, Robert Sedgewick,
and Douglas Wyatt at the Xerox Palo Alto Research Center.  Several hundred
refinements were introduced into \TeX82 based on the experiences gained with
the original implementations, so that essentially every part of the system
has been substantially improved. After the appearance of ``Version 0'' in
September 1982, this program benefited greatly from the comments of
many other people, notably David~R. Fuchs and Howard~W. Trickey.
A final revision in September 1989 extended the input character set to
eight-bit codes and introduced the ability to hyphenate words from
different languages, based on some ideas of Michael~J. Ferguson.
\xref[Fuchs, David Raymond]
\xref[Trickey, Howard Wellington]
\xref[Ferguson, Michael John]

No doubt there still is plenty of room for improvement, but the author
is firmly committed to keeping \TeX82 ``frozen'' from now on; stability
and reliability are to be its main virtues.

On the other hand, the \.[WEB] description can be extended without changing
the core of \TeX82 itself, and the program has been designed so that such
extensions are not extremely difficult to make.
The |banner| string defined here should be changed whenever \TeX\
undergoes any modifications, so that it will be clear which version of
\TeX\ might be the guilty party when a problem arises.
\xref[extensions to \TeX]
\xref[system dependencies]

If this program is changed, the resulting system should not be called
`\TeX'; the official name `\TeX' by itself is reserved
for software systems that are fully compatible with each other.
A special test suite called the ``\.[TRIP] test'' is available for
helping to determine whether a particular implementation deserves to be
known as `\TeX' [cf.~Stanford Computer Science report CS1027,
November 1984].

ML\TeX[] will add new primitives changing the behaviour of \TeX.  The
|banner| string has to be changed.  We do not change the |banner|
string, but will output an additional line to make clear that this is
a modified \TeX[] version. }

{ 3. }

{tangle:pos tex.web:195:3: }

{ Different \PASCAL s have slightly different conventions, and the present
 \xref[PASCAL H][\ph]
program expresses \TeX\ in terms of the \PASCAL\ that was
available to the author in 1982. Constructions that apply to
this particular compiler, which we shall call \ph, should help the
reader see how to make an appropriate interface for other systems
if necessary. (\ph\ is Charles Hedrick's modification of a compiler
\xref[Hedrick, Charles Locke]
for the DECsystem-10 that was originally developed at the University of
Hamburg; cf.\ [\sl Software---Practice and Experience \bf6] (1976),
29--42. The \TeX\ program below is intended to be adaptable, without
extensive changes, to most other versions of \PASCAL, so it does not fully
use the admirable features of \ph. Indeed, a conscious effort has been
made here to avoid using several idiosyncratic features of standard
\PASCAL\ itself, so that most of the code can be translated mechanically
into other high-level languages. For example, the `\&[with]' and `\\[new]'
features are not used, nor are pointer types, set types, or enumerated
scalar types; there are no `\&[var]' parameters, except in the case of files;
there are no tag fields on variant records; there are no assignments
|real:=integer|; no procedures are declared local to other procedures.)

The portions of this program that involve system-dependent code, where
changes might be necessary because of differences between \PASCAL\ compilers
and/or differences between
operating systems, can be identified by looking at the sections whose
numbers are listed under `system dependencies' in the index. Furthermore,
the index entries for `dirty \PASCAL' list all places where the restrictions
of \PASCAL\ have not been followed perfectly, for one reason or another.
 \xref[system dependencies]
 \xref[dirty \PASCAL]

Incidentally, \PASCAL's standard |round| function can be problematical,
because it disagrees with the IEEE floating-point standard.
Many implementors have
therefore chosen to substitute their own home-grown rounding procedure. }

{ 4. }

{tangle:pos tex.web:231:3: }

{ The program begins with a normal \PASCAL\ program heading, whose
components will be filled in later, using the conventions of \.[WEB].
\xref[WEB]
For example, the portion of the program called `\X\glob:Global
variables\X' below will be replaced by a sequence of variable declarations
that starts in $\section\glob$ of this documentation. In this way, we are able
to define each individual global variable when we are prepared to
understand what it means; we do not have to define all of the globals at
once.  Cross references in $\section\glob$, where it says ``See also
sections \gglob, \dots,'' also make it possible to look at the set of
all global variables, if desired.  Similar remarks apply to the other
portions of the program heading. } { \4 }
{ Compiler directives }
{$C-,A+,D-} {no range check, catch arithmetic overflow, no debug overhead}
 ifdef('TEXMF_DEBUG')  {$C+,D+}  endif('TEXMF_DEBUG')  {but turn everything on when debugging}



program TEX; {all file names are defined dynamically}

const

  start_of_TEX = 1 {go here when \TeX's variables are initialized} ;
  final_end = 9999 {this label marks the ending of the program} ;
  ssup_error_line =  255 ;
  max_font_max = 9000 {maximum number of internal fonts; this can be
                      increased, but |hash_size+max_font_max|
                      should not exceed 29000.} ;
  font_base = 0 {smallest internal font number; must be
                |>= min_quarterword|; do not change this without
                modifying the dynamic definition of the font arrays.} ;
  hash_size = 15000 {maximum number of control sequences; it should be at most
  about |(mem_max-mem_min)/10|; see also |font_max|} ;
  hash_prime = 8501 {a prime number equal to about 85\pct! of |hash_size|} ;
  hyph_prime = 607 {another prime for hashing \.[\\hyphenation] exceptions;
                if you change this, you should also change |iinf_hyphen_size|.}
{ \xref[system dependencies] } ;
  exit = 10 {go here to leave a procedure} ;
  restart = 20 {go here to start a procedure again} ;
  reswitch = 21 {go here to start a case statement again} ;
  continue = 22 {go here to resume a loop} ;
  done = 30 {go here to exit a loop} ;
  done1 = 31 {like |done|, when there is more than one loop} ;
  done2 = 32 {for exiting the second loop in a long block} ;
  done3 = 33 {for exiting the third loop in a very long block} ;
  done4 = 34 {for exiting the fourth loop in an extremely long block} ;
  done5 = 35 {for exiting the fifth loop in an immense block} ;
  done6 = 36 {for exiting the sixth loop in a block} ;
  found = 40 {go here when you've found it} ;
  found1 = 41 {like |found|, when there's more than one per routine} ;
  found2 = 42 {like |found|, when there's more than two per routine} ;
  not_found = 45 {go here when you've found nothing} ;
  common_ending = 50 {go here when you want to merge with another branch} ;
  empty = 0 {symbolic name for a null constant} ;
  first_text_char = 0 {ordinal number of the smallest element of |text_char|} ;
  last_text_char = 255 {ordinal number of the largest element of |text_char|} ;
  null_code = {00=}0 {ASCII code that might disappear} ;
  carriage_return = {015=}13 {ASCII code used at end of line} ;
  invalid_code = {0177=}127 {ASCII code that many systems prohibit in text files} ;
  no_print = 16 {|selector| setting that makes data disappear} ;
  term_only = 17 {printing is destined for the terminal only} ;
  log_only = 18 {printing is destined for the transcript file only} ;
  term_and_log = 19 {normal |selector| setting} ;
  pseudo = 20 {special |selector| setting for |show_context|} ;
  new_string = 21 {printing is deflected to the string pool} ;
  max_selector = 21 {highest selector setting} ;
  batch_mode = 0 {omits all stops and omits terminal output} ;
  nonstop_mode = 1 {omits all stops} ;
  scroll_mode = 2 {omits error stops} ;
  error_stop_mode = 3 {stops at every opportunity to interact} ;
  unspecified_mode = 4 {extra value for command-line switch} ;
  spotless = 0 {|history| value when nothing has been amiss yet} ;
  warning_issued = 1 {|history| value when |begin_diagnostic| has been called} ;
  error_message_issued = 2 {|history| value when |error| has been called} ;
  fatal_error_stop = 3 {|history| value when termination was premature} ;
  inf_bad =  10000 {infinitely bad value} ;
  min_quarterword = 0 {smallest allowable value in a |quarterword|} ;
  max_quarterword = 255 {largest allowable value in a |quarterword|} ;
  hlist_node = 0 {|type| of hlist nodes} ;
  box_node_size = 7 {number of words to allocate for a box node} ;
  width_offset = 1 {position of |width| field in a box node} ;
  depth_offset = 2 {position of |depth| field in a box node} ;
  height_offset = 3 {position of |height| field in a box node} ;
  list_offset = 5 {position of |list_ptr| field in a box node} ;
  normal = 0 {the most common case when several cases are named} ;
  stretching =  1 {glue setting applies to the stretch components} ;
  shrinking =  2 {glue setting applies to the shrink components} ;
  glue_offset =  6 {position of |glue_set| in a box node} ;
  vlist_node = 1 {|type| of vlist nodes} ;
  rule_node = 2 {|type| of rule nodes} ;
  rule_node_size = 4 {number of words to allocate for a rule node} ;
  ins_node = 3 {|type| of insertion nodes} ;
  ins_node_size = 5 {number of words to allocate for an insertion} ;
  mark_node = 4 {|type| of a mark node} ;
  small_node_size = 2 {number of words to allocate for most node types} ;
  adjust_node = 5 {|type| of an adjust node} ;
  ligature_node = 6 {|type| of a ligature node} ;
  disc_node = 7 {|type| of a discretionary node} ;
  whatsit_node = 8 {|type| of special extension nodes} ;
  math_node = 9 {|type| of a math node} ;
  before = 0 {|subtype| for math node that introduces a formula} ;
  after = 1 {|subtype| for math node that winds up a formula} ;
  glue_node = 10 {|type| of node that points to a glue specification} ;
  cond_math_glue = 98 {special |subtype| to suppress glue in the next node} ;
  mu_glue = 99 {|subtype| for math glue} ;
  a_leaders = 100 {|subtype| for aligned leaders} ;
  c_leaders = 101 {|subtype| for centered leaders} ;
  x_leaders = 102 {|subtype| for expanded leaders} ;
  glue_spec_size = 4 {number of words to allocate for a glue specification} ;
  fil = 1 {first-order infinity} ;
  fill = 2 {second-order infinity} ;
  filll = 3 {third-order infinity} ;
  kern_node = 11 {|type| of a kern node} ;
  explicit = 1 {|subtype| of kern nodes from \.[\\kern] and \.[\\/]} ;
  acc_kern = 2 {|subtype| of kern nodes from accents} ;
  penalty_node = 12 {|type| of a penalty node} ;
  inf_penalty = inf_bad {``infinite'' penalty value} ;
  eject_penalty = -inf_penalty {``negatively infinite'' penalty value} ;
  unset_node = 13 {|type| for an unset node} ;
  hi_mem_stat_usage = 14 {the number of one-word nodes always present} ;
  escape = 0 {escape delimiter (called \.\\ in [\sl The \TeX book\/])}
{ \xref[TeXbook][\sl The \TeX book] } ;
  relax = 0 {do nothing ( \.[\\relax] )} ;
  left_brace = 1 {beginning of a group ( \.\[ )} ;
  right_brace = 2 {ending of a group ( \.\] )} ;
  math_shift = 3 {mathematics shift character ( \.\$ )} ;
  tab_mark = 4 {alignment delimiter ( \.\&, \.[\\span] )} ;
  car_ret = 5 {end of line ( |carriage_return|, \.[\\cr], \.[\\crcr] )} ;
  out_param = 5 {output a macro parameter} ;
  mac_param = 6 {macro parameter symbol ( \.\# )} ;
  sup_mark = 7 {superscript ( \.[\char'136] )} ;
  sub_mark = 8 {subscript ( \.[\char'137] )} ;
  ignore = 9 {characters to ignore ( \.[\^\^@] )} ;
  endv = 9 {end of \<v_j> list in alignment template} ;
  spacer = 10 {characters equivalent to blank space ( \.[\ ] )} ;
  letter = 11 {characters regarded as letters ( \.[A..Z], \.[a..z] )} ;
  other_char = 12 {none of the special character types} ;
  active_char = 13 {characters that invoke macros ( \.[\char`\~] )} ;
  par_end = 13 {end of paragraph ( \.[\\par] )} ;
  match = 13 {match a macro parameter} ;
  comment = 14 {characters that introduce comments ( \.\% )} ;
  end_match = 14 {end of parameters to macro} ;
  stop = 14 {end of job ( \.[\\end], \.[\\dump] )} ;
  invalid_char = 15 {characters that shouldn't appear ( \.[\^\^?] )} ;
  delim_num = 15 {specify delimiter numerically ( \.[\\delimiter] )} ;
  max_char_code = 15 {largest catcode for individual characters} ;
  char_num = 16 {character specified numerically ( \.[\\char] )} ;
  math_char_num = 17 {explicit math code ( \.[\\mathchar] )} ;
  mark = 18 {mark definition ( \.[\\mark] )} ;
  xray = 19 {peek inside of \TeX\ ( \.[\\show], \.[\\showbox], etc.~)} ;
  make_box = 20 {make a box ( \.[\\box], \.[\\copy], \.[\\hbox], etc.~)} ;
  hmove = 21 {horizontal motion ( \.[\\moveleft], \.[\\moveright] )} ;
  vmove = 22 {vertical motion ( \.[\\raise], \.[\\lower] )} ;
  un_hbox = 23 {unglue a box ( \.[\\unhbox], \.[\\unhcopy] )} ;
  un_vbox = 24 {unglue a box ( \.[\\unvbox], \.[\\unvcopy] )} ;
  remove_item = 25 {nullify last item ( \.[\\unpenalty],
  \.[\\unkern], \.[\\unskip] )} ;
  hskip = 26 {horizontal glue ( \.[\\hskip], \.[\\hfil], etc.~)} ;
  vskip = 27 {vertical glue ( \.[\\vskip], \.[\\vfil], etc.~)} ;
  mskip = 28 {math glue ( \.[\\mskip] )} ;
  kern = 29 {fixed space ( \.[\\kern] )} ;
  mkern = 30 {math kern ( \.[\\mkern] )} ;
  leader_ship = 31 {use a box ( \.[\\shipout], \.[\\leaders], etc.~)} ;
  halign = 32 {horizontal table alignment ( \.[\\halign] )} ;
  valign = 33 {vertical table alignment ( \.[\\valign] )} ;
  no_align = 34 {temporary escape from alignment ( \.[\\noalign] )} ;
  vrule = 35 {vertical rule ( \.[\\vrule] )} ;
  hrule = 36 {horizontal rule ( \.[\\hrule] )} ;
  insert = 37 {vlist inserted in box ( \.[\\insert] )} ;
  vadjust = 38 {vlist inserted in enclosing paragraph ( \.[\\vadjust] )} ;
  ignore_spaces = 39 {gobble |spacer| tokens ( \.[\\ignorespaces] )} ;
  after_assignment = 40 {save till assignment is done ( \.[\\afterassignment] )} ;
  after_group = 41 {save till group is done ( \.[\\aftergroup] )} ;
  break_penalty = 42 {additional badness ( \.[\\penalty] )} ;
  start_par = 43 {begin paragraph ( \.[\\indent], \.[\\noindent] )} ;
  ital_corr = 44 {italic correction ( \.[\\/] )} ;
  accent = 45 {attach accent in text ( \.[\\accent] )} ;
  math_accent = 46 {attach accent in math ( \.[\\mathaccent] )} ;
  discretionary = 47 {discretionary texts ( \.[\\-], \.[\\discretionary] )} ;
  eq_no = 48 {equation number ( \.[\\eqno], \.[\\leqno] )} ;
  left_right = 49 {variable delimiter ( \.[\\left], \.[\\right] )} ;
  math_comp = 50 {component of formula ( \.[\\mathbin], etc.~)} ;
  limit_switch = 51 {diddle limit conventions ( \.[\\displaylimits], etc.~)} ;
  above = 52 {generalized fraction ( \.[\\above], \.[\\atop], etc.~)} ;
  math_style = 53 {style specification ( \.[\\displaystyle], etc.~)} ;
  math_choice = 54 {choice specification ( \.[\\mathchoice] )} ;
  non_script = 55 {conditional math glue ( \.[\\nonscript] )} ;
  vcenter = 56 {vertically center a vbox ( \.[\\vcenter] )} ;
  case_shift = 57 {force specific case ( \.[\\lowercase], \.[\\uppercase]~)} ;
  message = 58 {send to user ( \.[\\message], \.[\\errmessage] )} ;
  extension = 59 {extensions to \TeX\ ( \.[\\write], \.[\\special], etc.~)} ;
  in_stream = 60 {files for reading ( \.[\\openin], \.[\\closein] )} ;
  begin_group = 61 {begin local grouping ( \.[\\begingroup] )} ;
  end_group = 62 {end local grouping ( \.[\\endgroup] )} ;
  omit = 63 {omit alignment template ( \.[\\omit] )} ;
  ex_space = 64 {explicit space ( \.[\\\ ] )} ;
  no_boundary = 65 {suppress boundary ligatures ( \.[\\noboundary] )} ;
  radical = 66 {square root and similar signs ( \.[\\radical] )} ;
  end_cs_name = 67 {end control sequence ( \.[\\endcsname] )} ;
  min_internal = 68 {the smallest code that can follow \.[\\the]} ;
  char_given = 68 {character code defined by \.[\\chardef]} ;
  math_given = 69 {math code defined by \.[\\mathchardef]} ;
  last_item = 70 {most recent item ( \.[\\lastpenalty],
  \.[\\lastkern], \.[\\lastskip] )} ;
  max_non_prefixed_command = 70 {largest command code that can't be \.[\\global]} ;
  toks_register = 71 {token list register ( \.[\\toks] )} ;
  assign_toks = 72 {special token list ( \.[\\output], \.[\\everypar], etc.~)} ;
  assign_int = 73 {user-defined integer ( \.[\\tolerance], \.[\\day], etc.~)} ;
  assign_dimen = 74 {user-defined length ( \.[\\hsize], etc.~)} ;
  assign_glue = 75 {user-defined glue ( \.[\\baselineskip], etc.~)} ;
  assign_mu_glue = 76 {user-defined muglue ( \.[\\thinmuskip], etc.~)} ;
  assign_font_dimen = 77 {user-defined font dimension ( \.[\\fontdimen] )} ;
  assign_font_int = 78 {user-defined font integer ( \.[\\hyphenchar],
  \.[\\skewchar] )} ;
  set_aux = 79 {specify state info ( \.[\\spacefactor], \.[\\prevdepth] )} ;
  set_prev_graf = 80 {specify state info ( \.[\\prevgraf] )} ;
  set_page_dimen = 81 {specify state info ( \.[\\pagegoal], etc.~)} ;
  set_page_int = 82 {specify state info ( \.[\\deadcycles],
  \.[\\insertpenalties] )} ;
  set_box_dimen = 83 {change dimension of box ( \.[\\wd], \.[\\ht], \.[\\dp] )} ;
  set_shape = 84 {specify fancy paragraph shape ( \.[\\parshape] )} ;
  def_code = 85 {define a character code ( \.[\\catcode], etc.~)} ;
  def_family = 86 {declare math fonts ( \.[\\textfont], etc.~)} ;
  set_font = 87 {set current font ( font identifiers )} ;
  def_font = 88 {define a font file ( \.[\\font] )} ;
  register = 89 {internal register ( \.[\\count], \.[\\dimen], etc.~)} ;
  max_internal = 89 {the largest code that can follow \.[\\the]} ;
  advance = 90 {advance a register or parameter ( \.[\\advance] )} ;
  multiply = 91 {multiply a register or parameter ( \.[\\multiply] )} ;
  divide = 92 {divide a register or parameter ( \.[\\divide] )} ;
  prefix = 93 {qualify a definition ( \.[\\global], \.[\\long], \.[\\outer] )} ;
  let = 94 {assign a command code ( \.[\\let], \.[\\futurelet] )} ;
  shorthand_def = 95 {code definition ( \.[\\chardef], \.[\\countdef], etc.~)}
  {or \.[\\charsubdef]} ;
  read_to_cs = 96 {read into a control sequence ( \.[\\read] )} ;
  def = 97 {macro definition ( \.[\\def], \.[\\gdef], \.[\\xdef], \.[\\edef] )} ;
  set_box = 98 {set a box ( \.[\\setbox] )} ;
  hyph_data = 99 {hyphenation data ( \.[\\hyphenation], \.[\\patterns] )} ;
  set_interaction = 100 {define level of interaction ( \.[\\batchmode], etc.~)} ;
  max_command = 100 {the largest command code seen at |big_switch|} ;
  undefined_cs = max_command+1 {initial state of most |eq_type| fields} ;
  expand_after = max_command+2 {special expansion ( \.[\\expandafter] )} ;
  no_expand = max_command+3 {special nonexpansion ( \.[\\noexpand] )} ;
  input = max_command+4 {input a source file ( \.[\\input], \.[\\endinput] )} ;
  if_test = max_command+5 {conditional text ( \.[\\if], \.[\\ifcase], etc.~)} ;
  fi_or_else = max_command+6 {delimiters for conditionals ( \.[\\else], etc.~)} ;
  cs_name = max_command+7 {make a control sequence from tokens ( \.[\\csname] )} ;
  convert = max_command+8 {convert to text ( \.[\\number], \.[\\string], etc.~)} ;
  the = max_command+9 {expand an internal quantity ( \.[\\the] )} ;
  top_bot_mark = max_command+10 {inserted mark ( \.[\\topmark], etc.~)} ;
  call = max_command+11 {non-long, non-outer control sequence} ;
  long_call = max_command+12 {long, non-outer control sequence} ;
  outer_call = max_command+13 {non-long, outer control sequence} ;
  long_outer_call = max_command+14 {long, outer control sequence} ;
  end_template = max_command+15 {end of an alignment template} ;
  dont_expand = max_command+16 {the following token was marked by \.[\\noexpand]} ;
  glue_ref = max_command+17 {the equivalent points to a glue specification} ;
  shape_ref = max_command+18 {the equivalent points to a parshape specification} ;
  box_ref = max_command+19 {the equivalent points to a box node, or is |null|} ;
  data = max_command+20 {the equivalent is simply a halfword number} ;
  vmode = 1 {vertical mode} ;
  hmode = vmode+max_command+1 {horizontal mode} ;
  mmode = hmode+max_command+1 {math mode} ;
  level_zero = min_quarterword {level for undefined quantities} ;
  level_one = level_zero+1 {outermost level for defined quantities} ;
  active_base = 1 {beginning of region 1, for active character equivalents} ;
  single_base = active_base+256 {equivalents of one-character control sequences} ;
  null_cs = single_base+256 {equivalent of \.[\\csname\\endcsname]} ;
  hash_base = null_cs+1 {beginning of region 2, for the hash table} ;
  frozen_control_sequence = hash_base+hash_size {for error recovery} ;
  frozen_protection = frozen_control_sequence {inaccessible but definable} ;
  frozen_cr = frozen_control_sequence+1 {permanent `\.[\\cr]'} ;
  frozen_end_group = frozen_control_sequence+2 {permanent `\.[\\endgroup]'} ;
  frozen_right = frozen_control_sequence+3 {permanent `\.[\\right]'} ;
  frozen_fi = frozen_control_sequence+4 {permanent `\.[\\fi]'} ;
  frozen_end_template = frozen_control_sequence+5 {permanent `\.[\\endtemplate]'} ;
  frozen_endv = frozen_control_sequence+6 {second permanent `\.[\\endtemplate]'} ;
  frozen_relax = frozen_control_sequence+7 {permanent `\.[\\relax]'} ;
  end_write = frozen_control_sequence+8 {permanent `\.[\\endwrite]'} ;
  frozen_dont_expand = frozen_control_sequence+9
  {permanent `\.[\\notexpanded:]'} ;
  frozen_special = frozen_control_sequence+10
  {permanent `\.[\\special]'} ;
  frozen_null_font = frozen_control_sequence+11
  {permanent `\.[\\nullfont]'} ;
  font_id_base = frozen_null_font-font_base
  {begins table of 257 permanent font identifiers} ;
  undefined_control_sequence = frozen_null_font+max_font_max+1 {dummy location} ;
  glue_base = undefined_control_sequence+1 {beginning of region 3} ;
  line_skip_code = 0 {interline glue if |baseline_skip| is infeasible} ;
  baseline_skip_code = 1 {desired glue between baselines} ;
  par_skip_code = 2 {extra glue just above a paragraph} ;
  above_display_skip_code = 3 {extra glue just above displayed math} ;
  below_display_skip_code = 4 {extra glue just below displayed math} ;
  above_display_short_skip_code = 5
  {glue above displayed math following short lines} ;
  below_display_short_skip_code = 6
  {glue below displayed math following short lines} ;
  left_skip_code = 7 {glue at left of justified lines} ;
  right_skip_code = 8 {glue at right of justified lines} ;
  top_skip_code = 9 {glue at top of main pages} ;
  split_top_skip_code = 10 {glue at top of split pages} ;
  tab_skip_code = 11 {glue between aligned entries} ;
  space_skip_code = 12 {glue between words (if not |zero_glue|)} ;
  xspace_skip_code = 13 {glue after sentences (if not |zero_glue|)} ;
  par_fill_skip_code = 14 {glue on last line of paragraph} ;
  thin_mu_skip_code = 15 {thin space in math formula} ;
  med_mu_skip_code = 16 {medium space in math formula} ;
  thick_mu_skip_code = 17 {thick space in math formula} ;
  glue_pars = 18 {total number of glue parameters} ;
  skip_base = glue_base+glue_pars {table of 256 ``skip'' registers} ;
  mu_skip_base = skip_base+256 {table of 256 ``muskip'' registers} ;
  local_base = mu_skip_base+256 {beginning of region 4} ;
  par_shape_loc = local_base {specifies paragraph shape} ;
  output_routine_loc = local_base+1 {points to token list for \.[\\output]} ;
  every_par_loc = local_base+2 {points to token list for \.[\\everypar]} ;
  every_math_loc = local_base+3 {points to token list for \.[\\everymath]} ;
  every_display_loc = local_base+4 {points to token list for \.[\\everydisplay]} ;
  every_hbox_loc = local_base+5 {points to token list for \.[\\everyhbox]} ;
  every_vbox_loc = local_base+6 {points to token list for \.[\\everyvbox]} ;
  every_job_loc = local_base+7 {points to token list for \.[\\everyjob]} ;
  every_cr_loc = local_base+8 {points to token list for \.[\\everycr]} ;
  err_help_loc = local_base+9 {points to token list for \.[\\errhelp]} ;
  toks_base = local_base+10 {table of 256 token list registers} ;
  box_base = toks_base+256 {table of 256 box registers} ;
  cur_font_loc = box_base+256 {internal font number outside math mode} ;
  math_font_base = cur_font_loc+1 {table of 48 math font numbers} ;
  cat_code_base = math_font_base+48
  {table of 256 command codes (the ``catcodes'')} ;
  lc_code_base = cat_code_base+256 {table of 256 lowercase mappings} ;
  uc_code_base = lc_code_base+256 {table of 256 uppercase mappings} ;
  sf_code_base = uc_code_base+256 {table of 256 spacefactor mappings} ;
  math_code_base = sf_code_base+256 {table of 256 math mode mappings} ;
  char_sub_code_base = math_code_base+256 {table of character substitutions} ;
  int_base = char_sub_code_base+256 {beginning of region 5} ;
  pretolerance_code = 0 {badness tolerance before hyphenation} ;
  tolerance_code = 1 {badness tolerance after hyphenation} ;
  line_penalty_code = 2 {added to the badness of every line} ;
  hyphen_penalty_code = 3 {penalty for break after discretionary hyphen} ;
  ex_hyphen_penalty_code = 4 {penalty for break after explicit hyphen} ;
  club_penalty_code = 5 {penalty for creating a club line} ;
  widow_penalty_code = 6 {penalty for creating a widow line} ;
  display_widow_penalty_code = 7 {ditto, just before a display} ;
  broken_penalty_code = 8 {penalty for breaking a page at a broken line} ;
  bin_op_penalty_code = 9 {penalty for breaking after a binary operation} ;
  rel_penalty_code = 10 {penalty for breaking after a relation} ;
  pre_display_penalty_code = 11
  {penalty for breaking just before a displayed formula} ;
  post_display_penalty_code = 12
  {penalty for breaking just after a displayed formula} ;
  inter_line_penalty_code = 13 {additional penalty between lines} ;
  double_hyphen_demerits_code = 14 {demerits for double hyphen break} ;
  final_hyphen_demerits_code = 15 {demerits for final hyphen break} ;
  adj_demerits_code = 16 {demerits for adjacent incompatible lines} ;
  mag_code = 17 {magnification ratio} ;
  delimiter_factor_code = 18 {ratio for variable-size delimiters} ;
  looseness_code = 19 {change in number of lines for a paragraph} ;
  time_code = 20 {current time of day} ;
  day_code = 21 {current day of the month} ;
  month_code = 22 {current month of the year} ;
  year_code = 23 {current year of our Lord} ;
  show_box_breadth_code = 24 {nodes per level in |show_box|} ;
  show_box_depth_code = 25 {maximum level in |show_box|} ;
  hbadness_code = 26 {hboxes exceeding this badness will be shown by |hpack|} ;
  vbadness_code = 27 {vboxes exceeding this badness will be shown by |vpack|} ;
  pausing_code = 28 {pause after each line is read from a file} ;
  tracing_online_code = 29 {show diagnostic output on terminal} ;
  tracing_macros_code = 30 {show macros as they are being expanded} ;
  tracing_stats_code = 31 {show memory usage if \TeX\ knows it} ;
  tracing_paragraphs_code = 32 {show line-break calculations} ;
  tracing_pages_code = 33 {show page-break calculations} ;
  tracing_output_code = 34 {show boxes when they are shipped out} ;
  tracing_lost_chars_code = 35 {show characters that aren't in the font} ;
  tracing_commands_code = 36 {show command codes at |big_switch|} ;
  tracing_restores_code = 37 {show equivalents when they are restored} ;
  uc_hyph_code = 38 {hyphenate words beginning with a capital letter} ;
  output_penalty_code = 39 {penalty found at current page break} ;
  max_dead_cycles_code = 40 {bound on consecutive dead cycles of output} ;
  hang_after_code = 41 {hanging indentation changes after this many lines} ;
  floating_penalty_code = 42 {penalty for insertions held over after a split} ;
  global_defs_code = 43 {override \.[\\global] specifications} ;
  cur_fam_code = 44 {current family} ;
  escape_char_code = 45 {escape character for token output} ;
  default_hyphen_char_code = 46 {value of \.[\\hyphenchar] when a font is loaded} ;
  default_skew_char_code = 47 {value of \.[\\skewchar] when a font is loaded} ;
  end_line_char_code = 48 {character placed at the right end of the buffer} ;
  new_line_char_code = 49 {character that prints as |print_ln|} ;
  language_code = 50 {current hyphenation table} ;
  left_hyphen_min_code = 51 {minimum left hyphenation fragment size} ;
  right_hyphen_min_code = 52 {minimum right hyphenation fragment size} ;
  holding_inserts_code = 53 {do not remove insertion nodes from \.[\\box255]} ;
  error_context_lines_code = 54 {maximum intermediate line pairs shown} ;
  tex_int_pars = 55 {total number of \TeX's integer parameters} ;
  web2c_int_base = tex_int_pars {base for web2c's integer parameters} ;
  char_sub_def_min_code = web2c_int_base {smallest value in the charsubdef list} ;
  char_sub_def_max_code = web2c_int_base+1 {largest value in the charsubdef list} ;
  tracing_char_sub_def_code = web2c_int_base+2 {traces changes to a charsubdef def} ;
  web2c_int_pars = web2c_int_base+3 {total number of web2c's integer parameters} ;
  int_pars = web2c_int_pars {total number of integer parameters} ;
  count_base = int_base+int_pars {256 user \.[\\count] registers} ;
  del_code_base = count_base+256 {256 delimiter code mappings} ;
  dimen_base = del_code_base+256 {beginning of region 6} ;
  par_indent_code = 0 {indentation of paragraphs} ;
  math_surround_code = 1 {space around math in text} ;
  line_skip_limit_code = 2 {threshold for |line_skip| instead of |baseline_skip|} ;
  hsize_code = 3 {line width in horizontal mode} ;
  vsize_code = 4 {page height in vertical mode} ;
  max_depth_code = 5 {maximum depth of boxes on main pages} ;
  split_max_depth_code = 6 {maximum depth of boxes on split pages} ;
  box_max_depth_code = 7 {maximum depth of explicit vboxes} ;
  hfuzz_code = 8 {tolerance for overfull hbox messages} ;
  vfuzz_code = 9 {tolerance for overfull vbox messages} ;
  delimiter_shortfall_code = 10 {maximum amount uncovered by variable delimiters} ;
  null_delimiter_space_code = 11 {blank space in null delimiters} ;
  script_space_code = 12 {extra space after subscript or superscript} ;
  pre_display_size_code = 13 {length of text preceding a display} ;
  display_width_code = 14 {length of line for displayed equation} ;
  display_indent_code = 15 {indentation of line for displayed equation} ;
  overfull_rule_code = 16 {width of rule that identifies overfull hboxes} ;
  hang_indent_code = 17 {amount of hanging indentation} ;
  h_offset_code = 18 {amount of horizontal offset when shipping pages out} ;
  v_offset_code = 19 {amount of vertical offset when shipping pages out} ;
  emergency_stretch_code = 20 {reduces badnesses on final pass of line-breaking} ;
  dimen_pars = 21 {total number of dimension parameters} ;
  scaled_base = dimen_base+dimen_pars
  {table of 256 user-defined \.[\\dimen] registers} ;
  eqtb_size = scaled_base+255 {largest subscript of |eqtb|} ;
  restore_old_value = 0 {|save_type| when a value should be restored later} ;
  restore_zero = 1 {|save_type| when an undefined entry should be restored} ;
  insert_token = 2 {|save_type| when a token is being saved for later use} ;
  level_boundary = 3 {|save_type| corresponding to beginning of group} ;
  bottom_level = 0 {group code for the outside world} ;
  simple_group = 1 {group code for local structure only} ;
  hbox_group = 2 {code for `\.[\\hbox]\grp'} ;
  adjusted_hbox_group = 3 {code for `\.[\\hbox]\grp' in vertical mode} ;
  vbox_group = 4 {code for `\.[\\vbox]\grp'} ;
  vtop_group = 5 {code for `\.[\\vtop]\grp'} ;
  align_group = 6 {code for `\.[\\halign]\grp', `\.[\\valign]\grp'} ;
  no_align_group = 7 {code for `\.[\\noalign]\grp'} ;
  output_group = 8 {code for output routine} ;
  math_group = 9 {code for, e.g., `\.[\char'136]\grp'} ;
  disc_group = 10 {code for `\.[\\discretionary]\grp\grp\grp'} ;
  insert_group = 11 {code for `\.[\\insert]\grp', `\.[\\vadjust]\grp'} ;
  vcenter_group = 12 {code for `\.[\\vcenter]\grp'} ;
  math_choice_group = 13 {code for `\.[\\mathchoice]\grp\grp\grp\grp'} ;
  semi_simple_group = 14 {code for `\.[\\begingroup...\\endgroup]'} ;
  math_shift_group = 15 {code for `\.[\$...\$]'} ;
  math_left_group = 16 {code for `\.[\\left...\\right]'} ;
  max_group_code = 16 ;
  left_brace_token = {0400=}256 {$2^8\cdot|left_brace|$} ;
  left_brace_limit = {01000=}512 {$2^8\cdot(|left_brace|+1)$} ;
  right_brace_token = {01000=}512 {$2^8\cdot|right_brace|$} ;
  right_brace_limit = {01400=}768 {$2^8\cdot(|right_brace|+1)$} ;
  math_shift_token = {01400=}768 {$2^8\cdot|math_shift|$} ;
  tab_token = {02000=}1024 {$2^8\cdot|tab_mark|$} ;
  out_param_token = {02400=}1280 {$2^8\cdot|out_param|$} ;
  space_token = {05040=}2592 {$2^8\cdot|spacer|+|" "|$} ;
  letter_token = {05400=}2816 {$2^8\cdot|letter|$} ;
  other_token = {06000=}3072 {$2^8\cdot|other_char|$} ;
  match_token = {06400=}3328 {$2^8\cdot|match|$} ;
  end_match_token = {07000=}3584 {$2^8\cdot|end_match|$} ;
  mid_line = 1 {|state| code when scanning a line of characters} ;
  skip_blanks = 2+max_char_code {|state| code when ignoring blanks} ;
  new_line = 3+max_char_code+max_char_code {|state| code at start of line} ;
  skipping = 1 {|scanner_status| when passing conditional text} ;
  defining = 2 {|scanner_status| when reading a macro definition} ;
  matching = 3 {|scanner_status| when reading macro arguments} ;
  aligning = 4 {|scanner_status| when reading an alignment preamble} ;
  absorbing = 5 {|scanner_status| when reading a balanced text} ;
  token_list = 0 {|state| code when scanning a token list} ;
  parameter = 0 {|token_type| code for parameter} ;
  u_template = 1 {|token_type| code for \<u_j> template} ;
  v_template = 2 {|token_type| code for \<v_j> template} ;
  backed_up = 3 {|token_type| code for text to be reread} ;
  inserted = 4 {|token_type| code for inserted texts} ;
  macro = 5 {|token_type| code for defined control sequences} ;
  output_text = 6 {|token_type| code for output routines} ;
  every_par_text = 7 {|token_type| code for \.[\\everypar]} ;
  every_math_text = 8 {|token_type| code for \.[\\everymath]} ;
  every_display_text = 9 {|token_type| code for \.[\\everydisplay]} ;
  every_hbox_text = 10 {|token_type| code for \.[\\everyhbox]} ;
  every_vbox_text = 11 {|token_type| code for \.[\\everyvbox]} ;
  every_job_text = 12 {|token_type| code for \.[\\everyjob]} ;
  every_cr_text = 13 {|token_type| code for \.[\\everycr]} ;
  mark_text = 14 {|token_type| code for \.[\\topmark], etc.} ;
  write_text = 15 {|token_type| code for \.[\\write]} ;
  switch = 25 {a label in |get_next|} ;
  start_cs = 26 {another} ;
  no_expand_flag = 257 {this characterizes a special variant of |relax|} ;
  top_mark_code = 0 {the mark in effect at the previous page break} ;
  first_mark_code = 1 {the first mark between |top_mark| and |bot_mark|} ;
  bot_mark_code = 2 {the mark in effect at the current page break} ;
  split_first_mark_code = 3 {the first mark found by \.[\\vsplit]} ;
  split_bot_mark_code = 4 {the last mark found by \.[\\vsplit]} ;
  int_val = 0 {integer values} ;
  dimen_val = 1 {dimension values} ;
  glue_val = 2 {glue specifications} ;
  mu_val = 3 {math glue specifications} ;
  ident_val = 4 {font identifier} ;
  tok_val = 5 {token lists} ;
  input_line_no_code = glue_val+1 {code for \.[\\inputlineno]} ;
  badness_code = glue_val+2 {code for \.[\\badness]} ;
  octal_token = other_token+{"'"=}39 {apostrophe, indicates an octal constant} ;
  hex_token = other_token+{""""=}34 {double quote, indicates a hex constant} ;
  alpha_token = other_token+{"`"=}96 {reverse apostrophe, precedes alpha constants} ;
  point_token = other_token+{"."=}46 {decimal point} ;
  continental_point_token = other_token+{","=}44 {decimal point, Eurostyle} ;
  zero_token = other_token+{"0"=}48 {zero, the smallest digit} ;
  A_token = letter_token+{"A"=}65 {the smallest special hex digit} ;
  other_A_token = other_token+{"A"=}65 {special hex digit of type |other_char|} ;
  attach_fraction = 88 {go here to pack |cur_val| and |f| into |cur_val|} ;
  attach_sign = 89 {go here when |cur_val| is correct except perhaps for sign} ;
  default_rule = 26214 {0.4\thinspace pt} ;
  number_code = 0 {command code for \.[\\number]} ;
  roman_numeral_code = 1 {command code for \.[\\romannumeral]} ;
  string_code = 2 {command code for \.[\\string]} ;
  meaning_code = 3 {command code for \.[\\meaning]} ;
  font_name_code = 4 {command code for \.[\\fontname]} ;
  job_name_code = 5 {command code for \.[\\jobname]} ;
  closed = 2 {not open, or at end of file} ;
  just_open = 1 {newly opened, first line not yet read} ;
  if_char_code = 0 { `\.[\\if]' } ;
  if_cat_code = 1 { `\.[\\ifcat]' } ;
  if_int_code = 2 { `\.[\\ifnum]' } ;
  if_dim_code = 3 { `\.[\\ifdim]' } ;
  if_odd_code = 4 { `\.[\\ifodd]' } ;
  if_vmode_code = 5 { `\.[\\ifvmode]' } ;
  if_hmode_code = 6 { `\.[\\ifhmode]' } ;
  if_mmode_code = 7 { `\.[\\ifmmode]' } ;
  if_inner_code = 8 { `\.[\\ifinner]' } ;
  if_void_code = 9 { `\.[\\ifvoid]' } ;
  if_hbox_code = 10 { `\.[\\ifhbox]' } ;
  if_vbox_code = 11 { `\.[\\ifvbox]' } ;
  ifx_code = 12 { `\.[\\ifx]' } ;
  if_eof_code = 13 { `\.[\\ifeof]' } ;
  if_true_code = 14 { `\.[\\iftrue]' } ;
  if_false_code = 15 { `\.[\\iffalse]' } ;
  if_case_code = 16 { `\.[\\ifcase]' } ;
  if_node_size = 2 {number of words in stack entry for conditionals} ;
  if_code = 1 {code for \.[\\if...] being evaluated} ;
  fi_code = 2 {code for \.[\\fi]} ;
  else_code = 3 {code for \.[\\else]} ;
  or_code = 4 {code for \.[\\or]} ;
  format_area_length = 0 {length of its area part} ;
  format_ext_length = 4 {length of its `\.[.fmt]' part} ;
  format_extension = {".fmt"=}794 {the extension, as a \.[WEB] constant} ;
  no_tag = 0 {vanilla character} ;
  lig_tag = 1 {character has a ligature/kerning program} ;
  list_tag = 2 {character has a successor in a charlist} ;
  ext_tag = 3 {character is extensible} ;
  slant_code = 1 ;
  space_code = 2 ;
  space_stretch_code = 3 ;
  space_shrink_code = 4 ;
  x_height_code = 5 ;
  quad_code = 6 ;
  extra_space_code = 7 ;
  non_address = 0 {a spurious |bchar_label|} ;
  bad_tfm = 11 {label for |read_font_info|} ;
  set_char_0 = 0 {typeset character 0 and move right} ;
  set1 = 128 {typeset a character and move right} ;
  set_rule = 132 {typeset a rule and move right} ;
  put_rule = 137 {typeset a rule} ;
  nop = 138 {no operation} ;
  bop = 139 {beginning of page} ;
  eop = 140 {ending of page} ;
  push = 141 {save the current positions} ;
  pop = 142 {restore previous positions} ;
  right1 = 143 {move right} ;
  w0 = 147 {move right by |w|} ;
  w1 = 148 {move right and set |w|} ;
  x0 = 152 {move right by |x|} ;
  x1 = 153 {move right and set |x|} ;
  down1 = 157 {move down} ;
  y0 = 161 {move down by |y|} ;
  y1 = 162 {move down and set |y|} ;
  z0 = 166 {move down by |z|} ;
  z1 = 167 {move down and set |z|} ;
  fnt_num_0 = 171 {set current font to 0} ;
  fnt1 = 235 {set current font} ;
  xxx1 = 239 {extension to \.[DVI] primitives} ;
  xxx4 = 242 {potentially long extension to \.[DVI] primitives} ;
  fnt_def1 = 243 {define the meaning of a font number} ;
  pre = 247 {preamble} ;
  post = 248 {postamble beginning} ;
  post_post = 249 {postamble ending} ;
  id_byte = 2 {identifies the kind of \.[DVI] files described here} ;
  movement_node_size = 3 {number of words per entry in the down and right stacks} ;
  y_here = 1 {|info| when the movement entry points to a |y| command} ;
  z_here = 2 {|info| when the movement entry points to a |z| command} ;
  yz_OK = 3 {|info| corresponding to an unconstrained \\[down] command} ;
  y_OK = 4 {|info| corresponding to a \\[down] that can't become a |z|} ;
  z_OK = 5 {|info| corresponding to a \\[down] that can't become a |y|} ;
  d_fixed = 6 {|info| corresponding to a \\[down] that can't change} ;
  none_seen = 0 {no |y_here| or |z_here| nodes have been encountered yet} ;
  y_seen = 6 {we have seen |y_here| but not |z_here|} ;
  z_seen = 12 {we have seen |z_here| but not |y_here|} ;
  move_past = 13 {go to this label when advancing past glue or a rule} ;
  fin_rule = 14 {go to this label to finish processing a rule} ;
  next_p = 15 {go to this label when finished with node |p|} ;
  exactly = 0 {a box dimension is pre-specified} ;
  additional = 1 {a box dimension is increased from the natural one} ;
  noad_size = 4 {number of words in a normal noad} ;
  math_char = 1 {|math_type| when the attribute is simple} ;
  sub_box = 2 {|math_type| when the attribute is a box} ;
  sub_mlist = 3 {|math_type| when the attribute is a formula} ;
  math_text_char = 4 {|math_type| when italic correction is dubious} ;
  ord_noad = unset_node+3 {|type| of a noad classified Ord} ;
  op_noad = ord_noad+1 {|type| of a noad classified Op} ;
  bin_noad = ord_noad+2 {|type| of a noad classified Bin} ;
  rel_noad = ord_noad+3 {|type| of a noad classified Rel} ;
  open_noad = ord_noad+4 {|type| of a noad classified Open} ;
  close_noad = ord_noad+5 {|type| of a noad classified Close} ;
  punct_noad = ord_noad+6 {|type| of a noad classified Punct} ;
  inner_noad = ord_noad+7 {|type| of a noad classified Inner} ;
  limits = 1 {|subtype| of |op_noad| whose scripts are to be above, below} ;
  no_limits = 2 {|subtype| of |op_noad| whose scripts are to be normal} ;
  radical_noad = inner_noad+1 {|type| of a noad for square roots} ;
  radical_noad_size = 5 {number of |mem| words in a radical noad} ;
  fraction_noad = radical_noad+1 {|type| of a noad for generalized fractions} ;
  fraction_noad_size = 6 {number of |mem| words in a fraction noad} ;
  under_noad = fraction_noad+1 {|type| of a noad for underlining} ;
  over_noad = under_noad+1 {|type| of a noad for overlining} ;
  accent_noad = over_noad+1 {|type| of a noad for accented subformulas} ;
  accent_noad_size = 5 {number of |mem| words in an accent noad} ;
  vcenter_noad = accent_noad+1 {|type| of a noad for \.[\\vcenter]} ;
  left_noad = vcenter_noad+1 {|type| of a noad for \.[\\left]} ;
  right_noad = left_noad+1 {|type| of a noad for \.[\\right]} ;
  style_node = unset_node+1 {|type| of a style node} ;
  style_node_size = 3 {number of words in a style node} ;
  display_style = 0 {|subtype| for \.[\\displaystyle]} ;
  text_style = 2 {|subtype| for \.[\\textstyle]} ;
  script_style = 4 {|subtype| for \.[\\scriptstyle]} ;
  script_script_style = 6 {|subtype| for \.[\\scriptscriptstyle]} ;
  cramped = 1 {add this to an uncramped style if you want to cramp it} ;
  choice_node = unset_node+2 {|type| of a choice node} ;
  text_size = 0 {size code for the largest size in a family} ;
  script_size = 16 {size code for the medium size in a family} ;
  script_script_size = 32 {size code for the smallest size in a family} ;
  total_mathsy_params = 22 ;
  total_mathex_params = 13 ;
  done_with_noad = 80 {go here when a noad has been fully translated} ;
  done_with_node = 81 {go here when a node has been fully converted} ;
  check_dimensions = 82 {go here to update |max_h| and |max_d|} ;
  delete_q = 83 {go here to delete |q| and move to the next node} ;
  math_spacing =  

{ \hskip-35pt }
{"0234000122*4000133**3**344*0400400*000000234000111*1111112341011"=}906
{ $ \hskip-35pt$ } ;
  align_stack_node_size = 5 {number of |mem| words to save alignment states} ;
  span_code = 256 {distinct from any character} ;
  cr_code = 257 {distinct from |span_code| and from any character} ;
  cr_cr_code = cr_code+1 {this distinguishes \.[\\crcr] from \.[\\cr]} ;
  span_node_size = 2 {number of |mem| words for a span node} ;
  tight_fit = 3 {fitness classification for lines shrinking 0.5 to 1.0 of their
  shrinkability} ;
  loose_fit = 1 {fitness classification for lines stretching 0.5 to 1.0 of their
  stretchability} ;
  very_loose_fit = 0 {fitness classification for lines stretching more than
  their stretchability} ;
  decent_fit = 2 {fitness classification for all other lines} ;
  active_node_size = 3 {number of words in active nodes} ;
  unhyphenated = 0 {the |type| of a normal active break node} ;
  hyphenated = 1 {the |type| of an active node that breaks at a |disc_node|} ;
  passive_node_size = 2 {number of words in passive nodes} ;
  delta_node_size = 7 {number of words in a delta node} ;
  delta_node = 2 {|type| field in a delta node} ;
  deactivate = 60 {go here when node |r| should be deactivated} ;
  update_heights = 90 {go here to record glue in the |active_height| table} ;
  inserts_only = 1
  {|page_contents| when an insert node has been contributed, but no boxes} ;
  box_there = 2 {|page_contents| when a box or rule has been contributed} ;
  page_ins_node_size = 4 {number of words for a page insertion node} ;
  inserting = 0 {an insertion class that has not yet overflowed} ;
  split_up = 1 {an overflowed insertion class} ;
  contribute = 80 {go here to link a node into the current page} ;
  big_switch = 60 {go here to branch on the next token of input} ;
  main_loop = 70 {go here to typeset a string of consecutive characters} ;
  main_loop_wrapup = 80 {go here to finish a character or ligature} ;
  main_loop_move = 90 {go here to advance the ligature cursor} ;
  main_loop_move_lig = 95 {same, when advancing past a generated ligature} ;
  main_loop_lookahead = 100 {go here to bring in another character, if any} ;
  main_lig_loop = 110 {go here to check for ligatures or kerning} ;
  append_normal_space = 120 {go here to append a normal space between words} ;
  fil_code = 0 {identifies \.[\\hfil] and \.[\\vfil]} ;
  fill_code = 1 {identifies \.[\\hfill] and \.[\\vfill]} ;
  ss_code = 2 {identifies \.[\\hss] and \.[\\vss]} ;
  fil_neg_code = 3 {identifies \.[\\hfilneg] and \.[\\vfilneg]} ;
  skip_code = 4 {identifies \.[\\hskip] and \.[\\vskip]} ;
  mskip_code = 5 {identifies \.[\\mskip]} ;
  box_code = 0 {|chr_code| for `\.[\\box]'} ;
  copy_code = 1 {|chr_code| for `\.[\\copy]'} ;
  last_box_code = 2 {|chr_code| for `\.[\\lastbox]'} ;
  vsplit_code = 3 {|chr_code| for `\.[\\vsplit]'} ;
  vtop_code = 4 {|chr_code| for `\.[\\vtop]'} ;
  above_code = 0 { `\.[\\above]' } ;
  over_code = 1 { `\.[\\over]' } ;
  atop_code = 2 { `\.[\\atop]' } ;
  delimited_code = 3 { `\.[\\abovewithdelims]', etc.} ;
  char_def_code = 0 {|shorthand_def| for \.[\\chardef]} ;
  math_char_def_code = 1 {|shorthand_def| for \.[\\mathchardef]} ;
  count_def_code = 2 {|shorthand_def| for \.[\\countdef]} ;
  dimen_def_code = 3 {|shorthand_def| for \.[\\dimendef]} ;
  skip_def_code = 4 {|shorthand_def| for \.[\\skipdef]} ;
  mu_skip_def_code = 5 {|shorthand_def| for \.[\\muskipdef]} ;
  toks_def_code = 6 {|shorthand_def| for \.[\\toksdef]} ;
  char_sub_def_code = 7 {|shorthand_def| for \.[\\charsubdef]} ;
  show_code = 0 { \.[\\show] } ;
  show_box_code = 1 { \.[\\showbox] } ;
  show_the_code = 2 { \.[\\showthe] } ;
  show_lists_code = 3 { \.[\\showlists] } ;
  bad_fmt = 6666 {go here if the format file is unacceptable} ;
  breakpoint = 888 {place where a breakpoint is desirable} ;
  write_node_size = 2 {number of words in a write/whatsit node} ;
  open_node_size = 3 {number of words in an open/whatsit node} ;
  open_node = 0 {|subtype| in whatsits that represent files to \.[\\openout]} ;
  write_node = 1 {|subtype| in whatsits that represent things to \.[\\write]} ;
  close_node = 2 {|subtype| in whatsits that represent streams to \.[\\closeout]} ;
  special_node = 3 {|subtype| in whatsits that represent \.[\\special] things} ;
  language_node = 4 {|subtype| in whatsits that change the current language} ;
  immediate_code = 4 {command modifier for \.[\\immediate]} ;
  set_language_code = 5 {command modifier for \.[\\setlanguage]} ;
 
{ Constants in the outer block }
 hash_offset=514; {smallest index in hash array, i.e., |hash_base| }
  {Use |hash_offset=0| for compilers which cannot decrement pointers.}
 trie_op_size=35111; {space for ``opcodes'' in the hyphenation patterns;
  best if relatively prime to 313, 361, and 1009.}
 neg_trie_op_size=-35111; {for lower |trie_op_hash| array bound;
  must be equal to |-trie_op_size|.}
 min_trie_op=0; {first possible trie op code for any language}
 max_trie_op= 65535 ; {largest possible trie opcode for any language}
 pool_name=TEXMF_POOL_NAME; {this is configurable, for the sake of ML-\TeX}
  {string of length |file_name_size|; tells where the string pool appears}
 engine_name=TEXMF_ENGINE_NAME; {the name of this engine}


 inf_mem_bot = 0;
 sup_mem_bot = 1;

 inf_main_memory = 3000;
 sup_main_memory = 256000000;

 inf_trie_size = 8000;
 sup_trie_size =  {0x3fffff=}4194303 ;

 inf_max_strings = 3000;
 sup_max_strings =  2097151 ;
 inf_strings_free = 100;
 sup_strings_free = sup_max_strings;

 inf_buf_size = 500;
 sup_buf_size = 30000000;

 inf_nest_size = 40;
 sup_nest_size = 4000;

 inf_max_in_open = 6;
 sup_max_in_open = 127;

 inf_param_size = 60;
 sup_param_size = 32767;

 inf_save_size = 600;
 sup_save_size = 30000000;

 inf_stack_size = 200;
 sup_stack_size = 30000;

 inf_dvi_buf_size = 800;
 sup_dvi_buf_size = 65536;

 inf_font_mem_size = 20000;
 sup_font_mem_size = 147483647; {|integer|-limited, so 2 could be prepended?}

 sup_font_max = max_font_max;
 inf_font_max = 50; {could be smaller, but why?}

 inf_pool_size = 32000;
 sup_pool_size = 40000000;
 inf_pool_free = 1000;
 sup_pool_free = sup_pool_size;
 inf_string_vacancies = 8000;
 sup_string_vacancies = sup_pool_size - 23000;

 sup_hash_extra = sup_max_strings;
 inf_hash_extra = 0;

 sup_hyph_size =  65535 ;
 inf_hyph_size =  610 ; {Must be not less than |hyph_prime|!}

 inf_expand_depth = 10;
 sup_expand_depth = 10000000;
{ \xref[TeXformats] }



type  
{ Types in the outer block }
 ASCII_code=0..255; {eight-bit numbers}


 eight_bits=0..255; {unsigned one-byte quantity}
 alpha_file=packed file of  ASCII_code ; {files that contain textual data}
 byte_file=packed file of eight_bits; {files that contain binary data}


 pool_pointer = integer; {for variables that point into |str_pool|}
 str_number = 0.. 2097151 ; {for variables that point into |str_start|}
 packed_ASCII_code = 0..255; {elements of |str_pool| array}


 scaled = integer; {this type is used for scaled integers}
 nonnegative_integer=0..{017777777777=}2147483647; {$0\L x<2^[31]$}
 small_number=0..63; {this type is self-explanatory}





 quarterword = min_quarterword..max_quarterword; {1/4 of a word}
 halfword=-{0xfffffff=}268435455 ..{0xfffffff=}268435455 ; {1/2 of a word}
 two_choices = 1..2; {used when there are two variants in a record}
 four_choices = 1..4; {used when there are four variants in a record}
 two_halves = packed record 

   rh:halfword;
  case two_choices of
  1: ( lh:halfword);
  2: ( b0:quarterword;  b1:quarterword);
  end;
 four_quarters = packed record 

   b0:quarterword;
   b1:quarterword;
   b2:quarterword;
   b3:quarterword;
  end;
 memory_word = record 

  case four_choices of
  1: ( int:integer);
  2: ( gr:glue_ratio);
  3: ( hh:two_halves);
  4: ( qqqq:four_quarters);
  end;
 word_file = file of memory_word;


 glue_ord=normal..filll; {infinity to the 0, 1, 2, or 3 power}


 list_state_record=record mode_field:-mmode..mmode; 
   head_field, tail_field: halfword ;
   pg_field, ml_field: integer; 
   aux_field: memory_word;
  end;


 group_code=0..max_group_code; {|save_level| for a level boundary}


 in_state_record = record
   state_field,  index_field: quarterword;
   start_field, loc_field,  limit_field,  name_field: halfword;
  end;


 internal_font_number=integer; {|font| in a |char_node|}
 font_index=integer; {index into |font_info|}
 nine_bits=min_quarterword.. 256  ;


 dvi_index=0..dvi_buf_size; {an index into the output buffer}


 trie_pointer=0.. {0x3fffff=}4194303 ; {an index into |trie|}
 trie_opcode=0.. 65535 ;  {a trie opcode}


 hyph_pointer=0.. 65535 ; {index into hyphen exceptions hash table;
                     enlarging this requires changing (un)dump code}



var 
{ Global variables }
 bad:integer; {is some ``constant'' wrong?}


 xord: array [ ASCII_code ] of ASCII_code;
  {specifies conversion of input characters}
xchr: array [ASCII_code] of  ASCII_code ;
   { specifies conversion of output characters }
xprn: array [ASCII_code] of ASCII_code;
   { non zero iff character is printable }


 name_of_file:^ ASCII_code ;
 name_length:0.. maxint ;
{this many characters are actually
  relevant in |name_of_file| (the rest are blank)}


 buffer:^ASCII_code; {lines of characters being read}
 first:0..buf_size; {the first unused position in |buffer|}
 last:0..buf_size; {end of the line just input to |buffer|}
 max_buf_stack:0..buf_size; {largest index used in |buffer|}


 ifdef('INITEX') 
 ini_version:boolean; {are we \.[INITEX]?}
 dump_option:boolean; {was the dump name option used?}
 dump_line:boolean; {was a \.[\%\AM format] line seen?}
endif('INITEX') 



 dump_name:const_cstring; {format name for terminal display}


 bound_default:integer; {temporary for setup}
 bound_name:const_cstring; {temporary for setup}


 mem_bot:integer;{smallest index in the |mem| array dumped by \.[INITEX];
  must not be less than |mem_min|}
 main_memory:integer; {total memory words allocated in initex}
 extra_mem_bot:integer; {|mem_min:=mem_bot-extra_mem_bot| except in \.[INITEX]}
 mem_min:integer; {smallest index in \TeX's internal |mem| array;
  must be |min_halfword| or more;
  must be equal to |mem_bot| in \.[INITEX], otherwise |<=mem_bot|}
 mem_top:integer; {largest index in the |mem| array dumped by \.[INITEX];
  must be substantially larger than |mem_bot|,
  equal to |mem_max| in \.[INITEX], else not greater than |mem_max|}
 extra_mem_top:integer; {|mem_max:=mem_top+extra_mem_top| except in \.[INITEX]}
 mem_max:integer; {greatest index in \TeX's internal |mem| array;
  must be strictly less than |max_halfword|;
  must be equal to |mem_top| in \.[INITEX], otherwise |>=mem_top|}
 error_line:integer; {width of context lines on terminal error messages}
 half_error_line:integer; {width of first lines of contexts in terminal
  error messages; should be between 30 and |error_line-15|}
 max_print_line:integer;
  {width of longest text lines output; should be at least 60}
 max_strings:integer; {maximum number of strings; must not exceed |max_halfword|}
 strings_free:integer; {strings available after format loaded}
 string_vacancies:integer; {the minimum number of characters that should be
  available for the user's control sequences and font names,
  after \TeX's own error messages are stored}
 pool_size:integer; {maximum number of characters in strings, including all
  error messages and help texts, and the names of all fonts and
  control sequences; must exceed |string_vacancies| by the total
  length of \TeX's own strings, which is currently about 23000}
 pool_free:integer;{pool space free after format loaded}
 font_mem_size:integer; {number of words of |font_info| for all fonts}
 font_max:integer; {maximum internal font number; ok to exceed |max_quarterword|
  and must be at most |font_base|+|max_font_max|}
 font_k:integer; {loop variable for initialization}
 hyph_size:integer; {maximum number of hyphen exceptions}
 trie_size:integer; {space for hyphenation patterns; should be larger for
  \.[INITEX] than it is in production versions of \TeX.  50000 is
  needed for English, German, and Portuguese.}
 buf_size:integer; {maximum number of characters simultaneously present in
  current lines of open files and in control sequences between
  \.[\\csname] and \.[\\endcsname]; must not exceed |max_halfword|}
 stack_size:integer; {maximum number of simultaneous input sources}
 max_in_open:integer; {maximum number of input files and error insertions that
  can be going on simultaneously}
 param_size:integer; {maximum number of simultaneous macro parameters}
 nest_size:integer; {maximum number of semantic levels simultaneously active}
 save_size:integer; {space for saving values outside of current group; must be
  at most |max_halfword|}
 dvi_buf_size:integer; {size of the output buffer; must be a multiple of 8}
 expand_depth:integer; {limits recursive calls to the |expand| procedure}
 parse_first_line_p:cinttype; {parse the first line for options}
 file_line_error_style_p:cinttype; {format messages as file:line:error}
 eight_bit_p:cinttype; {make all characters printable by default}
 halt_on_error_p:cinttype; {stop at first error}
 halting_on_error_p:boolean; {already trying to halt?}
 quoted_filename:boolean; {current filename is quoted}
{Variables for source specials}
 src_specials_p : boolean;{Whether |src_specials| are enabled at all}
 insert_src_special_auto : boolean;
 insert_src_special_every_par : boolean;
 insert_src_special_every_parend : boolean;
 insert_src_special_every_cr : boolean;
 insert_src_special_every_math : boolean;
 insert_src_special_every_hbox : boolean;
 insert_src_special_every_vbox : boolean;
 insert_src_special_every_display : boolean;


 str_pool: ^packed_ASCII_code; {the characters}
 str_start : ^pool_pointer; {the starting pointers}
 pool_ptr : pool_pointer; {first unused position in |str_pool|}
 str_ptr : str_number; {number of the current string being created}
 init_pool_ptr : pool_pointer; {the starting value of |pool_ptr|}
 init_str_ptr : str_number; {the starting value of |str_ptr|}


 ifdef('INITEX')   pool_file:alpha_file; {the string-pool file output by \.[TANGLE]}
endif('INITEX') 


 log_file : alpha_file; {transcript of \TeX\ session}
 selector : 0..max_selector; {where to print a message}
 dig : array[0..22] of 0..15; {digits in a number being output}
 tally : integer; {the number of characters recently printed}
 term_offset : 0..max_print_line;
  {the number of characters on the current terminal line}
 file_offset : 0..max_print_line;
  {the number of characters on the current file line}
 trick_buf:array[0..ssup_error_line] of ASCII_code; {circular buffer for
  pseudoprinting}
 trick_count: integer; {threshold for pseudoprinting, explained later}
 first_count: integer; {another variable for pseudoprinting}


 interaction:batch_mode..error_stop_mode; {current level of interaction}
 interaction_option:batch_mode..unspecified_mode; {set from command line}


 deletions_allowed:boolean; {is it safe for |error| to call |get_token|?}
 set_box_allowed:boolean; {is it safe to do a \.[\\setbox] assignment?}
 history:spotless..fatal_error_stop; {has the source input been clean so far?}
 error_count:-1..100; {the number of scrolled errors since the
  last paragraph ended}


 help_line:array[0..5] of str_number; {helps for the next |error|}
 help_ptr:0..6; {the number of help lines present}
 use_err_help:boolean; {should the |err_help| list be shown?}


 interrupt:integer; {should \TeX\ pause for instructions?}
 OK_to_interrupt:boolean; {should interrupts be observed?}


 save_arith_error:boolean;
 arith_error:boolean; {has arithmetic overflow occurred recently?}
 tex_remainder :scaled; {amount subtracted to get an exact division}


 temp_ptr:halfword ; {a pointer variable for occasional emergency use}


 yzmem : ^memory_word; {the big dynamic storage area}
 zmem : ^memory_word; {the big dynamic storage area}
 lo_mem_max : halfword ; {the largest location of variable-size memory in use}
 hi_mem_min : halfword ; {the smallest location of one-word memory in use}


 var_used,  dyn_used : integer; {how much memory is in use}


 avail : halfword ; {head of the list of available one-word nodes}
 mem_end : halfword ; {the last one-word node used in |mem|}


 rover : halfword ; {points to some node in the list of empties}


{The debug memory arrays have not been mallocated yet.}
 ifdef('TEXMF_DEBUG')   free_arr : packed array [0..9] of boolean; {free cells}
{ \hskip10pt } was_free: packed array [0..9] of boolean;
  {previously free cells}
{ \hskip10pt } was_mem_end, was_lo_max, was_hi_min: halfword ;
  {previous |mem_end|, |lo_mem_max|, and |hi_mem_min|}
{ \hskip10pt } panicking:boolean; {do we want to check memory constantly?}
endif('TEXMF_DEBUG') 


 font_in_short_display:integer; {an internal font number}


 depth_threshold : integer; {maximum nesting depth in box displays}
 breadth_max : integer; {maximum number of items shown at the same list level}


 nest:^list_state_record;
 nest_ptr:0..nest_size; {first unused location of |nest|}
 max_nest_stack:0..nest_size; {maximum of |nest_ptr| when pushing}
 cur_list:list_state_record; {the ``top'' semantic state}
 shown_mode:-mmode..mmode; {most recent mode shown by \.[\\tracingcommands]}


 old_setting:0..max_selector;
 sys_time, sys_day, sys_month, sys_year:integer;
    {date and time supplied by external system}


 zeqtb:^memory_word;
 xeq_level:array[int_base..eqtb_size] of quarterword;


 hash: ^two_halves; {the hash table}
 yhash: ^two_halves; {auxiliary pointer for freeing hash}
 hash_used:halfword ; {allocation pointer for |hash|}
 hash_extra:halfword ; {|hash_extra=hash| above |eqtb_size|}
 hash_top:halfword ; {maximum of the hash array}
 eqtb_top:halfword ; {maximum of the |eqtb|}
 hash_high:halfword ; {pointer to next high hash location}
 no_new_control_sequence:boolean; {are new identifiers legal?}
 cs_count:integer; {total number of known identifiers}


 save_stack : ^memory_word;
 save_ptr : 0..save_size; {first unused entry on |save_stack|}
 max_save_stack:0..save_size; {maximum usage of save stack}
 cur_level: quarterword; {current nesting level for groups}
 cur_group: group_code; {current group type}
 cur_boundary: 0..save_size; {where the current level begins}


 mag_set:integer; {if nonzero, this magnification should be used henceforth}


 cur_cmd: eight_bits; {current command set by |get_next|}
 cur_chr: halfword; {operand of current command}
 cur_cs: halfword ; {control sequence found here, zero if none found}
 cur_tok: halfword; {packed representative of |cur_cmd| and |cur_chr|}


 input_stack : ^in_state_record;
 input_ptr : 0..stack_size; {first unused location of |input_stack|}
 max_in_stack: 0..stack_size; {largest value of |input_ptr| when pushing}
 cur_input : in_state_record;
  {the ``top'' input state, according to convention (1)}


 in_open : 0..max_in_open; {the number of lines in the buffer, less one}
 open_parens : 0..max_in_open; {the number of open text files}
 input_file : ^alpha_file;
 line : integer; {current line number in the current source file}
 line_stack : ^integer;
 source_filename_stack : ^str_number;
 full_source_filename_stack : ^str_number;


 scanner_status : normal..absorbing; {can a subfile end now?}
 warning_index : halfword ; {identifier relevant to non-|normal| scanner status}
 def_ref : halfword ; {reference count of token list being defined}


 param_stack: ^halfword ;
  {token list pointers for parameters}
 param_ptr:0..param_size; {first unused entry in |param_stack|}
 max_param_stack:integer;
  {largest value of |param_ptr|, will be |<=param_size+9|}


 align_state:integer; {group level with respect to current alignment}


 base_ptr:0..stack_size; {shallowest level shown by |show_context|}


 par_loc:halfword ; {location of `\.[\\par]' in |eqtb|}
 par_token:halfword; {token representing `\.[\\par]'}


 force_eof:boolean; {should the next \.[\\input] be aborted early?}


 cur_mark:array[top_mark_code..split_bot_mark_code] of halfword ;
  {token lists for marks}


 long_state:call..long_outer_call; {governs the acceptance of \.[\\par]}


 pstack:array[0..8] of halfword ; {arguments supplied to a macro}


 cur_val:integer; {value returned by numeric scanners}
 cur_val_level:int_val..tok_val; {the ``level'' of this value}


 radix:small_number; {|scan_int| sets this to 8, 10, 16, or zero}


 cur_order:glue_ord; {order of infinity found by |scan_dimen|}


 read_file:array[0..15] of alpha_file; {used for \.[\\read]}
 read_open:array[0..16] of normal..closed; {state of |read_file[n]|}


 cond_ptr:halfword ; {top of the condition stack}
 if_limit:normal..or_code; {upper bound on |fi_or_else| codes}
 cur_if:small_number; {type of conditional being worked on}
 if_line:integer; {line where that conditional began}


 skip_line:integer; {skipping began here}


 cur_name:str_number; {name of file just scanned}
 cur_area:str_number; {file area just scanned, or \.[""]}
 cur_ext:str_number; {file extension just scanned, or \.[""]}


 area_delimiter:pool_pointer; {the most recent `\./', if any}
 ext_delimiter:pool_pointer; {the most recent `\..', if any}


 format_default_length: integer;
 TEX_format_default: cstring;


 name_in_progress:boolean; {is a file name being scanned?}
 job_name:str_number; {principal file name}
 log_opened:boolean; {has the transcript file been opened?}


 dvi_file: byte_file; {the device-independent output goes here}
 output_file_name: str_number; {full name of the output file}
  texmf_log_name :str_number; {full name of the log file}


 tfm_file:byte_file;
buf:eight_bits;


 font_info: ^fmemory_word;
  {the big collection of font data}
 fmem_ptr:font_index; {first unused word of |font_info|}
 font_ptr:internal_font_number; {largest internal font number in use}
 font_check: ^four_quarters; {check sum}
 font_size: ^scaled; {``at'' size}
 font_dsize: ^scaled; {``design'' size}
 font_params: ^font_index; {how many font
  parameters are present}
 font_name: ^str_number; {name of the font}
 font_area: ^str_number; {area of the font}
 font_bc: ^eight_bits;
  {beginning (smallest) character code}
 font_ec: ^eight_bits;
  {ending (largest) character code}
 font_glue: ^halfword ;
  {glue specification for interword space, |null| if not allocated}
 font_used: ^boolean;
  {has a character from this font actually appeared in the output?}
 hyphen_char: ^integer;
  {current \.[\\hyphenchar] values}
 skew_char: ^integer;
  {current \.[\\skewchar] values}
 bchar_label: ^font_index;
  {start of |lig_kern| program for left boundary character,
  |non_address| if there is none}
 font_bchar: ^nine_bits;
  {boundary character, |non_char| if there is none}
 font_false_bchar: ^nine_bits;
  {|font_bchar| if it doesn't exist in the font, otherwise |non_char|}


 char_base: ^integer;
  {base addresses for |char_info|}
 width_base: ^integer;
  {base addresses for widths}
 height_base: ^integer;
  {base addresses for heights}
 depth_base: ^integer;
  {base addresses for depths}
 italic_base: ^integer;
  {base addresses for italic corrections}
 lig_kern_base: ^integer;
  {base addresses for ligature/kerning programs}
 kern_base: ^integer;
  {base addresses for kerns}
 exten_base: ^integer;
  {base addresses for extensible recipes}
 param_base: ^integer;
  {base addresses for font parameters}


 null_character:four_quarters; {nonexistent character information}


 total_pages:integer; {the number of pages that have been shipped out}
 max_v:scaled; {maximum height-plus-depth of pages shipped so far}
 max_h:scaled; {maximum width of pages shipped so far}
 max_push:integer; {deepest nesting of |push| commands encountered so far}
 last_bop:integer; {location of previous |bop| in the \.[DVI] output}
 dead_cycles:integer; {recent outputs that didn't ship anything out}
 doing_leaders:boolean; {are we inside a leader box?}


{character and font in current |char_node|}
 c:quarterword;
 f:internal_font_number;
 rule_ht, rule_dp, rule_wd:scaled; {size of current rule being output}
 g:halfword ; {current glue specification}
 lq, lr:integer; {quantities used in calculations for leaders}


 dvi_buf:^eight_bits; {buffer for \.[DVI] output}
 half_buf:integer; {half of |dvi_buf_size|}
 dvi_limit:integer; {end of the current half buffer}
 dvi_ptr:integer; {the next available buffer address}
 dvi_offset:integer; {|dvi_buf_size| times the number of times the
  output buffer has been fully emptied}
 dvi_gone:integer; {the number of bytes already output to |dvi_file|}


 down_ptr, right_ptr:halfword ; {heads of the down and right stacks}


 dvi_h, dvi_v:scaled; {a \.[DVI] reader program thinks we are here}
 cur_h, cur_v:scaled; {\TeX\ thinks we are here}
 dvi_f:internal_font_number; {the current font}
 cur_s:integer; {current depth of output box nesting, initially $-1$}


 total_stretch,  total_shrink: array[glue_ord] of scaled;
  {glue found by |hpack| or |vpack|}
 last_badness:integer; {badness of the most recently packaged box}


 adjust_tail:halfword ; {tail of adjustment list}


 pack_begin_line:integer; {source file line where the current paragraph
  or alignment began; a negative value denotes alignment}


 empty_field:two_halves;
 null_delimiter:four_quarters;


 cur_mlist:halfword ; {beginning of mlist to be translated}
 cur_style:small_number; {style code at current place in the list}
 cur_size:small_number; {size code corresponding to |cur_style|}
 cur_mu:scaled; {the math unit width corresponding to |cur_size|}
 mlist_penalties:boolean; {should |mlist_to_hlist| insert penalties?}


 cur_f:internal_font_number; {the |font| field of a |math_char|}
 cur_c:quarterword; {the |character| field of a |math_char|}
 cur_i:four_quarters; {the |char_info| of a |math_char|,
  or a lig/kern instruction}


 magic_offset:integer; {used to find inter-element spacing}


 cur_align:halfword ; {current position in preamble list}
 cur_span:halfword ; {start of currently spanned columns in preamble list}
 cur_loop:halfword ; {place to copy when extending a periodic preamble}
 align_ptr:halfword ; {most recently pushed-down alignment stack node}
 cur_head, cur_tail:halfword ; {adjustment list pointers}


 just_box:halfword ; {the |hlist_node| for the last line of the new paragraph}


 passive:halfword ; {most recent node on passive list}
 printed_node:halfword ; {most recent node that has been printed}
 pass_number:halfword; {the number of passive nodes allocated on this pass}


 active_width:array[1..6] of scaled;
  {distance from first active node to~|cur_p|}
 cur_active_width:array[1..6] of scaled; {distance from current active node}
 background:array[1..6] of scaled; {length of an ``empty'' line}
 break_width:array[1..6] of scaled; {length being computed after current break}


 no_shrink_error_yet:boolean; {have we complained about infinite shrinkage?}


 cur_p:halfword ; {the current breakpoint under consideration}
 second_pass:boolean; {is this our second attempt to break this paragraph?}
 final_pass:boolean; {is this our final attempt to break this paragraph?}
 threshold:integer; {maximum badness on feasible lines}


 minimal_demerits:array[very_loose_fit..tight_fit] of integer; {best total
  demerits known for current line class and position, given the fitness}
 minimum_demerits:integer; {best total demerits known for current line class
  and position}
 best_place:array[very_loose_fit..tight_fit] of halfword ; {how to achieve
  |minimal_demerits|}
 best_pl_line:array[very_loose_fit..tight_fit] of halfword; {corresponding
  line number}


 disc_width:scaled; {the length of discretionary material preceding a break}


 easy_line:halfword; {line numbers |>easy_line| are equivalent in break nodes}
 last_special_line:halfword; {line numbers |>last_special_line| all have
  the same width}
 first_width:scaled; {the width of all lines |<=last_special_line|, if
  no \.[\\parshape] has been specified}
 second_width:scaled; {the width of all lines |>last_special_line|}
 first_indent:scaled; {left margin to go with |first_width|}
 second_indent:scaled; {left margin to go with |second_width|}


 best_bet:halfword ; {use this passive node and its predecessors}
 fewest_demerits:integer; {the demerits associated with |best_bet|}
 best_line:halfword; {line number following the last line of the new paragraph}
 actual_looseness:integer; {the difference between |line_number(best_bet)|
  and the optimum |best_line|}
 line_diff:integer; {the difference between the current line number and
  the optimum |best_line|}


 hc:array[0..65] of 0..256; {word to be hyphenated}
 hn:0..64; {the number of positions occupied in |hc|;
                                  not always a |small_number|}
 ha, hb:halfword ; {nodes |ha..hb| should be replaced by the hyphenated result}
 hf:internal_font_number; {font number of the letters in |hc|}
 hu:array[0..63] of 0..256; {like |hc|, before conversion to lowercase}
 hyf_char:integer; {hyphen character of the relevant font}
 cur_lang, init_cur_lang:ASCII_code; {current hyphenation table of interest}
 l_hyf, r_hyf, init_l_hyf, init_r_hyf:integer; {limits on fragment sizes}
 hyf_bchar:halfword; {boundary character after $c_n$}


 hyf:array [0..64] of 0..9; {odd values indicate discretionary hyphens}
 init_list:halfword ; {list of punctuation characters preceding the word}
 init_lig:boolean; {does |init_list| represent a ligature?}
 init_lft:boolean; {if so, did the ligature involve a left boundary?}


 hyphen_passed:small_number; {first hyphen in a ligature, if any}


 cur_l, cur_r:halfword; {characters before and after the cursor}
 cur_q:halfword ; {where a ligature should be detached}
 lig_stack:halfword ; {unfinished business to the right of the cursor}
 ligature_present:boolean; {should a ligature node be made for |cur_l|?}
 lft_hit, rt_hit:boolean; {did we hit a ligature with a boundary character?}


{We will dynamically allocate these arrays.}
 trie_trl:^trie_pointer; {|trie_link|}
 trie_tro:^trie_pointer; {|trie_op|}
 trie_trc:^quarterword; {|trie_char|}
 hyf_distance:array[1..trie_op_size] of small_number; {position |k-j| of $n_j$}
 hyf_num:array[1..trie_op_size] of small_number; {value of $n_j$}
 hyf_next:array[1..trie_op_size] of trie_opcode; {continuation code}
 op_start:array[ASCII_code] of 0..trie_op_size; {offset for current language}


 hyph_word: ^str_number; {exception words}
 hyph_list: ^halfword ; {lists of hyphen positions}
 hyph_link: ^hyph_pointer; {link array for hyphen exceptions hash table}
 hyph_count:integer; {the number of words in the exception dictionary}
 hyph_next:integer; {next free slot in hyphen exceptions hash table}


 ifdef('INITEX')   trie_op_hash:array[neg_trie_op_size..trie_op_size] of 0..trie_op_size;
  {trie op codes for quadruples}
 trie_used:array[ASCII_code] of trie_opcode;
  {largest opcode used so far for this language}
 trie_op_lang:array[1..trie_op_size] of ASCII_code;
  {language part of a hashed quadruple}
 trie_op_val:array[1..trie_op_size] of trie_opcode;
  {opcode corresponding to a hashed quadruple}
 trie_op_ptr:0..trie_op_size; {number of stored ops so far}
endif('INITEX')  
 max_op_used:trie_opcode; {largest opcode used for any language}
 small_op:boolean; {flag used while dumping or undumping}


 ifdef('INITEX')   trie_c:^packed_ASCII_code;
  {characters to match}
{ \hskip10pt } trie_o:^trie_opcode;
  {operations to perform}
{ \hskip10pt } trie_l:^trie_pointer;
  {left subtrie links}
{ \hskip10pt } trie_r:^trie_pointer;
  {right subtrie links}
{ \hskip10pt } trie_ptr:trie_pointer; {the number of nodes in the trie}
{ \hskip10pt } trie_hash:^trie_pointer;
  {used to identify equivalent subtries}
endif('INITEX') 


 ifdef('INITEX')   trie_taken: ^boolean;
  {does a family start here?}
{ \hskip10pt } trie_min:array[ASCII_code] of trie_pointer;
  {the first possible slot for each character}
{ \hskip10pt } trie_max:trie_pointer; {largest location used in |trie|}
{ \hskip10pt } trie_not_ready:boolean; {is the trie still in linked form?}
endif('INITEX') 


 best_height_plus_depth:scaled; {height of the best box, without stretching or
  shrinking}


 page_tail:halfword ; {the final node on the current page}
 page_contents:empty..box_there; {what is on the current page so far?}
 page_max_depth:scaled; {maximum box depth on page being built}
 best_page_break:halfword ; {break here to get the best page known so far}
 least_page_cost:integer; {the score for this currently best page}
 best_size:scaled; {its |page_goal|}


 page_so_far:array [0..7] of scaled; {height and glue of the current page}
 last_glue:halfword ; {used to implement \.[\\lastskip]}
 last_penalty:integer; {used to implement \.[\\lastpenalty]}
 last_kern:scaled; {used to implement \.[\\lastkern]}
 insert_penalties:integer; {sum of the penalties for insertions
  that were held over}


 output_active:boolean; {are we in the midst of an output routine?}


 main_f:internal_font_number; {the current font}
 main_i:four_quarters; {character information bytes for |cur_l|}
 main_j:four_quarters; {ligature/kern command}
 main_k:font_index; {index into |font_info|}
 main_p:halfword ; {temporary register for list manipulation}
 main_s:integer; {space factor value}
 bchar:halfword; {boundary character of current font, or |non_char|}
 false_bchar:halfword; {nonexistent character matching |bchar|, or |non_char|}
 cancel_boundary:boolean; {should the left boundary be ignored?}
 ins_disc:boolean; {should we insert a discretionary node?}


 cur_box:halfword ; {box to be placed into its context}


 after_token:halfword; {zero, or a saved token}


 long_help_seen:boolean; {has the long \.[\\errmessage] help been used?}


 format_ident:str_number;


 fmt_file:word_file; {for input or output of format information}


 ready_already:integer; {a sacrifice of purity for economy}


 write_file:array[0..15] of alpha_file;
 write_open:array[0..17] of boolean;


 write_loc:halfword ; {|eqtb| address of \.[\\write]}


 edit_name_start: pool_pointer; {where the filename to switch to starts}
 edit_name_length, edit_line: integer; {what line to start editing at}
 ipc_on: cinttype; {level of IPC action, 0 for none [default]}
 stop_at_space: boolean; {whether |more_name| returns false for space}


 save_str_ptr: str_number;
 save_pool_ptr: pool_pointer;
 shellenabledp: cinttype;
 restrictedshell: cinttype;
 output_comment: ^char;
 k,l: 0..255; {used by `Make the first 256 strings', etc.}


 debug_format_file: boolean;



expand_depth_count:integer;


 mltex_p: boolean;


 mltex_enabled_p:boolean;  {enable character substitution}



 accent_c, base_c, replace_c:integer;
 ia_c, ib_c:four_quarters; {accent and base character information}
 base_slant, accent_slant:real; {amount of slant}
 base_x_height:scaled; {accent is designed for characters of this height}
 base_width, base_height:scaled; {height and width for base character}
 accent_width, accent_height:scaled; {height and width for accent}
 delta:scaled; {amount of right shift}






procedure initialize; {this procedure gets things started properly}
  var 
{ Local variables for initialization }
 i:integer;


 k:integer; {index into |mem|, |eqtb|, etc.}


 z:hyph_pointer; {runs through the exception dictionary}



  begin 
{ Initialize whatever \TeX\ might access }

{ Set initial values of key variables }
xchr[{040=}32]:=' ';
xchr[{041=}33]:='!';
xchr[{042=}34]:='"';
xchr[{043=}35]:='#';
xchr[{044=}36]:='$';
xchr[{045=}37]:='%';
xchr[{046=}38]:='&';
xchr[{047=}39]:='''';

xchr[{050=}40]:='(';
xchr[{051=}41]:=')';
xchr[{052=}42]:='*';
xchr[{053=}43]:='+';
xchr[{054=}44]:=',';
xchr[{055=}45]:='-';
xchr[{056=}46]:='.';
xchr[{057=}47]:='/';

xchr[{060=}48]:='0';
xchr[{061=}49]:='1';
xchr[{062=}50]:='2';
xchr[{063=}51]:='3';
xchr[{064=}52]:='4';
xchr[{065=}53]:='5';
xchr[{066=}54]:='6';
xchr[{067=}55]:='7';

xchr[{070=}56]:='8';
xchr[{071=}57]:='9';
xchr[{072=}58]:=':';
xchr[{073=}59]:=';';
xchr[{074=}60]:='<';
xchr[{075=}61]:='=';
xchr[{076=}62]:='>';
xchr[{077=}63]:='?';

xchr[{0100=}64]:='@';
xchr[{0101=}65]:='A';
xchr[{0102=}66]:='B';
xchr[{0103=}67]:='C';
xchr[{0104=}68]:='D';
xchr[{0105=}69]:='E';
xchr[{0106=}70]:='F';
xchr[{0107=}71]:='G';

xchr[{0110=}72]:='H';
xchr[{0111=}73]:='I';
xchr[{0112=}74]:='J';
xchr[{0113=}75]:='K';
xchr[{0114=}76]:='L';
xchr[{0115=}77]:='M';
xchr[{0116=}78]:='N';
xchr[{0117=}79]:='O';

xchr[{0120=}80]:='P';
xchr[{0121=}81]:='Q';
xchr[{0122=}82]:='R';
xchr[{0123=}83]:='S';
xchr[{0124=}84]:='T';
xchr[{0125=}85]:='U';
xchr[{0126=}86]:='V';
xchr[{0127=}87]:='W';

xchr[{0130=}88]:='X';
xchr[{0131=}89]:='Y';
xchr[{0132=}90]:='Z';
xchr[{0133=}91]:='[';
xchr[{0134=}92]:='\';
xchr[{0135=}93]:=']';
xchr[{0136=}94]:='^';
xchr[{0137=}95]:='_';

xchr[{0140=}96]:='`';
xchr[{0141=}97]:='a';
xchr[{0142=}98]:='b';
xchr[{0143=}99]:='c';
xchr[{0144=}100]:='d';
xchr[{0145=}101]:='e';
xchr[{0146=}102]:='f';
xchr[{0147=}103]:='g';

xchr[{0150=}104]:='h';
xchr[{0151=}105]:='i';
xchr[{0152=}106]:='j';
xchr[{0153=}107]:='k';
xchr[{0154=}108]:='l';
xchr[{0155=}109]:='m';
xchr[{0156=}110]:='n';
xchr[{0157=}111]:='o';

xchr[{0160=}112]:='p';
xchr[{0161=}113]:='q';
xchr[{0162=}114]:='r';
xchr[{0163=}115]:='s';
xchr[{0164=}116]:='t';
xchr[{0165=}117]:='u';
xchr[{0166=}118]:='v';
xchr[{0167=}119]:='w';

xchr[{0170=}120]:='x';
xchr[{0171=}121]:='y';
xchr[{0172=}122]:='z';
xchr[{0173=}123]:='{';
xchr[{0174=}124]:='|';
xchr[{0175=}125]:='}';
xchr[{0176=}126]:='~';



{Initialize |xchr| to the identity mapping.}
for i:=0 to {037=}31 do xchr[i]:=i;
for i:={0177=}127 to {0377=}255 do xchr[i]:=i;


for i:=first_text_char to last_text_char do xord[chr(i)]:=invalid_code;
for i:={0200=}128 to {0377=}255 do xord[xchr[i]]:=i;
for i:=0 to {0176=}126 do xord[xchr[i]]:=i;
{Set |xprn| for printable ASCII, unless |eight_bit_p| is set.}
for i:=0 to 255 do xprn[i]:=(eight_bit_p or ((i>={" "=}32)and(i<={"~"=}126)));

{The idea for this dynamic translation comes from the patch by
 Libor Skarvada \.[<libor@informatics.muni.cz>]
 and Petr Sojka \.[<sojka@informatics.muni.cz>]. I didn't use any of the
 actual code, though, preferring a more general approach.}

{This updates the |xchr|, |xord|, and |xprn| arrays from the provided
 |translate_filename|.  See the function definition in \.[texmfmp.c] for
 more comments.}
if translate_filename then read_tcx_file;

if interaction_option=unspecified_mode then
  interaction:=error_stop_mode
else
  interaction:=interaction_option;


deletions_allowed:=true; set_box_allowed:=true;
error_count:=0; {|history| is initialized elsewhere}


help_ptr:=0; use_err_help:=false;


interrupt:=0; OK_to_interrupt:=true;


 ifdef('TEXMF_DEBUG')  was_mem_end:=mem_min; {indicate that everything was previously free}
was_lo_max:=mem_min; was_hi_min:=mem_max;
panicking:=false;
endif('TEXMF_DEBUG') 


nest_ptr:=0; max_nest_stack:=0;
cur_list.mode_field :=vmode; cur_list.head_field :=mem_top-1 ; cur_list.tail_field :=mem_top-1 ;
cur_list.aux_field .int  :=-65536000 ; cur_list.ml_field :=0;
cur_list.pg_field :=0; shown_mode:=0;

{The following piece of code is a copy of module 991:}
page_contents:=empty; page_tail:=mem_top-2 ; {|link(page_head):=null;|}

last_glue:={0xfffffff=}268435455 ; last_penalty:=0; last_kern:=0;
page_so_far[7] :=0; page_max_depth:=0;


for k:=int_base to eqtb_size do xeq_level[k]:=level_one;


no_new_control_sequence:=true; {new identifiers are usually forbidden}



save_ptr:=0; cur_level:=level_one; cur_group:=bottom_level; cur_boundary:=0;
max_save_stack:=0;


mag_set:=0;


cur_mark[top_mark_code] :=-{0xfffffff=}268435455  ; cur_mark[first_mark_code] :=-{0xfffffff=}268435455  ; cur_mark[bot_mark_code] :=-{0xfffffff=}268435455  ;
cur_mark[split_first_mark_code] :=-{0xfffffff=}268435455  ; cur_mark[split_bot_mark_code] :=-{0xfffffff=}268435455  ;


cur_val:=0; cur_val_level:=int_val; radix:=0; cur_order:=normal;


for k:=0 to 16 do read_open[k]:=closed;


cond_ptr:=-{0xfffffff=}268435455  ; if_limit:=normal; cur_if:=0; if_line:=0;





null_character.b0:=min_quarterword; null_character.b1:=min_quarterword;
null_character.b2:=min_quarterword; null_character.b3:=min_quarterword;


total_pages:=0; max_v:=0; max_h:=0; max_push:=0; last_bop:=-1;
doing_leaders:=false; dead_cycles:=0; cur_s:=-1;


half_buf:=dvi_buf_size div 2; dvi_limit:=dvi_buf_size; dvi_ptr:=0;
dvi_offset:=0; dvi_gone:=0;


down_ptr:=-{0xfffffff=}268435455  ; right_ptr:=-{0xfffffff=}268435455  ;

adjust_tail:=-{0xfffffff=}268435455  ; last_badness:=0;


pack_begin_line:=0;


empty_field.rh:=empty; empty_field.lh:=-{0xfffffff=}268435455  ;

null_delimiter.b0:=0; null_delimiter.b1:=min_quarterword;

null_delimiter.b2:=0; null_delimiter.b3:=min_quarterword;


align_ptr:=-{0xfffffff=}268435455  ; cur_align:=-{0xfffffff=}268435455  ; cur_span:=-{0xfffffff=}268435455  ; cur_loop:=-{0xfffffff=}268435455  ;
cur_head:=-{0xfffffff=}268435455  ; cur_tail:=-{0xfffffff=}268435455  ;


for z:=0 to hyph_size do
  begin hyph_word[z]:=0; hyph_list[z]:=-{0xfffffff=}268435455  ; hyph_link[z]:=0;
  end;
hyph_count:=0;
hyph_next:=hyph_prime+1; if hyph_next>hyph_size then hyph_next:=hyph_prime;


output_active:=false; insert_penalties:=0;


ligature_present:=false; cancel_boundary:=false; lft_hit:=false; rt_hit:=false;
ins_disc:=false;


after_token:=0;

long_help_seen:=false;


format_ident:=0;


for k:=0 to 17 do write_open[k]:=false;


edit_name_start:=0;
stop_at_space:=true;
halting_on_error_p:=false;


expand_depth_count:=0;


mltex_enabled_p:=false;




 ifdef('INITEX')  if ini_version then begin  
{ Initialize table entries (done by \.[INITEX] only) }
for k:=mem_bot+1 to mem_bot +glue_spec_size +glue_spec_size +glue_spec_size +glue_spec_size +glue_spec_size-1  do mem[k].int :=0;
  {all glue dimensions are zeroed}
{ \xref[data structure assumptions] }
k:=mem_bot; while k<=mem_bot +glue_spec_size +glue_spec_size +glue_spec_size +glue_spec_size +glue_spec_size-1  do
    {set first words of glue specifications}
  begin   mem[  k].hh.rh  :=-{0xfffffff=}268435455  +1;
    mem[ k].hh.b0 :=normal;   mem[ k].hh.b1 :=normal;
  k:=k+glue_spec_size;
  end;
 mem[ mem_bot +glue_spec_size +2].int  := {0200000=}65536 ;   mem[ mem_bot +glue_spec_size ].hh.b0 :=fil;

 mem[ mem_bot +glue_spec_size +glue_spec_size +2].int  := {0200000=}65536 ;   mem[ mem_bot +glue_spec_size +glue_spec_size ].hh.b0 :=fill;

 mem[ mem_bot +glue_spec_size +glue_spec_size +glue_spec_size +2].int  := {0200000=}65536 ;   mem[ mem_bot +glue_spec_size +glue_spec_size +glue_spec_size ].hh.b0 :=fil;

 mem[ mem_bot +glue_spec_size +glue_spec_size +glue_spec_size +3].int  := {0200000=}65536 ;   mem[ mem_bot +glue_spec_size +glue_spec_size +glue_spec_size ].hh.b1 :=fil;

 mem[ mem_bot +glue_spec_size +glue_spec_size +glue_spec_size +glue_spec_size +2].int  :=- {0200000=}65536 ;   mem[ mem_bot +glue_spec_size +glue_spec_size +glue_spec_size +glue_spec_size ].hh.b0 :=fil;

rover:=mem_bot +glue_spec_size +glue_spec_size +glue_spec_size +glue_spec_size +glue_spec_size-1 +1;
 mem[ rover].hh.rh := {0xfffffff=}268435455  ; {now initialize the dynamic memory}
  mem[ rover].hh.lh :=1000; {which is a 1000-word available node}
  mem[  rover+ 1].hh.lh  :=rover;   mem[  rover+ 1].hh.rh  :=rover;

lo_mem_max:=rover+1000;  mem[ lo_mem_max].hh.rh :=-{0xfffffff=}268435455  ;  mem[ lo_mem_max].hh.lh :=-{0xfffffff=}268435455  ;

for k:=mem_top-13  to mem_top do
  mem[k]:=mem[lo_mem_max]; {clear list heads}

{ Initialize the special list heads and constant nodes }
 mem[ mem_top-10 ].hh.lh :={07777=}4095 +frozen_end_template ; {|link(omit_template)=null|}


 mem[ mem_top-9 ].hh.rh :=max_quarterword+1;  mem[ mem_top-9 ].hh.lh :=-{0xfffffff=}268435455  ;


 mem[ mem_top-7  ].hh.b0 :=hyphenated;   mem[  mem_top-7  + 1].hh.lh  :={0xfffffff=}268435455 ;
 mem[ mem_top-7  ].hh.b1 :=0; {the |subtype| is never examined by the algorithm}


 mem[ mem_top ].hh.b1 := 255 ;
 mem[ mem_top ].hh.b0 :=split_up;  mem[ mem_top ].hh.rh :=mem_top ;


 mem[ mem_top-2 ].hh.b0 :=glue_node;  mem[ mem_top-2 ].hh.b1 :=normal;

;
avail:=-{0xfffffff=}268435455  ; mem_end:=mem_top;
hi_mem_min:=mem_top-13 ; {initialize the one-word memory}
var_used:=mem_bot +glue_spec_size +glue_spec_size +glue_spec_size +glue_spec_size +glue_spec_size-1 +1-mem_bot; dyn_used:=hi_mem_stat_usage;
  {initialize statistics}


 eqtb[  undefined_control_sequence].hh.b0  :=undefined_cs;
 eqtb[  undefined_control_sequence].hh.rh  :=-{0xfffffff=}268435455  ;
 eqtb[  undefined_control_sequence].hh.b1  :=level_zero;
for k:=active_base to eqtb_top do
  eqtb[k]:=eqtb[undefined_control_sequence];


 eqtb[  glue_base].hh.rh  :=mem_bot ;  eqtb[  glue_base].hh.b1  :=level_one;
 eqtb[  glue_base].hh.b0  :=glue_ref;
for k:=glue_base+1 to local_base-1 do eqtb[k]:=eqtb[glue_base];
  mem[  mem_bot ].hh.rh  :=  mem[  mem_bot ].hh.rh  +local_base-glue_base;


 eqtb[  par_shape_loc].hh.rh   :=-{0xfffffff=}268435455  ;  eqtb[  par_shape_loc].hh.b0  :=shape_ref;
 eqtb[  par_shape_loc].hh.b1  :=level_one;

for k:=output_routine_loc to toks_base+255 do
  eqtb[k]:=eqtb[undefined_control_sequence];
 eqtb[  box_base+   0].hh.rh   :=-{0xfffffff=}268435455  ;  eqtb[  box_base].hh.b0  :=box_ref;  eqtb[  box_base].hh.b1  :=level_one;
for k:=box_base+1 to box_base+255 do eqtb[k]:=eqtb[box_base];
 eqtb[  cur_font_loc].hh.rh   :=font_base ;  eqtb[  cur_font_loc].hh.b0  :=data;
 eqtb[  cur_font_loc].hh.b1  :=level_one;

for k:=math_font_base to math_font_base+47 do eqtb[k]:=eqtb[cur_font_loc];
 eqtb[  cat_code_base].hh.rh  :=0;  eqtb[  cat_code_base].hh.b0  :=data;
 eqtb[  cat_code_base].hh.b1  :=level_one;

for k:=cat_code_base+1 to int_base-1 do eqtb[k]:=eqtb[cat_code_base];
for k:=0 to 255 do
  begin  eqtb[  cat_code_base+   k].hh.rh   :=other_char;  eqtb[  math_code_base+   k].hh.rh   := k ;  eqtb[  sf_code_base+   k].hh.rh   :=1000;
  end;
 eqtb[  cat_code_base+   carriage_return].hh.rh   :=car_ret;  eqtb[  cat_code_base+{" "=}   32].hh.rh   :=spacer;
 eqtb[  cat_code_base+{"\"=}   92].hh.rh   :=escape;  eqtb[  cat_code_base+{"%"=}   37].hh.rh   :=comment;
 eqtb[  cat_code_base+   invalid_code].hh.rh   :=invalid_char;  eqtb[  cat_code_base+   null_code].hh.rh   :=ignore;
for k:={"0"=}48 to {"9"=}57 do  eqtb[  math_code_base+   k].hh.rh   := k+ {070000=}28672  ;
for k:={"A"=}65 to {"Z"=}90 do
  begin  eqtb[  cat_code_base+   k].hh.rh   :=letter;  eqtb[  cat_code_base+   k+{"a"=}   97-{"A"=}   65].hh.rh   :=letter;

   eqtb[  math_code_base+   k].hh.rh   := k+ {070000=}28672 +{0x100=} 256 ;
   eqtb[  math_code_base+   k+{"a"=}   97-{"A"=}   65].hh.rh   := k+{"a"=} 97-{"A"=} 65+ {070000=}28672 +{0x100=} 256 ;

   eqtb[  lc_code_base+   k].hh.rh   :=k+{"a"=}97-{"A"=}65;  eqtb[  lc_code_base+   k+{"a"=}   97-{"A"=}   65].hh.rh   :=k+{"a"=}97-{"A"=}65;

   eqtb[  uc_code_base+   k].hh.rh   :=k;  eqtb[  uc_code_base+   k+{"a"=}   97-{"A"=}   65].hh.rh   :=k;

   eqtb[  sf_code_base+   k].hh.rh   :=999;
  end;


for k:=int_base to del_code_base-1 do eqtb[k].int:=0;
eqtb[int_base+ char_sub_def_min_code].int  :=256; eqtb[int_base+ char_sub_def_max_code].int  :=-1;
{allow \.[\\charsubdef] for char 0}

{|tracing_char_sub_def:=0| is already done}

eqtb[int_base+ mag_code].int  :=1000; eqtb[int_base+ tolerance_code].int  :=10000; eqtb[int_base+ hang_after_code].int  :=1; eqtb[int_base+ max_dead_cycles_code].int  :=25;
eqtb[int_base+ escape_char_code].int  :={"\"=}92; eqtb[int_base+ end_line_char_code].int  :=carriage_return;
for k:=0 to 255 do eqtb[del_code_base+ k].int :=-1;
eqtb[del_code_base+{"."=} 46].int :=0; {this null delimiter is used in error recovery}


for k:=dimen_base to eqtb_size do eqtb[k].int :=0;


hash_used:=frozen_control_sequence; {nothing is used}
hash_high:=0;
cs_count:=0;
 eqtb[  frozen_dont_expand].hh.b0  :=dont_expand;
 hash[ frozen_dont_expand].rh :={"notexpanded:"=}510;
{ \xref[notexpanded:] }





for k:=-trie_op_size to trie_op_size do trie_op_hash[k]:=0;
for k:=0 to 255 do trie_used[k]:=min_trie_op;
max_op_used:=min_trie_op;
trie_op_ptr:=0;


trie_not_ready:=true;


 hash[ frozen_protection].rh :={"inaccessible"=}1203;
{ \xref[inaccessible] }


if ini_version then format_ident:={" (INITEX)"=}1273;


 hash[ end_write].rh :={"endwrite"=}1315;  eqtb[  end_write].hh.b1  :=level_one;
 eqtb[  end_write].hh.b0  :=outer_call;  eqtb[  end_write].hh.rh  :=-{0xfffffff=}268435455  ;

  end; endif('INITEX')  

 
  end;

{ \4 }
{ Basic printing procedures }
procedure print_ln; {prints an end-of-line}
begin case selector of
term_and_log: begin writeln( stdout )  ; writeln( log_file)  ;
  term_offset:=0; file_offset:=0;
  end;
log_only: begin writeln( log_file)  ; file_offset:=0;
  end;
term_only: begin writeln( stdout )  ; term_offset:=0;
  end;
no_print,pseudo,new_string:  ;
 else  writeln( write_file[ selector]) 
 end ;

end; {|tally| is not affected}


procedure print_char( s:ASCII_code); {prints a single character}
label exit;
begin if 
{ Character |s| is the current new-line character }s=eqtb[int_base+ new_line_char_code].int  

 then
 if selector<pseudo then
  begin print_ln;  goto exit ;
  end;
case selector of
term_and_log: begin write(stdout , xchr[ s]) ; write(log_file, xchr[ s]) ;
  incr(term_offset); incr(file_offset);
  if term_offset=max_print_line then
    begin writeln( stdout )  ; term_offset:=0;
    end;
  if file_offset=max_print_line then
    begin writeln( log_file)  ; file_offset:=0;
    end;
  end;
log_only: begin write(log_file, xchr[ s]) ; incr(file_offset);
  if file_offset=max_print_line then print_ln;
  end;
term_only: begin write(stdout , xchr[ s]) ; incr(term_offset);
  if term_offset=max_print_line then print_ln;
  end;
no_print:  ;
pseudo: if tally<trick_count then trick_buf[tally mod error_line]:=s;
new_string: begin if pool_ptr<pool_size then  begin str_pool[pool_ptr]:=   s ; incr(pool_ptr); end ;
  end; {we drop characters if the string space is full}
 else  write(write_file[selector],xchr[s])
 end ;

incr(tally);
exit:end;


procedure print( s:integer); {prints string |s|}
label exit;
var j:pool_pointer; {current character code position}
 nl:integer; {new-line character to restore}
begin if s>=str_ptr then s:={"???"=}259 {this can't happen}
{ \xref[???] }
else if s<256 then
  if s<0 then s:={"???"=}259 {can't happen}
  else begin if selector>pseudo then
      begin print_char(s);  goto exit ; {internal strings are not expanded}
      end;
    if (
{ Character |s| is the current new-line character }s=eqtb[int_base+ new_line_char_code].int  

) then
      if selector<pseudo then
        begin print_ln;  goto exit ;
        end;
    nl:=eqtb[int_base+ new_line_char_code].int  ; eqtb[int_base+ new_line_char_code].int  :=-1;
      {temporarily disable new-line character}
    j:=str_start[s];
    while j<str_start[s+1] do
      begin print_char(  str_pool[ j] ); incr(j);
      end;
    eqtb[int_base+ new_line_char_code].int  :=nl;  goto exit ;
    end;
j:=str_start[s];
while j<str_start[s+1] do
  begin print_char(  str_pool[ j] ); incr(j);
  end;
exit:end;


procedure slow_print( s:integer); {prints string |s|}
var j:pool_pointer; {current character code position}
begin if (s>=str_ptr) or (s<256) then print(s)
else begin j:=str_start[s];
  while j<str_start[s+1] do
    begin print(  str_pool[ j] ); incr(j);
    end;
  end;
end;


procedure print_nl( s:str_number); {prints string |s| at beginning of line}
begin if ((term_offset>0)and(odd(selector)))or 
  ((file_offset>0)and(selector>=log_only)) then print_ln;
print(s);
end;


procedure print_esc( s:str_number); {prints escape character, then |s|}
var c:integer; {the escape character code}
begin  
{ Set variable |c| to the current escape character }c:=eqtb[int_base+ escape_char_code].int  

;
if c>=0 then if c<256 then print(c);
slow_print(s);
end;


procedure print_the_digs( k:eight_bits);
  {prints |dig[k-1]|$\,\ldots\,$|dig[0]|}
begin while k>0 do
  begin decr(k);
  if dig[k]<10 then print_char({"0"=}48+dig[k])
  else print_char({"A"=}65-10+dig[k]);
  end;
end;


procedure print_int( n:integer); {prints an integer in decimal form}
var k:0..23; {index to current digit; we assume that $\vert n\vert<10^[23]$}
 m:integer; {used to negate |n| in possibly dangerous cases}
begin k:=0;
if n<0 then
  begin print_char({"-"=}45);
  if n>-100000000 then   n:=- n 
  else  begin m:=-1-n; n:=m div 10; m:=(m mod 10)+1; k:=1;
    if m<10 then dig[0]:=m
    else  begin dig[0]:=0; incr(n);
      end;
    end;
  end;
repeat dig[k]:=n mod 10; n:=n div 10; incr(k);
until n=0;
print_the_digs(k);
end;


procedure print_cs( p:integer); {prints a purported control sequence}
begin if p<hash_base then {single character}
  if p>=single_base then
    if p=null_cs then
      begin print_esc({"csname"=}512); print_esc({"endcsname"=}513); print_char({" "=}32);
      end
    else  begin print_esc(p-single_base);
      if  eqtb[  cat_code_base+   p-   single_base].hh.rh   =letter then print_char({" "=}32);
      end
  else if p<active_base then print_esc({"IMPOSSIBLE."=}514)
{ \xref[IMPOSSIBLE] }
  else print(p-active_base)
else if ((p>=undefined_control_sequence)and(p<=eqtb_size))or(p>eqtb_top) then
  print_esc({"IMPOSSIBLE."=}514)
else if ( hash[ p].rh >=str_ptr) then print_esc({"NONEXISTENT."=}515)
{ \xref[NONEXISTENT] }
else  begin print_esc( hash[ p].rh );
  print_char({" "=}32);
  end;
end;


procedure sprint_cs( p:halfword ); {prints a control sequence}
begin if p<hash_base then
  if p<single_base then print(p-active_base)
  else  if p<null_cs then print_esc(p-single_base)
    else  begin print_esc({"csname"=}512); print_esc({"endcsname"=}513);
      end
else print_esc( hash[ p].rh );
end;


procedure print_file_name( n, a, e:integer);
var must_quote: boolean; {whether to quote the filename}
 j:pool_pointer; {index into |str_pool|}
begin
must_quote:=false;
 if  a<>0 then begin j:=str_start[ a]; while (not must_quote) and (j<str_start[ a+1]) do begin must_quote:=str_pool[j]={" "=}32; incr(j); end; end ;  if  n<>0 then begin j:=str_start[ n]; while (not must_quote) and (j<str_start[ n+1]) do begin must_quote:=str_pool[j]={" "=}32; incr(j); end; end ;  if  e<>0 then begin j:=str_start[ e]; while (not must_quote) and (j<str_start[ e+1]) do begin must_quote:=str_pool[j]={" "=}32; incr(j); end; end ;
{FIXME: Alternative is to assume that any filename that has to be quoted has
 at least one quoted component...if we pick this, a number of insertions
 of |print_file_name| should go away.
|must_quote|:=((|a|<>0)and(|str_pool|[|str_start|[|a|]]=""""))or
              ((|n|<>0)and(|str_pool|[|str_start|[|n|]]=""""))or
              ((|e|<>0)and(|str_pool|[|str_start|[|e|]]=""""));}
if must_quote then print_char({""""=}34);
 if  a<>0 then for j:=str_start[ a] to str_start[ a+1]-1 do if   str_pool[ j] <>{""""=}34 then print(  str_pool[ j] ) ;  if  n<>0 then for j:=str_start[ n] to str_start[ n+1]-1 do if   str_pool[ j] <>{""""=}34 then print(  str_pool[ j] ) ;  if  e<>0 then for j:=str_start[ e] to str_start[ e+1]-1 do if   str_pool[ j] <>{""""=}34 then print(  str_pool[ j] ) ;
if must_quote then print_char({""""=}34);
end;


procedure print_size( s:integer);
begin if s=text_size then print_esc({"textfont"=}417)
else if s=script_size then print_esc({"scriptfont"=}418)
else print_esc({"scriptscriptfont"=}419);
end;


procedure print_write_whatsit( s:str_number; p:halfword );
begin print_esc(s);
if   mem[  p+ 1].hh.lh  <16 then print_int(  mem[  p+ 1].hh.lh  )
else if   mem[  p+ 1].hh.lh  =16 then print_char({"*"=}42)
{ \xref[*\relax] }
else print_char({"-"=}45);
end;


procedure print_csnames (hstart:integer; hfinish:integer);
var c,h:integer;
begin
  writeln( stderr, 'fmtdebug:csnames from ',  hstart, ' to ',  hfinish, ':') ;
  for h := hstart to hfinish do begin
    if  hash[ h].rh  > 0 then begin {if have anything at this position}
      for c := str_start[ hash[ h].rh ] to str_start[ hash[ h].rh  + 1] - 1
      do begin
        put_byte(str_pool[c], stderr); {print the characters}
      end;
      writeln( stderr, '|') ;
    end;
  end;
end;


procedure print_file_line;
var level: 0..max_in_open;
begin
  level:=in_open;
  while (level>0) and (full_source_filename_stack[level]=0) do
    decr(level);
  if level=0 then
    print_nl({"! "=}262)
  else begin
    print_nl ({""=}335); print (full_source_filename_stack[level]); print ({":"=}58);
    if level=in_open then print_int (line)
    else print_int (line_stack[level+1]);
    print ({": "=}576);
  end;
end;



{ \4 }
{ Error handling procedures }
procedure normalize_selector; forward;{ \2 }

procedure get_token; forward;{ \2 }

procedure term_input; forward;{ \2 }

procedure show_context; forward;{ \2 }

procedure begin_file_reading; forward;{ \2 }

procedure open_log_file; forward;{ \2 }

procedure close_files_and_terminate; forward;{ \2 }

procedure clear_for_error_prompt; forward;{ \2 }

procedure give_err_help; forward;{ \2 }

{ \4\hskip-\fontdimen2\font }   ifdef('TEXMF_DEBUG')  procedure debug_help;
  forward;  endif('TEXMF_DEBUG') 


noreturn procedure jump_out;
begin
close_files_and_terminate;
begin  fflush (stdout ) ; ready_already:=0; if (history <> spotless) and (history <> warning_issued) then uexit(1) else uexit(0); end ;
panic(end_of_TEX);
end;


procedure error; {completes the job of error reporting}
label continue,exit;
var c:ASCII_code; {what the user types}
 s1, s2, s3, s4:integer;
  {used to save global variables when deleting tokens}
begin if history<error_message_issued then history:=error_message_issued;
print_char({"."=}46); show_context;
if (halt_on_error_p) then begin
  {If |close_files_and_terminate| generates an error, we'll end up back
   here; just give up in that case. If files are truncated, too bad.}
  if (halting_on_error_p) then begin  fflush (stdout ) ; ready_already:=0; if (history <> spotless) and (history <> warning_issued) then uexit(1) else uexit(0); end ; {quit immediately}
  halting_on_error_p:=true;

  {This module is executed at the end of the |error| procedure in
   \.[tex.web], but we'll never get there when |halt_on_error_p|, so the
   error help shouldn't get duplicated. It's potentially useful to see,
   especially if \.[\\errhelp] is being used. See thread at:
   \.[https://tug.org/pipermail/tex-live/2024-July/050741.html].}
  
{ Put help message on the transcript file }
if interaction>batch_mode then decr(selector); {avoid terminal output}
if use_err_help then
  begin print_ln; give_err_help;
  end
else while help_ptr>0 do
  begin decr(help_ptr); print_nl(help_line[help_ptr]);
  end;
print_ln;
if interaction>batch_mode then incr(selector); {re-enable terminal output}
print_ln

;

  {Proceed with normal exit.}
  history:=fatal_error_stop;
  jump_out;
end;
if interaction=error_stop_mode then
  
{ Get user's advice and |return| }
 while true do  begin continue: if interaction<>error_stop_mode then  goto exit ;
  clear_for_error_prompt; begin    ; print({"? "=} 264); term_input; end ;
{ \xref[?\relax] }
  if last=first then  goto exit ;
  c:=buffer[first];
  if c>={"a"=}97 then c:=c+{"A"=}65-{"a"=}97; {convert to uppercase}
  
{ Interpret code |c| and |return| if done }
case c of
{"0"=}48,{"1"=}49,{"2"=}50,{"3"=}51,{"4"=}52,{"5"=}53,{"6"=}54,{"7"=}55,{"8"=}56,{"9"=}57: if deletions_allowed then
  
{ Delete \(c)|c-"0"| tokens and |goto continue| }
begin s1:=cur_tok; s2:=cur_cmd; s3:=cur_chr; s4:=align_state;
align_state:=1000000; OK_to_interrupt:=false;
if (last>first+1) and (buffer[first+1]>={"0"=}48)and(buffer[first+1]<={"9"=}57) then
  c:=c*10+buffer[first+1]-{"0"=}48*11
else c:=c-{"0"=}48;
while c>0 do
  begin get_token; {one-level recursive call of |error| is possible}
  decr(c);
  end;
cur_tok:=s1; cur_cmd:=s2; cur_chr:=s3; align_state:=s4; OK_to_interrupt:=true;
 begin help_ptr:=2; help_line[1]:={"I have just deleted some text, as you asked."=} 277; help_line[0]:={"You can now delete more, or insert, or whatever."=} 278; end ;
show_context; goto continue;
end

;
{ \4\4 }   ifdef('TEXMF_DEBUG')  {"D"=}68: begin debug_help; goto continue; end; endif('TEXMF_DEBUG') 

{"E"=}69: if base_ptr>0 then if input_stack[base_ptr].name_field>=256 then
    begin edit_name_start:=str_start[input_stack[base_ptr] .name_field];
    edit_name_length:=str_start[input_stack[base_ptr] .name_field+1] -
                      str_start[input_stack[base_ptr] .name_field];
    edit_line:=line;
    jump_out;
  end;
{"H"=}72: 
{ Print the help information and |goto continue| }
begin if use_err_help then
  begin give_err_help; use_err_help:=false;
  end
else  begin if help_ptr=0 then
     begin help_ptr:=2; help_line[1]:={"Sorry, I don't know how to help in this situation."=} 279; help_line[0]:={"Maybe you should try asking a human?"=} 280; end ;
  repeat decr(help_ptr); print(help_line[help_ptr]); print_ln;
  until help_ptr=0;
  end;
 begin help_ptr:=4; help_line[3]:={"Sorry, I already gave what help I could..."=} 281; help_line[2]:={"Maybe you should try asking a human?"=} 280; help_line[1]:={"An error might have occurred before I noticed any problems."=} 282; help_line[0]:={"``If all else fails, read the instructions.''"=} 283; end ;

goto continue;
end

;
{"I"=}73:
{ Introduce new material from the terminal and |return| }
begin begin_file_reading; {enter a new syntactic level for terminal input}
{now |state=mid_line|, so an initial blank space will count as a blank}
if last>first+1 then
  begin cur_input.loc_field :=first+1; buffer[first]:={" "=}32;
  end
else  begin begin    ; print({"insert>"=} 276); term_input; end ; cur_input.loc_field :=first;
{ \xref[insert>] }
  end;
first:=last;
cur_input.limit_field:=last-1; {no |end_line_char| ends this line}
 goto exit ;
end

;
{"Q"=}81,{"R"=}82,{"S"=}83:
{ Change the interaction level and |return| }
begin error_count:=0; interaction:=batch_mode+c-{"Q"=}81;
print({"OK, entering "=}271);
case c of
{"Q"=}81:begin print_esc({"batchmode"=}272); decr(selector);
  end;
{"R"=}82:print_esc({"nonstopmode"=}273);
{"S"=}83:print_esc({"scrollmode"=}274);
end; {there are no other cases}
print({"..."=}275); print_ln;  fflush (stdout ) ;  goto exit ;
end

;
{"X"=}88:begin interaction:=scroll_mode; jump_out;
  end;
 else   
 end ;


{ Print the menu of available options }
begin print({"Type <return> to proceed, S to scroll future error messages,"=}265);

{ \xref[Type <return> to proceed...] }
print_nl({"R to run without stopping, Q to run quietly,"=}266);

print_nl({"I to insert something, "=}267);
if base_ptr>0 then if input_stack[base_ptr].name_field>=256 then
  print({"E to edit your file,"=}268);
if deletions_allowed then
  print_nl({"1 or ... or 9 to ignore the next 1 to 9 tokens of input,"=}269);
print_nl({"H for help, X to quit."=}270);
end



;
  end

;
incr(error_count);
if error_count=100 then
  begin print_nl({"(That makes 100 errors; please try again.)"=}263);
{ \xref[That makes 100 errors...] }
  history:=fatal_error_stop; jump_out;
  end;

{ Put help message on the transcript file }
if interaction>batch_mode then decr(selector); {avoid terminal output}
if use_err_help then
  begin print_ln; give_err_help;
  end
else while help_ptr>0 do
  begin decr(help_ptr); print_nl(help_line[help_ptr]);
  end;
print_ln;
if interaction>batch_mode then incr(selector); {re-enable terminal output}
print_ln

;
exit:end;


noreturn procedure fatal_error( s:str_number); {prints |s|, and that's it}
begin normalize_selector;

begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Emergency stop"=} 285); end ;  begin help_ptr:=1; help_line[0]:= s; end ; begin if interaction=error_stop_mode then interaction:=scroll_mode; if log_opened then error; ifdef('TEXMF_DEBUG')  if interaction>batch_mode then debug_help; endif('TEXMF_DEBUG')  history:=fatal_error_stop; jump_out; end ;
{ \xref[Emergency stop] }
end;


noreturn procedure overflow( s:str_number; n:integer); {stop due to finiteness}
begin normalize_selector;
begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"TeX capacity exceeded, sorry ["=} 286); end ;
{ \xref[TeX capacity exceeded ...] }
print(s); print_char({"="=}61); print_int(n); print_char({"]"=}93);
 begin help_ptr:=2; help_line[1]:={"If you really absolutely need more capacity,"=} 287; help_line[0]:={"you can ask a wizard to enlarge me."=} 288; end ;
begin if interaction=error_stop_mode then interaction:=scroll_mode; if log_opened then error; ifdef('TEXMF_DEBUG')  if interaction>batch_mode then debug_help; endif('TEXMF_DEBUG')  history:=fatal_error_stop; jump_out; end ;
end;


noreturn procedure confusion( s:str_number);
  {consistency check violated; |s| tells where}
begin normalize_selector;
if history<error_message_issued then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"This can't happen ("=} 289); end ; print(s); print_char({")"=}41);
{ \xref[This can't happen] }
   begin help_ptr:=1; help_line[0]:={"I'm broken. Please show this to someone who can fix can fix"=} 290; end ;
  end
else  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"I can't go on meeting you like this"=} 291); end ;
{ \xref[I can't go on...] }
   begin help_ptr:=2; help_line[1]:={"One of your faux pas seems to have wounded me deeply..."=} 292; help_line[0]:={"in fact, I'm barely conscious. Please fix it and try again."=} 293; end ;
  end;
begin if interaction=error_stop_mode then interaction:=scroll_mode; if log_opened then error; ifdef('TEXMF_DEBUG')  if interaction>batch_mode then debug_help; endif('TEXMF_DEBUG')  history:=fatal_error_stop; jump_out; end ;
end;






{ 5. }

{tangle:pos tex.web:268:3: }

{ The overall \TeX\ program begins with the heading just shown, after which
comes a bunch of procedure declarations and function declarations.
Finally we will get to the main program, which begins with the
comment `|start_here|'. If you want to skip down to the
main program now, you can look up `|start_here|' in the index.
But the author suggests that the best way to understand this program
is to follow pretty much the order of \TeX's components as they appear in the
\.[WEB] description you are now reading, since the present ordering is
intended to combine the advantages of the ``bottom up'' and ``top down''
approaches to the problem of understanding a somewhat complicated system. }

{ 6. }

{tangle:pos tex.ch:61:3: }

{ For Web2c, labels are not declared in the main program, but
we still have to declare the symbolic names. }

{ 7. }

{tangle:pos tex.web:290:3: }

{ Some of the code below is intended to be used only when diagnosing the
strange behavior that sometimes occurs when \TeX\ is being installed or
when system wizards are fooling around with \TeX\ without quite knowing
what they are doing. Such code will not normally be compiled; it is
delimited by the codewords `$|debug|\ldots|gubed|$', with apologies
to people who wish to preserve the purity of English.

Similarly, there is some conditional code delimited by
`$|stat|\ldots|tats|$' that is intended for use when statistics are to be
kept about \TeX's memory usage.  The |stat| $\ldots$ |tats| code also
implements diagnostic information for \.[\\tracingparagraphs],
\.[\\tracingpages], and \.[\\tracingrestores].
\xref[debugging] }

{ 10. }

{tangle:pos tex.web:348:3: }

{ This \TeX\ implementation conforms to the rules of the [\sl Pascal User
\xref[PASCAL][\PASCAL]
\xref[system dependencies]
Manual] published by Jensen and Wirth in 1975, except where system-dependent
\xref[Wirth, Niklaus]
\xref[Jensen, Kathleen]
code is necessary to make a useful system program, and except in another
respect where such conformity would unnecessarily obscure the meaning
and clutter up the code: We assume that |case| statements may include a
default case that applies if no matching label is found. Thus, we shall use
constructions like
$$\vbox[\halign[\ignorespaces#\hfil\cr
|case x of|\cr
1: $\langle\,$code for $x=1\,\rangle$;\cr
3: $\langle\,$code for $x=3\,\rangle$;\cr
|othercases| $\langle\,$code for |x<>1| and |x<>3|$\,\rangle$\cr
|endcases|\cr]]$$
since most \PASCAL\ compilers have plugged this hole in the language by
incorporating some sort of default mechanism. For example, the \ph\
compiler allows `|others|:' as a default label, and other \PASCAL s allow
syntaxes like `\&[else]' or `\&[otherwise]' or `\\[otherwise]:', etc. The
definitions of |othercases| and |endcases| should be changed to agree with
local conventions.  Note that no semicolon appears before |endcases| in
this program, so the definition of |endcases| should include a semicolon
if the compiler wants one. (Of course, if no default mechanism is
available, the |case| statements of \TeX\ will have to be laboriously
extended by listing all remaining cases. People who are stuck with such
\PASCAL s have, in fact, done this, successfully but not happily!)
\xref[PASCAL H][\ph] }

{ 12. }

{tangle:pos tex.web:430:3: }

{ Like the preceding parameters, the following quantities can be changed
at compile time to extend or reduce \TeX's capacity. But if they are changed,
it is necessary to rerun the initialization program \.[INITEX]
\xref[INITEX]
to generate new tables for the production \TeX\ program.
One can't simply make helter-skelter changes to the following constants,
since certain rather complex initialization
numbers are computed from them. They are defined here using
\.[WEB] macros, instead of being put into \PASCAL's |const| list, in order to
emphasize this distinction. }

{ 15. }

{tangle:pos tex.web:476:3: }

{ Labels are given symbolic names by the following definitions, so that
occasional |goto| statements will be meaningful. We insert the label
`|exit|' just before the `\ignorespaces|end|\unskip' of a procedure in
which we have used the `|return|' statement defined below; the label
`|restart|' is occasionally used at the very beginning of a procedure; and
the label `|reswitch|' is occasionally used just prior to a |case|
statement in which some cases change the conditions and we wish to branch
to the newly applicable case.  Loops that are set up with the |loop|
construction defined below are commonly exited by going to `|done|' or to
`|found|' or to `|not_found|', and they are sometimes repeated by going to
`|continue|'.  If two or more parts of a subroutine start differently but
end up the same, the shared code may be gathered together at
`|common_ending|'.

Incidentally, this program never declares a label that isn't actually used,
because some fussy \PASCAL\ compilers will complain about redundant labels. }

{ 16. }

{tangle:pos tex.web:510:3: }

{ Here are some macros for common programming idioms. }

{ 17. \[2] The character set }

{tangle:pos tex.web:523:27: }

{ In order to make \TeX\ readily portable to a wide variety of
computers, all of its input text is converted to an internal eight-bit
code that includes standard ASCII, the ``American Standard Code for
Information Interchange.''  This conversion is done immediately when each
character is read in. Conversely, characters are converted from ASCII to
the user's external representation just before they are output to a
text file.

Such an internal code is relevant to users of \TeX\ primarily because it
governs the positions of characters in the fonts. For example, the
character `\.A' has ASCII code $65=@'101$, and when \TeX\ typesets
this letter it specifies character number 65 in the current font.
If that font actually has `\.A' in a different position, \TeX\ doesn't
know what the real position is; the program that does the actual printing from
\TeX's device-independent files is responsible for converting from ASCII to
a particular font encoding.
\xref[ASCII code]

\TeX's internal code also defines the value of constants
that begin with a reverse apostrophe; and it provides an index to the
\.[\\catcode], \.[\\mathcode], \.[\\uccode], \.[\\lccode], and \.[\\delcode]
tables. }

{ 22. }

{tangle:pos tex.web:702:3: }

{ Some of the ASCII codes without visible characters have been given symbolic
names in this program because they are used with a special meaning. }

{ 27. }

{tangle:pos tex.ch:382:3: }

{ All of the file opening functions are defined in C. }

{ 28. }

{tangle:pos tex.ch:409:3: }

{ And all the file closing routines as well. }

{ 29. }

{tangle:pos tex.web:887:3: }

{ Binary input and output are done with \PASCAL's ordinary |get| and |put|
procedures, so we don't have to make any other special arrangements for
binary~I/O. Text output is also easy to do with standard \PASCAL\ routines.
The treatment of text input is more difficult, however, because
of the necessary translation to |ASCII_code| values.
\TeX's conventions should be efficient, and they should
blend nicely with the user's operating environment. }

{ 31. }

{tangle:pos tex.web:908:3: }

{ The |input_ln| function brings the next line of input from the specified
file into available positions of the buffer array and returns the value
|true|, unless the file has already been entirely read, in which case it
returns |false| and sets |last:=first|.  In general, the |ASCII_code|
numbers that represent the next line of the file are input into
|buffer[first]|, |buffer[first+1]|, \dots, |buffer[last-1]|; and the
global variable |last| is set equal to |first| plus the length of the
line. Trailing blanks are removed from the line; thus, either |last=first|
(in which case the line was entirely blank) or |buffer[last-1]<>" "|.

An overflow error is given, however, if the normal actions of |input_ln|
would make |last>=buf_size|; this is done so that other parts of \TeX\
can safely look at the contents of |buffer[last+1]| without overstepping
the bounds of the |buffer| array. Upon entry to |input_ln|, the condition
|first<buf_size| will always hold, so that there is always room for an
``empty'' line.

The variable |max_buf_stack|, which is used to keep track of how large
the |buf_size| parameter must be to accommodate the present job, is
also kept up to date by |input_ln|.

If the |bypass_eoln| parameter is |true|, |input_ln| will do a |get|
before looking at the first character of the line; this skips over
an |eoln| that was in |f^|. The procedure does not do a |get| when it
reaches the end of the line; therefore it can be used to acquire input
from the user's terminal as well as from ordinary text files.

Standard \PASCAL\ says that a file should have |eoln| immediately
before |eof|, but \TeX\ needs only a weaker restriction: If |eof|
occurs in the middle of a line, the system function |eoln| should return
a |true| result (even though |f^| will be undefined).

Since the inner loop of |input_ln| is part of \TeX's ``inner loop''---each
character of input comes in at this place---it is wise to reduce system
overhead by making use of special routines that read in an entire array
of characters at once, if such routines are available. The following
code uses standard \PASCAL\ to illustrate what needs to be done, but
finer tuning is often possible at well-developed \PASCAL\ sites.
\xref[inner loop]

We define |input_ln| in C, for efficiency. Nevertheless we quote the module
`Report overflow of the input buffer, and abort' here in order to make
\.[WEAVE] happy, since part of that module is needed by e-TeX. } { 
[ Report overflow of the input buffer, and abort ]
  begin cur_input.loc_field:=first; cur_input.limit_field:=last-1;
  overflow(["buffer size"=]256,buf_size);
[ \xref[TeX capacity exceeded buffer size][\quad buffer size] ]
  end

 }



{ 33. }

{tangle:pos tex.ch:543:3: }

{ Here is how to open the terminal files.  |t_open_out| does nothing.
|t_open_in|, on the other hand, does the work of ``rescanning,'' or getting
any command line arguments the user has provided.  It's defined in C. }

{ 34. }

{tangle:pos tex.web:987:3: }

{ Sometimes it is necessary to synchronize the input/output mixture that
happens on the user's terminal, and three system-dependent
procedures are used for this
purpose. The first of these, |update_terminal|, is called when we want
to make sure that everything we have output to the terminal so far has
actually left the computer's internal buffers and been sent.
The second, |clear_terminal|, is called when we wish to cancel any
input that the user may have typed ahead (since we are about to
issue an unexpected error message). The third, |wake_up_terminal|,
is supposed to revive the terminal if the user has disabled it by
some instruction to the operating system.  The following macros show how
these operations can be specified with [\mc UNIX].  |update_terminal|
does an |fflush|. |clear_terminal| is redefined
to do nothing, since the user should control the terminal.
\xref[system dependencies] }

{ 36. }

{tangle:pos tex.web:1044:3: }

{ Different systems have different ways to get started. But regardless of
what conventions are adopted, the routine that initializes the terminal
should satisfy the following specifications:

\yskip\textindent[1)]It should open file |term_in| for input from the
  terminal. (The file |term_out| will already be open for output to the
  terminal.)

\textindent[2)]If the user has given a command line, this line should be
  considered the first line of terminal input. Otherwise the
  user should be prompted with `\.[**]', and the first line of input
  should be whatever is typed in response.

\textindent[3)]The first line of input, which might or might not be a
  command line, should appear in locations |first| to |last-1| of the
  |buffer| array.

\textindent[4)]The global variable |loc| should be set so that the
  character to be read next by \TeX\ is in |buffer[loc]|. This
  character should not be blank, and we should have |loc<last|.

\yskip\noindent(It may be necessary to prompt the user several times
before a non-blank line comes in. The prompt is `\.[**]' instead of the
later `\.*' because the meaning is slightly different: `\.[\\input]' need
not be typed immediately after~`\.[**]'.) }

{ 37. }

{tangle:pos tex.ch:592:3: }

{ The following program does the required initialization.
Iff anything has been specified on the command line, then |t_open_in|
will return with |last > first|.
\xref[system dependencies] } function init_terminal:boolean; {gets the terminal input started}
label exit;
begin t_open_in;
if last > first then
  begin cur_input.loc_field  := first;
  while (cur_input.loc_field  < last) and (buffer[cur_input.loc_field ]=' ') do incr(cur_input.loc_field );
  if cur_input.loc_field  < last then
    begin init_terminal := true; goto exit;
    end;
  end;
 while true do  begin    ; write(stdout ,'**');  fflush (stdout ) ;
{ \xref[**] }
  if not input_ln(stdin ,true) then {this shouldn't happen}
    begin writeln( stdout ) ;
    writeln( stdout ,'! End of file on the terminal... why?') ;
{ \xref[End of file on the terminal] }
    init_terminal:=false;  goto exit ;
    end;
  cur_input.loc_field :=first;
  while (cur_input.loc_field <last)and(buffer[cur_input.loc_field ]={" "=}32) do incr(cur_input.loc_field );
  if cur_input.loc_field <last then
    begin init_terminal:=true;
     goto exit ; {return unless the line was all blank}
    end;
  writeln( stdout ,'Please type the name of your input file.') ;
  end;
exit:end;



{ 40. }

{tangle:pos tex.web:1156:3: }

{ Several of the elementary string operations are performed using \.[WEB]
macros instead of \PASCAL\ procedures, because many of the
operations are done quite frequently and we want to avoid the
overhead of procedure calls. For example, here is
a simple macro that computes the length of a string.
\xref[WEB] }

{ 41. }

{tangle:pos tex.web:1166:3: }

{ The length of the current string is called |cur_length|: }

{ 42. }

{tangle:pos tex.web:1170:3: }

{ Strings are created by appending character codes to |str_pool|.
The |append_char| macro, defined here, does not check to see if the
value of |pool_ptr| has gotten too high; this test is supposed to be
made before |append_char| is used. There is also a |flush_char|
macro, which erases the last character appended.

To test if there is room to append |l| more characters to |str_pool|,
we shall write |str_room(l)|, which aborts \TeX\ and gives an
apologetic error message if there isn't enough room. }

{ 43. }

{tangle:pos tex.web:1190:3: }

{ Once a sequence of characters has been appended to |str_pool|, it
officially becomes a string when the function |make_string| is called.
This function returns the identification number of the new string as its
value. } function make_string : str_number; {current string enters the pool}
begin if str_ptr=max_strings then
  overflow({"number of strings"=}258,max_strings-init_str_ptr);
{ \xref[TeX capacity exceeded number of strings][\quad number of strings] }
incr(str_ptr); str_start[str_ptr]:=pool_ptr;
make_string:=str_ptr-1;
end;



{ 44. }

{tangle:pos tex.web:1203:3: }

{ To destroy the most recently made string, we say |flush_string|. }

{ 45. }

{tangle:pos tex.web:1208:3: }

{ The following subroutine compares string |s| with another string of the
same length that appears in |buffer| starting at position |k|;
the result is |true| if and only if the strings are equal.
Empirical tests indicate that |str_eq_buf| is used in such a way that
it tends to return |true| about 80 percent of the time. } function str_eq_buf( s:str_number; k:integer):boolean;
  {test equality of strings}
label not_found; {loop exit}
var j: pool_pointer; {running index}
 result: boolean; {result of comparison}
begin j:=str_start[s];
while j<str_start[s+1] do
  begin if   str_pool[ j] <>buffer[k] then
    begin result:=false; goto not_found;
    end;
  incr(j); incr(k);
  end;
result:=true;
not_found: str_eq_buf:=result;
end;



{ 46. }

{tangle:pos tex.web:1230:3: }

{ Here is a similar routine, but it compares two strings in the string pool,
and it does not assume that they have the same length. } function str_eq_str( s, t:str_number):boolean;
  {test equality of strings}
label not_found; {loop exit}
var j, k: pool_pointer; {running indices}
 result: boolean; {result of comparison}
begin result:=false;
if (str_start[ s+1]-str_start[ s]) <>(str_start[ t+1]-str_start[ t])  then goto not_found;
j:=str_start[s]; k:=str_start[t];
while j<str_start[s+1] do
  begin if str_pool[j]<>str_pool[k] then goto not_found;
  incr(j); incr(k);
  end;
result:=true;
not_found: str_eq_str:=result;
end;



{ 47. }

{tangle:pos tex.web:1249:3: }

{ The initial values of |str_pool|, |str_start|, |pool_ptr|,
and |str_ptr| are computed by the \.[INITEX] program, based in part
on the information that \.[WEB] has output while processing \TeX.
\xref[INITEX]
\xref[string pool] } { \4 }
{ Declare additional routines for string recycling }
function search_string( search:str_number):str_number;
label found;
var result: str_number;
 s: str_number; {running index}
 len: integer; {length of searched string}
begin result:=0; len:=(str_start[ search+1]-str_start[ search]) ;
if len=0 then  {trivial case}
  begin result:={""=}335; goto found;
  end
else  begin s:=search-1;  {start search with newest string below |s|; |search>1|!}
  while s>255 do  {first 256 strings depend on implementation!!}
    begin if (str_start[ s+1]-str_start[ s]) =len then
      if str_eq_str(s,search) then
        begin result:=s; goto found;
        end;
    decr(s);
    end;
  end;
found:search_string:=result;
end;


function slow_make_string : str_number;
label exit;
var s: str_number; {result of |search_string|}
 t: str_number; {new string}
begin t:=make_string; s:=search_string(t);
if s>0 then
  begin begin decr(str_ptr); pool_ptr:=str_start[str_ptr]; end ; slow_make_string:=s;  goto exit ;
  end;
slow_make_string:=t;
exit:end;





 ifdef('INITEX')  function get_strings_started:boolean; {initializes the string pool,
  but returns |false| if something goes wrong}
label done,exit;
var k, l:0..255; {small indices or counters}
 m, n: ASCII_code ; {characters input from |pool_file|}
 g:str_number; {garbage}
 a:integer; {accumulator for check sum}
 c:boolean; {check sum has been checked}
begin
if g=0 then;
pool_ptr:=0; str_ptr:=0; str_start[0]:=0;

{ Make the first 256 strings }
for k:=0 to 255 do
  begin if (
{ Character |k| cannot be printed }
  (k<{" "=}32)or(k>{"~"=}126)

) then
    begin  begin str_pool[pool_ptr]:= {"^"=}  94 ; incr(pool_ptr); end ;  begin str_pool[pool_ptr]:= {"^"=}  94 ; incr(pool_ptr); end ;
    if k<{0100=}64 then  begin str_pool[pool_ptr]:=   k+{0100=}  64 ; incr(pool_ptr); end 
    else if k<{0200=}128 then  begin str_pool[pool_ptr]:=   k-{0100=}  64 ; incr(pool_ptr); end 
    else begin l:= k  div  16; if l<10 then  begin str_pool[pool_ptr]:=   l+{"0"=}  48 ; incr(pool_ptr); end  else  begin str_pool[pool_ptr]:=   l-  10+{"a"=}  97 ; incr(pool_ptr); end  ; l:= k  mod  16; if l<10 then  begin str_pool[pool_ptr]:=   l+{"0"=}  48 ; incr(pool_ptr); end  else  begin str_pool[pool_ptr]:=   l-  10+{"a"=}  97 ; incr(pool_ptr); end  ;
      end;
    end
  else  begin str_pool[pool_ptr]:=   k ; incr(pool_ptr); end ;
  g:=make_string;
  end

;

{ Read the other strings from the \.[TEX.POOL] file and return |true|, or give an error message and return |false| }
name_length := strlen (pool_name);
name_of_file := xmalloc_array (ASCII_code, name_length + 1);
strcpy (stringcast(name_of_file+1), pool_name); {copy the string}
if a_open_in (pool_file, kpse_texpool_format) then
  begin c:=false;
  repeat 
{ Read one string, but return |false| if the string memory space is getting too tight for comfort }
begin if eof(pool_file) then begin    ; writeln( stdout ,'! ',   pool_name, ' has no check sum.') ; a_close(pool_file); get_strings_started:=false;  goto exit ; end ;
{ \xref[TEX.POOL has no check sum] }
read(pool_file,m); read(pool_file,n); {read two digits of string length}
if m='*' then 
{ Check the pool check sum }
begin a:=0; k:=1;
 while true do    begin if (xord[n]<{"0"=}48)or(xord[n]>{"9"=}57) then
  begin    ; writeln( stdout ,'! ',   pool_name, ' check sum doesn''t have nine digits.') ; a_close(pool_file); get_strings_started:=false;  goto exit ; end ;
{ \xref[TEX.POOL check sum...] }
  a:=10*a+xord[n]-{"0"=}48;
  if k=9 then goto done;
  incr(k); read(pool_file,n);
  end;
done: if a<>@$ then
  begin    ; writeln( stdout ,'! ',   pool_name, ' doesn''t match; tangle me again (or fix the path).') ; a_close(pool_file); get_strings_started:=false;  goto exit ; end ;
{ \xref[TEX.POOL doesn't match] }
c:=true;
end


else  begin if (xord[m]<{"0"=}48)or(xord[m]>{"9"=}57)or 
      (xord[n]<{"0"=}48)or(xord[n]>{"9"=}57) then
    begin    ; writeln( stdout ,'! ',   pool_name, ' line doesn''t begin with two digits.') ; a_close(pool_file); get_strings_started:=false;  goto exit ; end ;
{ \xref[TEX.POOL line doesn't...] }
  l:=xord[m]*10+xord[n]-{"0"=}48*11; {compute the length}
  if pool_ptr+l+string_vacancies>pool_size then
    begin    ; writeln( stdout ,'! You have to increase POOLSIZE.') ; a_close(pool_file); get_strings_started:=false;  goto exit ; end ;
{ \xref[You have to increase POOLSIZE] }
  for k:=1 to l do
    begin if eoln(pool_file) then m:=' ' else read(pool_file,m);
     begin str_pool[pool_ptr]:=   xord[  m] ; incr(pool_ptr); end ;
    end;
  readln( pool_file) ; g:=make_string;
  end;
end

;
  until c;
  a_close(pool_file); get_strings_started:=true;
  end
else  begin    ; writeln( stdout ,'! I can''t read ',   pool_name, '; bad path?') ; a_close(pool_file); get_strings_started:=false;  goto exit ; end 
{ \xref[I can't read TEX.POOL] }

;
exit:end;
endif('INITEX') 



{ 56. }

{tangle:pos tex.web:1448:3: }

{ Macro abbreviations for output to the terminal and to the log file are
defined here for convenience. Some systems need special conventions
for terminal output, and it is possible to adhere to those conventions
by changing |wterm|, |wterm_ln|, and |wterm_cr| in this section.
\xref[system dependencies] }

{ 66. }

{tangle:pos tex.web:1638:3: }

{ Here is a trivial procedure to print two digits; it is usually called with
a parameter in the range |0<=n<=99|. } procedure print_two( n:integer); {prints two least significant digits}
begin n:=abs(n) mod 100; print_char({"0"=}48+(n div 10));
print_char({"0"=}48+(n mod 10));
end;



{ 67. }

{tangle:pos tex.web:1646:3: }

{ Hexadecimal printing of nonnegative integers is accomplished by |print_hex|. } procedure print_hex( n:integer);
  {prints a positive integer in hexadecimal form}
var k:0..22; {index to current digit; we assume that $0\L n<16^[22]$}
begin k:=0; print_char({""""=}34);
repeat dig[k]:=n mod 16; n:=n div 16; incr(k);
until n=0;
print_the_digs(k);
end;



{ 68. }

{tangle:pos tex.web:1657:3: }

{ Old versions of \TeX\ needed a procedure called |print_ASCII| whose function
is now subsumed by |print|. We retain the old name here as a possible aid to
future software arch\ae ologists. }

{ 69. }

{tangle:pos tex.web:1663:3: }

{ Roman numerals are produced by the |print_roman_int| routine.  Readers
who like puzzles might enjoy trying to figure out how this tricky code
works; therefore no explanation will be given. Notice that 1990 yields
\.[mcmxc], not \.[mxm]. } procedure print_roman_int( n:integer);
label exit;
var j, k: pool_pointer; {mysterious indices into |str_pool|}
 u, v: nonnegative_integer; {mysterious numbers}
begin j:=str_start[{"m2d5c2l5x2v5i"=}260]; v:=1000;
 while true do    begin while n>=v do
    begin print_char(  str_pool[ j] ); n:=n-v;
    end;
  if n<=0 then  goto exit ; {nonpositive input produces no output}
  k:=j+2; u:=v div (  str_pool[ k- 1] -{"0"=}48);
  if str_pool[k-1]= {"2"=} 50  then
    begin k:=k+2; u:=u div (  str_pool[ k- 1] -{"0"=}48);
    end;
  if n+u>=v then
    begin print_char(  str_pool[ k] ); n:=n+u;
    end
  else  begin j:=j+2; v:=v div (  str_pool[ j- 1] -{"0"=}48);
    end;
  end;
exit:end;



{ 70. }

{tangle:pos tex.web:1689:3: }

{ The |print| subroutine will not print a string that is still being
created. The following procedure will. } procedure print_current_string; {prints a yet-unmade string}
var j:pool_pointer; {points to current character code}
begin j:=str_start[str_ptr];
while j<pool_ptr do
  begin print_char(  str_pool[ j] ); incr(j);
  end;
end;



{ 71. }

{tangle:pos tex.web:1700:3: }

{ Here is a procedure that asks the user to type a line of input,
assuming that the |selector| setting is either |term_only| or |term_and_log|.
The input is placed into locations |first| through |last-1| of the
|buffer| array, and echoed on the transcript file if appropriate.

This procedure is never called when |interaction<scroll_mode|. } procedure term_input; {gets a line from the terminal}
var k:0..buf_size; {index into |buffer|}
begin  fflush (stdout ) ; {now the user sees the prompt for sure}
if not input_ln(stdin ,true) then begin
  cur_input.limit_field :=0; fatal_error({"End of file on the terminal!"=}261); end;
{ \xref[End of file on the terminal] }
term_offset:=0; {the user's line ended with \<\rm return>}
decr(selector); {prepare to echo the input}
if last<>first then for k:=first to last-1 do print(buffer[k]);
print_ln; incr(selector); {restore previous status}
end;



{ 72. \[6] Reporting errors }

{tangle:pos tex.web:1721:26: }

{ When something anomalous is detected, \TeX\ typically does something like this:
$$\vbox[\halign[#\hfil\cr
|print_err("Something anomalous has been detected");|\cr
|help3("This is the first line of my offer to help.")|\cr
|("This is the second line. I'm trying to")|\cr
|("explain the best way for you to proceed.");|\cr
|error;|\cr]]$$
A two-line help message would be given using |help2|, etc.; these informal
helps should use simple vocabulary that complements the words used in the
official error message that was printed. (Outside the U.S.A., the help
messages should preferably be translated into the local vernacular. Each
line of help is at most 60 characters long, in the present implementation,
so that |max_print_line| will not be exceeded.)

The |print_err| procedure supplies a `\.!' before the official message,
and makes sure that the terminal is awake if a stop is going to occur.
The |error| procedure supplies a `\..' after the official message, then it
shows the location of the error; and if |interaction=error_stop_mode|,
it also enters into a dialog with the user, during which time the help
message may be printed.
\xref[system dependencies] }

{ 91. }

{tangle:pos tex.web:2034:3: }

{ A dozen or so error messages end with a parenthesized integer, so we
save a teeny bit of program space by declaring the following procedure: } procedure int_error( n:integer);
begin print({" ("=}284); print_int(n); print_char({")"=}41); error;
end;



{ 92. }

{tangle:pos tex.web:2041:3: }

{ In anomalous cases, the print selector might be in an unknown state;
the following subroutine is called to fix things just enough to keep
running a bit longer. } procedure normalize_selector;
begin if log_opened then selector:=term_and_log
else selector:=term_only;
if job_name=0 then open_log_file;
if interaction=batch_mode then decr(selector);
end;



{ 98. }

{tangle:pos tex.web:2124:3: }

{ When an interrupt has been detected, the program goes into its
highest interaction level and lets the user have nearly the full flexibility of
the |error| routine.  \TeX\ checks for interrupts only at times when it is
safe to do this. } procedure pause_for_instructions;
begin if OK_to_interrupt then
  begin interaction:=error_stop_mode;
  if (selector=log_only)or(selector=no_print) then
    incr(selector);
  begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Interruption"=} 294); end ;
{ \xref[Interruption] }
   begin help_ptr:=3; help_line[2]:={"You rang?"=} 295; help_line[1]:={"Try to insert an instruction for me (e.g., `I\showlists'),"=} 296; help_line[0]:={"unless you just want to quit by typing `X'."=} 297; end ;
  deletions_allowed:=false; error; deletions_allowed:=true;
  interrupt:=0;
  end;
end;



{ 99. \[7] Arithmetic with scaled dimensions }

{tangle:pos tex.web:2144:43: }

{ The principal computations performed by \TeX\ are done entirely in terms of
integers less than $2^[31]$ in magnitude; and divisions are done only when both
dividend and divisor are nonnegative. Thus, the arithmetic specified in this
program can be carried out in exactly the same way on a wide variety of
computers, including some small ones. Why? Because the arithmetic
calculations need to be spelled out precisely in order to guarantee that
\TeX\ will produce identical output on different machines. If some
quantities were rounded differently in different implementations, we would
find that line breaks and even page breaks might occur in different places.
Hence the arithmetic of \TeX\ has been designed with care, and systems that
claim to be implementations of \TeX82 should follow precisely the
\xref[TeX82][\TeX82]
calculations as they appear in the present program.

(Actually there are three places where \TeX\ uses |div| with a possibly negative
numerator. These are harmless; see |div| in the index. Also if the user
sets the \.[\\time] or the \.[\\year] to a negative value, some diagnostic
information will involve negative-numerator division. The same remarks
apply for |mod| as well as for |div|.) }

{ 100. }

{tangle:pos tex.web:2165:3: }

{ Here is a routine that calculates half of an integer, using an
unambiguous convention with respect to signed odd numbers. } function half( x:integer):integer;
begin if odd(x) then half:=(x+1) div 2
else half:=x  div 2;
end;



{ 102. }

{tangle:pos tex.web:2185:3: }

{ The following function is used to create a scaled integer from a given decimal
fraction $(.d_0d_1\ldots d_[k-1])$, where |0<=k<=17|. The digit $d_i$ is
given in |dig[i]|, and the calculation produces a correctly rounded result. } function round_decimals( k:small_number) : scaled;
  {converts a decimal fraction}
var a:integer; {the accumulator}
begin a:=0;
while k>0 do
  begin decr(k); a:=(a+dig[k]* {0400000=}131072 ) div 10;
  end;
round_decimals:=(a+1) div 2;
end;



{ 103. }

{tangle:pos tex.web:2199:3: }

{ Conversely, here is a procedure analogous to |print_int|. If the output
of this procedure is subsequently read by \TeX\ and converted by the
|round_decimals| routine above, it turns out that the original value will
be reproduced exactly; the ``simplest'' such decimal number is output,
but there is always at least one digit following the decimal point.

The invariant relation in the \&[repeat] loop is that a sequence of
decimal digits yet to be printed will yield the original number if and only if
they form a fraction~$f$ in the range $s-\delta\L10\cdot2^[16]f<s$.
We can stop if and only if $f=0$ satisfies this condition; the loop will
terminate before $s$ can possibly become zero. } procedure print_scaled( s:scaled); {prints scaled real, rounded to five
  digits}
var delta:scaled; {amount of allowable inaccuracy}
begin if s<0 then
  begin print_char({"-"=}45);   s:=- s ; {print the sign, if negative}
  end;
print_int(s div  {0200000=}65536 ); {print the integer part}
print_char({"."=}46);
s:=10*(s mod  {0200000=}65536 )+5; delta:=10;
repeat if delta> {0200000=}65536  then s:=s+{0100000=}32768-50000; {round the last digit}
print_char({"0"=}48+(s div  {0200000=}65536 )); s:=10*(s mod  {0200000=}65536 ); delta:=delta*10;
until s<=delta;
end;



{ 105. }

{tangle:pos tex.web:2254:3: }

{ The first arithmetical subroutine we need computes $nx+y$, where |x|
and~|y| are |scaled| and |n| is an integer. We will also use it to
multiply integers. } function mult_and_add( n:integer; x, y, max_answer:scaled):scaled;
begin if n<0 then
  begin   x:=- x ;   n:=- n ;
  end;
if n=0 then mult_and_add:=y
else if ((x<=(max_answer-y) div n)and(-x<=(max_answer+y) div n)) then
  mult_and_add:=n*x+y
else  begin arith_error:=true; mult_and_add:=0;
  end;
end;



{ 106. }

{tangle:pos tex.web:2272:3: }

{ We also need to divide scaled dimensions by integers. } function x_over_n( x:scaled; n:integer):scaled;
var negative:boolean; {should |remainder| be negated?}
begin negative:=false;
if n=0 then
  begin arith_error:=true; x_over_n:=0; tex_remainder :=x;
  end
else  begin if n<0 then
    begin   x:=- x ;   n:=- n ; negative:=true;
    end;
  if x>=0 then
    begin x_over_n:=x div n; tex_remainder :=x mod n;
    end
  else  begin x_over_n:=-((-x) div n); tex_remainder :=-((-x) mod n);
    end;
  end;
if negative then   tex_remainder :=- tex_remainder  ;
end;



{ 107. }

{tangle:pos tex.web:2292:3: }

{ Then comes the multiplication of a scaled number by a fraction |n/d|,
where |n| and |d| are nonnegative integers |<=$2^[16]$| and |d| is
positive. It would be too dangerous to multiply by~|n| and then divide
by~|d|, in separate operations, since overflow might well occur; and it
would be too inaccurate to divide by |d| and then multiply by |n|. Hence
this subroutine simulates 1.5-precision arithmetic. } function xn_over_d( x:scaled;  n, d:integer):scaled;
var positive:boolean; {was |x>=0|?}
 t, u, v:nonnegative_integer; {intermediate quantities}
begin if x>=0 then positive:=true
else  begin   x:=- x ; positive:=false;
  end;
t:=(x mod {0100000=}32768)*n;
u:=(x div {0100000=}32768)*n+(t div {0100000=}32768);
v:=(u mod d)*{0100000=}32768 + (t mod {0100000=}32768);
if u div d>={0100000=}32768 then arith_error:=true
else u:={0100000=}32768*(u div d) + (v div d);
if positive then
  begin xn_over_d:=u; tex_remainder :=v mod d;
  end
else  begin xn_over_d:=-u; tex_remainder :=-(v mod d);
  end;
end;



{ 108. }

{tangle:pos tex.web:2317:3: }

{ The next subroutine is used to compute the ``badness'' of glue, when a
total~|t| is supposed to be made from amounts that sum to~|s|.  According
to [\sl The \TeX book], the badness of this situation is $100(t/s)^3$;
however, badness is simply a heuristic, so we need not squeeze out the
last drop of accuracy when computing it. All we really want is an
approximation that has similar properties.
\xref[TeXbook][\sl The \TeX book]

The actual method used to compute the badness is easier to read from the
program than to describe in words. It produces an integer value that is a
reasonably close approximation to $100(t/s)^3$, and all implementations
of \TeX\ should use precisely this method. Any badness of $2^[13]$ or more is
treated as infinitely bad, and represented by 10000.

It is not difficult to prove that $$\hbox[|badness(t+1,s)>=badness(t,s)
>=badness(t,s+1)|].$$ The badness function defined here is capable of
computing at most 1095 distinct values, but that is plenty. } function badness( t, s:scaled):halfword; {compute badness, given |t>=0|}
var r:integer; {approximation to $\alpha t/s$, where $\alpha^3\approx
  100\cdot2^[18]$}
begin if t=0 then badness:=0
else if s<=0 then badness:=inf_bad
else  begin if t<=7230584 then  r:=(t*297) div s {$297^3=99.94\times2^[18]$}
  else if s>=1663497 then r:=t div (s div 297)
  else r:=t;
  if r>1290 then badness:=inf_bad {$1290^3<2^[31]<1291^3$}
  else badness:=(r*r*r+{0400000=}131072) div {01000000=}262144;
  end; {that was $r^3/2^[18]$, rounded to the nearest integer}
end;



{ 110. \[8] Packed data }

{tangle:pos tex.web:2375:21: }

{ In order to make efficient use of storage space, \TeX\ bases its major data
structures on a |memory_word|, which contains either a (signed) integer,
possibly scaled, or a (signed) |glue_ratio|, or a small number of
fields that are one half or one quarter of the size used for storing
integers.

If |x| is a variable of type |memory_word|, it contains up to four
fields that can be referred to as follows:
$$\vbox[\halign[\hfil#&#\hfil&#\hfil\cr
|x|&.|int|&(an |integer|)\cr
|x|&.|sc|\qquad&(a |scaled| integer)\cr
|x|&.|gr|&(a |glue_ratio|)\cr
|x.hh.lh|, |x.hh|&.|rh|&(two halfword fields)\cr
|x.hh.b0|, |x.hh.b1|, |x.hh|&.|rh|&(two quarterword fields, one halfword
  field)\cr
|x.qqqq.b0|, |x.qqqq.b1|, |x.qqqq|&.|b2|, |x.qqqq.b3|\hskip-100pt
  &\qquad\qquad\qquad(four quarterword fields)\cr]]$$
This is somewhat cumbersome to write, and not very readable either, but
macros will be used to make the notation shorter and more transparent.
The \PASCAL\ code below gives a formal definition of |memory_word| and
its subsidiary types, using packed variant records. \TeX\ makes no
assumptions about the relative positions of the fields within a word.

Since we are assuming 32-bit integers, a halfword must contain at least
16 bits, and a quarterword must contain at least 8 bits.
\xref[system dependencies]
But it doesn't hurt to have more bits; for example, with enough 36-bit
words you might be able to have |mem_max| as large as 262142, which is
eight times as much memory as anybody had during the first four years of
\TeX's existence.

N.B.: Valuable memory space will be dreadfully wasted unless \TeX\ is compiled
by a \PASCAL\ that packs all of the |memory_word| variants into
the space of a single integer. This means, for example, that |glue_ratio|
words should be |short_real| instead of |real| on some computers. Some
\PASCAL\ compilers will pack an integer whose subrange is `|0..255|' into
an eight-bit field, but others insist on allocating space for an additional
sign bit; on such systems you can get 256 values into a quarterword only
if the subrange is `|-128..127|'.

The present implementation tries to accommodate as many variations as possible,
so it makes few assumptions. If integers having the subrange
`|min_quarterword..max_quarterword|' can be packed into a quarterword,
and if integers having the subrange `|min_halfword..max_halfword|'
can be packed into a halfword, everything should work satisfactorily.

It is usually most efficient to have |min_quarterword=min_halfword=0|,
so one should try to achieve this unless it causes a severe problem.
The values defined here are recommended for most 32-bit computers. }

{ 112. }

{tangle:pos tex.web:2449:3: }

{ The operation of adding or subtracting |min_quarterword| occurs quite
frequently in \TeX, so it is convenient to abbreviate this operation
by using the macros |qi| and |qo| for input and output to and from
quarterword format.

The inner loop of \TeX\ will run faster with respect to compilers
that don't optimize expressions like `|x+0|' and `|x-0|', if these
macros are simplified in the obvious way when |min_quarterword=0|.
So they have been simplified here in the obvious way.
\xref[inner loop]\xref[system dependencies]

The \.[WEB] source for \TeX\ defines |hi(#)==#+min_halfword| which can be
simplified when |min_halfword=0|.  The Web2C implementation of \TeX\ can use
|hi(#)==#| together with |min_halfword<0| as long as |max_halfword| is
sufficiently large. }

{ 114. }

{tangle:pos tex.web:2499:3: }

{ When debugging, we may want to print a |memory_word| without knowing
what type it is; so we print it in all modes.
\xref[dirty \PASCAL]\xref[debugging] }  ifdef('TEXMF_DEBUG')  procedure print_word( w:memory_word);
  {prints |w| in all ways}
begin print_int(w.int); print_char({" "=}32);

print_scaled(w.int ); print_char({" "=}32);

print_scaled(round( {0200000=}65536 *  w. gr )); print_ln;

{ \xref[real multiplication] }
print_int(w.hh.lh); print_char({"="=}61); print_int(w.hh.b0); print_char({":"=}58);
print_int(w.hh.b1); print_char({";"=}59); print_int(w.hh.rh); print_char({" "=}32);

print_int(w.qqqq.b0); print_char({":"=}58); print_int(w.qqqq.b1); print_char({":"=}58);
print_int(w.qqqq.b2); print_char({":"=}58); print_int(w.qqqq.b3);
end;
endif('TEXMF_DEBUG') 



{ 119. }

{tangle:pos tex.web:2596:3: }

{ If memory is exhausted, it might mean that the user has forgotten
a right brace. We will define some procedures later that try to help
pinpoint the trouble. } 
{ Declare the procedure called |show_token_list| }
procedure show_token_list( p, q:integer; l:integer);
label exit;
var m, c:integer; {pieces of a token}
 match_chr:ASCII_code; {character used in a `|match|'}
 n:ASCII_code; {the highest parameter number, as an ASCII digit}
begin match_chr:={"#"=}35; n:={"0"=}48; tally:=0;
while (p<>-{0xfffffff=}268435455  ) and (tally<l) do
  begin if p=q then 
{ Do magic computation } begin first_count:=tally; trick_count:=tally+1+error_line-half_error_line; if trick_count<error_line then trick_count:=error_line; end 

;
  
{ Display token |p|, and |return| if there are problems }
if (p<hi_mem_min) or (p>mem_end) then
  begin print_esc({"CLOBBERED."=}307);  goto exit ;
{ \xref[CLOBBERED] }
  end;
if  mem[ p].hh.lh >={07777=}4095  then print_cs( mem[ p].hh.lh -{07777=}4095 )
else  begin m:= mem[ p].hh.lh  div {0400=}256; c:= mem[ p].hh.lh  mod {0400=}256;
  if  mem[ p].hh.lh <0 then print_esc({"BAD."=}563)
{ \xref[BAD] }
  else 
{ Display the token $(|m|,|c|)$ }
case m of
left_brace,right_brace,math_shift,tab_mark,sup_mark,sub_mark,spacer,
  letter,other_char: print(c);
mac_param: begin print(c); print(c);
  end;
out_param: begin print(match_chr);
  if c<=9 then print_char(c+{"0"=}48)
  else  begin print_char({"!"=}33);  goto exit ;
    end;
  end;
match: begin match_chr:=c; print(c); incr(n); print_char(n);
  if n>{"9"=}57 then  goto exit ;
  end;
end_match: print({"->"=}564);
{ \xref[->] }
 else  print_esc({"BAD."=}563)
{ \xref[BAD] }
 end 

;
  end

;
  p:= mem[ p].hh.rh ;
  end;
if p<>-{0xfffffff=}268435455   then print_esc({"ETC."=}562);
{ \xref[ETC] }
exit:
end;




{ Declare the procedure called |runaway| }
procedure runaway;
var p:halfword ; {head of runaway list}
begin if scanner_status>skipping then
  begin
{ \xref[Runaway...] }
  case scanner_status of
  defining: begin print_nl({"Runaway definition"=}577); p:=def_ref;
    end;
  matching: begin print_nl({"Runaway argument"=}578); p:=mem_top-3 ;
    end;
  aligning: begin print_nl({"Runaway preamble"=}579); p:=mem_top-4 ;
    end;
  absorbing: begin print_nl({"Runaway text"=}580); p:=def_ref;
    end;
  end; {there are no other cases}
  print_char({"?"=}63);print_ln; show_token_list( mem[ p].hh.rh ,-{0xfffffff=}268435455  ,error_line-10);
  end;
end;





{ 120. }

{tangle:pos tex.web:2603:3: }

{ The function |get_avail| returns a pointer to a new one-word node whose
|link| field is null. However, \TeX\ will halt if there is no more room left.
\xref[inner loop]

If the available-space list is empty, i.e., if |avail=null|,
we try first to increase |mem_end|. If that cannot be done, i.e., if
|mem_end=mem_max|, we try to decrease |hi_mem_min|. If that cannot be
done, i.e., if |hi_mem_min=lo_mem_max+1|, we have to quit. } function get_avail : halfword ; {single-word node allocation}
var p:halfword ; {the new node being got}
begin p:=avail; {get top location in the |avail| stack}
if p<>-{0xfffffff=}268435455   then avail:= mem[ avail].hh.rh  {and pop it off}
else if mem_end<mem_max then {or go into virgin territory}
  begin incr(mem_end); p:=mem_end;
  end
else   begin decr(hi_mem_min); p:=hi_mem_min;
  if hi_mem_min<=lo_mem_max then
    begin runaway; {if memory is exhausted, display possible runaway text}
    overflow({"main memory size"=}298,mem_max+1-mem_min);
      {quit; all one-word nodes are busy}
{ \xref[TeX capacity exceeded main memory size][\quad main memory size] }
    end;
  end;
 mem[ p].hh.rh :=-{0xfffffff=}268435455  ; {provide an oft-desired initialization of the new node}
 ifdef('STAT')  incr(dyn_used); endif('STAT')  {maintain statistics}
get_avail:=p;
end;



{ 121. }

{tangle:pos tex.web:2632:3: }

{ Conversely, a one-word node is recycled by calling |free_avail|.
This routine is part of \TeX's ``inner loop,'' so we want it to be fast.
\xref[inner loop] }

{ 122. }

{tangle:pos tex.web:2641:3: }

{ There's also a |fast_get_avail| routine, which saves the procedure-call
overhead at the expense of extra programming. This routine is used in
the places that would otherwise account for the most calls of |get_avail|.
\xref[inner loop] }

{ 123. }

{tangle:pos tex.web:2654:3: }

{ The procedure |flush_list(p)| frees an entire linked list of
one-word nodes that starts at position |p|.
\xref[inner loop] } procedure flush_list( p:halfword ); {makes list of single-word nodes
  available}
var  q, r:halfword ; {list traversers}
begin if p<>-{0xfffffff=}268435455   then
  begin r:=p;
  repeat q:=r; r:= mem[ r].hh.rh ;  ifdef('STAT')  decr(dyn_used); endif('STAT') 

  until r=-{0xfffffff=}268435455  ; {now |q| is the last node on the list}
   mem[ q].hh.rh :=avail; avail:=p;
  end;
end;



{ 125. }

{tangle:pos tex.web:2694:3: }

{ A call to |get_node| with argument |s| returns a pointer to a new node
of size~|s|, which must be 2~or more. The |link| field of the first word
of this new node is set to null. An overflow stop occurs if no suitable
space exists.

If |get_node| is called with $s=2^[30]$, it simply merges adjacent free
areas and returns the value |max_halfword|. } function get_node( s:integer):halfword ; {variable-size node allocation}
label found,exit,restart;
var p:halfword ; {the node currently under inspection}
 q:halfword ; {the node physically after node |p|}
 r:integer; {the newly allocated node, or a candidate for this honor}
 t:integer; {temporary register}
begin restart: p:=rover; {start at some free node in the ring}
repeat 
{ Try to allocate within node |p| and its physical successors, and |goto found| if allocation was possible }
q:=p+  mem[ p].hh.lh ; {find the physical successor}
{ \xref[inner loop] }
while  ( mem[  q].hh.rh = {0xfffffff=}268435455  )  do {merge node |p| with node |q|}
  begin t:=  mem[  q+ 1].hh.rh  ;
  if q=rover then rover:=t;
    mem[  t+ 1].hh.lh  :=  mem[  q+ 1].hh.lh  ;   mem[    mem[    q+ 1].hh.lh  + 1].hh.rh  :=t;

  q:=q+  mem[ q].hh.lh ;
  end;
r:=q-s;
if r>p+1 then 
{ Allocate from the top of node |p| and |goto found| }
begin   mem[ p].hh.lh :=r-p; {store the remaining size}
{ \xref[inner loop] }
rover:=p; {start searching here next time}
goto found;
end

;
if r=p then if   mem[  p+ 1].hh.rh  <>p then
  
{ Allocate entire node |p| and |goto found| }
begin rover:=  mem[  p+ 1].hh.rh  ; t:=  mem[  p+ 1].hh.lh  ;
  mem[  rover+ 1].hh.lh  :=t;   mem[  t+ 1].hh.rh  :=rover;
goto found;
end

;
  mem[ p].hh.lh :=q-p {reset the size in case it grew}

;
{ \xref[inner loop] }
p:=  mem[  p+ 1].hh.rh  ; {move to the next node in the ring}
until p=rover; {repeat until the whole list has been traversed}
if s={010000000000=}1073741824 then
  begin get_node:={0xfffffff=}268435455 ;  goto exit ;
  end;
if lo_mem_max+2<hi_mem_min then if lo_mem_max+2<=mem_bot+{0xfffffff=}268435455  then
  
{ Grow more variable-size memory and |goto restart| }
begin if hi_mem_min-lo_mem_max>=1998 then t:=lo_mem_max+1000
else t:=lo_mem_max+1+(hi_mem_min-lo_mem_max) div 2;
  {|lo_mem_max+2<=t<hi_mem_min|}
p:=  mem[  rover+ 1].hh.lh  ; q:=lo_mem_max;   mem[  p+ 1].hh.rh  :=q;   mem[  rover+ 1].hh.lh  :=q;

if t>mem_bot+{0xfffffff=}268435455  then t:=mem_bot+{0xfffffff=}268435455 ;
  mem[  q+ 1].hh.rh  :=rover;   mem[  q+ 1].hh.lh  :=p;  mem[ q].hh.rh := {0xfffffff=}268435455  ;   mem[ q].hh.lh :=t-lo_mem_max;

lo_mem_max:=t;  mem[ lo_mem_max].hh.rh :=-{0xfffffff=}268435455  ;  mem[ lo_mem_max].hh.lh :=-{0xfffffff=}268435455  ;
rover:=q; goto restart;
end

;
overflow({"main memory size"=}298,mem_max+1-mem_min);
  {sorry, nothing satisfactory is left}
{ \xref[TeX capacity exceeded main memory size][\quad main memory size] }
found:  mem[ r].hh.rh :=-{0xfffffff=}268435455  ; {this node is now nonempty}
 ifdef('STAT')  var_used:=var_used+s; {maintain usage statistics}
endif('STAT')  

get_node:=r;
exit:end;



{ 130. }

{tangle:pos tex.web:2780:3: }

{ Conversely, when some variable-size node |p| of size |s| is no longer needed,
the operation |free_node(p,s)| will make its words available, by inserting
|p| as a new empty node just before where |rover| now points.
\xref[inner loop] } procedure free_node( p:halfword ;  s:halfword); {variable-size node
  liberation}
var q:halfword ; {|llink(rover)|}
begin   mem[ p].hh.lh :=s;  mem[ p].hh.rh := {0xfffffff=}268435455  ;
q:=  mem[  rover+ 1].hh.lh  ;   mem[  p+ 1].hh.lh  :=q;   mem[  p+ 1].hh.rh  :=rover; {set both links}
  mem[  rover+ 1].hh.lh  :=p;   mem[  q+ 1].hh.rh  :=p; {insert |p| into the ring}
 ifdef('STAT')  var_used:=var_used-s; endif('STAT')  {maintain statistics}
end;



{ 131. }

{tangle:pos tex.web:2794:3: }

{ Just before \.[INITEX] writes out the memory, it sorts the doubly linked
available space list. The list is probably very short at such times, so a
simple insertion sort is used. The smallest available location will be
pointed to by |rover|, the next-smallest by |rlink(rover)|, etc. }  ifdef('INITEX')  procedure sort_avail; {sorts the available variable-size nodes
  by location}
var p, q, r: halfword ; {indices into |mem|}
 old_rover:halfword ; {initial |rover| setting}
begin p:=get_node({010000000000=}1073741824); {merge adjacent free areas}
p:=  mem[  rover+ 1].hh.rh  ;   mem[  rover+ 1].hh.rh  :={0xfffffff=}268435455 ; old_rover:=rover;
while p<>old_rover do 
{ Sort \(p)|p| into the list starting at |rover| and advance |p| to |rlink(p)| }
if p<rover then
  begin q:=p; p:=  mem[  q+ 1].hh.rh  ;   mem[  q+ 1].hh.rh  :=rover; rover:=q;
  end
else  begin q:=rover;
  while   mem[  q+ 1].hh.rh  <p do q:=  mem[  q+ 1].hh.rh  ;
  r:=  mem[  p+ 1].hh.rh  ;   mem[  p+ 1].hh.rh  :=  mem[  q+ 1].hh.rh  ;   mem[  q+ 1].hh.rh  :=p; p:=r;
  end

;
p:=rover;
while   mem[  p+ 1].hh.rh  <>{0xfffffff=}268435455  do
  begin   mem[    mem[    p+ 1].hh.rh  + 1].hh.lh  :=p; p:=  mem[  p+ 1].hh.rh  ;
  end;
  mem[  p+ 1].hh.rh  :=rover;   mem[  rover+ 1].hh.lh  :=p;
end;
endif('INITEX') 



{ 133. \[10] Data structures for boxes and their friends }

{tangle:pos tex.web:2828:54: }

{ From the computer's standpoint, \TeX's chief mission is to create
horizontal and vertical lists. We shall now investigate how the elements
of these lists are represented internally as nodes in the dynamic memory.

A horizontal or vertical list is linked together by |link| fields in
the first word of each node. Individual nodes represent boxes, glue,
penalties, or special things like discretionary hyphens; because of this
variety, some nodes are longer than others, and we must distinguish different
kinds of nodes. We do this by putting a `|type|' field in the first word,
together with the link and an optional `|subtype|'. }

{ 134. }

{tangle:pos tex.web:2843:3: }

{ A | char_node|, which represents a single character, is the most important
kind of node because it accounts for the vast majority of all boxes.
Special precautions are therefore taken to ensure that a |char_node| does
not take up much memory space. Every such node is one word long, and in fact
it is identifiable by this property, since other kinds of nodes have at least
two words, and they appear in |mem| locations less than |hi_mem_min|.
This makes it possible to omit the |type| field in a |char_node|, leaving
us room for two bytes that identify a |font| and a |character| within
that font.

Note that the format of a |char_node| allows for up to 256 different
fonts and up to 256 characters per font; but most implementations will
probably limit the total number of fonts to fewer than 75 per job,
and most fonts will stick to characters whose codes are
less than 128 (since higher codes
are more difficult to access on most keyboards).

Extensions of \TeX\ intended for oriental languages will need even more
than $256\times256$ possible characters, when we consider different sizes
\xref[oriental characters]\xref[Chinese characters]\xref[Japanese characters]
and styles of type.  It is suggested that Chinese and Japanese fonts be
handled by representing such characters in two consecutive |char_node|
entries: The first of these has |font=font_base|, and its |link| points
to the second;
the second identifies the font and the character dimensions.
The saving feature about oriental characters is that most of them have
the same box dimensions. The |character| field of the first |char_node|
is a ``\\[charext]'' that distinguishes between graphic symbols whose
dimensions are identical for typesetting purposes. (See the \MF\ manual.)
Such an extension of \TeX\ would not be difficult; further details are
left to the reader.

In order to make sure that the |character| code fits in a quarterword,
\TeX\ adds the quantity |min_quarterword| to the actual code.

Character nodes appear only in horizontal lists, never in vertical lists. }

{ 135. }

{tangle:pos tex.web:2885:3: }

{ An |hlist_node| stands for a box that was made from a horizontal list.
Each |hlist_node| is seven words long, and contains the following fields
(in addition to the mandatory |type| and |link|, which we shall not
mention explicitly when discussing the other node types): The |height| and
|width| and |depth| are scaled integers denoting the dimensions of the
box.  There is also a |shift_amount| field, a scaled integer indicating
how much this box should be lowered (if it appears in a horizontal list),
or how much it should be moved to the right (if it appears in a vertical
list). There is a |list_ptr| field, which points to the beginning of the
list from which this box was fabricated; if |list_ptr| is |null|, the box
is empty. Finally, there are three fields that represent the setting of
the glue:  |glue_set(p)| is a word of type |glue_ratio| that represents
the proportionality constant for glue setting; |glue_sign(p)| is
|stretching| or |shrinking| or |normal| depending on whether or not the
glue should stretch or shrink or remain rigid; and |glue_order(p)|
specifies the order of infinity to which glue setting applies (|normal|,
|fil|, |fill|, or |filll|). The |subtype| field is not used. }

{ 136. }

{tangle:pos tex.web:2923:3: }

{ The |new_null_box| function returns a pointer to an |hlist_node| in
which all subfields have the values corresponding to `\.[\\hbox\[\]]'.
(The |subtype| field is set to |min_quarterword|, for historic reasons
that are no longer relevant.) } function new_null_box:halfword ; {creates a new box node}
var p:halfword ; {the new node}
begin p:=get_node(box_node_size);  mem[ p].hh.b0 :=hlist_node;
 mem[ p].hh.b1 :=min_quarterword;
 mem[ p+width_offset].int  :=0;  mem[ p+depth_offset].int  :=0;  mem[ p+height_offset].int  :=0;  mem[ p+4].int  :=0;   mem[  p+ list_offset].hh.rh  :=-{0xfffffff=}268435455  ;
  mem[  p+ list_offset].hh.b0  :=normal;   mem[  p+ list_offset].hh.b1  :=normal;    mem[  p+glue_offset].gr :=0.0 ;
new_null_box:=p;
end;



{ 137. }

{tangle:pos tex.web:2937:3: }

{ A |vlist_node| is like an |hlist_node| in all respects except that it
contains a vertical list. }

{ 138. }

{tangle:pos tex.web:2942:3: }

{ A |rule_node| stands for a solid black rectangle; it has |width|,
|depth|, and |height| fields just as in an |hlist_node|. However, if
any of these dimensions is $-2^[30]$, the actual value will be determined
by running the rule up to the boundary of the innermost enclosing box.
This is called a ``running dimension.'' The |width| is never running in
an hlist; the |height| and |depth| are never running in a~vlist. }

{ 139. }

{tangle:pos tex.web:2954:3: }

{ A new rule node is delivered by the |new_rule| function. It
makes all the dimensions ``running,'' so you have to change the
ones that are not allowed to run. } function new_rule:halfword ;
var p:halfword ; {the new node}
begin p:=get_node(rule_node_size);  mem[ p].hh.b0 :=rule_node;
 mem[ p].hh.b1 :=0; {the |subtype| is not used}
 mem[ p+width_offset].int  :=-{010000000000=}1073741824 ;  mem[ p+depth_offset].int  :=-{010000000000=}1073741824 ;  mem[ p+height_offset].int  :=-{010000000000=}1073741824 ;
new_rule:=p;
end;



{ 140. }

{tangle:pos tex.web:2966:3: }

{ Insertions are represented by |ins_node| records, where the |subtype|
indicates the corresponding box number. For example, `\.[\\insert 250]'
leads to an |ins_node| whose |subtype| is |250+min_quarterword|.
The |height| field of an |ins_node| is slightly misnamed; it actually holds
the natural height plus depth of the vertical list being inserted.
The |depth| field holds the |split_max_depth| to be used in case this
insertion is split, and the |split_top_ptr| points to the corresponding
|split_top_skip|. The |float_cost| field holds the |floating_penalty| that
will be used if this insertion floats to a subsequent page after a
split insertion of the same class.  There is one more field, the
|ins_ptr|, which points to the beginning of the vlist for the insertion. }

{ 141. }

{tangle:pos tex.web:2984:3: }

{ A |mark_node| has a |mark_ptr| field that points to the reference count
of a token list that contains the user's \.[\\mark] text.
This field occupies a full word instead of a halfword, because
there's nothing to put in the other halfword; it is easier in \PASCAL\ to
use the full word than to risk leaving garbage in the unused half. }

{ 142. }

{tangle:pos tex.web:2994:3: }

{ An |adjust_node|, which occurs only in horizontal lists,
specifies material that will be moved out into the surrounding
vertical list; i.e., it is used to implement \TeX's `\.[\\vadjust]'
operation.  The |adjust_ptr| field points to the vlist containing this
material. }

{ 143. }

{tangle:pos tex.web:3003:3: }

{ A |ligature_node|, which occurs only in horizontal lists, specifies
a character that was fabricated from the interaction of two or more
actual characters.  The second word of the node, which is called the
|lig_char| word, contains |font| and |character| fields just as in a
|char_node|. The characters that generated the ligature have not been
forgotten, since they are needed for diagnostic messages and for
hyphenation; the |lig_ptr| field points to a linked list of character
nodes for all original characters that have been deleted. (This list
might be empty if the characters that generated the ligature were
retained in other nodes.)

The |subtype| field is 0, plus 2 and/or 1 if the original source of the
ligature included implicit left and/or right boundaries. }

{ 144. }

{tangle:pos tex.web:3021:3: }

{ The |new_ligature| function creates a ligature node having given
contents of the |font|, |character|, and |lig_ptr| fields. We also have
a |new_lig_item| function, which returns a two-word node having a given
|character| field. Such nodes are used for temporary processing as ligatures
are being created. } function new_ligature( f:internal_font_number;  c:quarterword;
                          q:halfword ):halfword ;
var p:halfword ; {the new node}
begin p:=get_node(small_node_size);  mem[ p].hh.b0 :=ligature_node;
  mem[   p+1 ].hh.b0 :=f;   mem[   p+1 ].hh.b1 :=c;  mem[    p+1 ].hh.rh  :=q;
 mem[ p].hh.b1 :=0; new_ligature:=p;
end;


function new_lig_item( c:quarterword):halfword ;
var p:halfword ; {the new node}
begin p:=get_node(small_node_size);   mem[ p].hh.b1 :=c;  mem[    p+1 ].hh.rh  :=-{0xfffffff=}268435455  ;
new_lig_item:=p;
end;



{ 145. }

{tangle:pos tex.web:3040:3: }

{ A |disc_node|, which occurs only in horizontal lists, specifies a
``dis\-cretion\-ary'' line break. If such a break occurs at node |p|, the text
that starts at |pre_break(p)| will precede the break, the text that starts at
|post_break(p)| will follow the break, and text that appears in the next
|replace_count(p)| nodes will be ignored. For example, an ordinary
discretionary hyphen, indicated by `\.[\\-]', yields a |disc_node| with
|pre_break| pointing to a |char_node| containing a hyphen, |post_break=null|,
and |replace_count=0|. All three of the discretionary texts must be
lists that consist entirely of character, kern, box, rule, and ligature nodes.

If |pre_break(p)=null|, the |ex_hyphen_penalty| will be charged for this
break.  Otherwise the |hyphen_penalty| will be charged.  The texts will
actually be substituted into the list by the line-breaking algorithm if it
decides to make the break, and the discretionary node will disappear at
that time; thus, the output routine sees only discretionaries that were
not chosen. } function new_disc:halfword ; {creates an empty |disc_node|}
var p:halfword ; {the new node}
begin p:=get_node(small_node_size);  mem[ p].hh.b0 :=disc_node;
 mem[ p].hh.b1 :=0;   mem[  p+ 1].hh.lh  :=-{0xfffffff=}268435455  ;   mem[  p+ 1].hh.rh  :=-{0xfffffff=}268435455  ;
new_disc:=p;
end;



{ 146. }

{tangle:pos tex.web:3069:3: }

{ A |whatsit_node| is a wild card reserved for extensions to \TeX. The
|subtype| field in its first word says what `\\[whatsit]' it is, and
implicitly determines the node size (which must be 2 or more) and the
format of the remaining words. When a |whatsit_node| is encountered
in a list, special actions are invoked; knowledgeable people who are
careful not to mess up the rest of \TeX\ are able to make \TeX\ do new
things by adding code at the end of the program. For example, there
might be a `\TeX nicolor' extension to specify different colors of ink,
\xref[extensions to \TeX]
and the whatsit node might contain the desired parameters.

The present implementation of \TeX\ treats the features associated with
`\.[\\write]' and `\.[\\special]' as if they were extensions, in order to
illustrate how such routines might be coded. We shall defer further
discussion of extensions until the end of this program. }

{ 147. }

{tangle:pos tex.web:3087:3: }

{ A |math_node|, which occurs only in horizontal lists, appears before and
after mathematical formulas. The |subtype| field is |before| before the
formula and |after| after it. There is a |width| field, which represents
the amount of surrounding space inserted by \.[\\mathsurround]. } function new_math( w:scaled; s:small_number):halfword ;
var p:halfword ; {the new node}
begin p:=get_node(small_node_size);  mem[ p].hh.b0 :=math_node;
 mem[ p].hh.b1 :=s;  mem[ p+width_offset].int  :=w; new_math:=p;
end;



{ 148. }

{tangle:pos tex.web:3102:3: }

{ \TeX\ makes use of the fact that |hlist_node|, |vlist_node|,
|rule_node|, |ins_node|, |mark_node|, |adjust_node|, |ligature_node|,
|disc_node|, |whatsit_node|, and |math_node| are at the low end of the
type codes, by permitting a break at glue in a list if and only if the
|type| of the previous node is less than |math_node|. Furthermore, a
node is discarded after a break if its type is |math_node| or~more. }

{ 149. }

{tangle:pos tex.web:3112:3: }

{ A |glue_node| represents glue in a list. However, it is really only
a pointer to a separate glue specification, since \TeX\ makes use of the
fact that many essentially identical nodes of glue are usually present.
If |p| points to a |glue_node|, |glue_ptr(p)| points to
another packet of words that specify the stretch and shrink components, etc.

Glue nodes also serve to represent leaders; the |subtype| is used to
distinguish between ordinary glue (which is called |normal|) and the three
kinds of leaders (which are called |a_leaders|, |c_leaders|, and |x_leaders|).
The |leader_ptr| field points to a rule node or to a box node containing the
leaders; it is set to |null| in ordinary glue nodes.

Many kinds of glue are computed from \TeX's ``skip'' parameters, and
it is helpful to know which parameter has led to a particular glue node.
Therefore the |subtype| is set to indicate the source of glue, whenever
it originated as a parameter. We will be defining symbolic names for the
parameter numbers later (e.g., |line_skip_code=0|, |baseline_skip_code=1|,
etc.); it suffices for now to say that the |subtype| of parametric glue
will be the same as the parameter number, plus~one.

In math formulas there are two more possibilities for the |subtype| in a
glue node: |mu_glue| denotes an \.[\\mskip] (where the units are scaled \.[mu]
instead of scaled \.[pt]); and |cond_math_glue| denotes the `\.[\\nonscript]'
feature that cancels the glue node immediately following if it appears
in a subscript. }

{ 151. }

{tangle:pos tex.web:3174:3: }

{ Here is a function that returns a pointer to a copy of a glue spec.
The reference count in the copy is |null|, because there is assumed
to be exactly one reference to the new specification. } function new_spec( p:halfword ):halfword ; {duplicates a glue specification}
var q:halfword ; {the new spec}
begin q:=get_node(glue_spec_size);

mem[q]:=mem[p];   mem[  q].hh.rh  :=-{0xfffffff=}268435455  ;

 mem[ q+width_offset].int  := mem[ p+width_offset].int  ;  mem[ q+2].int  := mem[ p+2].int  ;  mem[ q+3].int  := mem[ p+3].int  ;
new_spec:=q;
end;



{ 152. }

{tangle:pos tex.web:3186:3: }

{ And here's a function that creates a glue node for a given parameter
identified by its code number; for example,
|new_param_glue(line_skip_code)| returns a pointer to a glue node for the
current \.[\\lineskip]. } function new_param_glue( n:small_number):halfword ;
var p:halfword ; {the new node}
 q:halfword ; {the glue specification}
begin p:=get_node(small_node_size);  mem[ p].hh.b0 :=glue_node;  mem[ p].hh.b1 :=n+1;
  mem[  p+ 1].hh.rh  :=-{0xfffffff=}268435455  ;

q:=
{ Current |mem| equivalent of glue parameter number |n| } eqtb[  glue_base+   n].hh.rh   

{  };
  mem[  p+ 1].hh.lh  :=q; incr(  mem[  q].hh.rh  );
new_param_glue:=p;
end;



{ 153. }

{tangle:pos tex.web:3201:3: }

{ Glue nodes that are more or less anonymous are created by |new_glue|,
whose argument points to a glue specification. } function new_glue( q:halfword ):halfword ;
var p:halfword ; {the new node}
begin p:=get_node(small_node_size);  mem[ p].hh.b0 :=glue_node;  mem[ p].hh.b1 :=normal;
  mem[  p+ 1].hh.rh  :=-{0xfffffff=}268435455  ;   mem[  p+ 1].hh.lh  :=q; incr(  mem[  q].hh.rh  );
new_glue:=p;
end;



{ 154. }

{tangle:pos tex.web:3211:3: }

{ Still another subroutine is needed: This one is sort of a combination
of |new_param_glue| and |new_glue|. It creates a glue node for one of
the current glue parameters, but it makes a fresh copy of the glue
specification, since that specification will probably be subject to change,
while the parameter will stay put. The global variable |temp_ptr| is
set to the address of the new spec. } function new_skip_param( n:small_number):halfword ;
var p:halfword ; {the new node}
begin temp_ptr:=new_spec(
{ Current |mem| equivalent of glue parameter... } eqtb[  glue_base+   n].hh.rh   

);
p:=new_glue(temp_ptr);   mem[  temp_ptr].hh.rh  :=-{0xfffffff=}268435455  ;  mem[ p].hh.b1 :=n+1;
new_skip_param:=p;
end;



{ 155. }

{tangle:pos tex.web:3225:3: }

{ A |kern_node| has a |width| field to specify a (normally negative)
amount of spacing. This spacing correction appears in horizontal lists
between letters like A and V when the font designer said that it looks
better to move them closer together or further apart. A kern node can
also appear in a vertical list, when its `|width|' denotes additional
spacing in the vertical direction. The |subtype| is either |normal| (for
kerns inserted from font information or math mode calculations) or |explicit|
(for kerns inserted from \.[\\kern] and \.[\\/] commands) or |acc_kern|
(for kerns inserted from non-math accents) or |mu_glue| (for kerns
inserted from \.[\\mkern] specifications in math formulas). }

{ 156. }

{tangle:pos tex.web:3240:3: }

{ The |new_kern| function creates a kern node having a given width. } function new_kern( w:scaled):halfword ;
var p:halfword ; {the new node}
begin p:=get_node(small_node_size);  mem[ p].hh.b0 :=kern_node;
 mem[ p].hh.b1 :=normal;
 mem[ p+width_offset].int  :=w;
new_kern:=p;
end;



{ 157. }

{tangle:pos tex.web:3250:3: }

{ A |penalty_node| specifies the penalty associated with line or page
breaking, in its |penalty| field. This field is a fullword integer, but
the full range of integer values is not used: Any penalty |>=10000| is
treated as infinity, and no break will be allowed for such high values.
Similarly, any penalty |<=-10000| is treated as negative infinity, and a
break will be forced. }

{ 158. }

{tangle:pos tex.web:3262:3: }

{ Anyone who has been reading the last few sections of the program will
be able to guess what comes next. } function new_penalty( m:integer):halfword ;
var p:halfword ; {the new node}
begin p:=get_node(small_node_size);  mem[ p].hh.b0 :=penalty_node;
 mem[ p].hh.b1 :=0; {the |subtype| is not used}
 mem[ p+1].int :=m; new_penalty:=p;
end;



{ 159. }

{tangle:pos tex.web:3272:3: }

{ You might think that we have introduced enough node types by now. Well,
almost, but there is one more: An |unset_node| has nearly the same format
as an |hlist_node| or |vlist_node|; it is used for entries in \.[\\halign]
or \.[\\valign] that are not yet in their final form, since the box
dimensions are their ``natural'' sizes before any glue adjustment has been
made. The |glue_set| word is not present; instead, we have a |glue_stretch|
field, which contains the total stretch of order |glue_order| that is
present in the hlist or vlist being boxed.
Similarly, the |shift_amount| field is replaced by a |glue_shrink| field,
containing the total shrink of order |glue_sign| that is present.
The |subtype| field is called |span_count|; an unset box typically
contains the data for |qo(span_count)+1| columns.
Unset nodes will be changed to box nodes when alignment is completed. }

{ 160. }

{tangle:pos tex.web:3291:3: }

{ In fact, there are still more types coming. When we get to math formula
processing we will see that a |style_node| has |type=14|; and a number
of larger type codes will also be defined, for use in math mode only. }

{ 161. }

{tangle:pos tex.web:3295:3: }

{ Warning: If any changes are made to these data structure layouts, such as
changing any of the node sizes or even reordering the words of nodes,
the |copy_node_list| procedure and the memory initialization code
below may have to be changed. Such potentially dangerous parts of the
program are listed in the index under `data structure assumptions'.
 \xref[data structure assumptions]
However, other references to the nodes are made symbolically in terms of
the \.[WEB] macro definitions above, so that format changes will leave
\TeX's other algorithms intact.
\xref[system dependencies] }

{ 162. \[11] Memory layout }

{tangle:pos tex.web:3306:24: }

{ Some areas of |mem| are dedicated to fixed usage, since static allocation is
more efficient than dynamic allocation when we can get away with it. For
example, locations |mem_bot| to |mem_bot+3| are always used to store the
specification for glue that is `\.[0pt plus 0pt minus 0pt]'. The
following macro definitions accomplish the static allocation by giving
symbolic names to the fixed positions. Static variable-size nodes appear
in locations |mem_bot| through |lo_mem_stat_max|, and static single-word nodes
appear in locations |hi_mem_stat_min| through |mem_top|, inclusive. It is
harmless to let |lig_trick| and |garbage| share the same location of |mem|. }

{ 167. }

{tangle:pos tex.web:3402:3: }

{ Procedure |check_mem| makes sure that the available space lists of
|mem| are well formed, and it optionally prints out all locations
that are reserved now but were free the last time this procedure was called. }  ifdef('TEXMF_DEBUG')  procedure check_mem( print_locs : boolean);
label done1,done2; {loop exits}
var p, q:halfword ; {current locations of interest in |mem|}
 clobbered:boolean; {is something amiss?}
begin for p:=mem_min to lo_mem_max do free_arr [p]:=false; {you can probably
  do this faster}
for p:=hi_mem_min to mem_end do free_arr [p]:=false; {ditto}

{ Check single-word |avail| list }
p:=avail; q:=-{0xfffffff=}268435455  ; clobbered:=false;
while p<>-{0xfffffff=}268435455   do
  begin if (p>mem_end)or(p<hi_mem_min) then clobbered:=true
  else if free_arr [p] then clobbered:=true;
  if clobbered then
    begin print_nl({"AVAIL list clobbered at "=}299);
{ \xref[AVAIL list clobbered...] }
    print_int(q); goto done1;
    end;
  free_arr [p]:=true; q:=p; p:= mem[ q].hh.rh ;
  end;
done1:

;

{ Check variable-size |avail| list }
p:=rover; q:=-{0xfffffff=}268435455  ; clobbered:=false;
repeat if (p>=lo_mem_max)or(p<mem_min) then clobbered:=true
  else if (  mem[  p+ 1].hh.rh  >=lo_mem_max)or(  mem[  p+ 1].hh.rh  <mem_min) then clobbered:=true
  else if  not( ( mem[  p].hh.rh = {0xfffffff=}268435455  ) )or(  mem[ p].hh.lh <2)or 
   (p+  mem[ p].hh.lh >lo_mem_max)or  (  mem[    mem[    p+ 1].hh.rh  + 1].hh.lh  <>p) then clobbered:=true;
  if clobbered then
  begin print_nl({"Double-AVAIL list clobbered at "=}300);
  print_int(q); goto done2;
  end;
for q:=p to p+  mem[ p].hh.lh -1 do {mark all locations free}
  begin if free_arr [q] then
    begin print_nl({"Doubly free location at "=}301);
{ \xref[Doubly free location...] }
    print_int(q); goto done2;
    end;
  free_arr [q]:=true;
  end;
q:=p; p:=  mem[  p+ 1].hh.rh  ;
until p=rover;
done2:

;

{ Check flags of unavailable nodes }
p:=mem_min;
while p<=lo_mem_max do {node |p| should not be empty}
  begin if  ( mem[  p].hh.rh = {0xfffffff=}268435455  )  then
    begin print_nl({"Bad flag at "=}302); print_int(p);
{ \xref[Bad flag...] }
    end;
  while (p<=lo_mem_max) and not free_arr [p] do incr(p);
  while (p<=lo_mem_max) and free_arr [p] do incr(p);
  end

;
if print_locs then 
{ Print newly busy locations }
begin print_nl({"New busy locs:"=}303);
for p:=mem_min to lo_mem_max do
  if not free_arr [p] and ((p>was_lo_max) or was_free[p]) then
    begin print_char({" "=}32); print_int(p);
    end;
for p:=hi_mem_min to mem_end do
  if not free_arr [p] and
   ((p<was_hi_min) or (p>was_mem_end) or was_free[p]) then
    begin print_char({" "=}32); print_int(p);
    end;
end

;
for p:=mem_min to lo_mem_max do was_free[p]:=free_arr [p];
for p:=hi_mem_min to mem_end do was_free[p]:=free_arr [p];
  {|was_free:=free| might be faster}
was_mem_end:=mem_end; was_lo_max:=lo_mem_max; was_hi_min:=hi_mem_min;
end;
endif('TEXMF_DEBUG') 



{ 172. }

{tangle:pos tex.web:3484:3: }

{ The |search_mem| procedure attempts to answer the question ``Who points
to node~|p|?'' In doing so, it fetches |link| and |info| fields of |mem|
that might not be of type |two_halves|. Strictly speaking, this is
\xref[dirty \PASCAL]
undefined in \PASCAL, and it can lead to ``false drops'' (words that seem to
point to |p| purely by coincidence). But for debugging purposes, we want
to rule out the places that do [\sl not\/] point to |p|, so a few false
drops are tolerable. }  ifdef('TEXMF_DEBUG')  procedure search_mem( p:halfword ); {look for pointers to |p|}
var q:integer; {current position being searched}
begin for q:=mem_min to lo_mem_max do
  begin if  mem[ q].hh.rh =p then
    begin print_nl({"LINK("=}304); print_int(q); print_char({")"=}41);
    end;
  if  mem[ q].hh.lh =p then
    begin print_nl({"INFO("=}305); print_int(q); print_char({")"=}41);
    end;
  end;
for q:=hi_mem_min to mem_end do
  begin if  mem[ q].hh.rh =p then
    begin print_nl({"LINK("=}304); print_int(q); print_char({")"=}41);
    end;
  if  mem[ q].hh.lh =p then
    begin print_nl({"INFO("=}305); print_int(q); print_char({")"=}41);
    end;
  end;

{ Search |eqtb| for equivalents equal to |p| }
for q:=active_base to box_base+255 do
  begin if  eqtb[  q].hh.rh  =p then
    begin print_nl({"EQUIV("=}509); print_int(q); print_char({")"=}41);
    end;
  end

;

{ Search |save_stack| for equivalents that point to |p| }
if save_ptr>0 then for q:=0 to save_ptr-1 do
  begin if  save_stack[ q].hh.rh =p then
    begin print_nl({"SAVE("=}554); print_int(q); print_char({")"=}41);
    end;
  end

;

{ Search |hyph_list| for pointers to |p| }
for q:=0 to hyph_size do
  begin if hyph_list[q]=p then
    begin print_nl({"HYPH("=}954); print_int(q); print_char({")"=}41);
    end;
  end

;
end;
endif('TEXMF_DEBUG') 



{ 174. }

{tangle:pos tex.web:3538:3: }

{ Boxes, rules, inserts, whatsits, marks, and things in general that are
sort of ``complicated'' are indicated only by printing `\.[[]]'. } procedure short_display( p:integer); {prints highlights of list |p|}
var n:integer; {for replacement counts}
begin while p>mem_min do
  begin if  ( p>=hi_mem_min)  then
    begin if p<=mem_end then
      begin if   mem[ p].hh.b0 <>font_in_short_display then
        begin if (  mem[ p].hh.b0 >font_max) then
          print_char({"*"=}42)
{ \xref[*\relax] }
        else 
{ Print the font identifier for |font(p)| }
print_esc(  hash[ font_id_base+    mem[   p].hh.b0 ].rh  )

;
        print_char({" "=}32); font_in_short_display:=  mem[ p].hh.b0 ;
        end;
       print (   mem[  p].hh.b1  );
      end;
    end
  else 
{ Print a short indication of the contents of node |p| }
case  mem[ p].hh.b0  of
hlist_node,vlist_node,ins_node,whatsit_node,mark_node,adjust_node,
  unset_node: print({"[]"=}306);
rule_node: print_char({"|"=}124);
glue_node: if   mem[  p+ 1].hh.lh  <>mem_bot  then print_char({" "=}32);
math_node: print_char({"$"=}36);
ligature_node: short_display( mem[    p+1 ].hh.rh  );
disc_node: begin short_display(  mem[  p+ 1].hh.lh  );
  short_display(  mem[  p+ 1].hh.rh  );

  n:= mem[ p].hh.b1 ;
  while n>0 do
    begin if  mem[ p].hh.rh <>-{0xfffffff=}268435455   then p:= mem[ p].hh.rh ;
    decr(n);
    end;
  end;
 else   
 end 

;
  p:= mem[ p].hh.rh ;
  end;
end;



{ 176. }

{tangle:pos tex.web:3580:3: }

{ The |show_node_list| routine requires some auxiliary subroutines: one to
print a font-and-character combination, one to print a token list without
its reference count, and one to print a rule dimension. } procedure print_font_and_char( p:integer); {prints |char_node| data}
begin if p>mem_end then print_esc({"CLOBBERED."=}307)
else  begin if (  mem[ p].hh.b0 >font_max) then print_char({"*"=}42)
{ \xref[*\relax] }
  else 
{ Print the font identifier for |font(p)| }
print_esc(  hash[ font_id_base+    mem[   p].hh.b0 ].rh  )

;
  print_char({" "=}32);  print (   mem[  p].hh.b1  );
  end;
end;


procedure print_mark( p:integer); {prints token list data in braces}
begin print_char({"["=}123);
if (p<hi_mem_min)or(p>mem_end) then print_esc({"CLOBBERED."=}307)
else show_token_list( mem[ p].hh.rh ,-{0xfffffff=}268435455  ,max_print_line-10);
print_char({"]"=}125);
end;


procedure print_rule_dimen( d:scaled); {prints dimension in rule node}
begin if  ( d=-{010000000000=}1073741824 )  then print_char({"*"=}42) else print_scaled(d);
{ \xref[*\relax] }
end;



{ 177. }

{tangle:pos tex.web:3605:3: }

{ Then there is a subroutine that prints glue stretch and shrink, possibly
followed by the name of finite units: } procedure print_glue( d:scaled; order:integer; s:str_number);
  {prints a glue component}
begin print_scaled(d);
if (order<normal)or(order>filll) then print({"foul"=}308)
else if order>normal then
  begin print({"fil"=}309);
  while order>fil do
    begin print_char({"l"=}108); decr(order);
    end;
  end
else if s<>0 then print(s);
end;



{ 178. }

{tangle:pos tex.web:3621:3: }

{ The next subroutine prints a whole glue specification. } procedure print_spec( p:integer; s:str_number);
  {prints a glue specification}
begin if (p<mem_min)or(p>=lo_mem_max) then print_char({"*"=}42)
{ \xref[*\relax] }
else  begin print_scaled( mem[ p+width_offset].int  );
  if s<>0 then print(s);
  if  mem[ p+2].int  <>0 then
    begin print({" plus "=}310); print_glue( mem[ p+2].int  ,  mem[ p].hh.b0 ,s);
    end;
  if  mem[ p+3].int  <>0 then
    begin print({" minus "=}311); print_glue( mem[ p+3].int  ,  mem[ p].hh.b1 ,s);
    end;
  end;
end;



{ 179. }

{tangle:pos tex.web:3638:3: }

{ We also need to declare some procedures that appear later in this
documentation. } 
{ Declare procedures needed for displaying the elements of mlists }
procedure print_fam_and_char( p:halfword ); {prints family and character}
begin print_esc({"fam"=}469); print_int(  mem[ p].hh.b0 ); print_char({" "=}32);
 print (   mem[  p].hh.b1  );
end;


procedure print_delimiter( p:halfword ); {prints a delimiter as 24-bit hex value}
var a:integer; {accumulator}
begin a:=mem[ p].qqqq.b0 *256+ mem[  p].qqqq.b1  ;
a:=a*{0x1000=}4096+mem[ p].qqqq.b2 *256+ mem[  p].qqqq.b3  ;
if a<0 then print_int(a) {this should never happen}
else print_hex(a);
end;


procedure show_info; forward;{ \2 } {|show_node_list(info(temp_ptr))|}
procedure print_subsidiary_data( p:halfword ; c:ASCII_code);
  {display a noad field}
begin if  (pool_ptr - str_start[str_ptr]) >=depth_threshold then
  begin if  mem[ p].hh.rh <>empty then print({" []"=}312);
  end
else  begin  begin str_pool[pool_ptr]:=   c ; incr(pool_ptr); end ; {include |c| in the recursion history}
  temp_ptr:=p; {prepare for |show_info| if recursion is needed}
  case  mem[ p].hh.rh  of
  math_char: begin print_ln; print_current_string; print_fam_and_char(p);
    end;
  sub_box: show_info; {recursive call}
  sub_mlist: if  mem[ p].hh.lh =-{0xfffffff=}268435455   then
      begin print_ln; print_current_string; print({"[]"=}874);
      end
    else show_info; {recursive call}
   else    {|empty|}
   end ;

   decr(pool_ptr) ; {remove |c| from the recursion history}
  end;
end;


procedure print_style( c:integer);
begin case c div 2 of
0: print_esc({"displaystyle"=}875); {|display_style=0|}
1: print_esc({"textstyle"=}876); {|text_style=2|}
2: print_esc({"scriptstyle"=}877); {|script_style=4|}
3: print_esc({"scriptscriptstyle"=}878); {|script_script_style=6|}
 else  print({"Unknown style!"=}879)
 end ;
end;

 

{ Declare the procedure called |print_skip_param| }
procedure print_skip_param( n:integer);
begin case n of
line_skip_code: print_esc({"lineskip"=}381);
baseline_skip_code: print_esc({"baselineskip"=}382);
par_skip_code: print_esc({"parskip"=}383);
above_display_skip_code: print_esc({"abovedisplayskip"=}384);
below_display_skip_code: print_esc({"belowdisplayskip"=}385);
above_display_short_skip_code: print_esc({"abovedisplayshortskip"=}386);
below_display_short_skip_code: print_esc({"belowdisplayshortskip"=}387);
left_skip_code: print_esc({"leftskip"=}388);
right_skip_code: print_esc({"rightskip"=}389);
top_skip_code: print_esc({"topskip"=}390);
split_top_skip_code: print_esc({"splittopskip"=}391);
tab_skip_code: print_esc({"tabskip"=}392);
space_skip_code: print_esc({"spaceskip"=}393);
xspace_skip_code: print_esc({"xspaceskip"=}394);
par_fill_skip_code: print_esc({"parfillskip"=}395);
thin_mu_skip_code: print_esc({"thinmuskip"=}396);
med_mu_skip_code: print_esc({"medmuskip"=}397);
thick_mu_skip_code: print_esc({"thickmuskip"=}398);
 else  print({"[unknown glue parameter!]"=}399)
 end ;
end;





{ 180. }

{tangle:pos tex.web:3644:3: }

{ Since boxes can be inside of boxes, |show_node_list| is inherently recursive,
\xref[recursion]
up to a given maximum number of levels.  The history of nesting is indicated
by the current string, which will be printed at the beginning of each line;
the length of this string, namely |cur_length|, is the depth of nesting.

Recursive calls on |show_node_list| therefore use the following pattern: }

{ 182. }

{tangle:pos tex.web:3667:3: }

{ Now we are ready for |show_node_list| itself. This procedure has been
written to be ``extra robust'' in the sense that it should not crash or get
into a loop even if the data structures have been messed up by bugs in
the rest of the program. You can safely call its parent routine
|show_box(p)| for arbitrary values of |p| when you are debugging \TeX.
However, in the presence of bad data, the procedure may
\xref[dirty \PASCAL]\xref[debugging]
fetch a |memory_word| whose variant is different from the way it was stored;
for example, it might try to read |mem[p].hh| when |mem[p]|
contains a scaled integer, if |p| is a pointer that has been
clobbered or chosen at random. } procedure show_node_list( p:integer); {prints a node list symbolically}
label exit;
var n:integer; {the number of items already printed at this level}
 g:real; {a glue ratio, as a floating point number}
begin if  (pool_ptr - str_start[str_ptr]) >depth_threshold then
  begin if p>-{0xfffffff=}268435455   then print({" []"=}312);
    {indicate that there's been some truncation}
   goto exit ;
  end;
n:=0;
while p>mem_min do
  begin print_ln; print_current_string; {display the nesting history}
  if p>mem_end then {pointer out of range}
    begin print({"Bad link, display aborted."=}313);  goto exit ;
{ \xref[Bad link...] }
    end;
  incr(n); if n>breadth_max then {time to stop}
    begin print({"etc."=}314);  goto exit ;
{ \xref[etc] }
    end;
  
{ Display node |p| }
if  ( p>=hi_mem_min)  then print_font_and_char(p)
else  case  mem[ p].hh.b0  of
  hlist_node,vlist_node,unset_node: 
{ Display box |p| }
begin if  mem[ p].hh.b0 =hlist_node then print_esc({"h"=}104)
else if  mem[ p].hh.b0 =vlist_node then print_esc({"v"=}118)
else print_esc({"unset"=}316);
print({"box("=}317); print_scaled( mem[ p+height_offset].int  ); print_char({"+"=}43);
print_scaled( mem[ p+depth_offset].int  ); print({")x"=}318); print_scaled( mem[ p+width_offset].int  );
if  mem[ p].hh.b0 =unset_node then
  
{ Display special fields of the unset node |p| }
begin if  mem[ p].hh.b1 <>min_quarterword then
  begin print({" ("=}284); print_int(  mem[  p].hh.b1  +1);
  print({" columns)"=}320);
  end;
if mem[ p+glue_offset].int  <>0 then
  begin print({", stretch "=}321); print_glue(mem[ p+glue_offset].int  ,  mem[  p+ list_offset].hh.b1  ,0);
  end;
if  mem[ p+4].int  <>0 then
  begin print({", shrink "=}322); print_glue( mem[ p+4].int  ,  mem[  p+ list_offset].hh.b0  ,0);
  end;
end


else  begin 
{ Display the value of |glue_set(p)| }
g:=   mem[  p+glue_offset].gr  ;
if (g<>  0.0 )and(  mem[  p+ list_offset].hh.b0  <>normal) then
  begin print({", glue set "=}323);
  if   mem[  p+ list_offset].hh.b0  =shrinking then print({"- "=}324);
  { The Unix |pc| folks removed this restriction with a remark that
    invalid bit patterns were vanishingly improbable, so we follow
    their example without really understanding it.
  |if abs(mem[p+glue_offset].int)<@'4000000 then print('?.?')|
  |else| }
  if fabs(g)>  20000.0  then
    begin if g>  0.0  then print_char({">"=}62)
    else print({"< -"=}325);
    print_glue(20000* {0200000=}65536 ,  mem[  p+ list_offset].hh.b1  ,0);
    end
  else print_glue(round( {0200000=}65536 *g),  mem[  p+ list_offset].hh.b1  ,0);
{ \xref[real multiplication] }
  end

;
  if  mem[ p+4].int  <>0 then
    begin print({", shifted "=}319); print_scaled( mem[ p+4].int  );
    end;
  end;
 begin  begin str_pool[pool_ptr]:= {"."=}  46 ; incr(pool_ptr); end ; show_node_list(   mem[   p+ list_offset].hh.rh  );  decr(pool_ptr) ; end ; {recursive call}
end

;
  rule_node: 
{ Display rule |p| }
begin print_esc({"rule("=}326); print_rule_dimen( mem[ p+height_offset].int  ); print_char({"+"=}43);
print_rule_dimen( mem[ p+depth_offset].int  ); print({")x"=}318); print_rule_dimen( mem[ p+width_offset].int  );
end

;
  ins_node: 
{ Display insertion |p| }
begin print_esc({"insert"=}327); print_int(  mem[  p].hh.b1  );
print({", natural size "=}328); print_scaled( mem[ p+height_offset].int  );
print({"; split("=}329); print_spec( mem[  p+ 4].hh.rh  ,0);
print_char({","=}44); print_scaled( mem[ p+depth_offset].int  );
print({"); float cost "=}330); print_int(mem[ p+1].int );
 begin  begin str_pool[pool_ptr]:= {"."=}  46 ; incr(pool_ptr); end ; show_node_list(  mem[   p+ 4].hh.lh  );  decr(pool_ptr) ; end ; {recursive call}
end

;
  whatsit_node: 
{ Display the whatsit node |p| }
case  mem[ p].hh.b1  of
open_node:begin print_write_whatsit({"openout"=}1304,p);
  print_char({"="=}61); print_file_name(  mem[  p+ 1].hh.rh  ,  mem[  p+ 2].hh.lh  ,  mem[  p+ 2].hh.rh  );
  end;
write_node:begin print_write_whatsit({"write"=}601,p);
  print_mark(  mem[  p+ 1].hh.rh  );
  end;
close_node:print_write_whatsit({"closeout"=}1305,p);
special_node:begin print_esc({"special"=}1306);
  print_mark(  mem[  p+ 1].hh.rh  );
  end;
language_node:begin print_esc({"setlanguage"=}1308);
  print_int( mem[  p+ 1].hh.rh  ); print({" (hyphenmin "=}1311);
  print_int( mem[  p+ 1].hh.b0  ); print_char({","=}44);
  print_int( mem[  p+ 1].hh.b1  ); print_char({")"=}41);
  end;
 else  print({"whatsit?"=}1312)
 end 

;
  glue_node: 
{ Display glue |p| }
if  mem[ p].hh.b1 >=a_leaders then 
{ Display leaders |p| }
begin print_esc({""=}335);
if  mem[ p].hh.b1 =c_leaders then print_char({"c"=}99)
else if  mem[ p].hh.b1 =x_leaders then print_char({"x"=}120);
print({"leaders "=}336); print_spec(  mem[  p+ 1].hh.lh  ,0);
 begin  begin str_pool[pool_ptr]:= {"."=}  46 ; incr(pool_ptr); end ; show_node_list(   mem[   p+ 1].hh.rh  );  decr(pool_ptr) ; end ; {recursive call}
end


else  begin print_esc({"glue"=}331);
  if  mem[ p].hh.b1 <>normal then
    begin print_char({"("=}40);
    if  mem[ p].hh.b1 <cond_math_glue then
      print_skip_param( mem[ p].hh.b1 -1)
    else if  mem[ p].hh.b1 =cond_math_glue then print_esc({"nonscript"=}332)
    else print_esc({"mskip"=}333);
    print_char({")"=}41);
    end;
  if  mem[ p].hh.b1 <>cond_math_glue then
    begin print_char({" "=}32);
    if  mem[ p].hh.b1 <cond_math_glue then print_spec(  mem[  p+ 1].hh.lh  ,0)
    else print_spec(  mem[  p+ 1].hh.lh  ,{"mu"=}334);
    end;
  end

;
  kern_node: 
{ Display kern |p| }
if  mem[ p].hh.b1 <>mu_glue then
  begin print_esc({"kern"=}337);
  if  mem[ p].hh.b1 <>normal then print_char({" "=}32);
  print_scaled( mem[ p+width_offset].int  );
  if  mem[ p].hh.b1 =acc_kern then print({" (for accent)"=}338);
{ \xref[for accent] }
  end
else  begin print_esc({"mkern"=}339); print_scaled( mem[ p+width_offset].int  ); print({"mu"=}334);
  end

;
  math_node: 
{ Display math node |p| }
begin print_esc({"math"=}340);
if  mem[ p].hh.b1 =before then print({"on"=}341)
else print({"off"=}342);
if  mem[ p+width_offset].int  <>0 then
  begin print({", surrounded "=}343); print_scaled( mem[ p+width_offset].int  );
  end;
end

;
  ligature_node: 
{ Display ligature |p| }
begin print_font_and_char( p+1 ); print({" (ligature "=}344);
if  mem[ p].hh.b1 >1 then print_char({"|"=}124);
font_in_short_display:=  mem[   p+1 ].hh.b0 ; short_display( mem[    p+1 ].hh.rh  );
if odd( mem[ p].hh.b1 ) then print_char({"|"=}124);
print_char({")"=}41);
end

;
  penalty_node: 
{ Display penalty |p| }
begin print_esc({"penalty "=}345); print_int( mem[ p+1].int );
end

;
  disc_node: 
{ Display discretionary |p| }
begin print_esc({"discretionary"=}346);
if  mem[ p].hh.b1 >0 then
  begin print({" replacing "=}347); print_int( mem[ p].hh.b1 );
  end;
 begin  begin str_pool[pool_ptr]:= {"."=}  46 ; incr(pool_ptr); end ; show_node_list(   mem[   p+ 1].hh.lh  );  decr(pool_ptr) ; end ; {recursive call}
 begin str_pool[pool_ptr]:= {"|"=}  124 ; incr(pool_ptr); end ; show_node_list(  mem[  p+ 1].hh.rh  );  decr(pool_ptr) ; {recursive call}
end

;
  mark_node: 
{ Display mark |p| }
begin print_esc({"mark"=}348); print_mark(mem[ p+1].int );
end

;
  adjust_node: 
{ Display adjustment |p| }
begin print_esc({"vadjust"=}349);  begin  begin str_pool[pool_ptr]:= {"."=}  46 ; incr(pool_ptr); end ; show_node_list( mem[  p+1].int );  decr(pool_ptr) ; end ; {recursive call}
end

;
  { \4 }
{ Cases of |show_node_list| that arise in mlists only }
style_node:print_style( mem[ p].hh.b1 );
choice_node:
{ Display choice node |p| }
begin print_esc({"mathchoice"=}533);
 begin str_pool[pool_ptr]:= {"D"=}  68 ; incr(pool_ptr); end ; show_node_list( mem[  p+ 1].hh.lh  );  decr(pool_ptr) ;
 begin str_pool[pool_ptr]:= {"T"=}  84 ; incr(pool_ptr); end ; show_node_list( mem[  p+ 1].hh.rh  );  decr(pool_ptr) ;
 begin str_pool[pool_ptr]:= {"S"=}  83 ; incr(pool_ptr); end ; show_node_list( mem[  p+ 2].hh.lh  );  decr(pool_ptr) ;
 begin str_pool[pool_ptr]:= {"s"=}  115 ; incr(pool_ptr); end ; show_node_list( mem[  p+ 2].hh.rh  );  decr(pool_ptr) ;
end

;
ord_noad,op_noad,bin_noad,rel_noad,open_noad,close_noad,punct_noad,inner_noad,
  radical_noad,over_noad,under_noad,vcenter_noad,accent_noad,
  left_noad,right_noad:
{ Display normal noad |p| }
begin case  mem[ p].hh.b0  of
ord_noad: print_esc({"mathord"=}880);
op_noad: print_esc({"mathop"=}881);
bin_noad: print_esc({"mathbin"=}882);
rel_noad: print_esc({"mathrel"=}883);
open_noad: print_esc({"mathopen"=}884);
close_noad: print_esc({"mathclose"=}885);
punct_noad: print_esc({"mathpunct"=}886);
inner_noad: print_esc({"mathinner"=}887);
over_noad: print_esc({"overline"=}888);
under_noad: print_esc({"underline"=}889);
vcenter_noad: print_esc({"vcenter"=}547);
radical_noad: begin print_esc({"radical"=}541); print_delimiter( p+4 );
  end;
accent_noad: begin print_esc({"accent"=}516); print_fam_and_char( p+4 );
  end;
left_noad: begin print_esc({"left"=}890); print_delimiter( p+1 );
  end;
right_noad: begin print_esc({"right"=}891); print_delimiter( p+1 );
  end;
end;
if  mem[ p].hh.b1 <>normal then
  if  mem[ p].hh.b1 =limits then print_esc({"limits"=}892)
  else print_esc({"nolimits"=}893);
if  mem[ p].hh.b0 <left_noad then print_subsidiary_data( p+1 ,{"."=}46);
print_subsidiary_data( p+2 ,{"^"=}94);
print_subsidiary_data( p+3 ,{"_"=}95);
end

;
fraction_noad:
{ Display fraction noad |p| }
begin print_esc({"fraction, thickness "=}894);
if  mem[ p+width_offset].int  ={010000000000=}1073741824  then print({"= default"=}895)
else print_scaled( mem[ p+width_offset].int  );
if (mem[   p+4 ].qqqq.b0 <>0)or 
  (mem[   p+4 ].qqqq.b1 <>min_quarterword)or 
  (mem[   p+4 ].qqqq.b2 <>0)or 
  (mem[   p+4 ].qqqq.b3 <>min_quarterword) then
  begin print({", left-delimiter "=}896); print_delimiter( p+4 );
  end;
if (mem[   p+5 ].qqqq.b0 <>0)or 
  (mem[   p+5 ].qqqq.b1 <>min_quarterword)or 
  (mem[   p+5 ].qqqq.b2 <>0)or 
  (mem[   p+5 ].qqqq.b3 <>min_quarterword) then
  begin print({", right-delimiter "=}897); print_delimiter( p+5 );
  end;
print_subsidiary_data( p+2 ,{"\"=}92);
print_subsidiary_data( p+3 ,{"/"=}47);
end

;

 
   else  print({"Unknown node type!"=}315)
   end 

;
  p:= mem[ p].hh.rh ;
  end;
exit:
end;



{ 198. }

{tangle:pos tex.web:3872:3: }

{ The recursive machinery is started by calling |show_box|.
\xref[recursion] } procedure show_box( p:halfword );
begin 
{ Assign the values |depth_threshold:=show_box_depth| and |breadth_max:=show_box_breadth| }
depth_threshold:=eqtb[int_base+ show_box_depth_code].int  ;
breadth_max:=eqtb[int_base+ show_box_breadth_code].int  

;
if breadth_max<=0 then breadth_max:=5;
if pool_ptr+depth_threshold>=pool_size then
  depth_threshold:=pool_size-pool_ptr-1;
  {now there's enough room for prefix string}
show_node_list(p); {the show starts at |p|}
print_ln;
end;



{ 199. \[13] Destroying boxes }

{tangle:pos tex.web:3886:27: }

{ When we are done with a node list, we are obliged to return it to free
storage, including all of its sublists. The recursive procedure
|flush_node_list| does this for us. }

{ 200. }

{tangle:pos tex.web:3891:3: }

{ First, however, we shall consider two non-recursive procedures that do
simpler tasks. The first of these, |delete_token_ref|, is called when
a pointer to a token list's reference count is being removed. This means
that the token list should disappear if the reference count was |null|,
otherwise the count should be decreased by one.
\xref[reference counts] } procedure delete_token_ref( p:halfword ); {|p| points to the reference count
  of a token list that is losing one reference}
begin if   mem[  p].hh.lh  =-{0xfffffff=}268435455   then flush_list(p)
else decr(  mem[  p].hh.lh  );
end;



{ 201. }

{tangle:pos tex.web:3906:3: }

{ Similarly, |delete_glue_ref| is called when a pointer to a glue
specification is being withdrawn.
\xref[reference counts] } procedure delete_glue_ref( p:halfword ); {|p| points to a glue specification}
{  } begin if   mem[   p].hh.rh  =-{0xfffffff=}268435455   then free_node( p,glue_spec_size) else decr(  mem[   p].hh.rh  ); end ;



{ 202. }

{tangle:pos tex.web:3917:3: }

{ Now we are ready to delete any node list, recursively.
In practice, the nodes deleted are usually charnodes (about 2/3 of the time),
and they are glue nodes in about half of the remaining cases.
\xref[recursion] } procedure flush_node_list( p:halfword ); {erase list of nodes starting at |p|}
label done; {go here when node |p| has been freed}
var q:halfword ; {successor to node |p|}
begin while p<>-{0xfffffff=}268435455   do
{ \xref[inner loop] }
  begin q:= mem[ p].hh.rh ;
  if  ( p>=hi_mem_min)  then  begin  mem[  p].hh.rh :=avail; avail:= p; ifdef('STAT')  decr(dyn_used); endif('STAT')  end 
  else  begin case  mem[ p].hh.b0  of
    hlist_node,vlist_node,unset_node: begin flush_node_list(  mem[  p+ list_offset].hh.rh  );
      free_node(p,box_node_size); goto done;
      end;
    rule_node: begin free_node(p,rule_node_size); goto done;
      end;
    ins_node: begin flush_node_list( mem[  p+ 4].hh.lh  );
      delete_glue_ref( mem[  p+ 4].hh.rh  );
      free_node(p,ins_node_size); goto done;
      end;
    whatsit_node: 
{ Wipe out the whatsit node |p| and |goto done| }
begin case  mem[ p].hh.b1  of
open_node: free_node(p,open_node_size);
write_node,special_node: begin delete_token_ref(  mem[  p+ 1].hh.rh  );
  free_node(p,write_node_size); goto done;
  end;
close_node,language_node: free_node(p,small_node_size);
 else  confusion({"ext3"=}1314)
{ \xref[this can't happen ext3][\quad ext3] }
 end ;

goto done;
end

;
    glue_node: begin {  } begin if   mem[     mem[     p+ 1].hh.lh  ].hh.rh  =-{0xfffffff=}268435455   then free_node(   mem[   p+ 1].hh.lh  ,glue_spec_size) else decr(  mem[     mem[     p+ 1].hh.lh  ].hh.rh  ); end ;
      if   mem[  p+ 1].hh.rh  <>-{0xfffffff=}268435455   then flush_node_list(  mem[  p+ 1].hh.rh  );
      end;
    kern_node,math_node,penalty_node:  ;
    ligature_node: flush_node_list( mem[    p+1 ].hh.rh  );
    mark_node: delete_token_ref(mem[ p+1].int );
    disc_node: begin flush_node_list(  mem[  p+ 1].hh.lh  );
      flush_node_list(  mem[  p+ 1].hh.rh  );
      end;
    adjust_node: flush_node_list(mem[ p+1].int );
    { \4 }
{ Cases of |flush_node_list| that arise in mlists only }
style_node: begin free_node(p,style_node_size); goto done;
  end;
choice_node:begin flush_node_list( mem[  p+ 1].hh.lh  );
  flush_node_list( mem[  p+ 1].hh.rh  );
  flush_node_list( mem[  p+ 2].hh.lh  );
  flush_node_list( mem[  p+ 2].hh.rh  );
  free_node(p,style_node_size); goto done;
  end;
ord_noad,op_noad,bin_noad,rel_noad,open_noad,close_noad,punct_noad,inner_noad,
  radical_noad,over_noad,under_noad,vcenter_noad,accent_noad:{  } 

  begin if  mem[   p+1 ].hh.rh >=sub_box then
    flush_node_list( mem[   p+1 ].hh.lh );
  if  mem[   p+2 ].hh.rh >=sub_box then
    flush_node_list( mem[   p+2 ].hh.lh );
  if  mem[   p+3 ].hh.rh >=sub_box then
    flush_node_list( mem[   p+3 ].hh.lh );
  if  mem[ p].hh.b0 =radical_noad then free_node(p,radical_noad_size)
  else if  mem[ p].hh.b0 =accent_noad then free_node(p,accent_noad_size)
  else free_node(p,noad_size);
  goto done;
  end;
left_noad,right_noad: begin free_node(p,noad_size); goto done;
  end;
fraction_noad: begin flush_node_list( mem[   p+2 ].hh.lh );
  flush_node_list( mem[   p+3 ].hh.lh );
  free_node(p,fraction_noad_size); goto done;
  end;

 
     else  confusion({"flushing"=}350)
{ \xref[this can't happen flushing][\quad flushing] }
     end ;

    free_node(p,small_node_size);
    done:end;
  p:=q;
  end;
end;



{ 203. \[14] Copying boxes }

{tangle:pos tex.web:3960:24: }

{ Another recursive operation that acts on boxes is sometimes needed: The
procedure |copy_node_list| returns a pointer to another node list that has
the same structure and meaning as the original. Note that since glue
specifications and token lists have reference counts, we need not make
copies of them. Reference counts can never get too large to fit in a
halfword, since each pointer to a node is in a different memory address,
and the total number of memory addresses fits in a halfword.
\xref[recursion]
\xref[reference counts]

(Well, there actually are also references from outside |mem|; if the
|save_stack| is made arbitrarily large, it would theoretically be possible
to break \TeX\ by overflowing a reference count. But who would want to do that?) }

{ 204. }

{tangle:pos tex.web:3978:3: }

{ The copying procedure copies words en masse without bothering
to look at their individual fields. If the node format changes---for
example, if the size is altered, or if some link field is moved to another
relative position---then this code may need to be changed too.
\xref[data structure assumptions] } function copy_node_list( p:halfword ):halfword ; {makes a duplicate of the
  node list that starts at |p| and returns a pointer to the new list}
var h:halfword ; {temporary head of copied list}
 q:halfword ; {previous position in new list}
 r:halfword ; {current node being fabricated for new list}
 words:0..5; {number of words remaining to be copied}
begin h:=get_avail; q:=h;
while p<>-{0xfffffff=}268435455   do
  begin 
{ Make a copy of node |p| in node |r| }
words:=1; {this setting occurs in more branches than any other}
if  ( p>=hi_mem_min)  then r:=get_avail
else 
{ Case statement to copy different types and set |words| to the number of initial words not yet copied }
case  mem[ p].hh.b0  of
hlist_node,vlist_node,unset_node: begin r:=get_node(box_node_size);
  mem[r+6]:=mem[p+6]; mem[r+5]:=mem[p+5]; {copy the last two words}
    mem[  r+ list_offset].hh.rh  :=copy_node_list(  mem[  p+ list_offset].hh.rh  ); {this affects |mem[r+5]|}
  words:=5;
  end;
rule_node: begin r:=get_node(rule_node_size); words:=rule_node_size;
  end;
ins_node: begin r:=get_node(ins_node_size); mem[r+4]:=mem[p+4];
  incr(  mem[    mem[     p+ 4].hh.rh  ].hh.rh  ) ;
   mem[  r+ 4].hh.lh  :=copy_node_list( mem[  p+ 4].hh.lh  ); {this affects |mem[r+4]|}
  words:=ins_node_size-1;
  end;
whatsit_node:
{ Make a partial copy of the whatsit node |p| and make |r| point to it; set |words| to the number of initial words not yet copied }
case  mem[ p].hh.b1  of
open_node: begin r:=get_node(open_node_size); words:=open_node_size;
  end;
write_node,special_node: begin r:=get_node(write_node_size);
  incr(  mem[     mem[     p+ 1].hh.rh  ].hh.lh  ) ; words:=write_node_size;
  end;
close_node,language_node: begin r:=get_node(small_node_size);
  words:=small_node_size;
  end;
 else  confusion({"ext2"=}1313)
{ \xref[this can't happen ext2][\quad ext2] }
 end 

;
glue_node: begin r:=get_node(small_node_size); incr(  mem[     mem[     p+ 1].hh.lh  ].hh.rh  ) ;
    mem[  r+ 1].hh.lh  :=  mem[  p+ 1].hh.lh  ;   mem[  r+ 1].hh.rh  :=copy_node_list(  mem[  p+ 1].hh.rh  );
  end;
kern_node,math_node,penalty_node: begin r:=get_node(small_node_size);
  words:=small_node_size;
  end;
ligature_node: begin r:=get_node(small_node_size);
  mem[ r+1 ]:=mem[ p+1 ]; {copy |font| and |character|}
   mem[    r+1 ].hh.rh  :=copy_node_list( mem[    p+1 ].hh.rh  );
  end;
disc_node: begin r:=get_node(small_node_size);
    mem[  r+ 1].hh.lh  :=copy_node_list(  mem[  p+ 1].hh.lh  );
    mem[  r+ 1].hh.rh  :=copy_node_list(  mem[  p+ 1].hh.rh  );
  end;
mark_node: begin r:=get_node(small_node_size); incr(  mem[   mem[    p+1].int ].hh.lh  ) ;
  words:=small_node_size;
  end;
adjust_node: begin r:=get_node(small_node_size);
  mem[ r+1].int :=copy_node_list(mem[ p+1].int );
  end; {|words=1=small_node_size-1|}
 else  confusion({"copying"=}351)
{ \xref[this can't happen copying][\quad copying] }
 end 

;
while words>0 do
  begin decr(words); mem[r+words]:=mem[p+words];
  end

;
   mem[ q].hh.rh :=r; q:=r; p:= mem[ p].hh.rh ;
  end;
 mem[ q].hh.rh :=-{0xfffffff=}268435455  ; q:= mem[ h].hh.rh ;  begin  mem[  h].hh.rh :=avail; avail:= h; ifdef('STAT')  decr(dyn_used); endif('STAT')  end ;
copy_node_list:=q;
end;



{ 207. \[15] The command codes }

{tangle:pos tex.web:4048:28: }

{ Before we can go any further, we need to define symbolic names for the internal
code numbers that represent the various commands obeyed by \TeX. These codes
are somewhat arbitrary, but not completely so. For example, the command
codes for character types are fixed by the language, since a user says,
e.g., `\.[\\catcode \`\\\$[] = 3]' to make \.[\char'44] a math delimiter,
and the command code |math_shift| is equal to~3. Some other codes have
been made adjacent so that |case| statements in the program need not consider
cases that are widely spaced, or so that |case| statements can be replaced
by |if| statements.

At any rate, here is the list, for future reference. First come the
``catcode'' commands, several of which share their numeric codes with
ordinary commands when the catcode cannot emerge from \TeX's scanning routine. }

{ 208. }

{tangle:pos tex.web:4090:3: }

{ Next are the ordinary run-of-the-mill command codes.  Codes that are
|min_internal| or more represent internal quantities that might be
expanded by `\.[\\the]'. }

{ 209. }

{tangle:pos tex.web:4154:3: }

{ The next codes are special; they all relate to mode-independent
assignment of values to \TeX's internal registers or tables.
Codes that are |max_internal| or less represent internal quantities
that might be expanded by `\.[\\the]'. }

{ 210. }

{tangle:pos tex.web:4194:3: }

{ The remaining command codes are extra special, since they cannot get through
\TeX's scanner to the main control routine. They have been given values higher
than |max_command| so that their special nature is easily discernible.
The ``expandable'' commands come first. }

{ 211. \[16] The semantic nest }

{tangle:pos tex.web:4220:28: }

{ \TeX\ is typically in the midst of building many lists at once. For example,
when a math formula is being processed, \TeX\ is in math mode and
working on an mlist; this formula has temporarily interrupted \TeX\ from
being in horizontal mode and building the hlist of a paragraph; and this
paragraph has temporarily interrupted \TeX\ from being in vertical mode
and building the vlist for the next page of a document. Similarly, when a
\.[\\vbox] occurs inside of an \.[\\hbox], \TeX\ is temporarily
interrupted from working in restricted horizontal mode, and it enters
internal vertical mode.  The ``semantic nest'' is a stack that
keeps track of what lists and modes are currently suspended.

At each level of processing we are in one of six modes:

\yskip\hang|vmode| stands for vertical mode (the page builder);

\hang|hmode| stands for horizontal mode (the paragraph builder);

\hang|mmode| stands for displayed formula mode;

\hang|-vmode| stands for internal vertical mode (e.g., in a \.[\\vbox]);

\hang|-hmode| stands for restricted horizontal mode (e.g., in an \.[\\hbox]);

\hang|-mmode| stands for math formula mode (not displayed).

\yskip\noindent The mode is temporarily set to zero while processing \.[\\write]
texts.

Numeric values are assigned to |vmode|, |hmode|, and |mmode| so that
\TeX's ``big semantic switch'' can select the appropriate thing to
do by computing the value |abs(mode)+cur_cmd|, where |mode| is the current
mode and |cur_cmd| is the current command code. } procedure print_mode( m:integer); {prints the mode represented by |m|}
begin if m>0 then
  case m div (max_command+1) of
  0:print({"vertical mode"=}352);
  1:print({"horizontal mode"=}353);
  2:print({"display math mode"=}354);
  end
else if m=0 then print({"no mode"=}355)
else  case (-m) div (max_command+1) of
  0:print({"internal vertical mode"=}356);
  1:print({"restricted horizontal mode"=}357);
  2:print({"math mode"=}358);
  end;
end;

procedure print_in_mode( m:integer); {prints the mode represented by |m|}
begin if m>0 then
  case m div (max_command+1) of
  0:print({"' in vertical mode"=}359);
  1:print({"' in horizontal mode"=}360);
  2:print({"' in display math mode"=}361);
  end
else if m=0 then print({"' in no mode"=}362)
else  case (-m) div (max_command+1) of
  0:print({"' in internal vertical mode"=}363);
  1:print({"' in restricted horizontal mode"=}364);
  2:print({"' in math mode"=}365);
  end;
end;



{ 214. }

{tangle:pos tex.web:4348:3: }

{ Here is a common way to make the current list grow: }

{ 216. }

{tangle:pos tex.web:4368:3: }

{ When \TeX's work on one level is interrupted, the state is saved by
calling |push_nest|. This routine changes |head| and |tail| so that
a new (empty) list is begun; it does not change |mode| or |aux|. } procedure push_nest; {enter a new semantic level, save the old}
begin if nest_ptr>max_nest_stack then
  begin max_nest_stack:=nest_ptr;
  if nest_ptr=nest_size then overflow({"semantic nest size"=}366,nest_size);
{ \xref[TeX capacity exceeded semantic nest size][\quad semantic nest size] }
  end;
nest[nest_ptr]:=cur_list; {stack the record}
incr(nest_ptr); cur_list.head_field :=get_avail; cur_list.tail_field :=cur_list.head_field ; cur_list.pg_field :=0; cur_list.ml_field :=line;
end;



{ 217. }

{tangle:pos tex.web:4382:3: }

{ Conversely, when \TeX\ is finished on the current level, the former
state is restored by calling |pop_nest|. This routine will never be
called at the lowest semantic level, nor will it be called unless |head|
is a node that should be returned to free memory. } procedure pop_nest; {leave a semantic level, re-enter the old}
begin  begin  mem[  cur_list.head_field ].hh.rh :=avail; avail:= cur_list.head_field ; ifdef('STAT')  decr(dyn_used); endif('STAT')  end ; decr(nest_ptr); cur_list:=nest[nest_ptr];
end;



{ 218. }

{tangle:pos tex.web:4391:3: }

{ Here is a procedure that displays what \TeX\ is working on, at all levels. } procedure print_totals; forward;{ \2 }
procedure show_activities;
var p:0..nest_size; {index into |nest|}
 m:-mmode..mmode; {mode}
 a:memory_word; {auxiliary}
 q, r:halfword ; {for showing the current page}
 t:integer; {ditto}
begin nest[nest_ptr]:=cur_list; {put the top level into the array}
print_nl({""=}335); print_ln;
for p:=nest_ptr downto 0 do
  begin m:=nest[p].mode_field; a:=nest[p].aux_field;
  print_nl({"### "=}367); print_mode(m);
  print({" entered at line "=}368); print_int(abs(nest[p].ml_field));
  if m=hmode then if nest[p].pg_field <> {040600000=}8585216 then
    begin print({" (language"=}369); print_int(nest[p].pg_field mod {0200000=}65536);
    print({":hyphenmin"=}370); print_int(nest[p].pg_field div {020000000=}4194304);
    print_char({","=}44); print_int((nest[p].pg_field div {0200000=}65536) mod {0100=}64);
    print_char({")"=}41);
    end;
  if nest[p].ml_field<0 then print({" (\output routine)"=}371);
  if p=0 then
    begin 
{ Show the status of the current page }
if mem_top-2 <>page_tail then
  begin print_nl({"### current page:"=}994);
  if output_active then print({" (held over for next output)"=}995);
{ \xref[held over for next output] }
  show_box( mem[ mem_top-2 ].hh.rh );
  if page_contents>empty then
    begin print_nl({"total height "=}996); print_totals;
{ \xref[total_height][\.[total height]] }
    print_nl({" goal height "=}997); print_scaled(page_so_far[0] );
{ \xref[goal height] }
    r:= mem[ mem_top ].hh.rh ;
    while r<>mem_top  do
      begin print_ln; print_esc({"insert"=}327); t:=  mem[  r].hh.b1  ;
      print_int(t); print({" adds "=}998);
      if eqtb[count_base+ t].int =1000 then t:= mem[ r+height_offset].int  
      else t:=x_over_n( mem[ r+height_offset].int  ,1000)*eqtb[count_base+ t].int ;
      print_scaled(t);
      if  mem[ r].hh.b0 =split_up then
        begin q:=mem_top-2 ; t:=0;
        repeat q:= mem[ q].hh.rh ;
        if ( mem[ q].hh.b0 =ins_node)and( mem[ q].hh.b1 = mem[ r].hh.b1 ) then incr(t);
        until q= mem[  r+ 1].hh.lh  ;
        print({", #"=}999); print_int(t); print({" might split"=}1000);
        end;
      r:= mem[ r].hh.rh ;
      end;
    end;
  end

;
    if  mem[ mem_top-1 ].hh.rh <>-{0xfffffff=}268435455   then
      print_nl({"### recent contributions:"=}372);
    end;
  show_box( mem[ nest[ p]. head_field].hh.rh );
  
{ Show the auxiliary field, |a| }
case abs(m) div (max_command+1) of
0: begin print_nl({"prevdepth "=}373);
  if a.int <=-65536000  then print({"ignored"=}374)
  else print_scaled(a.int );
  if nest[p].pg_field<>0 then
    begin print({", prevgraf "=}375);
    print_int(nest[p].pg_field);
    if nest[p].pg_field<>1 then print({" lines"=}376)
    else print({" line"=}377);
    end;
  end;
1: begin print_nl({"spacefactor "=}378); print_int(a.hh.lh);
  if m>0 then  if a.hh.rh>0 then
    begin print({", current language "=}379); print_int(a.hh.rh); 
    end;
  end;
2: if a.int<>-{0xfffffff=}268435455   then
  begin print({"this will begin denominator of:"=}380); show_box(a.int); 
  end;
end {there are no other cases}

;
  end;
end;



{ 220. \[17] The table of equivalents }

{tangle:pos tex.web:4444:35: }

{ Now that we have studied the data structures for \TeX's semantic routines,
we ought to consider the data structures used by its syntactic routines. In
other words, our next concern will be
the tables that \TeX\ looks at when it is scanning
what the user has written.

The biggest and most important such table is called |eqtb|. It holds the
current ``equivalents'' of things; i.e., it explains what things mean
or what their current values are, for all quantities that are subject to
the nesting structure provided by \TeX's grouping mechanism. There are six
parts to |eqtb|:

\yskip\hangg 1) |eqtb[active_base..(hash_base-1)]| holds the current
equivalents of single-character control sequences.

\yskip\hangg 2) |eqtb[hash_base..(glue_base-1)]| holds the current
equivalents of multiletter control sequences.

\yskip\hangg 3) |eqtb[glue_base..(local_base-1)]| holds the current
equivalents of glue parameters like the current baselineskip.

\yskip\hangg 4) |eqtb[local_base..(int_base-1)]| holds the current
equivalents of local halfword quantities like the current box registers,
the current ``catcodes,'' the current font, and a pointer to the current
paragraph shape.
Additionally region~4 contains the table with ML\TeX's character
substitution definitions.

\yskip\hangg 5) |eqtb[int_base..(dimen_base-1)]| holds the current
equivalents of fullword integer parameters like the current hyphenation
penalty.

\yskip\hangg 6) |eqtb[dimen_base..eqtb_size]| holds the current equivalents
of fullword dimension parameters like the current hsize or amount of
hanging indentation.

\yskip\noindent Note that, for example, the current amount of
baselineskip glue is determined by the setting of a particular location
in region~3 of |eqtb|, while the current meaning of the control sequence
`\.[\\baselineskip]' (which might have been changed by \.[\\def] or
\.[\\let]) appears in region~2. }

{ 221. }

{tangle:pos tex.web:4485:3: }

{ Each entry in |eqtb| is a |memory_word|. Most of these words are of type
|two_halves|, and subdivided into three fields:

\yskip\hangg 1) The |eq_level| (a quarterword) is the level of grouping at
which this equivalent was defined. If the level is |level_zero|, the
equivalent has never been defined; |level_one| refers to the outer level
(outside of all groups), and this level is also used for global
definitions that never go away. Higher levels are for equivalents that
will disappear at the end of their group.  \xref[global definitions]

\yskip\hangg 2) The |eq_type| (another quarterword) specifies what kind of
entry this is. There are many types, since each \TeX\ primitive like
\.[\\hbox], \.[\\def], etc., has its own special code. The list of
command codes above includes all possible settings of the |eq_type| field.

\yskip\hangg 3) The |equiv| (a halfword) is the current equivalent value.
This may be a font number, a pointer into |mem|, or a variety of other
things. }

{ 237. }

{tangle:pos tex.web:5043:3: }

{ We can print the symbolic name of an integer parameter as follows. } procedure print_param( n:integer);
begin case n of
pretolerance_code:print_esc({"pretolerance"=}425);
tolerance_code:print_esc({"tolerance"=}426);
line_penalty_code:print_esc({"linepenalty"=}427);
hyphen_penalty_code:print_esc({"hyphenpenalty"=}428);
ex_hyphen_penalty_code:print_esc({"exhyphenpenalty"=}429);
club_penalty_code:print_esc({"clubpenalty"=}430);
widow_penalty_code:print_esc({"widowpenalty"=}431);
display_widow_penalty_code:print_esc({"displaywidowpenalty"=}432);
broken_penalty_code:print_esc({"brokenpenalty"=}433);
bin_op_penalty_code:print_esc({"binoppenalty"=}434);
rel_penalty_code:print_esc({"relpenalty"=}435);
pre_display_penalty_code:print_esc({"predisplaypenalty"=}436);
post_display_penalty_code:print_esc({"postdisplaypenalty"=}437);
inter_line_penalty_code:print_esc({"interlinepenalty"=}438);
double_hyphen_demerits_code:print_esc({"doublehyphendemerits"=}439);
final_hyphen_demerits_code:print_esc({"finalhyphendemerits"=}440);
adj_demerits_code:print_esc({"adjdemerits"=}441);
mag_code:print_esc({"mag"=}442);
delimiter_factor_code:print_esc({"delimiterfactor"=}443);
looseness_code:print_esc({"looseness"=}444);
time_code:print_esc({"time"=}445);
day_code:print_esc({"day"=}446);
month_code:print_esc({"month"=}447);
year_code:print_esc({"year"=}448);
show_box_breadth_code:print_esc({"showboxbreadth"=}449);
show_box_depth_code:print_esc({"showboxdepth"=}450);
hbadness_code:print_esc({"hbadness"=}451);
vbadness_code:print_esc({"vbadness"=}452);
pausing_code:print_esc({"pausing"=}453);
tracing_online_code:print_esc({"tracingonline"=}454);
tracing_macros_code:print_esc({"tracingmacros"=}455);
tracing_stats_code:print_esc({"tracingstats"=}456);
tracing_paragraphs_code:print_esc({"tracingparagraphs"=}457);
tracing_pages_code:print_esc({"tracingpages"=}458);
tracing_output_code:print_esc({"tracingoutput"=}459);
tracing_lost_chars_code:print_esc({"tracinglostchars"=}460);
tracing_commands_code:print_esc({"tracingcommands"=}461);
tracing_restores_code:print_esc({"tracingrestores"=}462);
uc_hyph_code:print_esc({"uchyph"=}463);
output_penalty_code:print_esc({"outputpenalty"=}464);
max_dead_cycles_code:print_esc({"maxdeadcycles"=}465);
hang_after_code:print_esc({"hangafter"=}466);
floating_penalty_code:print_esc({"floatingpenalty"=}467);
global_defs_code:print_esc({"globaldefs"=}468);
cur_fam_code:print_esc({"fam"=}469);
escape_char_code:print_esc({"escapechar"=}470);
default_hyphen_char_code:print_esc({"defaulthyphenchar"=}471);
default_skew_char_code:print_esc({"defaultskewchar"=}472);
end_line_char_code:print_esc({"endlinechar"=}473);
new_line_char_code:print_esc({"newlinechar"=}474);
language_code:print_esc({"language"=}475);
left_hyphen_min_code:print_esc({"lefthyphenmin"=}476);
right_hyphen_min_code:print_esc({"righthyphenmin"=}477);
holding_inserts_code:print_esc({"holdinginserts"=}478);
error_context_lines_code:print_esc({"errorcontextlines"=}479);
char_sub_def_min_code:print_esc({"charsubdefmin"=}480);
char_sub_def_max_code:print_esc({"charsubdefmax"=}481);
tracing_char_sub_def_code:print_esc({"tracingcharsubdef"=}482);
 else  print({"[unknown integer parameter!]"=}483)
 end ;
end;



{ 241. }

{tangle:pos tex.ch:1292:3: }

{ The following procedure, which is called just before \TeX\ initializes
its input and output, establishes the initial values of the date and
time. It calls a |date_and_time| C macro (a.k.a.\ |dateandtime|), which
calls the C function |get_date_and_time|, passing it the addresses of
|sys_time|, etc., so they can be set by the routine. |get_date_and_time|
also sets up interrupt catching if that is conditionally compiled in the
C code.

We have to initialize the |sys_| variables because that is what gets
output on the first line of the log file. (New in 2021.)
\xref[system dependencies] } procedure fix_date_and_time;
begin date_and_time(sys_time,sys_day,sys_month,sys_year);
eqtb[int_base+ time_code].int  :=sys_time; {minutes since midnight}
eqtb[int_base+ day_code].int  :=sys_day; {day of the month}
eqtb[int_base+ month_code].int  :=sys_month; {month of the year}
eqtb[int_base+ year_code].int  :=sys_year; {Anno Domini}
end;



{ 245. }

{tangle:pos tex.web:5272:3: }

{ \TeX\ is occasionally supposed to print diagnostic information that
goes only into the transcript file, unless |tracing_online| is positive.
Here are two routines that adjust the destination of print commands: } procedure begin_diagnostic; {prepare to do some tracing}
begin old_setting:=selector;
if (eqtb[int_base+ tracing_online_code].int  <=0)and(selector=term_and_log) then
  begin decr(selector);
  if history=spotless then history:=warning_issued;
  end;
end;


procedure end_diagnostic( blank_line:boolean);
  {restore proper conditions after tracing}
begin print_nl({""=}335);
if blank_line then print_ln;
selector:=old_setting;
end;



{ 247. }

{tangle:pos tex.web:5299:3: }

{ The final region of |eqtb| contains the dimension parameters defined
here, and the 256 \.[\\dimen] registers. } procedure print_length_param( n:integer);
begin case n of
par_indent_code:print_esc({"parindent"=}486);
math_surround_code:print_esc({"mathsurround"=}487);
line_skip_limit_code:print_esc({"lineskiplimit"=}488);
hsize_code:print_esc({"hsize"=}489);
vsize_code:print_esc({"vsize"=}490);
max_depth_code:print_esc({"maxdepth"=}491);
split_max_depth_code:print_esc({"splitmaxdepth"=}492);
box_max_depth_code:print_esc({"boxmaxdepth"=}493);
hfuzz_code:print_esc({"hfuzz"=}494);
vfuzz_code:print_esc({"vfuzz"=}495);
delimiter_shortfall_code:print_esc({"delimitershortfall"=}496);
null_delimiter_space_code:print_esc({"nulldelimiterspace"=}497);
script_space_code:print_esc({"scriptspace"=}498);
pre_display_size_code:print_esc({"predisplaysize"=}499);
display_width_code:print_esc({"displaywidth"=}500);
display_indent_code:print_esc({"displayindent"=}501);
overfull_rule_code:print_esc({"overfullrule"=}502);
hang_indent_code:print_esc({"hangindent"=}503);
h_offset_code:print_esc({"hoffset"=}504);
v_offset_code:print_esc({"voffset"=}505);
emergency_stretch_code:print_esc({"emergencystretch"=}506);
 else  print({"[unknown dimen parameter!]"=}507)
 end ;
end;



{ 252. }

{tangle:pos tex.web:5441:3: }

{ Here is a procedure that displays the contents of |eqtb[n]|
symbolically. }{ \4 }
{ Declare the procedure called |print_cmd_chr| }
procedure print_cmd_chr( cmd:quarterword; chr_code:halfword);
begin case cmd of
left_brace: begin print({"begin-group character "=} 565);  print (chr_code); end ;
right_brace: begin print({"end-group character "=} 566);  print (chr_code); end ;
math_shift: begin print({"math shift character "=} 567);  print (chr_code); end ;
mac_param: begin print({"macro parameter character "=} 568);  print (chr_code); end ;
sup_mark: begin print({"superscript character "=} 569);  print (chr_code); end ;
sub_mark: begin print({"subscript character "=} 570);  print (chr_code); end ;
endv: print({"end of alignment template"=}571);
spacer: begin print({"blank space "=} 572);  print (chr_code); end ;
letter: begin print({"the letter "=} 573);  print (chr_code); end ;
other_char: begin print({"the character "=} 574);  print (chr_code); end ;
{ \4 }
{ Cases of |print_cmd_chr| for symbolic printing of primitives }
assign_glue,assign_mu_glue: if chr_code<skip_base then
    print_skip_param(chr_code-glue_base)
  else if chr_code<mu_skip_base then
    begin print_esc({"skip"=}400); print_int(chr_code-skip_base);
    end
  else  begin print_esc({"muskip"=}401); print_int(chr_code-mu_skip_base);
    end;


assign_toks: if chr_code>=toks_base then
  begin print_esc({"toks"=}412); print_int(chr_code-toks_base);
  end
else  case chr_code of
  output_routine_loc: print_esc({"output"=}403);
  every_par_loc: print_esc({"everypar"=}404);
  every_math_loc: print_esc({"everymath"=}405);
  every_display_loc: print_esc({"everydisplay"=}406);
  every_hbox_loc: print_esc({"everyhbox"=}407);
  every_vbox_loc: print_esc({"everyvbox"=}408);
  every_job_loc: print_esc({"everyjob"=}409);
  every_cr_loc: print_esc({"everycr"=}410);
   else  print_esc({"errhelp"=}411)
   end ;


assign_int: if chr_code<count_base then print_param(chr_code-int_base)
  else  begin print_esc({"count"=}484); print_int(chr_code-count_base);
    end;


assign_dimen: if chr_code<scaled_base then
    print_length_param(chr_code-dimen_base)
  else  begin print_esc({"dimen"=}508); print_int(chr_code-scaled_base);
    end;


accent: print_esc({"accent"=}516);
advance: print_esc({"advance"=}517);
after_assignment: print_esc({"afterassignment"=}518);
after_group: print_esc({"aftergroup"=}519);
assign_font_dimen: print_esc({"fontdimen"=}527);
begin_group: print_esc({"begingroup"=}520);
break_penalty: print_esc({"penalty"=}539);
char_num: print_esc({"char"=}521);
cs_name: print_esc({"csname"=}512);
def_font: print_esc({"font"=}526);
delim_num: print_esc({"delimiter"=}522);
divide: print_esc({"divide"=}523);
end_cs_name: print_esc({"endcsname"=}513);
end_group: print_esc({"endgroup"=}524);
ex_space: print_esc({" "=}32);
expand_after: print_esc({"expandafter"=}525);
halign: print_esc({"halign"=}528);
hrule: print_esc({"hrule"=}529);
ignore_spaces: print_esc({"ignorespaces"=}530);
insert: print_esc({"insert"=}327);
ital_corr: print_esc({"/"=}47);
mark: print_esc({"mark"=}348);
math_accent: print_esc({"mathaccent"=}531);
math_char_num: print_esc({"mathchar"=}532);
math_choice: print_esc({"mathchoice"=}533);
multiply: print_esc({"multiply"=}534);
no_align: print_esc({"noalign"=}535);
no_boundary:print_esc({"noboundary"=}536);
no_expand: print_esc({"noexpand"=}537);
non_script: print_esc({"nonscript"=}332);
omit: print_esc({"omit"=}538);
radical: print_esc({"radical"=}541);
read_to_cs: print_esc({"read"=}542);
relax: print_esc({"relax"=}543);
set_box: print_esc({"setbox"=}544);
set_prev_graf: print_esc({"prevgraf"=}540);
set_shape: print_esc({"parshape"=}413);
the: print_esc({"the"=}545);
toks_register: print_esc({"toks"=}412);
vadjust: print_esc({"vadjust"=}349);
valign: print_esc({"valign"=}546);
vcenter: print_esc({"vcenter"=}547);
vrule: print_esc({"vrule"=}548);


par_end:print_esc({"par"=}604);


input: if chr_code=0 then print_esc({"input"=}639) else print_esc({"endinput"=}640);


top_bot_mark: case chr_code of
  first_mark_code: print_esc({"firstmark"=}642);
  bot_mark_code: print_esc({"botmark"=}643);
  split_first_mark_code: print_esc({"splitfirstmark"=}644);
  split_bot_mark_code: print_esc({"splitbotmark"=}645);
   else  print_esc({"topmark"=}641)
   end ;


register: if chr_code=int_val then print_esc({"count"=}484)
  else if chr_code=dimen_val then print_esc({"dimen"=}508)
  else if chr_code=glue_val then print_esc({"skip"=}400)
  else print_esc({"muskip"=}401);


set_aux: if chr_code=vmode then print_esc({"prevdepth"=}679)
 else print_esc({"spacefactor"=}678);
set_page_int: if chr_code=0 then print_esc({"deadcycles"=}680)
 else print_esc({"insertpenalties"=}681);
set_box_dimen: if chr_code=width_offset then print_esc({"wd"=}682)
else if chr_code=height_offset then print_esc({"ht"=}683)
else print_esc({"dp"=}684);
last_item: case chr_code of
  int_val: print_esc({"lastpenalty"=}685);
  dimen_val: print_esc({"lastkern"=}686);
  glue_val: print_esc({"lastskip"=}687);
  input_line_no_code: print_esc({"inputlineno"=}688);
   else  print_esc({"badness"=}689)
   end ;


convert: case chr_code of
  number_code: print_esc({"number"=}745);
  roman_numeral_code: print_esc({"romannumeral"=}746);
  string_code: print_esc({"string"=}747);
  meaning_code: print_esc({"meaning"=}748);
  font_name_code: print_esc({"fontname"=}749);
   else  print_esc({"jobname"=}750)
   end ;


if_test: case chr_code of
  if_cat_code:print_esc({"ifcat"=}768);
  if_int_code:print_esc({"ifnum"=}769);
  if_dim_code:print_esc({"ifdim"=}770);
  if_odd_code:print_esc({"ifodd"=}771);
  if_vmode_code:print_esc({"ifvmode"=}772);
  if_hmode_code:print_esc({"ifhmode"=}773);
  if_mmode_code:print_esc({"ifmmode"=}774);
  if_inner_code:print_esc({"ifinner"=}775);
  if_void_code:print_esc({"ifvoid"=}776);
  if_hbox_code:print_esc({"ifhbox"=}777);
  if_vbox_code:print_esc({"ifvbox"=}778);
  ifx_code:print_esc({"ifx"=}779);
  if_eof_code:print_esc({"ifeof"=}780);
  if_true_code:print_esc({"iftrue"=}781);
  if_false_code:print_esc({"iffalse"=}782);
  if_case_code:print_esc({"ifcase"=}783);
   else  print_esc({"if"=}767)
   end ;


fi_or_else: if chr_code=fi_code then print_esc({"fi"=}784)
  else if chr_code=or_code then print_esc({"or"=}785)
  else print_esc({"else"=}786);


tab_mark: if chr_code=span_code then print_esc({"span"=}912)
  else begin print({"alignment tab character "=} 916);  print (chr_code); end ;
car_ret: if chr_code=cr_code then print_esc({"cr"=}913)
  else print_esc({"crcr"=}914);


set_page_dimen: case chr_code of
0: print_esc({"pagegoal"=}984);
1: print_esc({"pagetotal"=}985);
2: print_esc({"pagestretch"=}986);
3: print_esc({"pagefilstretch"=}987);
4: print_esc({"pagefillstretch"=}988);
5: print_esc({"pagefilllstretch"=}989);
6: print_esc({"pageshrink"=}990);
 else  print_esc({"pagedepth"=}991)
 end ;


stop:if chr_code=1 then print_esc({"dump"=}1039) else print_esc({"end"=}1038);


hskip: case chr_code of
  skip_code:print_esc({"hskip"=}1040);
  fil_code:print_esc({"hfil"=}1041);
  fill_code:print_esc({"hfill"=}1042);
  ss_code:print_esc({"hss"=}1043);
   else  print_esc({"hfilneg"=}1044)
   end ;
vskip: case chr_code of
  skip_code:print_esc({"vskip"=}1045);
  fil_code:print_esc({"vfil"=}1046);
  fill_code:print_esc({"vfill"=}1047);
  ss_code:print_esc({"vss"=}1048);
   else  print_esc({"vfilneg"=}1049)
   end ;
mskip: print_esc({"mskip"=}333);
kern: print_esc({"kern"=}337);
mkern: print_esc({"mkern"=}339);


hmove: if chr_code=1 then print_esc({"moveleft"=}1067) else print_esc({"moveright"=}1068);
vmove: if chr_code=1 then print_esc({"raise"=}1069) else print_esc({"lower"=}1070);
make_box: case chr_code of
  box_code: print_esc({"box"=}414);
  copy_code: print_esc({"copy"=}1071);
  last_box_code: print_esc({"lastbox"=}1072);
  vsplit_code: print_esc({"vsplit"=}979);
  vtop_code: print_esc({"vtop"=}1073);
  vtop_code+vmode: print_esc({"vbox"=}981);
   else  print_esc({"hbox"=}1074)
   end ;
leader_ship: if chr_code=a_leaders then print_esc({"leaders"=}1076)
  else if chr_code=c_leaders then print_esc({"cleaders"=}1077)
  else if chr_code=x_leaders then print_esc({"xleaders"=}1078)
  else print_esc({"shipout"=}1075);


start_par: if chr_code=0 then print_esc({"noindent"=}1094)  else print_esc({"indent"=}1093);


remove_item: if chr_code=glue_node then print_esc({"unskip"=}1105)
  else if chr_code=kern_node then print_esc({"unkern"=}1104)
  else print_esc({"unpenalty"=}1103);
un_hbox: if chr_code=copy_code then print_esc({"unhcopy"=}1107)
  else print_esc({"unhbox"=}1106);
un_vbox: if chr_code=copy_code then print_esc({"unvcopy"=}1109)
  else print_esc({"unvbox"=}1108);


discretionary: if chr_code=1 then
  print_esc({"-"=}45) else print_esc({"discretionary"=}346);


eq_no:if chr_code=1 then print_esc({"leqno"=}1141) else print_esc({"eqno"=}1140);


math_comp: case chr_code of
  ord_noad: print_esc({"mathord"=}880);
  op_noad: print_esc({"mathop"=}881);
  bin_noad: print_esc({"mathbin"=}882);
  rel_noad: print_esc({"mathrel"=}883);
  open_noad: print_esc({"mathopen"=}884);
  close_noad: print_esc({"mathclose"=}885);
  punct_noad: print_esc({"mathpunct"=}886);
  inner_noad: print_esc({"mathinner"=}887);
  under_noad: print_esc({"underline"=}889);
   else  print_esc({"overline"=}888)
   end ;
limit_switch: if chr_code=limits then print_esc({"limits"=}892)
  else if chr_code=no_limits then print_esc({"nolimits"=}893)
  else print_esc({"displaylimits"=}1142);


math_style: print_style(chr_code);


above: case chr_code of
  over_code:print_esc({"over"=}1161);
  atop_code:print_esc({"atop"=}1162);
  delimited_code+above_code:print_esc({"abovewithdelims"=}1163);
  delimited_code+over_code:print_esc({"overwithdelims"=}1164);
  delimited_code+atop_code:print_esc({"atopwithdelims"=}1165);
   else  print_esc({"above"=}1160)
   end ;


left_right: if chr_code=left_noad then print_esc({"left"=}890)
else print_esc({"right"=}891);


prefix: if chr_code=1 then print_esc({"long"=}1184)
  else if chr_code=2 then print_esc({"outer"=}1185)
  else print_esc({"global"=}1186);
def: if chr_code=0 then print_esc({"def"=}1187)
  else if chr_code=1 then print_esc({"gdef"=}1188)
  else if chr_code=2 then print_esc({"edef"=}1189)
  else print_esc({"xdef"=}1190);


let: if chr_code<>normal then print_esc({"futurelet"=}1205) else print_esc({"let"=}1204);


shorthand_def: case chr_code of
  char_def_code: print_esc({"chardef"=}1206);
  math_char_def_code: print_esc({"mathchardef"=}1207);
  count_def_code: print_esc({"countdef"=}1208);
  dimen_def_code: print_esc({"dimendef"=}1209);
  skip_def_code: print_esc({"skipdef"=}1210);
  mu_skip_def_code: print_esc({"muskipdef"=}1211);
  char_sub_def_code: print_esc({"charsubdef"=}1213);
   else  print_esc({"toksdef"=}1212)
   end ;
char_given: begin print_esc({"char"=}521); print_hex(chr_code);
  end;
math_given: begin print_esc({"mathchar"=}532); print_hex(chr_code);
  end;


def_code: if chr_code=cat_code_base then print_esc({"catcode"=}420)
  else if chr_code=math_code_base then print_esc({"mathcode"=}424)
  else if chr_code=lc_code_base then print_esc({"lccode"=}421)
  else if chr_code=uc_code_base then print_esc({"uccode"=}422)
  else if chr_code=sf_code_base then print_esc({"sfcode"=}423)
  else print_esc({"delcode"=}485);
def_family: print_size(chr_code-math_font_base);


hyph_data: if chr_code=1 then print_esc({"patterns"=}967)
  else print_esc({"hyphenation"=}955);


assign_font_int: if chr_code=0 then print_esc({"hyphenchar"=}1233)
  else print_esc({"skewchar"=}1234);


set_font:begin print({"select font "=}1242); slow_print(font_name[chr_code]);
  if font_size[chr_code]<>font_dsize[chr_code] then
    begin print({" at "=}751); print_scaled(font_size[chr_code]);
    print({"pt"=}402);
    end;
  end;


set_interaction: case chr_code of
  batch_mode: print_esc({"batchmode"=}272);
  nonstop_mode: print_esc({"nonstopmode"=}273);
  scroll_mode: print_esc({"scrollmode"=}274);
   else  print_esc({"errorstopmode"=}1243)
   end ;


in_stream: if chr_code=0 then print_esc({"closein"=}1245)
  else print_esc({"openin"=}1244);


message: if chr_code=0 then print_esc({"message"=}1246)
  else print_esc({"errmessage"=}1247);


case_shift:if chr_code=lc_code_base then print_esc({"lowercase"=}1253)
  else print_esc({"uppercase"=}1254);


xray: case chr_code of
  show_box_code:print_esc({"showbox"=}1256);
  show_the_code:print_esc({"showthe"=}1257);
  show_lists_code:print_esc({"showlists"=}1258);
   else  print_esc({"show"=}1255)
   end ;


undefined_cs: print({"undefined"=}1265);
call: print({"macro"=}1266);
long_call: print_esc({"long macro"=}1267);
outer_call: print_esc({"outer macro"=}1268);
long_outer_call: begin print_esc({"long"=}1184); print_esc({"outer macro"=}1268);
  end;
end_template: print_esc({"outer endtemplate"=}1269);


extension: case chr_code of
  open_node:print_esc({"openout"=}1304);
  write_node:print_esc({"write"=}601);
  close_node:print_esc({"closeout"=}1305);
  special_node:print_esc({"special"=}1306);
  immediate_code:print_esc({"immediate"=}1307);
  set_language_code:print_esc({"setlanguage"=}1308);
   else  print({"[unknown extension!]"=}1309)
   end ;



 else  print({"[unknown command code!]"=}575)
 end ;
end;

 

 ifdef('STAT')  procedure show_eqtb( n:halfword );
begin if n<active_base then print_char({"?"=}63) {this can't happen}
else if (n<glue_base) or ((n>eqtb_size)and(n<=eqtb_top)) then
  
{ Show equivalent |n|, in region 1 or 2 }
begin sprint_cs(n); print_char({"="=}61); print_cmd_chr( eqtb[  n].hh.b0  , eqtb[  n].hh.rh  );
if  eqtb[  n].hh.b0  >=call then
  begin print_char({":"=}58); show_token_list( mem[  eqtb[   n].hh.rh  ].hh.rh ,-{0xfffffff=}268435455  ,32);
  end;
end


else if n<local_base then 
{ Show equivalent |n|, in region 3 }
if n<skip_base then
  begin print_skip_param(n-glue_base); print_char({"="=}61);
  if n<glue_base+thin_mu_skip_code then print_spec( eqtb[  n].hh.rh  ,{"pt"=}402)
  else print_spec( eqtb[  n].hh.rh  ,{"mu"=}334);
  end
else if n<mu_skip_base then
  begin print_esc({"skip"=}400); print_int(n-skip_base); print_char({"="=}61);
  print_spec( eqtb[  n].hh.rh  ,{"pt"=}402);
  end
else  begin print_esc({"muskip"=}401); print_int(n-mu_skip_base); print_char({"="=}61);
  print_spec( eqtb[  n].hh.rh  ,{"mu"=}334);
  end


else if n<int_base then 
{ Show equivalent |n|, in region 4 }
if n=par_shape_loc then
  begin print_esc({"parshape"=}413); print_char({"="=}61);
  if  eqtb[  par_shape_loc].hh.rh   =-{0xfffffff=}268435455   then print_char({"0"=}48)
  else print_int( mem[  eqtb[  par_shape_loc].hh.rh   ].hh.lh );
  end
else if n<toks_base then
  begin print_cmd_chr(assign_toks,n); print_char({"="=}61);
  if  eqtb[  n].hh.rh  <>-{0xfffffff=}268435455   then show_token_list( mem[  eqtb[   n].hh.rh  ].hh.rh ,-{0xfffffff=}268435455  ,32);
  end
else if n<box_base then
  begin print_esc({"toks"=}412); print_int(n-toks_base); print_char({"="=}61);
  if  eqtb[  n].hh.rh  <>-{0xfffffff=}268435455   then show_token_list( mem[  eqtb[   n].hh.rh  ].hh.rh ,-{0xfffffff=}268435455  ,32);
  end
else if n<cur_font_loc then
  begin print_esc({"box"=}414); print_int(n-box_base); print_char({"="=}61);
  if  eqtb[  n].hh.rh  =-{0xfffffff=}268435455   then print({"void"=}415)
  else  begin depth_threshold:=0; breadth_max:=1; show_node_list( eqtb[  n].hh.rh  );
    end;
  end
else if n<cat_code_base then 
{ Show the font identifier in |eqtb[n]| }
begin if n=cur_font_loc then print({"current font"=}416)
else if n<math_font_base+16 then
  begin print_esc({"textfont"=}417); print_int(n-math_font_base);
  end
else if n<math_font_base+32 then
  begin print_esc({"scriptfont"=}418); print_int(n-math_font_base-16);
  end
else  begin print_esc({"scriptscriptfont"=}419); print_int(n-math_font_base-32);
  end;
print_char({"="=}61);

print_esc(hash[font_id_base+ eqtb[  n].hh.rh  ].rh);
  {that's |font_id_text(equiv(n))|}
end


else 
{ Show the halfword code in |eqtb[n]| }
if n<math_code_base then
  begin if n<lc_code_base then
    begin print_esc({"catcode"=}420); print_int(n-cat_code_base);
    end
  else if n<uc_code_base then
    begin print_esc({"lccode"=}421); print_int(n-lc_code_base);
    end
  else if n<sf_code_base then
    begin print_esc({"uccode"=}422); print_int(n-uc_code_base);
    end
  else  begin print_esc({"sfcode"=}423); print_int(n-sf_code_base);
    end;
  print_char({"="=}61); print_int( eqtb[  n].hh.rh  );
  end
else  begin print_esc({"mathcode"=}424); print_int(n-math_code_base);
  print_char({"="=}61); print_int(  eqtb[   n].hh.rh   );
  end




else if n<dimen_base then 
{ Show equivalent |n|, in region 5 }
begin if n<count_base then print_param(n-int_base)
else if  n<del_code_base then
  begin print_esc({"count"=}484); print_int(n-count_base);
  end
else  begin print_esc({"delcode"=}485); print_int(n-del_code_base);
  end;
print_char({"="=}61); print_int(eqtb[n].int);
end


else if n<=eqtb_size then 
{ Show equivalent |n|, in region 6 }
begin if n<scaled_base then print_length_param(n-dimen_base)
else  begin print_esc({"dimen"=}508); print_int(n-scaled_base);
  end;
print_char({"="=}61); print_scaled(eqtb[n].int ); print({"pt"=}402);
end


else print_char({"?"=}63); {this can't happen either}
end;
endif('STAT') 



{ 259. }

{tangle:pos tex.web:5528:3: }

{ Here is the subroutine that searches the hash table for an identifier
that matches a given string of length |l>1| appearing in |buffer[j..
(j+l-1)]|. If the identifier is found, the corresponding hash table address
is returned. Otherwise, if the global variable |no_new_control_sequence|
is |true|, the dummy address |undefined_control_sequence| is returned.
Otherwise the identifier is inserted into the hash table and its location
is returned. } function id_lookup( j, l:integer):halfword ; {search the hash table}
label found; {go here if you found it}
var h:integer; {hash code}
 d:integer; {number of characters in incomplete current string}
 p:halfword ; {index in |hash| array}
 k:halfword ; {index in |buffer| array}
begin 
{ Compute the hash code |h| }
h:=buffer[j];
for k:=j+1 to j+l-1 do
  begin h:=h+h+buffer[k];
  while h>=hash_prime do h:=h-hash_prime;
  end

;
p:=h+hash_base; {we start searching here; note that |0<=h<hash_prime|}
 while true do  begin if  hash[ p].rh >0 then if (str_start[  hash[  p].rh +1]-str_start[  hash[  p].rh ]) =l then
    if str_eq_buf( hash[ p].rh ,j) then goto found;
  if  hash[ p].lh =0 then
    begin if no_new_control_sequence then
      p:=undefined_control_sequence
    else 
{ Insert a new control sequence after |p|, then make |p| point to it }
begin if  hash[ p].rh >0 then
  begin if hash_high<hash_extra then
      begin incr(hash_high);
       hash[ p].lh :=hash_high+eqtb_size; p:=hash_high+eqtb_size;
      end
    else begin
      repeat if  (hash_used=hash_base)  then overflow({"hash size"=}511,hash_size+hash_extra);
{ \xref[TeX capacity exceeded hash size][\quad hash size] }
      decr(hash_used);
      until  hash[ hash_used].rh =0; {search for an empty location in |hash|}
     hash[ p].lh :=hash_used; p:=hash_used;
    end;
  end;
 begin if pool_ptr+ l > pool_size then overflow({"pool size"=}257,pool_size-init_pool_ptr); { \xref[TeX capacity exceeded pool size][\quad pool size] } end ; d:= (pool_ptr - str_start[str_ptr]) ;
while pool_ptr>str_start[str_ptr] do
  begin decr(pool_ptr); str_pool[pool_ptr+l]:=str_pool[pool_ptr];
  end; {move current string up to make room for another}
for k:=j to j+l-1 do  begin str_pool[pool_ptr]:=   buffer[  k] ; incr(pool_ptr); end ;
 hash[ p].rh :=make_string; pool_ptr:=pool_ptr+d;
 ifdef('STAT')  incr(cs_count); endif('STAT')  

end

;
    goto found;
    end;
  p:= hash[ p].lh ;
  end;
found: id_lookup:=p;
end;



{ 264. }

{tangle:pos tex.web:5631:3: }

{ We need to put \TeX's ``primitive'' control sequences into the hash
table, together with their command code (which will be the |eq_type|)
and an operand (which will be the |equiv|). The |primitive| procedure
does this, in a way that no \TeX\ user can. The global value |cur_val|
contains the new |eqtb| pointer after |primitive| has acted. }  ifdef('INITEX')  procedure primitive( s:str_number; c:quarterword; o:halfword);
var k:pool_pointer; {index into |str_pool|}
 j:small_number; {index into |buffer|}
 l:small_number; {length of the string}
begin if s<256 then cur_val:=s+single_base
else  begin k:=str_start[s]; l:=str_start[s+1]-k;
    {we will move |s| into the (empty) |buffer|}
  for j:=0 to l-1 do buffer[j]:=  str_pool[ k+ j] ;
  cur_val:=id_lookup(0,l); {|no_new_control_sequence| is |false|}
  begin decr(str_ptr); pool_ptr:=str_start[str_ptr]; end ;  hash[ cur_val].rh :=s; {we don't want to have the string twice}
  end;
 eqtb[  cur_val].hh.b1  :=level_one;  eqtb[  cur_val].hh.b0  :=c;  eqtb[  cur_val].hh.rh  :=o;
end;
endif('INITEX') 



{ 268. \[19] Saving and restoring equivalents }

{tangle:pos tex.web:5813:43: }

{ The nested structure provided by `$\.[\char'173]\ldots\.[\char'175]$' groups
in \TeX\ means that |eqtb| entries valid in outer groups should be saved
and restored later if they are overridden inside the braces. When a new |eqtb|
value is being assigned, the program therefore checks to see if the previous
entry belongs to an outer level. In such a case, the old value is placed
on the |save_stack| just before the new value enters |eqtb|. At the
end of a grouping level, i.e., when the right brace is sensed, the
|save_stack| is used to restore the outer values, and the inner ones are
destroyed.

Entries on the |save_stack| are of type |memory_word|. The top item on
this stack is |save_stack[p]|, where |p=save_ptr-1|; it contains three
fields called |save_type|, |save_level|, and |save_index|, and it is
interpreted in one of four ways:

\yskip\hangg 1) If |save_type(p)=restore_old_value|, then
|save_index(p)| is a location in |eqtb| whose current value should
be destroyed at the end of the current group and replaced by |save_stack[p-1]|.
Furthermore if |save_index(p)>=int_base|, then |save_level(p)|
should replace the corresponding entry in |xeq_level|.

\yskip\hangg 2) If |save_type(p)=restore_zero|, then |save_index(p)|
is a location in |eqtb| whose current value should be destroyed at the end
of the current group, when it should be
replaced by the value of |eqtb[undefined_control_sequence]|.

\yskip\hangg 3) If |save_type(p)=insert_token|, then |save_index(p)|
is a token that should be inserted into \TeX's input when the current
group ends.

\yskip\hangg 4) If |save_type(p)=level_boundary|, then |save_level(p)|
is a code explaining what kind of group we were previously in, and
|save_index(p)| points to the level boundary word at the bottom of
the entries for that group. }

{ 270. }

{tangle:pos tex.web:5892:3: }

{ The global variable |cur_group| keeps track of what sort of group we are
currently in. Another global variable, |cur_boundary|, points to the
topmost |level_boundary| word.  And |cur_level| is the current depth of
nesting. The routines are designed to preserve the condition that no entry
in the |save_stack| or in |eqtb| ever has a level greater than |cur_level|. }

{ 273. }

{tangle:pos tex.web:5915:3: }

{ The following macro is used to test if there is room for up to six more
entries on |save_stack|. By making a conservative test like this, we can
get by with testing for overflow in only a few places. }

{ 274. }

{tangle:pos tex.web:5925:3: }

{ Procedure |new_save_level| is called when a group begins. The
argument is a group identification code like `|hbox_group|'. After
calling this routine, it is safe to put five more entries on |save_stack|.

In some cases integer-valued items are placed onto the
|save_stack| just below a |level_boundary| word, because this is a
convenient place to keep information that is supposed to ``pop up'' just
when the group has finished.
For example, when `\.[\\hbox to 100pt]\grp' is being treated, the 100pt
dimension is stored on |save_stack| just before |new_save_level| is
called.

We use the notation |saved(k)| to stand for an integer item that
appears in location |save_ptr+k| of the save stack. } procedure new_save_level( c:group_code); {begin a new level of grouping}
begin if save_ptr>max_save_stack then begin max_save_stack:=save_ptr; if max_save_stack>save_size-6 then overflow({"save size"=}549,save_size); { \xref[TeX capacity exceeded save size][\quad save size] } end ;
save_stack[ save_ptr].hh.b0 :=level_boundary; save_stack[ save_ptr].hh.b1 :=cur_group;
save_stack[ save_ptr].hh.rh :=cur_boundary;
if cur_level=max_quarterword then overflow({"grouping levels"=}550,
{ \xref[TeX capacity exceeded grouping levels][\quad grouping levels] }
  max_quarterword-min_quarterword);
  {quit if |(cur_level+1)| is too big to be stored in |eqtb|}
cur_boundary:=save_ptr; incr(cur_level); incr(save_ptr); cur_group:=c;
end;



{ 275. }

{tangle:pos tex.web:5953:3: }

{ Just before an entry of |eqtb| is changed, the following procedure should
be called to update the other data structures properly. It is important
to keep in mind that reference counts in |mem| include references from
within |save_stack|, so these counts must be handled carefully.
\xref[reference counts] } procedure eq_destroy( w:memory_word); {gets ready to forget |w|}
var q:halfword ; {|equiv| field of |w|}
begin case  w.hh.b0  of
call,long_call,outer_call,long_outer_call: delete_token_ref( w.hh.rh );
glue_ref: delete_glue_ref( w.hh.rh );
shape_ref: begin q:= w.hh.rh ; {we need to free a \.[\\parshape] block}
  if q<>-{0xfffffff=}268435455   then free_node(q, mem[ q].hh.lh + mem[ q].hh.lh +1);
  end; {such a block is |2n+1| words long, where |n=info(q)|}
box_ref: flush_node_list( w.hh.rh );
 else   
 end ;
end;



{ 276. }

{tangle:pos tex.web:5972:3: }

{ To save a value of |eqtb[p]| that was established at level |l|, we
can use the following subroutine. } procedure eq_save( p:halfword ; l:quarterword); {saves |eqtb[p]|}
begin if save_ptr>max_save_stack then begin max_save_stack:=save_ptr; if max_save_stack>save_size-6 then overflow({"save size"=}549,save_size); { \xref[TeX capacity exceeded save size][\quad save size] } end ;
if l=level_zero then save_stack[ save_ptr].hh.b0 :=restore_zero
else  begin save_stack[save_ptr]:=eqtb[p]; incr(save_ptr);
  save_stack[ save_ptr].hh.b0 :=restore_old_value;
  end;
save_stack[ save_ptr].hh.b1 :=l; save_stack[ save_ptr].hh.rh :=p; incr(save_ptr);
end;



{ 277. }

{tangle:pos tex.web:5984:3: }

{ The procedure |eq_define| defines an |eqtb| entry having specified
|eq_type| and |equiv| fields, and saves the former value if appropriate.
This procedure is used only for entries in the first four regions of |eqtb|,
i.e., only for entries that have |eq_type| and |equiv| fields.
After calling this routine, it is safe to put four more entries on
|save_stack|, provided that there was room for four more entries before
the call, since |eq_save| makes the necessary test. } procedure eq_define( p:halfword ; t:quarterword; e:halfword);
  {new data for |eqtb|}
begin if  eqtb[  p].hh.b1  =cur_level then eq_destroy(eqtb[p])
else if cur_level>level_one then eq_save(p, eqtb[  p].hh.b1  );
 eqtb[  p].hh.b1  :=cur_level;  eqtb[  p].hh.b0  :=t;  eqtb[  p].hh.rh  :=e;
end;



{ 278. }

{tangle:pos tex.web:5999:3: }

{ The counterpart of |eq_define| for the remaining (fullword) positions in
|eqtb| is called |eq_word_define|. Since |xeq_level[p]>=level_one| for all
|p|, a `|restore_zero|' will never be used in this case. } procedure eq_word_define( p:halfword ; w:integer);
begin if xeq_level[p]<>cur_level then
  begin eq_save(p,xeq_level[p]); xeq_level[p]:=cur_level;
  end;
eqtb[p].int:=w;
end;



{ 279. }

{tangle:pos tex.web:6010:3: }

{ The |eq_define| and |eq_word_define| routines take care of local definitions.
\xref[global definitions]
Global definitions are done in almost the same way, but there is no need
to save old values, and the new value is associated with |level_one|. } procedure geq_define( p:halfword ; t:quarterword; e:halfword);
  {global |eq_define|}
begin eq_destroy(eqtb[p]);
 eqtb[  p].hh.b1  :=level_one;  eqtb[  p].hh.b0  :=t;  eqtb[  p].hh.rh  :=e;
end;


procedure geq_word_define( p:halfword ; w:integer); {global |eq_word_define|}
begin eqtb[p].int:=w; xeq_level[p]:=level_one;
end;



{ 280. }

{tangle:pos tex.web:6025:3: }

{ Subroutine |save_for_after| puts a token on the stack for save-keeping. } procedure save_for_after( t:halfword);
begin if cur_level>level_one then
  begin if save_ptr>max_save_stack then begin max_save_stack:=save_ptr; if max_save_stack>save_size-6 then overflow({"save size"=}549,save_size); { \xref[TeX capacity exceeded save size][\quad save size] } end ;
  save_stack[ save_ptr].hh.b0 :=insert_token; save_stack[ save_ptr].hh.b1 :=level_zero;
  save_stack[ save_ptr].hh.rh :=t; incr(save_ptr);
  end;
end;



{ 281. }

{tangle:pos tex.web:6035:3: }

{ The |unsave| routine goes the other way, taking items off of |save_stack|.
This routine takes care of restoration when a level ends; everything
belonging to the topmost group is cleared off of the save stack. }{ \4 }
{ Declare the procedure called |restore_trace| }
 ifdef('STAT')  procedure restore_trace( p:halfword ; s:str_number);
  {|eqtb[p]| has just been restored or retained}
begin begin_diagnostic; print_char({"["=}123); print(s); print_char({" "=}32);
show_eqtb(p); print_char({"]"=}125);
end_diagnostic(false);
end;
endif('STAT') 

 

procedure back_input; forward; { \2 }
procedure unsave; {pops the top level off the save stack}
label done;
var p:halfword ; {position to be restored}
 l:quarterword; {saved level, if in fullword regions of |eqtb|}
 t:halfword; {saved value of |cur_tok|}
begin if cur_level>level_one then
  begin decr(cur_level);
  
{ Clear off top level from |save_stack| }
 while true do  begin decr(save_ptr);
  if save_stack[ save_ptr].hh.b0 =level_boundary then goto done;
  p:=save_stack[ save_ptr].hh.rh ;
  if save_stack[ save_ptr].hh.b0 =insert_token then
    
{ Insert token |p| into \TeX's input }
begin t:=cur_tok; cur_tok:=p; back_input; cur_tok:=t;
end


  else  begin if save_stack[ save_ptr].hh.b0 =restore_old_value then
      begin l:=save_stack[ save_ptr].hh.b1 ; decr(save_ptr);
      end
    else save_stack[save_ptr]:=eqtb[undefined_control_sequence];
    
{ Store \(s)|save_stack[save_ptr]| in |eqtb[p]|, unless |eqtb[p]| holds a global value }
if (p<int_base)or(p>eqtb_size) then
  if  eqtb[  p].hh.b1  =level_one then
    begin eq_destroy(save_stack[save_ptr]); {destroy the saved value}
     ifdef('STAT')  if eqtb[int_base+ tracing_restores_code].int  >0 then restore_trace(p,{"retaining"=}552); endif('STAT')  

    end
  else  begin eq_destroy(eqtb[p]); {destroy the current value}
    eqtb[p]:=save_stack[save_ptr]; {restore the saved value}
     ifdef('STAT')  if eqtb[int_base+ tracing_restores_code].int  >0 then restore_trace(p,{"restoring"=}553); endif('STAT')  

    end
else if xeq_level[p]<>level_one then
  begin eqtb[p]:=save_stack[save_ptr]; xeq_level[p]:=l;
   ifdef('STAT')  if eqtb[int_base+ tracing_restores_code].int  >0 then restore_trace(p,{"restoring"=}553); endif('STAT')  

  end
else  begin
   ifdef('STAT')  if eqtb[int_base+ tracing_restores_code].int  >0 then restore_trace(p,{"retaining"=}552); endif('STAT')  

  end

;
    end;
  end;
done: cur_group:=save_stack[ save_ptr].hh.b1 ; cur_boundary:=save_stack[ save_ptr].hh.rh 

;
  end
else confusion({"curlevel"=}551); {|unsave| is not used when |cur_group=bottom_level|}
{ \xref[this can't happen curlevel][\quad curlevel] }
end;



{ 288. }

{tangle:pos tex.web:6129:3: }

{ The |prepare_mag| subroutine is called whenever \TeX\ wants to use |mag|
for magnification. } procedure prepare_mag;
begin if (mag_set>0)and(eqtb[int_base+ mag_code].int  <>mag_set) then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Incompatible magnification ("=} 555); end ; print_int(eqtb[int_base+ mag_code].int  );
{ \xref[Incompatible magnification] }
  print({");"=}556); print_nl({" the previous value will be retained"=}557);
   begin help_ptr:=2; help_line[1]:={"I can handle only one magnification ratio per job. So I've"=} 558; help_line[0]:={"reverted to the magnification you used earlier on this run."=} 559; end ;

  int_error(mag_set);
  geq_word_define(int_base+mag_code,mag_set); {|mag:=mag_set|}
  end;
if (eqtb[int_base+ mag_code].int  <=0)or(eqtb[int_base+ mag_code].int  >32768) then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Illegal magnification has been changed to 1000"=} 560); end ;

{ \xref[Illegal magnification...] }
   begin help_ptr:=1; help_line[0]:={"The magnification ratio must be between 1 and 32768."=} 561; end ;
  int_error(eqtb[int_base+ mag_code].int  ); geq_word_define(int_base+mag_code,1000);
  end;
mag_set:=eqtb[int_base+ mag_code].int  ;
end;



{ 289. \[20] Token lists }

{tangle:pos tex.web:6151:22: }

{ A \TeX\ token is either a character or a control sequence, and it is
\xref[token]
represented internally in one of two ways: (1)~A character whose ASCII
code number is |c| and whose command code is |m| is represented as the
number $2^8m+c$; the command code is in the range |1<=m<=14|. (2)~A control
sequence whose |eqtb| address is |p| is represented as the number
|cs_token_flag+p|. Here |cs_token_flag=$2^[12]-1$| is larger than
$2^8m+c$, yet it is small enough that |cs_token_flag+p< max_halfword|;
thus, a token fits comfortably in a halfword.

A token |t| represents a |left_brace| command if and only if
|t<left_brace_limit|; it represents a |right_brace| command if and only if
we have |left_brace_limit<=t<right_brace_limit|; and it represents a |match| or
|end_match| command if and only if |match_token<=t<=end_match_token|.
The following definitions take care of these token-oriented constants
and a few others. }

{ 291. }

{tangle:pos tex.web:6187:3: }

{ A token list is a singly linked list of one-word nodes in |mem|, where
each word contains a token and a link. Macro definitions, output-routine
definitions, marks, \.[\\write] texts, and a few other things
are remembered by \TeX\ in the form
of token lists, usually preceded by a node with a reference count in its
|token_ref_count| field. The token stored in location |p| is called
|info(p)|.

Three special commands appear in the token lists of macro definitions.
When |m=match|, it means that \TeX\ should scan a parameter
for the current macro; when |m=end_match|, it means that parameter
matching should end and \TeX\ should start reading the macro text; and
when |m=out_param|, it means that \TeX\ should insert parameter
number |c| into the text at this point.

The enclosing \.[\char'173] and \.[\char'175] characters of a macro
definition are omitted, but an output routine
will be enclosed in braces.

Here is an example macro definition that illustrates these conventions.
After \TeX\ processes the text
$$\.[\\def\\mac a\#1\#2 \\b \[\#1\\-a \#\#1\#2 \#2\]]$$
the definition of \.[\\mac] is represented as a token list containing
$$\def\,[\hskip2pt]
\vbox[\halign[\hfil#\hfil\cr
(reference count), |letter|\,\.a, |match|\,\#, |match|\,\#, |spacer|\,\.\ ,
\.[\\b], |end_match|,\cr
|out_param|\,1, \.[\\-], |letter|\,\.a, |spacer|\,\.\ , |mac_param|\,\#,
|other_char|\,\.1,\cr
|out_param|\,2, |spacer|\,\.\ , |out_param|\,2.\cr]]$$
The procedure |scan_toks| builds such token lists, and |macro_call|
does the parameter matching.
\xref[reference counts]

Examples such as
$$\.[\\def\\m\[\\def\\m\[a\]\ b\]]$$
explain why reference counts would be needed even if \TeX\ had no \.[\\let]
operation: When the token list for \.[\\m] is being read, the redefinition of
\.[\\m] changes the |eqtb| entry before the token list has been fully
consumed, so we dare not simply destroy a token list when its
control sequence is being redefined.

If the parameter-matching part of a definition ends with `\.[\#\[]',
the corresponding token list will have `\.\[' just before the `|end_match|'
and also at the very end. The first `\.\[' is used to delimit the parameter; the
second one keeps the first from disappearing. }

{ 295. }

{tangle:pos tex.web:6313:3: }

{ Here's the way we sometimes want to display a token list, given a pointer
to its reference count; the pointer may be null. } procedure token_show( p:halfword );
begin if p<>-{0xfffffff=}268435455   then show_token_list( mem[ p].hh.rh ,-{0xfffffff=}268435455  ,10000000);
end;



{ 296. }

{tangle:pos tex.web:6320:3: }

{ The |print_meaning| subroutine displays |cur_cmd| and |cur_chr| in
symbolic form, including the expansion of a macro or mark. } procedure print_meaning;
begin print_cmd_chr(cur_cmd,cur_chr);
if cur_cmd>=call then
  begin print_char({":"=}58); print_ln; token_show(cur_chr);
  end
else if cur_cmd=top_bot_mark then
  begin print_char({":"=}58); print_ln;
  token_show(cur_mark[cur_chr]);
  end;
end;



{ 299. }

{tangle:pos tex.web:6419:3: }

{ Here is a procedure that displays the current command. } procedure show_cur_cmd_chr;
begin begin_diagnostic; print_nl({"["=}123);
if cur_list.mode_field <>shown_mode then
  begin print_mode(cur_list.mode_field ); print({": "=}576); shown_mode:=cur_list.mode_field ;
  end;
print_cmd_chr(cur_cmd,cur_chr); print_char({"]"=}125);
end_diagnostic(false);
end;



{ 302. }

{tangle:pos tex.web:6465:3: }

{ We've already defined the special variable |loc==cur_input.loc_field|
in our discussion of basic input-output routines. The other components of
|cur_input| are defined in the same way: }

{ 303. }

{tangle:pos tex.web:6475:3: }

{ Let's look more closely now at the control variables
(|state|,~|index|,~|start|,~|loc|,~|limit|,~|name|),
assuming that \TeX\ is reading a line of characters that have been input
from some file or from the user's terminal. There is an array called
|buffer| that acts as a stack of all lines of characters that are
currently being read from files, including all lines on subsidiary
levels of the input stack that are not yet completed. \TeX\ will return to
the other lines when it is finished with the present input file.

(Incidentally, on a machine with byte-oriented addressing, it might be
appropriate to combine |buffer| with the |str_pool| array,
letting the buffer entries grow downward from the top of the string pool
and checking that these two tables don't bump into each other.)

The line we are currently working on begins in position |start| of the
buffer; the next character we are about to read is |buffer[loc]|; and
|limit| is the location of the last character present.  If |loc>limit|,
the line has been completely read. Usually |buffer[limit]| is the
|end_line_char|, denoting the end of a line, but this is not
true if the current line is an insertion that was entered on the user's
terminal in response to an error message.

The |name| variable is a string number that designates the name of
the current file, if we are reading a text file. It is zero if we
are reading from the terminal; it is |n+1| if we are reading from
input stream |n|, where |0<=n<=16|. (Input stream 16 stands for
an invalid stream number; in such cases the input is actually from
the terminal, under control of the procedure |read_toks|.)

The |state| variable has one of three values, when we are scanning such
files:
$$\baselineskip 15pt\vbox[\halign[#\hfil\cr
1) |state=mid_line| is the normal state.\cr
2) |state=skip_blanks| is like |mid_line|, but blanks are ignored.\cr
3) |state=new_line| is the state at the beginning of a line.\cr]]$$
These state values are assigned numeric codes so that if we add the state
code to the next character's command code, we get distinct values. For
example, `|mid_line+spacer|' stands for the case that a blank
space character occurs in the middle of a line when it is not being
ignored; after this case is processed, the next value of |state| will
be |skip_blanks|. }

{ 307. }

{tangle:pos tex.web:6631:3: }

{ However, all this discussion about input state really applies only to the
case that we are inputting from a file. There is another important case,
namely when we are currently getting input from a token list. In this case
|state=token_list|, and the conventions about the other state variables
are different:

\yskip\hang|loc| is a pointer to the current node in the token list, i.e.,
the node that will be read next. If |loc=null|, the token list has been
fully read.

\yskip\hang|start| points to the first node of the token list; this node
may or may not contain a reference count, depending on the type of token
list involved.

\yskip\hang|token_type|, which takes the place of |index| in the
discussion above, is a code number that explains what kind of token list
is being scanned.

\yskip\hang|name| points to the |eqtb| address of the control sequence
being expanded, if the current token list is a macro.

\yskip\hang|param_start|, which takes the place of |limit|, tells where
the parameters of the current macro begin in the |param_stack|, if the
current token list is a macro.

\yskip\noindent The |token_type| can take several values, depending on
where the current token list came from:

\yskip\hang|parameter|, if a parameter is being scanned;

\hang|u_template|, if the \<u_j> part of an alignment
template is being scanned;

\hang|v_template|, if the \<v_j> part of an alignment
template is being scanned;

\hang|backed_up|, if the token list being scanned has been inserted as
`to be read again';

\hang|inserted|, if the token list being scanned has been inserted as
the text expansion of a \.[\\count] or similar variable;

\hang|macro|, if a user-defined control sequence is being scanned;

\hang|output_text|, if an \.[\\output] routine is being scanned;

\hang|every_par_text|, if the text of \.[\\everypar] is being scanned;

\hang|every_math_text|, if the text of \.[\\everymath] is being scanned;

\hang|every_display_text|, if the text of \.[\\everydisplay] is being scanned;

\hang|every_hbox_text|, if the text of \.[\\everyhbox] is being scanned;

\hang|every_vbox_text|, if the text of \.[\\everyvbox] is being scanned;

\hang|every_job_text|, if the text of \.[\\everyjob] is being scanned;

\hang|every_cr_text|, if the text of \.[\\everycr] is being scanned;

\hang|mark_text|, if the text of a \.[\\mark] is being scanned;

\hang|write_text|, if the text of a \.[\\write] is being scanned.

\yskip\noindent
The codes for |output_text|, |every_par_text|, etc., are equal to a constant
plus the corresponding codes for token list parameters |output_routine_loc|,
|every_par_loc|, etc.  The token list begins with a reference count if and
only if |token_type>=macro|.
\xref[reference counts] }

{ 311. }

{tangle:pos tex.web:6758:3: }

{ The status at each level is indicated by printing two lines, where the first
line indicates what was read so far and the second line shows what remains
to be read. The context is cropped, if necessary, so that the first line
contains at most |half_error_line| characters, and the second contains
at most |error_line|. Non-current input levels whose |token_type| is
`|backed_up|' are shown only if they have not been fully read. } procedure show_context; {prints where the scanner is}
label done;
var old_setting:0..max_selector; {saved |selector| setting}
 nn:integer; {number of contexts shown so far, less one}
 bottom_line:boolean; {have we reached the final context to be shown?}

{ Local variables for formatting calculations }
 i:0..buf_size; {index into |buffer|}
 j:0..buf_size; {end of current line in |buffer|}
 l:0..half_error_line; {length of descriptive information on line 1}
 m:integer; {context information gathered for line 2}
 n:0..error_line; {length of line 1}
 p: integer; {starting or ending place in |trick_buf|}
 q: integer; {temporary index}



begin base_ptr:=input_ptr; input_stack[base_ptr]:=cur_input;
  {store current state}
nn:=-1; bottom_line:=false;
 while true do  begin cur_input:=input_stack[base_ptr]; {enter into the context}
  if (cur_input.state_field <>token_list) then
    if (cur_input.name_field >17) or (base_ptr=0) then bottom_line:=true;
  if (base_ptr=input_ptr)or bottom_line or(nn<eqtb[int_base+ error_context_lines_code].int  ) then
    
{ Display the current context }
begin if (base_ptr=input_ptr) or (cur_input.state_field <>token_list) or
   (cur_input.index_field  <>backed_up) or (cur_input.loc_field <>-{0xfffffff=}268435455  ) then
    {we omit backed-up token lists that have already been read}
  begin tally:=0; {get ready to count characters}
  old_setting:=selector;
  if cur_input.state_field <>token_list then
    begin 
{ Print location of current line }
if cur_input.name_field <=17 then
  if (cur_input.name_field =0)  then
    if base_ptr=0 then print_nl({"<*>"=}581) else print_nl({"<insert> "=}582)
  else  begin print_nl({"<read "=}583);
    if cur_input.name_field =17 then print_char({"*"=}42) else print_int(cur_input.name_field -1);
{ \xref[*\relax] }
    print_char({">"=}62);
    end
else  begin print_nl({"l."=}584); print_int(line);
  end;
print_char({" "=}32)

;
    
{ Pseudoprint the line }
 begin l:=tally; tally:=0; selector:=pseudo; trick_count:=1000000; end ;
if buffer[cur_input.limit_field ]=eqtb[int_base+ end_line_char_code].int   then j:=cur_input.limit_field 
else j:=cur_input.limit_field +1; {determine the effective end of the line}
if j>0 then for i:=cur_input.start_field  to j-1 do
  begin if i=cur_input.loc_field  then  begin first_count:=tally; trick_count:=tally+1+error_line-half_error_line; if trick_count<error_line then trick_count:=error_line; end ;
  print(buffer[i]);
  end

;
    end
  else  begin 
{ Print type of token list }
case cur_input.index_field   of
parameter: print_nl({"<argument> "=}585);
u_template,v_template: print_nl({"<template> "=}586);
backed_up: if cur_input.loc_field =-{0xfffffff=}268435455   then print_nl({"<recently read> "=}587)
  else print_nl({"<to be read again> "=}588);
inserted: print_nl({"<inserted text> "=}589);
macro: begin print_ln; print_cs(cur_input.name_field );
  end;
output_text: print_nl({"<output> "=}590);
every_par_text: print_nl({"<everypar> "=}591);
every_math_text: print_nl({"<everymath> "=}592);
every_display_text: print_nl({"<everydisplay> "=}593);
every_hbox_text: print_nl({"<everyhbox> "=}594);
every_vbox_text: print_nl({"<everyvbox> "=}595);
every_job_text: print_nl({"<everyjob> "=}596);
every_cr_text: print_nl({"<everycr> "=}597);
mark_text: print_nl({"<mark> "=}598);
write_text: print_nl({"<write> "=}599);
 else  print_nl({"?"=}63) {this should never happen}
 end 

;
    
{ Pseudoprint the token list }
 begin l:=tally; tally:=0; selector:=pseudo; trick_count:=1000000; end ;
if cur_input.index_field  <macro then show_token_list(cur_input.start_field ,cur_input.loc_field ,100000)
else show_token_list( mem[ cur_input.start_field ].hh.rh ,cur_input.loc_field ,100000) {avoid reference count}

;
    end;
  selector:=old_setting; {stop pseudoprinting}
  
{ Print two lines using the tricky pseudoprinted information }
if trick_count=1000000 then  begin first_count:=tally; trick_count:=tally+1+error_line-half_error_line; if trick_count<error_line then trick_count:=error_line; end ;
  {|set_trick_count| must be performed}
if tally<trick_count then m:=tally-first_count
else m:=trick_count-first_count; {context on line 2}
if l+first_count<=half_error_line then
  begin p:=0; n:=l+first_count;
  end
else  begin print({"..."=}275); p:=l+first_count-half_error_line+3;
  n:=half_error_line;
  end;
for q:=p to first_count-1 do print_char(trick_buf[q mod error_line]);
print_ln;
for q:=1 to n do print_char({" "=}32); {print |n| spaces to begin line~2}
if m+n<=error_line then p:=first_count+m else p:=first_count+(error_line-n-3);
for q:=first_count to p-1 do print_char(trick_buf[q mod error_line]);
if m+n>error_line then print({"..."=}275)

;
  incr(nn);
  end;
end


  else if nn=eqtb[int_base+ error_context_lines_code].int   then
    begin print_nl({"..."=}275); incr(nn); {omitted if |error_context_lines<0|}
    end;
  if bottom_line then goto done;
  decr(base_ptr);
  end;
done: cur_input:=input_stack[input_ptr]; {restore original state}
end;



{ 316. }

{tangle:pos tex.web:6886:3: }

{ The following code sets up the print routines so that they will gather
the desired information. }

{ 321. \[23] Maintaining the input stacks }

{tangle:pos tex.web:6942:39: }

{ The following subroutines change the input status in commonly needed ways.

First comes |push_input|, which stores the current state and creates a
new level (having, initially, the same properties as the old). }

{ 322. }

{tangle:pos tex.web:6958:3: }

{ And of course what goes up must come down. }

{ 323. }

{tangle:pos tex.web:6964:3: }

{ Here is a procedure that starts a new level of token-list input, given
a token list |p| and its type |t|. If |t=macro|, the calling routine should
set |name| and |loc|. } procedure begin_token_list( p:halfword ; t:quarterword);
begin {  } begin if input_ptr>max_in_stack then begin max_in_stack:=input_ptr; if input_ptr=stack_size then overflow({"input stack size"=}600,stack_size); { \xref[TeX capacity exceeded input stack size][\quad input stack size] } end; input_stack[input_ptr]:=cur_input; incr(input_ptr); end ; cur_input.state_field :=token_list; cur_input.start_field :=p; cur_input.index_field  :=t;
if t>=macro then {the token list starts with a reference count}
  begin incr(  mem[   p].hh.lh  ) ;
  if t=macro then cur_input.limit_field  :=param_ptr
  else  begin cur_input.loc_field := mem[ p].hh.rh ;
    if eqtb[int_base+ tracing_macros_code].int  >1 then
      begin begin_diagnostic; print_nl({""=}335);
      case t of
      mark_text:print_esc({"mark"=}348);
      write_text:print_esc({"write"=}601);
       else  print_cmd_chr(assign_toks,t-output_text+output_routine_loc)
       end ;

      print({"->"=}564); token_show(p); end_diagnostic(false);
      end;
    end;
  end
else cur_input.loc_field :=p;
end;



{ 324. }

{tangle:pos tex.web:6991:3: }

{ When a token list has been fully scanned, the following computations
should be done as we leave that level of input. The |token_type| tends
to be equal to either |backed_up| or |inserted| about 2/3 of the time.
\xref[inner loop] } procedure end_token_list; {leave a token-list input level}
begin if cur_input.index_field  >=backed_up then {token list to be deleted}
  begin if cur_input.index_field  <=inserted then flush_list(cur_input.start_field )
  else  begin delete_token_ref(cur_input.start_field ); {update reference count}
    if cur_input.index_field  =macro then {parameters must be flushed}
      while param_ptr>cur_input.limit_field   do
        begin decr(param_ptr);
        flush_list(param_stack[param_ptr]);
        end;
    end;
  end
else if cur_input.index_field  =u_template then
  if align_state>500000 then align_state:=0
  else fatal_error({"(interwoven alignment preambles are not allowed)"=}602);
{ \xref[interwoven alignment preambles...] }
{  } begin decr(input_ptr); cur_input:=input_stack[input_ptr]; end ;
begin if interrupt<>0 then pause_for_instructions; end ;
end;



{ 325. }

{tangle:pos tex.web:7015:3: }

{ Sometimes \TeX\ has read too far and wants to ``unscan'' what it has
seen. The |back_input| procedure takes care of this by putting the token
just scanned back into the input stream, ready to be read again. This
procedure can be used only if |cur_tok| represents the token to be
replaced. Some applications of \TeX\ use this procedure a lot,
so it has been slightly optimized for speed.
\xref[inner loop] } procedure back_input; {undoes one token of input}
var p:halfword ; {a token list of length one}
begin while (cur_input.state_field =token_list)and(cur_input.loc_field =-{0xfffffff=}268435455  )and(cur_input.index_field  <>v_template) do
  end_token_list; {conserve stack space}
p:=get_avail;  mem[ p].hh.lh :=cur_tok;
if cur_tok<right_brace_limit then
  if cur_tok<left_brace_limit then decr(align_state)
  else incr(align_state);
{  } begin if input_ptr>max_in_stack then begin max_in_stack:=input_ptr; if input_ptr=stack_size then overflow({"input stack size"=}600,stack_size); { \xref[TeX capacity exceeded input stack size][\quad input stack size] } end; input_stack[input_ptr]:=cur_input; incr(input_ptr); end ; cur_input.state_field :=token_list; cur_input.start_field :=p; cur_input.index_field  :=backed_up;
cur_input.loc_field :=p; {that was |back_list(p)|, without procedure overhead}
end;



{ 327. }

{tangle:pos tex.web:7039:3: }

{ The |back_error| routine is used when we want to replace an offending token
just before issuing an error message. This routine, like |back_input|,
requires that |cur_tok| has been set. We disable interrupts during the
call of |back_input| so that the help message won't be lost. } procedure back_error; {back up one token and call |error|}
begin OK_to_interrupt:=false; back_input; OK_to_interrupt:=true; error;
end;


procedure ins_error; {back up one inserted token and call |error|}
begin OK_to_interrupt:=false; back_input; cur_input.index_field  :=inserted;
OK_to_interrupt:=true; error;
end;



{ 328. }

{tangle:pos tex.web:7053:3: }

{ The |begin_file_reading| procedure starts a new level of input for lines
of characters to be read from a file, or as an insertion from the
terminal. It does not take care of opening the file, nor does it set |loc|
or |limit| or |line|.
\xref[system dependencies] } procedure begin_file_reading;
begin if in_open=max_in_open then overflow({"text input levels"=}603,max_in_open);
{ \xref[TeX capacity exceeded text input levels][\quad text input levels] }
if first=buf_size then overflow({"buffer size"=}256,buf_size);
{ \xref[TeX capacity exceeded buffer size][\quad buffer size] }
incr(in_open); {  } begin if input_ptr>max_in_stack then begin max_in_stack:=input_ptr; if input_ptr=stack_size then overflow({"input stack size"=}600,stack_size); { \xref[TeX capacity exceeded input stack size][\quad input stack size] } end; input_stack[input_ptr]:=cur_input; incr(input_ptr); end ; cur_input.index_field :=in_open;
source_filename_stack[cur_input.index_field ]:=0;full_source_filename_stack[cur_input.index_field ]:=0;
line_stack[cur_input.index_field ]:=line; cur_input.start_field :=first; cur_input.state_field :=mid_line;
cur_input.name_field :=0; {|terminal_input| is now |true|}
end;



{ 329. }

{tangle:pos tex.web:7069:3: }

{ Conversely, the variables must be downdated when such a level of input
is finished: } procedure end_file_reading;
begin first:=cur_input.start_field ; line:=line_stack[cur_input.index_field ];
if cur_input.name_field >17 then a_close(input_file[cur_input.index_field ] ); {forget it}
{  } begin decr(input_ptr); cur_input:=input_stack[input_ptr]; end ; decr(in_open);
end;



{ 330. }

{tangle:pos tex.web:7078:3: }

{ In order to keep the stack from overflowing during a long sequence of
inserted `\.[\\show]' commands, the following routine removes completed
error-inserted lines from memory. } procedure clear_for_error_prompt;
begin while (cur_input.state_field <>token_list)and (cur_input.name_field =0)  and 
  (input_ptr>0)and(cur_input.loc_field >cur_input.limit_field ) do end_file_reading;
print_ln;    ;
end;



{ 332. \[24] Getting the next token }

{tangle:pos tex.web:7104:33: }

{ The heart of \TeX's input mechanism is the |get_next| procedure, which
we shall develop in the next few sections of the program. Perhaps we
shouldn't actually call it the ``heart,'' however, because it really acts
as \TeX's eyes and mouth, reading the source files and gobbling them up.
And it also helps \TeX\ to regurgitate stored token lists that are to be
processed again.
\xref[eyes and mouth]

The main duty of |get_next| is to input one token and to set |cur_cmd|
and |cur_chr| to that token's command code and modifier. Furthermore, if
the input token is a control sequence, the |eqtb| location of that control
sequence is stored in |cur_cs|; otherwise |cur_cs| is set to zero.

Underlying this simple description is a certain amount of complexity
because of all the cases that need to be handled.
However, the inner loop of |get_next| is reasonably short and fast.

When |get_next| is asked to get the next token of a \.[\\read] line,
it sets |cur_cmd=cur_chr=cur_cs=0| in the case that no more tokens
appear on that line. (There might not be any tokens at all, if the
|end_line_char| has |ignore| as its catcode.) }

{ 336. }

{tangle:pos tex.web:7144:3: }

{ Before getting into |get_next|, let's consider the subroutine that
is called when an `\.[\\outer]' control sequence has been scanned or
when the end of a file has been reached. These two cases are distinguished
by |cur_cs|, which is zero at the end of a file. } procedure check_outer_validity;
var p:halfword ; {points to inserted token list}
 q:halfword ; {auxiliary pointer}
begin if scanner_status<>normal then
  begin deletions_allowed:=false;
  
{ Back up an outer control sequence so that it can be reread }
if cur_cs<>0 then
  begin if (cur_input.state_field =token_list)or(cur_input.name_field <1)or(cur_input.name_field >17) then
    begin p:=get_avail;  mem[ p].hh.lh :={07777=}4095 +cur_cs;
    begin_token_list( p,backed_up) ; {prepare to read the control sequence again}
    end;
  cur_cmd:=spacer; cur_chr:={" "=}32; {replace it by a space}
  end

;
  if scanner_status>skipping then
    
{ Tell the user what has run away and try to recover }
begin runaway; {print a definition, argument, or preamble}
if cur_cs=0 then begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"File ended"=} 611); end 
{ \xref[File ended while scanning...] }
else  begin cur_cs:=0; begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Forbidden control sequence found"=} 612); end ;
{ \xref[Forbidden control sequence...] }
  end;


{ Print either `\.[definition]' or `\.[use]' or `\.[preamble]' or `\.[text]', and insert tokens that should lead to recovery }
p:=get_avail;
case scanner_status of
defining:begin print({" while scanning definition"=}618);  mem[ p].hh.lh :=right_brace_token+{"]"=}125;
  end;
matching:begin print({" while scanning use"=}619);  mem[ p].hh.lh :=par_token; long_state:=outer_call;
  end;
aligning:begin print({" while scanning preamble"=}620);  mem[ p].hh.lh :=right_brace_token+{"]"=}125; q:=p;
  p:=get_avail;  mem[ p].hh.rh :=q;  mem[ p].hh.lh :={07777=}4095 +frozen_cr;
  align_state:=-1000000;
  end;
absorbing:begin print({" while scanning text"=}621);  mem[ p].hh.lh :=right_brace_token+{"]"=}125;
  end;
end; {there are no other cases}
begin_token_list( p,inserted) 

;
print({" of "=}613); sprint_cs(warning_index);
 begin help_ptr:=4; help_line[3]:={"I suspect you have forgotten a `]', causing me"=} 614; help_line[2]:={"to read past where you wanted me to stop."=} 615; help_line[1]:={"I'll try to recover; but if the error is serious,"=} 616; help_line[0]:={"you'd better type `E' or `X' now and fix your file."=} 617; end ;

error;
end


  else  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Incomplete "=} 605); end ; print_cmd_chr(if_test,cur_if);
{ \xref[Incomplete \\if...] }
    print({"; all text was ignored after line "=}606); print_int(skip_line);
     begin help_ptr:=3; help_line[2]:={"A forbidden control sequence occurred in skipped text."=} 607; help_line[1]:={"This kind of error happens when you say `\if...' and forget"=} 608; help_line[0]:={"the matching `\fi'. I've inserted a `\fi'; this might work."=} 609; end ;
    if cur_cs<>0 then cur_cs:=0
    else help_line[2]:= 
      {"The file ended while I was skipping conditional text."=}610;
    cur_tok:={07777=}4095 +frozen_fi; ins_error;
    end;
  deletions_allowed:=true;
  end;
end;



{ 340. }

{tangle:pos tex.web:7226:3: }

{ We need to mention a procedure here that may be called by |get_next|. } procedure firm_up_the_line; forward;



{ 341. }

{tangle:pos tex.web:7230:3: }

{ Now we're ready to take the plunge into |get_next| itself. Parts of
this routine are executed more often than any other instructions of \TeX.
\xref[mastication]\xref[inner loop] } procedure get_next; {sets |cur_cmd|, |cur_chr|, |cur_cs| to next token}
label restart, {go here to get the next input token}
  switch, {go here to eat the next character from a file}
  reswitch, {go here to digest it again}
  start_cs, {go here to start looking for a control sequence}
  found, {go here when a control sequence has been found}
  exit; {go here when the next input token has been got}
var k:0..buf_size; {an index into |buffer|}
 t:halfword; {a token}
 cat:0..max_char_code; {|cat_code(cur_chr)|, usually}
 c, cc:ASCII_code; {constituents of a possible expanded code}
 d:2..3; {number of excess characters in an expanded code}
begin restart: cur_cs:=0;
if cur_input.state_field <>token_list then

{ Input from external file, |goto restart| if no input found }
{ \xref[inner loop] }
begin switch: if cur_input.loc_field <=cur_input.limit_field  then {current line not yet finished}
  begin cur_chr:=buffer[cur_input.loc_field ]; incr(cur_input.loc_field );
  reswitch: cur_cmd:= eqtb[  cat_code_base+   cur_chr].hh.rh   ;
  
{ Change state if necessary, and |goto switch| if the current character should be ignored, or |goto reswitch| if the current character changes to another }
case cur_input.state_field +cur_cmd of

{ Cases where character is ignored }
 mid_line+ ignore,skip_blanks+ ignore,new_line+ ignore ,skip_blanks+spacer,new_line+spacer

: goto switch;
 mid_line+ escape,skip_blanks+ escape,new_line+ escape : 
{ Scan a control sequence and set |state:=skip_blanks| or |mid_line| }
begin if cur_input.loc_field >cur_input.limit_field  then cur_cs:=null_cs {|state| is irrelevant in this case}
else  begin start_cs: k:=cur_input.loc_field ; cur_chr:=buffer[k]; cat:= eqtb[  cat_code_base+   cur_chr].hh.rh   ;
  incr(k);
  if cat=letter then cur_input.state_field :=skip_blanks
  else if cat=spacer then cur_input.state_field :=skip_blanks
  else cur_input.state_field :=mid_line;
  if (cat=letter)and(k<=cur_input.limit_field ) then
    
{ Scan ahead in the buffer until finding a nonletter; if an expanded code is encountered, reduce it and |goto start_cs|; otherwise if a multiletter control sequence is found, adjust |cur_cs| and |loc|, and |goto found| }
begin repeat cur_chr:=buffer[k]; cat:= eqtb[  cat_code_base+   cur_chr].hh.rh   ; incr(k);
until (cat<>letter)or(k>cur_input.limit_field );

{ If an expanded... }
begin if buffer[k]=cur_chr then  if cat=sup_mark then  if k<cur_input.limit_field  then
  begin c:=buffer[k+1];  if c<{0200=}128 then {yes, one is indeed present}
    begin d:=2;
    if ((( c>={"0"=}48)and( c<={"9"=}57))or(( c>={"a"=}97)and( c<={"f"=}102)))  then  if k+2<=cur_input.limit_field  then
      begin cc:=buffer[k+2];  if ((( cc>={"0"=}48)and( cc<={"9"=}57))or(( cc>={"a"=}97)and( cc<={"f"=}102)))  then incr(d);
      end;
    if d>2 then
      begin  if c<={"9"=}57 then cur_chr:=c-{"0"=}48 else cur_chr:=c-{"a"=}97+10; if cc<={"9"=}57 then cur_chr:=16*cur_chr+cc-{"0"=}48 else cur_chr:=16*cur_chr+cc-{"a"=}97+10 ; buffer[k-1]:=cur_chr;
      end
    else if c<{0100=}64 then buffer[k-1]:=c+{0100=}64
    else buffer[k-1]:=c-{0100=}64;
    cur_input.limit_field :=cur_input.limit_field -d; first:=first-d;
    while k<=cur_input.limit_field  do
      begin buffer[k]:=buffer[k+d]; incr(k);
      end;
    goto start_cs;
    end;
  end;
end

;
if cat<>letter then decr(k);
  {now |k| points to first nonletter}
if k>cur_input.loc_field +1 then {multiletter control sequence has been scanned}
  begin cur_cs:=id_lookup(cur_input.loc_field ,k-cur_input.loc_field ); cur_input.loc_field :=k; goto found;
  end;
end


  else 
{ If an expanded code is present, reduce it and |goto start_cs| }
begin if buffer[k]=cur_chr then  if cat=sup_mark then  if k<cur_input.limit_field  then
  begin c:=buffer[k+1];  if c<{0200=}128 then {yes, one is indeed present}
    begin d:=2;
    if ((( c>={"0"=}48)and( c<={"9"=}57))or(( c>={"a"=}97)and( c<={"f"=}102)))  then  if k+2<=cur_input.limit_field  then
      begin cc:=buffer[k+2];  if ((( cc>={"0"=}48)and( cc<={"9"=}57))or(( cc>={"a"=}97)and( cc<={"f"=}102)))  then incr(d);
      end;
    if d>2 then
      begin  if c<={"9"=}57 then cur_chr:=c-{"0"=}48 else cur_chr:=c-{"a"=}97+10; if cc<={"9"=}57 then cur_chr:=16*cur_chr+cc-{"0"=}48 else cur_chr:=16*cur_chr+cc-{"a"=}97+10 ; buffer[k-1]:=cur_chr;
      end
    else if c<{0100=}64 then buffer[k-1]:=c+{0100=}64
    else buffer[k-1]:=c-{0100=}64;
    cur_input.limit_field :=cur_input.limit_field -d; first:=first-d;
    while k<=cur_input.limit_field  do
      begin buffer[k]:=buffer[k+d]; incr(k);
      end;
    goto start_cs;
    end;
  end;
end

;
  cur_cs:=single_base+buffer[cur_input.loc_field ]; incr(cur_input.loc_field );
  end;
found: cur_cmd:= eqtb[  cur_cs].hh.b0  ; cur_chr:= eqtb[  cur_cs].hh.rh  ;
if cur_cmd>=outer_call then check_outer_validity;
end

;
 mid_line+ active_char,skip_blanks+ active_char,new_line+ active_char : 
{ Process an active-character control sequence and set |state:=mid_line| }
begin cur_cs:=cur_chr+active_base;
cur_cmd:= eqtb[  cur_cs].hh.b0  ; cur_chr:= eqtb[  cur_cs].hh.rh  ; cur_input.state_field :=mid_line;
if cur_cmd>=outer_call then check_outer_validity;
end

;
 mid_line+ sup_mark,skip_blanks+ sup_mark,new_line+ sup_mark : 
{ If this |sup_mark| starts an expanded character like~\.[\^\^A] or~\.[\^\^df], then |goto reswitch|, otherwise set |state:=mid_line| }
begin if cur_chr=buffer[cur_input.loc_field ] then if cur_input.loc_field <cur_input.limit_field  then
  begin c:=buffer[cur_input.loc_field +1];  if c<{0200=}128 then {yes we have an expanded char}
    begin cur_input.loc_field :=cur_input.loc_field +2;
    if ((( c>={"0"=}48)and( c<={"9"=}57))or(( c>={"a"=}97)and( c<={"f"=}102)))  then if cur_input.loc_field <=cur_input.limit_field  then
      begin cc:=buffer[cur_input.loc_field ];  if ((( cc>={"0"=}48)and( cc<={"9"=}57))or(( cc>={"a"=}97)and( cc<={"f"=}102)))  then
        begin incr(cur_input.loc_field );  if c<={"9"=}57 then cur_chr:=c-{"0"=}48 else cur_chr:=c-{"a"=}97+10; if cc<={"9"=}57 then cur_chr:=16*cur_chr+cc-{"0"=}48 else cur_chr:=16*cur_chr+cc-{"a"=}97+10 ; goto reswitch;
        end;
      end;
    if c<{0100=}64 then cur_chr:=c+{0100=}64  else cur_chr:=c-{0100=}64;
    goto reswitch;
    end;
  end;
cur_input.state_field :=mid_line;
end

;
 mid_line+ invalid_char,skip_blanks+ invalid_char,new_line+ invalid_char : 
{ Decry the invalid character and |goto restart| }
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Text line contains an invalid character"=} 622); end ;
{ \xref[Text line contains...] }
 begin help_ptr:=2; help_line[1]:={"A funny symbol that I can't read has just been input."=} 623; help_line[0]:={"Continue, and I'll forget that it ever happened."=} 624; end ;

deletions_allowed:=false; error; deletions_allowed:=true;
goto restart;
end

;
{ \4 }
{ Handle situations involving spaces, braces, changes of state }
mid_line+spacer:
{ Enter |skip_blanks| state, emit a space }
begin cur_input.state_field :=skip_blanks; cur_chr:={" "=}32;
end

;
mid_line+car_ret:
{ Finish line, emit a space }
begin cur_input.loc_field :=cur_input.limit_field +1; cur_cmd:=spacer; cur_chr:={" "=}32;
end

;
skip_blanks+car_ret, mid_line+ comment,skip_blanks+ comment,new_line+ comment :
  
{ Finish line, |goto switch| }
begin cur_input.loc_field :=cur_input.limit_field +1; goto switch;
end

;
new_line+car_ret:
{ Finish line, emit a \.[\\par] }
begin cur_input.loc_field :=cur_input.limit_field +1; cur_cs:=par_loc; cur_cmd:= eqtb[  cur_cs].hh.b0  ;
cur_chr:= eqtb[  cur_cs].hh.rh  ;
if cur_cmd>=outer_call then check_outer_validity;
end

;
mid_line+left_brace: incr(align_state);
skip_blanks+left_brace,new_line+left_brace: begin
  cur_input.state_field :=mid_line; incr(align_state);
  end;
mid_line+right_brace: decr(align_state);
skip_blanks+right_brace,new_line+right_brace: begin
  cur_input.state_field :=mid_line; decr(align_state);
  end;
 skip_blanks+math_shift, skip_blanks+tab_mark, skip_blanks+mac_param,  skip_blanks+sub_mark, skip_blanks+letter, skip_blanks+other_char , new_line+math_shift, new_line+tab_mark, new_line+mac_param,  new_line+sub_mark, new_line+letter, new_line+other_char : cur_input.state_field :=mid_line;

 
 else   
 end 

;
  end
else  begin cur_input.state_field :=new_line;

  
{ Move to next line of file, or |goto restart| if there is no next line, or |return| if a \.[\\read] line has finished }
if cur_input.name_field >17 then 
{ Read next line of file into |buffer|, or |goto restart| if the file has ended }
begin incr(line); first:=cur_input.start_field ;
if not force_eof then
  begin if input_ln(input_file[cur_input.index_field ] ,true) then {not end of file}
    firm_up_the_line {this sets |limit|}
  else force_eof:=true;
  end;
if force_eof then
  begin print_char({")"=}41); decr(open_parens);
   fflush (stdout ) ; {show user that file has been read}
  force_eof:=false;
  end_file_reading; {resume previous level}
  check_outer_validity; goto restart;
  end;
if  (eqtb[int_base+ end_line_char_code].int  <0)or(eqtb[int_base+ end_line_char_code].int  >255)  then decr(cur_input.limit_field )
else  buffer[cur_input.limit_field ]:=eqtb[int_base+ end_line_char_code].int  ;
first:=cur_input.limit_field +1; cur_input.loc_field :=cur_input.start_field ; {ready to read}
end


else  begin if not (cur_input.name_field =0)  then {\.[\\read] line has ended}
    begin cur_cmd:=0; cur_chr:=0;  goto exit ;
    end;
  if input_ptr>0 then {text was inserted during error recovery}
    begin end_file_reading; goto restart; {resume previous level}
    end;
  if selector<log_only then open_log_file;
  if interaction>nonstop_mode then
    begin if  (eqtb[int_base+ end_line_char_code].int  <0)or(eqtb[int_base+ end_line_char_code].int  >255)  then incr(cur_input.limit_field );
    if cur_input.limit_field =cur_input.start_field  then {previous line was empty}
      print_nl({"(Please type a command or say `\end')"=}625);
{ \xref[Please type...] }
    print_ln; first:=cur_input.start_field ;
    begin    ; print({"*"=} 42); term_input; end ; {input on-line into |buffer|}
{ \xref[*\relax] }
    cur_input.limit_field :=last;
    if  (eqtb[int_base+ end_line_char_code].int  <0)or(eqtb[int_base+ end_line_char_code].int  >255)  then decr(cur_input.limit_field )
    else  buffer[cur_input.limit_field ]:=eqtb[int_base+ end_line_char_code].int  ;
    first:=cur_input.limit_field +1;
    cur_input.loc_field :=cur_input.start_field ;
    end
  else fatal_error({"*** (job aborted, no legal \end found)"=}626);
{ \xref[job aborted] }
    {nonstop mode, which is intended for overnight batch processing,
    never waits for on-line input}
  end

;
  begin if interrupt<>0 then pause_for_instructions; end ;
  goto switch;
  end;
end


else 
{ Input from token list, |goto restart| if end of list or if a parameter needs to be expanded }
if cur_input.loc_field <>-{0xfffffff=}268435455   then {list not exhausted}
{ \xref[inner loop] }
  begin t:= mem[ cur_input.loc_field ].hh.lh ; cur_input.loc_field := mem[ cur_input.loc_field ].hh.rh ; {move to next}
  if t>={07777=}4095  then {a control sequence token}
    begin cur_cs:=t-{07777=}4095 ;
    cur_cmd:= eqtb[  cur_cs].hh.b0  ; cur_chr:= eqtb[  cur_cs].hh.rh  ;
    if cur_cmd>=outer_call then
      if cur_cmd=dont_expand then
        
{ Get the next token, suppressing expansion }
begin cur_cs:= mem[ cur_input.loc_field ].hh.lh -{07777=}4095 ; cur_input.loc_field :=-{0xfffffff=}268435455  ;

cur_cmd:= eqtb[  cur_cs].hh.b0  ; cur_chr:= eqtb[  cur_cs].hh.rh  ;
if cur_cmd>max_command then
  begin cur_cmd:=relax; cur_chr:=no_expand_flag;
  end;
end


      else check_outer_validity;
    end
  else  begin cur_cmd:=t div {0400=}256; cur_chr:=t mod {0400=}256;
    case cur_cmd of
    left_brace: incr(align_state);
    right_brace: decr(align_state);
    out_param: 
{ Insert macro parameter and |goto restart| }
begin begin_token_list(param_stack[cur_input.limit_field  +cur_chr-1],parameter);
goto restart;
end

;
     else   
     end ;
    end;
  end
else  begin {we are done with this token list}
  end_token_list; goto restart; {resume previous level}
  end

;

{ If an alignment entry has just ended, take appropriate action }
if cur_cmd<=car_ret then if cur_cmd>=tab_mark then if align_state=0 then
  
{ Insert the \(v)\<v_j> template and |goto restart| }
begin if (scanner_status=aligning) or (cur_align=-{0xfffffff=}268435455  ) then
  fatal_error({"(interwoven alignment preambles are not allowed)"=}602);
{ \xref[interwoven alignment preambles...] }
cur_cmd:= mem[  cur_align+ list_offset].hh.lh  ;  mem[  cur_align+ list_offset].hh.lh  :=cur_chr;
if cur_cmd=omit then begin_token_list(mem_top-10 ,v_template)
else begin_token_list(mem[ cur_align+depth_offset].int ,v_template);
align_state:=1000000; goto restart;
end



;
exit:end;



{ 363. }

{tangle:pos tex.web:7583:1: }

{ If the user has set the |pausing| parameter to some positive value,
and if nonstop mode has not been selected, each line of input is displayed
on the terminal and the transcript file, followed by `\.[=>]'.
\TeX\ waits for a response. If the response is simply |carriage_return|, the
line is accepted as it stands, otherwise the line typed is
used instead of the line in the file. } procedure firm_up_the_line;
var k:0..buf_size; {an index into |buffer|}
begin cur_input.limit_field :=last;
if eqtb[int_base+ pausing_code].int  >0 then if interaction>nonstop_mode then
  begin    ; print_ln;
  if cur_input.start_field <cur_input.limit_field  then for k:=cur_input.start_field  to cur_input.limit_field -1 do print(buffer[k]);
  first:=cur_input.limit_field ; begin    ; print({"=>"=} 627); term_input; end ; {wait for user response}
{ \xref[=>] }
  if last>first then
    begin for k:=first to last-1 do {move line down in buffer}
      buffer[k+cur_input.start_field -first]:=buffer[k];
    cur_input.limit_field :=cur_input.start_field +last-first;
    end;
  end;
end;



{ 364. }

{tangle:pos tex.web:7606:1: }

{ Since |get_next| is used so frequently in \TeX, it is convenient
to define three related procedures that do a little more:

\yskip\hang|get_token| not only sets |cur_cmd| and |cur_chr|, it
also sets |cur_tok|, a packed halfword version of the current token.

\yskip\hang|get_x_token|, meaning ``get an expanded token,'' is like
|get_token|, but if the current token turns out to be a user-defined
control sequence (i.e., a macro call), or a conditional,
or something like \.[\\topmark] or \.[\\expandafter] or \.[\\csname],
it is eliminated from the input by beginning the expansion of the macro
or the evaluation of the conditional.

\yskip\hang|x_token| is like |get_x_token| except that it assumes that
|get_next| has already been called.

\yskip\noindent
In fact, these three procedures account for almost every use of |get_next|. }

{ 365. }

{tangle:pos tex.web:7625:1: }

{ No new control sequences will be defined except during a call of
|get_token|, or when \.[\\csname] compresses a token list, because
|no_new_control_sequence| is always |true| at other times. } procedure get_token; {sets |cur_cmd|, |cur_chr|, |cur_tok|}
begin no_new_control_sequence:=false; get_next; no_new_control_sequence:=true;
{ \xref[inner loop] }
if cur_cs=0 then cur_tok:=(cur_cmd*{0400=}256)+cur_chr
else cur_tok:={07777=}4095 +cur_cs;
end;



{ 366. \[25] Expanding the next token }

{tangle:pos tex.web:7636:33: }

{ Only a dozen or so command codes |>max_command| can possibly be returned by
|get_next|; in increasing order, they are |undefined_cs|, |expand_after|,
|no_expand|, |input|, |if_test|, |fi_or_else|, |cs_name|, |convert|, |the|,
|top_bot_mark|, |call|, |long_call|, |outer_call|, |long_outer_call|, and
|end_template|.[\emergencystretch=40pt\par]

The |expand| subroutine is used when |cur_cmd>max_command|. It removes a
``call'' or a conditional or one of the other special operations just
listed.  It follows that |expand| might invoke itself recursively. In all
cases, |expand| destroys the current token, but it sets things up so that
the next |get_next| will deliver the appropriate next token. The value of
|cur_tok| need not be known when |expand| is called.

Since several of the basic scanning routines communicate via global variables,
their values are saved as local variables of |expand| so that
recursive calls don't invalidate them.
\xref[recursion] }{ \4 }
{ Declare the procedure called |macro_call| }
procedure macro_call; {invokes a user-defined control sequence}
label exit, continue, done, done1, found;
var r:halfword ; {current node in the macro's token list}
 p:halfword ; {current node in parameter token list being built}
 q:halfword ; {new node being put into the token list}
 s:halfword ; {backup pointer for parameter matching}
 t:halfword ; {cycle pointer for backup recovery}
 u, v:halfword ; {auxiliary pointers for backup recovery}
 rbrace_ptr:halfword ; {one step before the last |right_brace| token}
 n:small_number; {the number of parameters scanned}
 unbalance:halfword; {unmatched left braces in current parameter}
 m:halfword; {the number of tokens or groups (usually)}
 ref_count:halfword ; {start of the token list}
 save_scanner_status:small_number; {|scanner_status| upon entry}
 save_warning_index:halfword ; {|warning_index| upon entry}
 match_chr:ASCII_code; {character used in parameter}
begin save_scanner_status:=scanner_status; save_warning_index:=warning_index;
warning_index:=cur_cs; ref_count:=cur_chr; r:= mem[ ref_count].hh.rh ; n:=0;
if eqtb[int_base+ tracing_macros_code].int  >0 then 
{ Show the text of the macro being expanded }
begin begin_diagnostic; print_ln; print_cs(warning_index);
token_show(ref_count); end_diagnostic(false);
end

;
if  mem[ r].hh.lh <>end_match_token then
  
{ Scan the parameters and make |link(r)| point to the macro body; but |return| if an illegal \.[\\par] is detected }
begin scanner_status:=matching; unbalance:=0;
long_state:= eqtb[  cur_cs].hh.b0  ;
if long_state>=outer_call then long_state:=long_state-2;
repeat  mem[ mem_top-3 ].hh.rh :=-{0xfffffff=}268435455  ;
if ( mem[ r].hh.lh >match_token+255)or( mem[ r].hh.lh <match_token) then s:=-{0xfffffff=}268435455  
else  begin match_chr:= mem[ r].hh.lh -match_token; s:= mem[ r].hh.rh ; r:=s;
  p:=mem_top-3 ; m:=0;
  end;

{ Scan a parameter until its delimiter string has been found; or, if |s=null|, simply scan the delimiter string }
continue: get_token; {set |cur_tok| to the next token of input}
if cur_tok= mem[ r].hh.lh  then
  
{ Advance \(r)|r|; |goto found| if the parameter delimiter has been fully matched, otherwise |goto continue| }
begin r:= mem[ r].hh.rh ;
if ( mem[ r].hh.lh >=match_token)and( mem[ r].hh.lh <=end_match_token) then
  begin if cur_tok<left_brace_limit then decr(align_state);
  goto found;
  end
else goto continue;
end

;

{ Contribute the recently matched tokens to the current parameter, and |goto continue| if a partial match is still in effect; but abort if |s=null| }
if s<>r then
  if s=-{0xfffffff=}268435455   then 
{ Report an improper use of the macro and abort }
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Use of "=} 660); end ; sprint_cs(warning_index);
{ \xref[Use of x doesn't match...] }
print({" doesn't match its definition"=}661);
 begin help_ptr:=4; help_line[3]:={"If you say, e.g., `\def\a1[...]', then you must always"=} 662; help_line[2]:={"put `1' after `\a', since control sequence names are"=} 663; help_line[1]:={"made up of letters only. The macro here has not been"=} 664; help_line[0]:={"followed by the required stuff, so I'm ignoring it."=} 665; end ;
error;  goto exit ;
end


  else  begin t:=s;
    repeat begin q:=get_avail;  mem[ p].hh.rh :=q;  mem[ q].hh.lh :=  mem[  t].hh.lh ; p:=q; end ; incr(m); u:= mem[ t].hh.rh ; v:=s;
     while true do    begin if u=r then
        if cur_tok<> mem[ v].hh.lh  then goto done
        else  begin r:= mem[ v].hh.rh ; goto continue;
          end;
      if  mem[ u].hh.lh <> mem[ v].hh.lh  then goto done;
      u:= mem[ u].hh.rh ; v:= mem[ v].hh.rh ;
      end;
    done: t:= mem[ t].hh.rh ;
    until t=r;
    r:=s; {at this point, no tokens are recently matched}
    end

;
if cur_tok=par_token then if long_state<>long_call then
  
{ Report a runaway argument and abort }
begin if long_state=call then
  begin runaway; begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Paragraph ended before "=} 655); end ;
{ \xref[Paragraph ended before...] }
  sprint_cs(warning_index); print({" was complete"=}656);
   begin help_ptr:=3; help_line[2]:={"I suspect you've forgotten a `]', causing me to apply this"=} 657; help_line[1]:={"control sequence to too much text. How can we recover?"=} 658; help_line[0]:={"My plan is to forget the whole thing and hope for the best."=} 659; end ;
  back_error;
  end;
pstack[n]:= mem[ mem_top-3 ].hh.rh ; align_state:=align_state-unbalance;
for m:=0 to n do flush_list(pstack[m]);
 goto exit ;
end

;
if cur_tok<right_brace_limit then
  if cur_tok<left_brace_limit then
    
{ Contribute an entire group to the current parameter }
begin unbalance:=1;
{ \xref[inner loop] }
 while true do    begin begin {  } begin  q:=avail; if  q=-{0xfffffff=}268435455   then  q:=get_avail else begin avail:= mem[  q].hh.rh ;  mem[  q].hh.rh :=-{0xfffffff=}268435455  ; ifdef('STAT')  incr(dyn_used); endif('STAT')  end; end ;  mem[ p].hh.rh :=q;  mem[ q].hh.lh := cur_tok; p:=q; end ; get_token;
  if cur_tok=par_token then if long_state<>long_call then
    
{ Report a runaway argument and abort }
begin if long_state=call then
  begin runaway; begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Paragraph ended before "=} 655); end ;
{ \xref[Paragraph ended before...] }
  sprint_cs(warning_index); print({" was complete"=}656);
   begin help_ptr:=3; help_line[2]:={"I suspect you've forgotten a `]', causing me to apply this"=} 657; help_line[1]:={"control sequence to too much text. How can we recover?"=} 658; help_line[0]:={"My plan is to forget the whole thing and hope for the best."=} 659; end ;
  back_error;
  end;
pstack[n]:= mem[ mem_top-3 ].hh.rh ; align_state:=align_state-unbalance;
for m:=0 to n do flush_list(pstack[m]);
 goto exit ;
end

;
  if cur_tok<right_brace_limit then
    if cur_tok<left_brace_limit then incr(unbalance)
    else  begin decr(unbalance);
      if unbalance=0 then goto done1;
      end;
  end;
done1: rbrace_ptr:=p; begin q:=get_avail;  mem[ p].hh.rh :=q;  mem[ q].hh.lh := cur_tok; p:=q; end ;
end


  else 
{ Report an extra right brace and |goto continue| }
begin back_input; begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Argument of "=} 647); end ; sprint_cs(warning_index);
{ \xref[Argument of \\x has...] }
print({" has an extra ]"=}648);
 begin help_ptr:=6; help_line[5]:={"I've run across a `]' that doesn't seem to match anything."=} 649; help_line[4]:={"For example, `\def\a#1[...]' and `\a]' would produce"=} 650; help_line[3]:={"this error. If you simply proceed now, the `\par' that"=} 651; help_line[2]:={"I've just inserted will cause me to report a runaway"=} 652; help_line[1]:={"argument that might be the root of the problem. But if"=} 653; help_line[0]:={"your `]' was spurious, just type `2' and it will go away."=} 654; end ;
incr(align_state); long_state:=call; cur_tok:=par_token; ins_error;
goto continue;
end {a white lie; the \.[\\par] won't always trigger a runaway}


else 
{ Store the current token, but |goto continue| if it is a blank space that would become an undelimited parameter }
begin if cur_tok=space_token then
  if  mem[ r].hh.lh <=end_match_token then
    if  mem[ r].hh.lh >=match_token then goto continue;
begin q:=get_avail;  mem[ p].hh.rh :=q;  mem[ q].hh.lh := cur_tok; p:=q; end ;
end

;
incr(m);
if  mem[ r].hh.lh >end_match_token then goto continue;
if  mem[ r].hh.lh <match_token then goto continue;
found: if s<>-{0xfffffff=}268435455   then 
{ Tidy up the parameter just scanned, and tuck it away }
begin if (m=1)and( mem[ p].hh.lh <right_brace_limit) then
  begin  mem[ rbrace_ptr].hh.rh :=-{0xfffffff=}268435455  ;  begin  mem[  p].hh.rh :=avail; avail:= p; ifdef('STAT')  decr(dyn_used); endif('STAT')  end ;
  p:= mem[ mem_top-3 ].hh.rh ; pstack[n]:= mem[ p].hh.rh ;  begin  mem[  p].hh.rh :=avail; avail:= p; ifdef('STAT')  decr(dyn_used); endif('STAT')  end ;
  end
else pstack[n]:= mem[ mem_top-3 ].hh.rh ;
incr(n);
if eqtb[int_base+ tracing_macros_code].int  >0 then
  begin begin_diagnostic; print_nl(match_chr); print_int(n);
  print({"<-"=}666); show_token_list(pstack[n-1],-{0xfffffff=}268435455  ,1000);
  end_diagnostic(false);
  end;
end



;

{now |info(r)| is a token whose command code is either |match| or |end_match|}
until  mem[ r].hh.lh =end_match_token;
end

;

{ Feed the macro body and its parameters to the scanner }
while (cur_input.state_field =token_list)and(cur_input.loc_field =-{0xfffffff=}268435455  )and(cur_input.index_field  <>v_template) do
  end_token_list; {conserve stack space}
begin_token_list(ref_count,macro); cur_input.name_field :=warning_index; cur_input.loc_field := mem[ r].hh.rh ;
if n>0 then
  begin if param_ptr+n>max_param_stack then
    begin max_param_stack:=param_ptr+n;
    if max_param_stack>param_size then
      overflow({"parameter stack size"=}646,param_size);
{ \xref[TeX capacity exceeded parameter stack size][\quad parameter stack size] }
    end;
  for m:=0 to n-1 do param_stack[param_ptr+m]:=pstack[m];
  param_ptr:=param_ptr+n;
  end

;
exit:scanner_status:=save_scanner_status; warning_index:=save_warning_index;
end;

 

{ \4 }
{ Declare the procedure called |insert_relax| }
procedure insert_relax;
begin cur_tok:={07777=}4095 +cur_cs; back_input;
cur_tok:={07777=}4095 +frozen_relax; back_input; cur_input.index_field  :=inserted;
end;

 

procedure pass_text; forward;{ \2 }
procedure start_input; forward;{ \2 }
procedure conditional; forward;{ \2 }
procedure get_x_token; forward;{ \2 }
procedure conv_toks; forward;{ \2 }
procedure ins_the_toks; forward;{ \2 }
procedure expand;
var t:halfword; {token that is being ``expanded after''}
 p, q, r:halfword ; {for list manipulation}
 j:0..buf_size; {index into |buffer|}
 cv_backup:integer; {to save the global quantity |cur_val|}
 cvl_backup, radix_backup, co_backup:small_number;
  {to save |cur_val_level|, etc.}
 backup_backup:halfword ; {to save |link(backup_head)|}
 save_scanner_status:small_number; {temporary storage of |scanner_status|}
begin
incr(expand_depth_count);
if expand_depth_count>=expand_depth then overflow({"expansion depth"=}628,expand_depth);
cv_backup:=cur_val; cvl_backup:=cur_val_level; radix_backup:=radix;
co_backup:=cur_order; backup_backup:= mem[ mem_top-13 ].hh.rh ;
if cur_cmd<call then 
{ Expand a nonmacro }
begin if eqtb[int_base+ tracing_commands_code].int  >1 then show_cur_cmd_chr;
case cur_cmd of
top_bot_mark:
{ Insert the \(a)appropriate mark text into the scanner }
begin if cur_mark[cur_chr]<>-{0xfffffff=}268435455   then
  begin_token_list(cur_mark[cur_chr],mark_text);
end

;
expand_after:
{ Expand the token after the next token }
begin get_token; t:=cur_tok; get_token;
if cur_cmd>max_command then expand else back_input;
cur_tok:=t; back_input;
end

;
no_expand:
{ Suppress expansion of the next token }
begin save_scanner_status:=scanner_status; scanner_status:=normal;
get_token; scanner_status:=save_scanner_status; t:=cur_tok;
back_input; {now |start| and |loc| point to the backed-up token |t|}
if (t>={07777=}4095 )and(t<>{07777=}4095 +end_write ) then
  begin p:=get_avail;  mem[ p].hh.lh :={07777=}4095 +frozen_dont_expand;
   mem[ p].hh.rh :=cur_input.loc_field ; cur_input.start_field :=p; cur_input.loc_field :=p;
  end;
end

;
cs_name:
{ Manufacture a control sequence name }
begin r:=get_avail; p:=r; {head of the list of characters}
repeat get_x_token;
if cur_cs=0 then begin q:=get_avail;  mem[ p].hh.rh :=q;  mem[ q].hh.lh := cur_tok; p:=q; end ;
until cur_cs<>0;
if cur_cmd<>end_cs_name then 
{ Complain about missing \.[\\endcsname] }
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Missing "=} 635); end ; print_esc({"endcsname"=}513); print({" inserted"=}636);
{ \xref[Missing \\endcsname...] }
 begin help_ptr:=2; help_line[1]:={"The control sequence marked <to be read again> should"=} 637; help_line[0]:={"not appear between \csname and \endcsname."=} 638; end ;
back_error;
end

;

{ Look up the characters of list |r| in the hash table, and set |cur_cs| }
j:=first; p:= mem[ r].hh.rh ;
while p<>-{0xfffffff=}268435455   do
  begin if j>=max_buf_stack then
    begin max_buf_stack:=j+1;
    if max_buf_stack=buf_size then
      overflow({"buffer size"=}256,buf_size);
{ \xref[TeX capacity exceeded buffer size][\quad buffer size] }
    end;
  buffer[j]:= mem[ p].hh.lh  mod {0400=}256; incr(j); p:= mem[ p].hh.rh ;
  end;
if j>first+1 then
  begin no_new_control_sequence:=false; cur_cs:=id_lookup(first,j-first);
  no_new_control_sequence:=true;
  end
else if j=first then cur_cs:=null_cs {the list is empty}
else cur_cs:=single_base+buffer[first] {the list has length one}

;
flush_list(r);
if  eqtb[  cur_cs].hh.b0  =undefined_cs then
  begin eq_define(cur_cs,relax,256); {N.B.: The |save_stack| might change}
  end; {the control sequence will now match `\.[\\relax]'}
cur_tok:=cur_cs+{07777=}4095 ; back_input;
end

;
convert:conv_toks; {this procedure is discussed in Part 27 below}
the:ins_the_toks; {this procedure is discussed in Part 27 below}
if_test:conditional; {this procedure is discussed in Part 28 below}
fi_or_else:
{ Terminate the current conditional and skip to \.[\\fi] }
if cur_chr>if_limit then
  if if_limit=if_code then insert_relax {condition not yet evaluated}
  else  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Extra "=} 787); end ; print_cmd_chr(fi_or_else,cur_chr);
{ \xref[Extra \\or] }
{ \xref[Extra \\else] }
{ \xref[Extra \\fi] }
     begin help_ptr:=1; help_line[0]:={"I'm ignoring this; it doesn't match any \if."=} 788; end ;
    error;
    end
else  begin while cur_chr<>fi_code do pass_text; {skip to \.[\\fi]}
  
{ Pop the condition stack }
begin p:=cond_ptr; if_line:=mem[ p+1].int ;
cur_if:= mem[ p].hh.b1 ; if_limit:= mem[ p].hh.b0 ; cond_ptr:= mem[ p].hh.rh ;
free_node(p,if_node_size);
end

;
  end

;
input:
{ Initiate or terminate input from a file }
if cur_chr>0 then force_eof:=true
else if name_in_progress then insert_relax
else start_input

;
 else  
{ Complain about an undefined macro }
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Undefined control sequence"=} 629); end ;
{ \xref[Undefined control sequence] }
 begin help_ptr:=5; help_line[4]:={"The control sequence at the end of the top line"=} 630; help_line[3]:={"of your error message was never \def'ed. If you have"=} 631; help_line[2]:={"misspelled it (e.g., `\hobx'), type `I' and the correct"=} 632; help_line[1]:={"spelling (e.g., `I\hbox'). Otherwise just continue,"=} 633; help_line[0]:={"and I'll forget about whatever was undefined."=} 634; end ;
error;
end


 end ;
end


else if cur_cmd<end_template then macro_call
else 
{ Insert a token containing |frozen_endv| }
begin cur_tok:={07777=}4095 +frozen_endv; back_input;
end

;
cur_val:=cv_backup; cur_val_level:=cvl_backup; radix:=radix_backup;
cur_order:=co_backup;  mem[ mem_top-13 ].hh.rh :=backup_backup;
decr(expand_depth_count);
end;



{ 371. }

{tangle:pos tex.web:7734:1: }

{ The |expand| procedure and some other routines that construct token
lists find it convenient to use the following macros, which are valid only if
the variables |p| and |q| are reserved for token-list building. }

{ 380. }

{tangle:pos tex.web:7821:1: }

{ Here is a recursive procedure that is \TeX's usual way to get the
next token of input. It has been slightly optimized to take account of
common cases. } procedure get_x_token; {sets |cur_cmd|, |cur_chr|, |cur_tok|,
  and expands macros}
label restart,done;
begin restart: get_next;
{ \xref[inner loop] }
if cur_cmd<=max_command then goto done;
if cur_cmd>=call then
  if cur_cmd<end_template then macro_call
  else  begin cur_cs:=frozen_endv; cur_cmd:=endv;
    goto done; {|cur_chr=null_list|}
    end
else expand;
goto restart;
done: if cur_cs=0 then cur_tok:=(cur_cmd*{0400=}256)+cur_chr
else cur_tok:={07777=}4095 +cur_cs;
end;



{ 381. }

{tangle:pos tex.web:7842:1: }

{ The |get_x_token| procedure is essentially equivalent to two consecutive
procedure calls: |get_next; x_token|. } procedure x_token; {|get_x_token| without the initial |get_next|}
begin while cur_cmd>max_command do
  begin expand;
  get_next;
  end;
if cur_cs=0 then cur_tok:=(cur_cmd*{0400=}256)+cur_chr
else cur_tok:={07777=}4095 +cur_cs;
end;



{ 402. \[26] Basic scanning subroutines }

{tangle:pos tex.web:8181:35: }

{ Let's turn now to some procedures that \TeX\ calls upon frequently to digest
certain kinds of patterns in the input. Most of these are quite simple;
some are quite elaborate. Almost all of the routines call |get_x_token|,
which can cause them to be invoked recursively.
\xref[stomach]
\xref[recursion] }

{ 403. }

{tangle:pos tex.web:8189:1: }

{ The |scan_left_brace| routine is called when a left brace is supposed to be
the next non-blank token. (The term ``left brace'' means, more precisely,
a character whose catcode is |left_brace|.) \TeX\ allows \.[\\relax] to
appear before the |left_brace|. } procedure scan_left_brace; {reads a mandatory |left_brace|}
begin 
{ Get the next non-blank non-relax non-call token }
repeat get_x_token;
until (cur_cmd<>spacer)and(cur_cmd<>relax)

;
if cur_cmd<>left_brace then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Missing [ inserted"=} 667); end ;
{ \xref[Missing \[ inserted] }
   begin help_ptr:=4; help_line[3]:={"A left brace was mandatory here, so I've put one in."=} 668; help_line[2]:={"You might want to delete and/or insert some corrections"=} 669; help_line[1]:={"so that I will find a matching right brace soon."=} 670; help_line[0]:={"(If you're confused by all this, try typing `I]' now.)"=} 671; end ;
  back_error; cur_tok:=left_brace_token+{"["=}123; cur_cmd:=left_brace;
  cur_chr:={"["=}123; incr(align_state);
  end;
end;



{ 405. }

{tangle:pos tex.web:8212:1: }

{ The |scan_optional_equals| routine looks for an optional `\.=' sign preceded
by optional spaces; `\.[\\relax]' is not ignored here. } procedure scan_optional_equals;
begin  
{ Get the next non-blank non-call token }
repeat get_x_token;
until cur_cmd<>spacer

;
if cur_tok<>other_token+{"="=}61 then back_input;
end;



{ 407. }

{tangle:pos tex.web:8224:1: }

{ In case you are getting bored, here is a slightly less trivial routine:
Given a string of lowercase letters, like `\.[pt]' or `\.[plus]' or
`\.[width]', the |scan_keyword| routine checks to see whether the next
tokens of input match this string. The match must be exact, except that
uppercase letters will match their lowercase counterparts; uppercase
equivalents are determined by subtracting |"a"-"A"|, rather than using the
|uc_code| table, since \TeX\ uses this routine only for its own limited
set of keywords.

If a match is found, the characters are effectively removed from the input
and |true| is returned. Otherwise |false| is returned, and the input
is left essentially unchanged (except for the fact that some macros
may have been expanded, etc.).
\xref[inner loop] } function scan_keyword( s:str_number):boolean; {look for a given string}
label exit;
var p:halfword ; {tail of the backup list}
 q:halfword ; {new node being added to the token list via |store_new_token|}
 k:pool_pointer; {index into |str_pool|}
begin p:=mem_top-13 ;  mem[ p].hh.rh :=-{0xfffffff=}268435455  ; k:=str_start[s];
while k<str_start[s+1] do
  begin get_x_token; {recursion is possible here}
{ \xref[recursion] }
  if (cur_cs=0)and 
   ((cur_chr=  str_pool[ k] )or(cur_chr=  str_pool[ k] -{"a"=}97+{"A"=}65)) then
    begin begin q:=get_avail;  mem[ p].hh.rh :=q;  mem[ q].hh.lh := cur_tok; p:=q; end ; incr(k);
    end
  else if (cur_cmd<>spacer)or(p<>mem_top-13 ) then
    begin back_input;
    if p<>mem_top-13  then begin_token_list(  mem[  mem_top-13 ].hh.rh ,backed_up) ;
    scan_keyword:=false;  goto exit ;
    end;
  end;
flush_list( mem[ mem_top-13 ].hh.rh ); scan_keyword:=true;
exit:end;



{ 408. }

{tangle:pos tex.web:8261:1: }

{ Here is a procedure that sounds an alarm when mu and non-mu units
are being switched. } procedure mu_error;
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Incompatible glue units"=} 672); end ;
{ \xref[Incompatible glue units] }
 begin help_ptr:=1; help_line[0]:={"I'm going to assume that 1mu=1pt when they're mixed."=} 673; end ;
error;
end;



{ 409. }

{tangle:pos tex.web:8271:1: }

{ The next routine `|scan_something_internal|' is used to fetch internal
numeric quantities like `\.[\\hsize]', and also to handle the `\.[\\the]'
when expanding constructions like `\.[\\the\\toks0]' and
`\.[\\the\\baselineskip]'. Soon we will be considering the |scan_int|
procedure, which calls |scan_something_internal|; on the other hand,
|scan_something_internal| also calls |scan_int|, for constructions like
`\.[\\catcode\`\\\$]' or `\.[\\fontdimen] \.3 \.[\\ff]'. So we
have to declare |scan_int| as a |forward| procedure. A few other
procedures are also declared at this point. } procedure scan_int; forward; {scans an integer value}
{ \4\4 }
{ Declare procedures that scan restricted classes of integers }
procedure scan_eight_bit_int;
begin scan_int;
if (cur_val<0)or(cur_val>255) then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Bad register code"=} 697); end ;
{ \xref[Bad register code] }
   begin help_ptr:=2; help_line[1]:={"A register number must be between 0 and 255."=} 698; help_line[0]:={"I changed this one to zero."=} 699; end ; int_error(cur_val); cur_val:=0;
  end;
end;


procedure scan_char_num;
begin scan_int;
if (cur_val<0)or(cur_val>255) then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Bad character code"=} 700); end ;
{ \xref[Bad character code] }
   begin help_ptr:=2; help_line[1]:={"A character number must be between 0 and 255."=} 701; help_line[0]:={"I changed this one to zero."=} 699; end ; int_error(cur_val); cur_val:=0;
  end;
end;


procedure scan_four_bit_int;
begin scan_int;
if (cur_val<0)or(cur_val>15) then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Bad number"=} 702); end ;
{ \xref[Bad number] }
   begin help_ptr:=2; help_line[1]:={"Since I expected to read a number between 0 and 15,"=} 703; help_line[0]:={"I changed this one to zero."=} 699; end ; int_error(cur_val); cur_val:=0;
  end;
end;


procedure scan_fifteen_bit_int;
begin scan_int;
if (cur_val<0)or(cur_val>{077777=}32767) then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Bad mathchar"=} 704); end ;
{ \xref[Bad mathchar] }
   begin help_ptr:=2; help_line[1]:={"A mathchar number must be between 0 and 32767."=} 705; help_line[0]:={"I changed this one to zero."=} 699; end ; int_error(cur_val); cur_val:=0;
  end;
end;


procedure scan_twenty_seven_bit_int;
begin scan_int;
if (cur_val<0)or(cur_val>{0777777777=}134217727) then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Bad delimiter code"=} 706); end ;
{ \xref[Bad delimiter code] }
   begin help_ptr:=2; help_line[1]:={"A numeric delimiter code must be between 0 and 2^[27]-1."=} 707; help_line[0]:={"I changed this one to zero."=} 699; end ; int_error(cur_val); cur_val:=0;
  end;
end;


procedure scan_four_bit_int_or_18;
begin scan_int;
if (cur_val<0)or((cur_val>15)and(cur_val<>18)) then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Bad number"=} 702); end ;
{ \xref[Bad number] }
   begin help_ptr:=2; help_line[1]:={"Since I expected to read a number between 0 and 15,"=} 703; help_line[0]:={"I changed this one to zero."=} 699; end ; int_error(cur_val); cur_val:=0;
  end;
end;

 
{ \4\4 }
{ Declare procedures that scan font-related stuff }
procedure scan_font_ident;
var f:internal_font_number;
 m:halfword;
begin 
{ Get the next non-blank non-call... }
repeat get_x_token;
until cur_cmd<>spacer

;
if cur_cmd=def_font then f:= eqtb[  cur_font_loc].hh.rh   
else if cur_cmd=set_font then f:=cur_chr
else if cur_cmd=def_family then
  begin m:=cur_chr; scan_four_bit_int; f:= eqtb[  m+  cur_val].hh.rh  ;
  end
else  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Missing font identifier"=} 829); end ;
{ \xref[Missing font identifier] }
   begin help_ptr:=2; help_line[1]:={"I was looking for a control sequence whose"=} 830; help_line[0]:={"current meaning has been defined by \font."=} 831; end ;
  back_error; f:=font_base ;
  end;
cur_val:=f;
end;


procedure find_font_dimen( writing:boolean);
  {sets |cur_val| to |font_info| location}
var f:internal_font_number;
 n:integer; {the parameter number}
begin scan_int; n:=cur_val; scan_font_ident; f:=cur_val;
if n<=0 then cur_val:=fmem_ptr
else  begin if writing and(n<=space_shrink_code)and 
    (n>=space_code)and(font_glue[f]<>-{0xfffffff=}268435455  ) then
    begin delete_glue_ref(font_glue[f]);
    font_glue[f]:=-{0xfffffff=}268435455  ;
    end;
  if n>font_params[f] then
    if f<font_ptr then cur_val:=fmem_ptr
    else 
{ Increase the number of parameters in the last font }
begin repeat if fmem_ptr=font_mem_size then
  overflow({"font memory"=}836,font_mem_size);
{ \xref[TeX capacity exceeded font memory][\quad font memory] }
font_info[fmem_ptr].int :=0; incr(fmem_ptr); incr(font_params[f]);
until n=font_params[f];
cur_val:=fmem_ptr-1; {this equals |param_base[f]+font_params[f]|}
end


  else cur_val:=n+param_base[f];
  end;

{ Issue an error message if |cur_val=fmem_ptr| }
if cur_val=fmem_ptr then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Font "=} 812); end ; print_esc(  hash[ font_id_base+  f].rh  );
  print({" has only "=}832); print_int(font_params[f]);
  print({" fontdimen parameters"=}833);
{ \xref[Font x has only...] }
   begin help_ptr:=2; help_line[1]:={"To increase the number of font parameters, you must"=} 834; help_line[0]:={"use \fontdimen immediately after the \font is loaded."=} 835; end ;
  error;
  end

;
end;





{ 413. }

{tangle:pos tex.web:8345:1: }

{ OK, we're ready for |scan_something_internal| itself. A second parameter,
|negative|, is set |true| if the value that is found should be negated.
It is assumed that |cur_cmd| and |cur_chr| represent the first token of
the internal quantity to be scanned; an error will be signalled if
|cur_cmd<min_internal| or |cur_cmd>max_internal|. } procedure scan_something_internal( level:small_number; negative:boolean);
  {fetch an internal parameter}
var m:halfword; {|chr_code| part of the operand token}
 p:0..nest_size; {index into |nest|}
begin m:=cur_chr;
case cur_cmd of
def_code: 
{ Fetch a character code from some table }
begin scan_char_num;
if m=math_code_base then  begin cur_val:=    eqtb[  math_code_base+     cur_val].hh.rh    ;cur_val_level:= int_val; end 
else if m<math_code_base then  begin cur_val:=  eqtb[   m+   cur_val].hh.rh  ;cur_val_level:= int_val; end 
else  begin cur_val:= eqtb[ m+ cur_val]. int;cur_val_level:= int_val; end ;
end

;
toks_register,assign_toks,def_family,set_font,def_font: 
{ Fetch a token list or font identifier, provided that |level=tok_val| }
if level<>tok_val then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Missing number, treated as zero"=} 674); end ;
{ \xref[Missing number...] }
   begin help_ptr:=3; help_line[2]:={"A number should have been here; I inserted `0'."=} 675; help_line[1]:={"(If you can't figure out why I needed to see a number,"=} 676; help_line[0]:={"look up `weird error' in the index to The TeXbook.)"=} 677; end ;
{ \xref[TeXbook][\sl The \TeX book] }
  back_error;  begin cur_val:= 0;cur_val_level:= dimen_val; end ;
  end
else if cur_cmd<=assign_toks then
  begin if cur_cmd<assign_toks then {|cur_cmd=toks_register|}
    begin scan_eight_bit_int; m:=toks_base+cur_val;
    end;
   begin cur_val:=  eqtb[   m].hh.rh  ;cur_val_level:= tok_val; end ;
  end
else  begin back_input; scan_font_ident;
   begin cur_val:= font_id_base+ cur_val;cur_val_level:= ident_val; end ;
  end

;
assign_int:  begin cur_val:= eqtb[ m]. int;cur_val_level:= int_val; end ;
assign_dimen:  begin cur_val:= eqtb[ m]. int ;cur_val_level:= dimen_val; end ;
assign_glue:  begin cur_val:=  eqtb[   m].hh.rh  ;cur_val_level:= glue_val; end ;
assign_mu_glue:  begin cur_val:=  eqtb[   m].hh.rh  ;cur_val_level:= mu_val; end ;
set_aux: 
{ Fetch the |space_factor| or the |prev_depth| }
if abs(cur_list.mode_field )<>m then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Improper "=} 690); end ; print_cmd_chr(set_aux,m);
{ \xref[Improper \\spacefactor] }
{ \xref[Improper \\prevdepth] }
   begin help_ptr:=4; help_line[3]:={"You can refer to \spacefactor only in horizontal mode;"=} 691; help_line[2]:={"you can refer to \prevdepth only in vertical mode; and"=} 692; help_line[1]:={"neither of these is meaningful inside \write. So"=} 693; help_line[0]:={"I'm forgetting what you said and using zero instead."=} 694; end ;
  error;
  if level<>tok_val then  begin cur_val:= 0;cur_val_level:= dimen_val; end 
  else  begin cur_val:= 0;cur_val_level:= int_val; end ;
  end
else if m=vmode then  begin cur_val:= cur_list.aux_field .int  ;cur_val_level:= dimen_val; end 
else  begin cur_val:= cur_list.aux_field .hh.lh ;cur_val_level:= int_val; end 

;
set_prev_graf: 
{ Fetch the |prev_graf| }
if cur_list.mode_field =0 then  begin cur_val:= 0;cur_val_level:= int_val; end  {|prev_graf=0| within \.[\\write]}
else begin nest[nest_ptr]:=cur_list; p:=nest_ptr;
  while abs(nest[p].mode_field)<>vmode do decr(p);
   begin cur_val:= nest[ p]. pg_field;cur_val_level:= int_val; end ;
  end

;
set_page_int:
{ Fetch the |dead_cycles| or the |insert_penalties| }
begin if m=0 then cur_val:=dead_cycles else cur_val:=insert_penalties;
cur_val_level:=int_val;
end

;
set_page_dimen: 
{ Fetch something on the |page_so_far| }
begin if (page_contents=empty) and (not output_active) then
  if m=0 then cur_val:={07777777777=}1073741823  else cur_val:=0
else cur_val:=page_so_far[m];
cur_val_level:=dimen_val;
end

;
set_shape: 
{ Fetch the |par_shape| size }
begin if  eqtb[  par_shape_loc].hh.rh   =-{0xfffffff=}268435455   then cur_val:=0
else cur_val:= mem[  eqtb[  par_shape_loc].hh.rh   ].hh.lh ;
cur_val_level:=int_val;
end

;
set_box_dimen: 
{ Fetch a box dimension }
begin scan_eight_bit_int;
if  eqtb[  box_base+   cur_val].hh.rh   =-{0xfffffff=}268435455   then cur_val:=0  else cur_val:=mem[ eqtb[  box_base+   cur_val].hh.rh   +m].int ;
cur_val_level:=dimen_val;
end

;
char_given,math_given:  begin cur_val:= cur_chr;cur_val_level:= int_val; end ;
assign_font_dimen: 
{ Fetch a font dimension }
begin find_font_dimen(false); font_info[fmem_ptr].int :=0;
 begin cur_val:= font_info[ cur_val]. int ;cur_val_level:= dimen_val; end ;
end

;
assign_font_int: 
{ Fetch a font integer }
begin scan_font_ident;
if m=0 then  begin cur_val:= hyphen_char[ cur_val];cur_val_level:= int_val; end 
else  begin cur_val:= skew_char[ cur_val];cur_val_level:= int_val; end ;
end

;
register: 
{ Fetch a register }
begin scan_eight_bit_int;
case m of
int_val:cur_val:=eqtb[count_base+ cur_val].int ;
dimen_val:cur_val:=eqtb[scaled_base+ cur_val].int  ;
glue_val: cur_val:= eqtb[  skip_base+   cur_val].hh.rh   ;
mu_val: cur_val:= eqtb[  mu_skip_base+   cur_val].hh.rh   ;
end; {there are no other cases}
cur_val_level:=m;
end

;
last_item: 
{ Fetch an item in the current node, if appropriate }
if cur_chr>glue_val then
  begin if cur_chr=input_line_no_code then cur_val:=line
  else cur_val:=last_badness; {|cur_chr=badness_code|}
  cur_val_level:=int_val;
  end
else begin if cur_chr=glue_val then cur_val:=mem_bot  else cur_val:=0;
  cur_val_level:=cur_chr;
  if not  ( cur_list.tail_field >=hi_mem_min) and(cur_list.mode_field <>0) then
    case cur_chr of
    int_val: if  mem[ cur_list.tail_field ].hh.b0 =penalty_node then cur_val:= mem[ cur_list.tail_field +1].int ;
    dimen_val: if  mem[ cur_list.tail_field ].hh.b0 =kern_node then cur_val:= mem[ cur_list.tail_field +width_offset].int  ;
    glue_val: if  mem[ cur_list.tail_field ].hh.b0 =glue_node then
      begin cur_val:=  mem[  cur_list.tail_field + 1].hh.lh  ;
      if  mem[ cur_list.tail_field ].hh.b1 =mu_glue then cur_val_level:=mu_val;
      end;
    end {there are no other cases}
  else if (cur_list.mode_field =vmode)and(cur_list.tail_field =cur_list.head_field ) then
    case cur_chr of
    int_val: cur_val:=last_penalty;
    dimen_val: cur_val:=last_kern;
    glue_val: if last_glue<>{0xfffffff=}268435455  then cur_val:=last_glue;
    end; {there are no other cases}
  end

;
 else  
{ Complain that \.[\\the] can't do this; give zero result }
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"You can't use `"=} 695); end ; print_cmd_chr(cur_cmd,cur_chr);
{ \xref[You can't use x after ...] }
print({"' after "=}696); print_esc({"the"=}545);
 begin help_ptr:=1; help_line[0]:={"I'm forgetting what you said and using zero instead."=} 694; end ;
error;
if level<>tok_val then  begin cur_val:= 0;cur_val_level:= dimen_val; end 
else  begin cur_val:= 0;cur_val_level:= int_val; end ;
end


 end ;

while cur_val_level>level do 
{ Convert \(c)|cur_val| to a lower level }
begin if cur_val_level=glue_val then cur_val:= mem[ cur_val+width_offset].int  
else if cur_val_level=mu_val then mu_error;
decr(cur_val_level);
end

;

{ Fix the reference count, if any, and negate |cur_val| if |negative| }
if negative then
  if cur_val_level>=glue_val then
    begin cur_val:=new_spec(cur_val);
    
{ Negate all three glue components of |cur_val| }
begin    mem[  cur_val+width_offset].int  :=-  mem[  cur_val+width_offset].int   ;
   mem[  cur_val+2].int  :=-  mem[  cur_val+2].int   ;
   mem[  cur_val+3].int  :=-  mem[  cur_val+3].int   ;
end

;
    end
  else   cur_val:=- cur_val 
else if (cur_val_level>=glue_val)and(cur_val_level<=mu_val) then
  incr(  mem[   cur_val].hh.rh  ) 

;
end;



{ 432. }

{tangle:pos tex.web:8612:1: }

{ Our next goal is to write the |scan_int| procedure, which scans anything that
\TeX\ treats as an integer. But first we might as well look at some simple
applications of |scan_int| that have already been made inside of
|scan_something_internal|. }

{ 440. }

{tangle:pos tex.web:8701:1: }

{ The |scan_int| routine is used also to scan the integer part of a
fraction; for example, the `\.3' in `\.[3.14159]' will be found by
|scan_int|. The |scan_dimen| routine assumes that |cur_tok=point_token|
after the integer part of such a fraction has been scanned by |scan_int|,
and that the decimal point has been backed up to be scanned again. } procedure scan_int; {sets |cur_val| to an integer}
label done;
var negative:boolean; {should the answer be negated?}
 m:integer; {|$2^[31]$ div radix|, the threshold of danger}
 d:small_number; {the digit just scanned}
 vacuous:boolean; {have no digits appeared?}
 OK_so_far:boolean; {has an error message been issued?}
begin radix:=0; OK_so_far:=true;


{ Get the next non-blank non-sign token; set |negative| appropriately }
negative:=false;
repeat 
{ Get the next non-blank non-call token }
repeat get_x_token;
until cur_cmd<>spacer

;
if cur_tok=other_token+{"-"=}45 then
  begin negative := not negative; cur_tok:=other_token+{"+"=}43;
  end;
until cur_tok<>other_token+{"+"=}43

;
if cur_tok=alpha_token then 
{ Scan an alphabetic character code into |cur_val| }
begin get_token; {suppress macro expansion}
if cur_tok<{07777=}4095  then
  begin cur_val:=cur_chr;
  if cur_cmd<=right_brace then
    if cur_cmd=right_brace then incr(align_state)
    else decr(align_state);
  end
else if cur_tok<{07777=}4095 +single_base then
  cur_val:=cur_tok-{07777=}4095 -active_base
else cur_val:=cur_tok-{07777=}4095 -single_base;
if cur_val>255 then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Improper alphabetic constant"=} 708); end ;
{ \xref[Improper alphabetic constant] }
   begin help_ptr:=2; help_line[1]:={"A one-character control sequence belongs after a ` mark."=} 709; help_line[0]:={"So I'm essentially inserting \0 here."=} 710; end ;
  cur_val:={"0"=}48; back_error;
  end
else 
{ Scan an optional space }
begin get_x_token; if cur_cmd<>spacer then back_input;
end

;
end


else if (cur_cmd>=min_internal)and(cur_cmd<=max_internal) then
  scan_something_internal(int_val,false)
else 
{ Scan a numeric constant }
begin radix:=10; m:=214748364;
if cur_tok=octal_token then
  begin radix:=8; m:={02000000000=}268435456; get_x_token;
  end
else if cur_tok=hex_token then
  begin radix:=16; m:={01000000000=}134217728; get_x_token;
  end;
vacuous:=true; cur_val:=0;


{ Accumulate the constant until |cur_tok| is not a suitable digit }
 while true do    begin if (cur_tok<zero_token+radix)and(cur_tok>=zero_token)and
    (cur_tok<=zero_token+9) then d:=cur_tok-zero_token
  else if radix=16 then
    if (cur_tok<=A_token+5)and(cur_tok>=A_token) then d:=cur_tok-A_token+10
    else if (cur_tok<=other_A_token+5)and(cur_tok>=other_A_token) then
      d:=cur_tok-other_A_token+10
    else goto done
  else goto done;
  vacuous:=false;
  if (cur_val>=m)and((cur_val>m)or(d>7)or(radix<>10)) then
    begin if OK_so_far then
      begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Number too big"=} 711); end ;
{ \xref[Number too big] }
       begin help_ptr:=2; help_line[1]:={"I can only go up to 2147483647='17777777777=""7FFFFFFF,"=} 712; help_line[0]:={"so I'm using that number instead of yours."=} 713; end ;
      error; cur_val:={017777777777=}2147483647 ; OK_so_far:=false;
      end;
    end
  else cur_val:=cur_val*radix+d;
  get_x_token;
  end;
done:

;
if vacuous then 
{ Express astonishment that no number was here }
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Missing number, treated as zero"=} 674); end ;
{ \xref[Missing number...] }
 begin help_ptr:=3; help_line[2]:={"A number should have been here; I inserted `0'."=} 675; help_line[1]:={"(If you can't figure out why I needed to see a number,"=} 676; help_line[0]:={"look up `weird error' in the index to The TeXbook.)"=} 677; end ;
{ \xref[TeXbook][\sl The \TeX book] }
back_error;
end


else if cur_cmd<>spacer then back_input;
end

;
if negative then   cur_val:=- cur_val ;
end;



{ 448. }

{tangle:pos tex.web:8829:1: }

{ Constructions like `\.[-\'77 pt]' are legal dimensions, so |scan_dimen|
may begin with |scan_int|. This explains why it is convenient to use
|scan_int| also for the integer part of a decimal fraction.

Several branches of |scan_dimen| work with |cur_val| as an integer and
with an auxiliary fraction |f|, so that the actual quantity of interest is
$|cur_val|+|f|/2^[16]$. At the end of the routine, this ``unpacked''
representation is put into the single word |cur_val|, which suddenly
switches significance from |integer| to |scaled|. } procedure scan_dimen( mu, inf, shortcut:boolean);
  {sets |cur_val| to a dimension}
label done, done1, done2, found, not_found, attach_fraction, attach_sign;
var negative:boolean; {should the answer be negated?}
 f:integer; {numerator of a fraction whose denominator is $2^[16]$}

{ Local variables for dimension calculations }
 num, denom:1..65536; {conversion ratio for the scanned units}
 k, kk:small_number; {number of digits in a decimal fraction}
 p, q:halfword ; {top of decimal digit stack}
 v:scaled; {an internal dimension}
 save_cur_val:integer; {temporary storage of |cur_val|}

 
begin f:=0; arith_error:=false; cur_order:=normal; negative:=false;
if not shortcut then
  begin 
{ Get the next non-blank non-sign... }
negative:=false;
repeat 
{ Get the next non-blank non-call token }
repeat get_x_token;
until cur_cmd<>spacer

;
if cur_tok=other_token+{"-"=}45 then
  begin negative := not negative; cur_tok:=other_token+{"+"=}43;
  end;
until cur_tok<>other_token+{"+"=}43

;
  if (cur_cmd>=min_internal)and(cur_cmd<=max_internal) then
    
{ Fetch an internal dimension and |goto attach_sign|, or fetch an internal integer }
if mu then
  begin scan_something_internal(mu_val,false);
  
{ Coerce glue to a dimension }
if cur_val_level>=glue_val then
  begin v:= mem[ cur_val+width_offset].int  ; delete_glue_ref(cur_val); cur_val:=v;
  end

;
  if cur_val_level=mu_val then goto attach_sign;
  if cur_val_level<>int_val then mu_error;
  end
else  begin scan_something_internal(dimen_val,false);
  if cur_val_level=dimen_val then goto attach_sign;
  end


  else  begin back_input;
    if cur_tok=continental_point_token then cur_tok:=point_token;
    if cur_tok<>point_token then scan_int
    else  begin radix:=10; cur_val:=0;
      end;
    if cur_tok=continental_point_token then cur_tok:=point_token;
    if (radix=10)and(cur_tok=point_token) then 
{ Scan decimal fraction }
begin k:=0; p:=-{0xfffffff=}268435455  ; get_token; {|point_token| is being re-scanned}
 while true do    begin get_x_token;
  if (cur_tok>zero_token+9)or(cur_tok<zero_token) then goto done1;
  if k<17 then {digits for |k>=17| cannot affect the result}
    begin q:=get_avail;  mem[ q].hh.rh :=p;  mem[ q].hh.lh :=cur_tok-zero_token;
    p:=q; incr(k);
    end;
  end;
done1: for kk:=k downto 1 do
  begin dig[kk-1]:= mem[ p].hh.lh ; q:=p; p:= mem[ p].hh.rh ;  begin  mem[  q].hh.rh :=avail; avail:= q; ifdef('STAT')  decr(dyn_used); endif('STAT')  end ;
  end;
f:=round_decimals(k);
if cur_cmd<>spacer then back_input;
end

;
    end;
  end;
if cur_val<0 then {in this case |f=0|}
  begin negative := not negative;   cur_val:=- cur_val ;
  end;

{ Scan units and set |cur_val| to $x\cdot(|cur_val|+f/2^[16])$, where there are |x| sp per unit; |goto attach_sign| if the units are internal }
if inf then 
{ Scan for \(f)\.[fil] units; |goto attach_fraction| if found }
if scan_keyword({"fil"=}309) then
{ \xref[fil] }
  begin cur_order:=fil;
  while scan_keyword({"l"=}108) do
    begin if cur_order=filll then
      begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Illegal unit of measure ("=} 715); end ;
{ \xref[Illegal unit of measure] }
      print({"replaced by filll)"=}716);
       begin help_ptr:=1; help_line[0]:={"I dddon't go any higher than filll."=} 717; end ; error;
      end
    else incr(cur_order);
    end;
  goto attach_fraction;
  end

;

{ Scan for \(u)units that are internal dimensions; |goto attach_sign| with |cur_val| set if found }
save_cur_val:=cur_val;

{ Get the next non-blank non-call... }
repeat get_x_token;
until cur_cmd<>spacer

;
if (cur_cmd<min_internal)or(cur_cmd>max_internal) then back_input
else  begin if mu then
    begin scan_something_internal(mu_val,false); 
{ Coerce glue... }
if cur_val_level>=glue_val then
  begin v:= mem[ cur_val+width_offset].int  ; delete_glue_ref(cur_val); cur_val:=v;
  end

;
    if cur_val_level<>mu_val then mu_error;
    end
  else scan_something_internal(dimen_val,false);
  v:=cur_val; goto found;
  end;
if mu then goto not_found;
if scan_keyword({"em"=}718) then v:=(
{ The em width for |cur_font| }font_info[ quad_code+param_base[  eqtb[  cur_font_loc].hh.rh   ]].int  

)
{ \xref[em] }
else if scan_keyword({"ex"=}719) then v:=(
{ The x-height for |cur_font| }font_info[ x_height_code+param_base[  eqtb[  cur_font_loc].hh.rh   ]].int  

)
{ \xref[ex] }
else goto not_found;

{ Scan an optional space }
begin get_x_token; if cur_cmd<>spacer then back_input;
end

;
found:cur_val:=mult_and_add( save_cur_val, v, xn_over_d( v, f,{0200000=} 65536),{07777777777=}1073741823) ;
goto attach_sign;
not_found:

;
if mu then 
{ Scan for \(m)\.[mu] units and |goto attach_fraction| }
if scan_keyword({"mu"=}334) then goto attach_fraction
{ \xref[mu] }
else  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Illegal unit of measure ("=} 715); end ; print({"mu inserted)"=}720);
{ \xref[Illegal unit of measure] }
   begin help_ptr:=4; help_line[3]:={"The unit of measurement in math glue must be mu."=} 721; help_line[2]:={"To recover gracefully from this error, it's best to"=} 722; help_line[1]:={"delete the erroneous units; e.g., type `2' to delete"=} 723; help_line[0]:={"two letters. (See Chapter 27 of The TeXbook.)"=} 724; end ;
{ \xref[TeXbook][\sl The \TeX book] }
  error; goto attach_fraction;
  end

;
if scan_keyword({"true"=}714) then 
{ Adjust \(f)for the magnification ratio }
begin prepare_mag;
if eqtb[int_base+ mag_code].int  <>1000 then
  begin cur_val:=xn_over_d(cur_val,1000,eqtb[int_base+ mag_code].int  );
  f:=(1000*f+{0200000=}65536*tex_remainder ) div eqtb[int_base+ mag_code].int  ;
  cur_val:=cur_val+(f div {0200000=}65536); f:=f mod {0200000=}65536;
  end;
end

;
{ \xref[true] }
if scan_keyword({"pt"=}402) then goto attach_fraction; {the easy case}
{ \xref[pt] }

{ Scan for \(a)all other units and adjust |cur_val| and |f| accordingly; |goto done| in the case of scaled points }
if scan_keyword({"in"=}725) then  begin num:= 7227;  denom:= 100; end 
{ \xref[in] }
else if scan_keyword({"pc"=}726) then  begin num:= 12;  denom:= 1; end 
{ \xref[pc] }
else if scan_keyword({"cm"=}727) then  begin num:= 7227;  denom:= 254; end 
{ \xref[cm] }
else if scan_keyword({"mm"=}728) then  begin num:= 7227;  denom:= 2540; end 
{ \xref[mm] }
else if scan_keyword({"bp"=}729) then  begin num:= 7227;  denom:= 7200; end 
{ \xref[bp] }
else if scan_keyword({"dd"=}730) then  begin num:= 1238;  denom:= 1157; end 
{ \xref[dd] }
else if scan_keyword({"cc"=}731) then  begin num:= 14856;  denom:= 1157; end 
{ \xref[cc] }
else if scan_keyword({"sp"=}732) then goto done
{ \xref[sp] }
else 
{ Complain about unknown unit and |goto done2| }
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Illegal unit of measure ("=} 715); end ; print({"pt inserted)"=}733);
{ \xref[Illegal unit of measure] }
 begin help_ptr:=6; help_line[5]:={"Dimensions can be in units of em, ex, in, pt, pc,"=} 734; help_line[4]:={"cm, mm, dd, cc, bp, or sp; but yours is a new one!"=} 735; help_line[3]:={"I'll assume that you meant to say pt, for printer's points."=} 736; help_line[2]:={"To recover gracefully from this error, it's best to"=} 722; help_line[1]:={"delete the erroneous units; e.g., type `2' to delete"=} 723; help_line[0]:={"two letters. (See Chapter 27 of The TeXbook.)"=} 724; end ;
{ \xref[TeXbook][\sl The \TeX book] }
error; goto done2;
end


;
cur_val:=xn_over_d(cur_val,num,denom);
f:=(num*f+{0200000=}65536*tex_remainder ) div denom;

cur_val:=cur_val+(f div {0200000=}65536); f:=f mod {0200000=}65536;
done2:

;
attach_fraction: if cur_val>={040000=}16384 then arith_error:=true
else cur_val:=cur_val* {0200000=}65536 +f;
done:

;

{ Scan an optional space }
begin get_x_token; if cur_cmd<>spacer then back_input;
end

;
attach_sign: if arith_error or(abs(cur_val)>={010000000000=}1073741824) then
  
{ Report that this dimension is out of range }
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Dimension too large"=} 737); end ;
{ \xref[Dimension too large] }
 begin help_ptr:=2; help_line[1]:={"I can't work with sizes bigger than about 19 feet."=} 738; help_line[0]:={"Continue and I'll use the largest value I can."=} 739; end ;

error; cur_val:={07777777777=}1073741823 ; arith_error:=false;
end

;
if negative then   cur_val:=- cur_val ;
end;



{ 461. }

{tangle:pos tex.web:9064:1: }

{ The final member of \TeX's value-scanning trio is |scan_glue|, which
makes |cur_val| point to a glue specification. The reference count of that
glue spec will take account of the fact that |cur_val| is pointing to~it.

The |level| parameter should be either |glue_val| or |mu_val|.

Since |scan_dimen| was so much more complex than |scan_int|, we might expect
|scan_glue| to be even worse. But fortunately, it is very simple, since
most of the work has already been done. } procedure scan_glue( level:small_number);
  {sets |cur_val| to a glue spec pointer}
label exit;
var negative:boolean; {should the answer be negated?}
 q:halfword ; {new glue specification}
 mu:boolean; {does |level=mu_val|?}
begin mu:=(level=mu_val); 
{ Get the next non-blank non-sign... }
negative:=false;
repeat 
{ Get the next non-blank non-call token }
repeat get_x_token;
until cur_cmd<>spacer

;
if cur_tok=other_token+{"-"=}45 then
  begin negative := not negative; cur_tok:=other_token+{"+"=}43;
  end;
until cur_tok<>other_token+{"+"=}43

;
if (cur_cmd>=min_internal)and(cur_cmd<=max_internal) then
  begin scan_something_internal(level,negative);
  if cur_val_level>=glue_val then
    begin if cur_val_level<>level then mu_error;
     goto exit ;
    end;
  if cur_val_level=int_val then scan_dimen(mu,false,true)
  else if level=mu_val then mu_error;
  end
else  begin back_input; scan_dimen(mu,false,false);
  if negative then   cur_val:=- cur_val ;
  end;

{ Create a new glue specification whose width is |cur_val|; scan for its stretch and shrink components }
q:=new_spec(mem_bot );  mem[ q+width_offset].int  :=cur_val;
if scan_keyword({"plus"=}740) then
{ \xref[plus] }
  begin scan_dimen(mu,true,false);
   mem[ q+2].int  :=cur_val;   mem[ q].hh.b0 :=cur_order;
  end;
if scan_keyword({"minus"=}741) then
{ \xref[minus] }
  begin scan_dimen(mu,true,false);
   mem[ q+3].int  :=cur_val;   mem[ q].hh.b1 :=cur_order;
  end;
cur_val:=q

;
exit:end;



{ 463. }

{tangle:pos tex.web:9111:1: }

{ Here's a similar procedure that returns a pointer to a rule node. This
routine is called just after \TeX\ has seen \.[\\hrule] or \.[\\vrule];
therefore |cur_cmd| will be either |hrule| or |vrule|. The idea is to store
the default rule dimensions in the node, then to override them if
`\.[height]' or `\.[width]' or `\.[depth]' specifications are
found (in any order). } function scan_rule_spec:halfword ;
label reswitch;
var q:halfword ; {the rule node being created}
begin q:=new_rule; {|width|, |depth|, and |height| all equal |null_flag| now}
if cur_cmd=vrule then  mem[ q+width_offset].int  :=default_rule
else  begin  mem[ q+height_offset].int  :=default_rule;  mem[ q+depth_offset].int  :=0;
  end;
reswitch: if scan_keyword({"width"=}742) then
{ \xref[width] }
  begin scan_dimen(false,false,false) ;  mem[ q+width_offset].int  :=cur_val; goto reswitch;
  end;
if scan_keyword({"height"=}743) then
{ \xref[height] }
  begin scan_dimen(false,false,false) ;  mem[ q+height_offset].int  :=cur_val; goto reswitch;
  end;
if scan_keyword({"depth"=}744) then
{ \xref[depth] }
  begin scan_dimen(false,false,false) ;  mem[ q+depth_offset].int  :=cur_val; goto reswitch;
  end;
scan_rule_spec:=q;
end;



{ 464. \[27] Building token lists }

{tangle:pos tex.web:9142:29: }

{ The token lists for macros and for other things like \.[\\mark] and \.[\\output]
and \.[\\write] are produced by a procedure called |scan_toks|.

Before we get into the details of |scan_toks|, let's consider a much
simpler task, that of converting the current string into a token list.
The |str_toks| function does this; it classifies spaces as type |spacer|
and everything else as type |other_char|.

The token list created by |str_toks| begins at |link(temp_head)| and ends
at the value |p| that is returned. (If |p=temp_head|, the list is empty.) } function str_toks( b:pool_pointer):halfword ;
  {converts |str_pool[b..pool_ptr-1]| to a token list}
var p:halfword ; {tail of the token list}
 q:halfword ; {new node being added to the token list via |store_new_token|}
 t:halfword; {token being appended}
 k:pool_pointer; {index into |str_pool|}
begin  begin if pool_ptr+ 1 > pool_size then overflow({"pool size"=}257,pool_size-init_pool_ptr); { \xref[TeX capacity exceeded pool size][\quad pool size] } end ;
p:=mem_top-3 ;  mem[ p].hh.rh :=-{0xfffffff=}268435455  ; k:=b;
while k<pool_ptr do
  begin t:=  str_pool[ k] ;
  if t={" "=}32 then t:=space_token
  else t:=other_token+t;
  begin {  } begin  q:=avail; if  q=-{0xfffffff=}268435455   then  q:=get_avail else begin avail:= mem[  q].hh.rh ;  mem[  q].hh.rh :=-{0xfffffff=}268435455  ; ifdef('STAT')  incr(dyn_used); endif('STAT')  end; end ;  mem[ p].hh.rh :=q;  mem[ q].hh.lh := t; p:=q; end ;
  incr(k);
  end;
pool_ptr:=b; str_toks:=p;
end;



{ 465. }

{tangle:pos tex.web:9172:1: }

{ The main reason for wanting |str_toks| is the next function,
|the_toks|, which has similar input/output characteristics.

This procedure is supposed to scan something like `\.[\\skip\\count12]',
i.e., whatever can follow `\.[\\the]', and it constructs a token list
containing something like `\.[-3.0pt minus 0.5fill]'. } function the_toks:halfword ;
var old_setting:0..max_selector; {holds |selector| setting}
 p, q, r:halfword ; {used for copying a token list}
 b:pool_pointer; {base of temporary string}
begin get_x_token; scan_something_internal(tok_val,false);
if cur_val_level>=ident_val then 
{ Copy the token list }
begin p:=mem_top-3 ;  mem[ p].hh.rh :=-{0xfffffff=}268435455  ;
if cur_val_level=ident_val then begin q:=get_avail;  mem[ p].hh.rh :=q;  mem[ q].hh.lh := {07777=}4095 + cur_val; p:=q; end 
else if cur_val<>-{0xfffffff=}268435455   then
  begin r:= mem[ cur_val].hh.rh ; {do not copy the reference count}
  while r<>-{0xfffffff=}268435455   do
    begin begin {  } begin  q:=avail; if  q=-{0xfffffff=}268435455   then  q:=get_avail else begin avail:= mem[  q].hh.rh ;  mem[  q].hh.rh :=-{0xfffffff=}268435455  ; ifdef('STAT')  incr(dyn_used); endif('STAT')  end; end ;  mem[ p].hh.rh :=q;  mem[ q].hh.lh :=  mem[  r].hh.lh ; p:=q; end ; r:= mem[ r].hh.rh ;
    end;
  end;
the_toks:=p;
end


else begin old_setting:=selector; selector:=new_string; b:=pool_ptr;
  case cur_val_level of
  int_val:print_int(cur_val);
  dimen_val:begin print_scaled(cur_val); print({"pt"=}402);
    end;
  glue_val: begin print_spec(cur_val,{"pt"=}402); delete_glue_ref(cur_val);
    end;
  mu_val: begin print_spec(cur_val,{"mu"=}334); delete_glue_ref(cur_val);
    end;
  end; {there are no other cases}
  selector:=old_setting; the_toks:=str_toks(b);
  end;
end;



{ 467. }

{tangle:pos tex.web:9211:1: }

{ Here's part of the |expand| subroutine that we are now ready to complete: } procedure ins_the_toks;
begin  mem[ mem_top-12 ].hh.rh :=the_toks; begin_token_list(  mem[  mem_top-3 ].hh.rh ,inserted) ;
end;



{ 470. }

{tangle:pos tex.web:9251:1: }

{ The procedure |conv_toks| uses |str_toks| to insert the token list
for |convert| functions into the scanner; `\.[\\outer]' control sequences
are allowed to follow `\.[\\string]' and `\.[\\meaning]'. } procedure conv_toks;
var old_setting:0..max_selector; {holds |selector| setting}
 c:number_code..job_name_code; {desired type of conversion}
 save_scanner_status:small_number; {|scanner_status| upon entry}
 b:pool_pointer; {base of temporary string}
begin c:=cur_chr; 
{ Scan the argument for command |c| }
case c of
number_code,roman_numeral_code: scan_int;
string_code, meaning_code: begin save_scanner_status:=scanner_status;
  scanner_status:=normal; get_token; scanner_status:=save_scanner_status;
  end;
font_name_code: scan_font_ident;
job_name_code: if job_name=0 then open_log_file;
end {there are no other cases}

;
old_setting:=selector; selector:=new_string; b:=pool_ptr;

{ Print the result of command |c| }
case c of
number_code: print_int(cur_val);
roman_numeral_code: print_roman_int(cur_val);
string_code:if cur_cs<>0 then sprint_cs(cur_cs)
  else print_char(cur_chr);
meaning_code: print_meaning;
font_name_code: begin print(font_name[cur_val]);
  if font_size[cur_val]<>font_dsize[cur_val] then
    begin print({" at "=}751); print_scaled(font_size[cur_val]);
    print({"pt"=}402);
    end;
  end;
job_name_code: print(job_name);
end {there are no other cases}

;
selector:=old_setting;  mem[ mem_top-12 ].hh.rh :=str_toks(b); begin_token_list(  mem[  mem_top-3 ].hh.rh ,inserted) ;
end;



{ 473. }

{tangle:pos tex.web:9292:1: }

{ Now we can't postpone the difficulties any longer; we must bravely tackle
|scan_toks|. This function returns a pointer to the tail of a new token
list, and it also makes |def_ref| point to the reference count at the
head of that list.

There are two boolean parameters, |macro_def| and |xpand|. If |macro_def|
is true, the goal is to create the token list for a macro definition;
otherwise the goal is to create the token list for some other \TeX\
primitive: \.[\\mark], \.[\\output], \.[\\everypar], \.[\\lowercase],
\.[\\uppercase], \.[\\message], \.[\\errmessage], \.[\\write], or
\.[\\special]. In the latter cases a left brace must be scanned next; this
left brace will not be part of the token list, nor will the matching right
brace that comes at the end. If |xpand| is false, the token list will
simply be copied from the input using |get_token|. Otherwise all expandable
tokens will be expanded until unexpandable tokens are left, except that
the results of expanding `\.[\\the]' are not expanded further.
If both |macro_def| and |xpand| are true, the expansion applies
only to the macro body (i.e., to the material following the first
|left_brace| character).

The value of |cur_cs| when |scan_toks| begins should be the |eqtb|
address of the control sequence to display in ``runaway'' error
messages. } function scan_toks( macro_def, xpand:boolean):halfword ;
label found,continue,done,done1,done2;
var t:halfword; {token representing the highest parameter number}
 s:halfword; {saved token}
 p:halfword ; {tail of the token list being built}
 q:halfword ; {new node being added to the token list via |store_new_token|}
 unbalance:halfword; {number of unmatched left braces}
 hash_brace:halfword; {possible `\.[\#\[]' token}
begin if macro_def then scanner_status:=defining
 else scanner_status:=absorbing;
warning_index:=cur_cs; def_ref:=get_avail;   mem[  def_ref].hh.lh  :=-{0xfffffff=}268435455  ;
p:=def_ref; hash_brace:=0; t:=zero_token;
if macro_def then 
{ Scan and build the parameter part of the macro definition }
begin  while true do  begin continue: get_token; {set |cur_cmd|, |cur_chr|, |cur_tok|}
  if cur_tok<right_brace_limit then goto done1;
  if cur_cmd=mac_param then
    
{ If the next character is a parameter number, make |cur_tok| a |match| token; but if it is a left brace, store `|left_brace|, |end_match|', set |hash_brace|, and |goto done| }
begin s:=match_token+cur_chr; get_token;
if cur_tok<left_brace_limit then
  begin hash_brace:=cur_tok;
  begin q:=get_avail;  mem[ p].hh.rh :=q;  mem[ q].hh.lh := cur_tok; p:=q; end ; begin q:=get_avail;  mem[ p].hh.rh :=q;  mem[ q].hh.lh := end_match_token; p:=q; end ;
  goto done;
  end;
if t=zero_token+9 then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"You already have nine parameters"=} 754); end ;
{ \xref[You already have nine...] }
   begin help_ptr:=2; help_line[1]:={"I'm going to ignore the # sign you just used,"=} 755; help_line[0]:={"as well as the token that followed it."=} 756; end ; error; goto continue;
  end
else  begin incr(t);
  if cur_tok<>t then
    begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Parameters must be numbered consecutively"=} 757); end ;
{ \xref[Parameters...consecutively] }
     begin help_ptr:=2; help_line[1]:={"I've inserted the digit you should have used after the #."=} 758; help_line[0]:={"Type `1' to delete what you did use."=} 759; end ; back_error;
    end;
  cur_tok:=s;
  end;
end

;
  begin q:=get_avail;  mem[ p].hh.rh :=q;  mem[ q].hh.lh := cur_tok; p:=q; end ;
  end;
done1: begin q:=get_avail;  mem[ p].hh.rh :=q;  mem[ q].hh.lh := end_match_token; p:=q; end ;
if cur_cmd=right_brace then
  
{ Express shock at the missing left brace; |goto found| }
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Missing [ inserted"=} 667); end ; incr(align_state);
{ \xref[Missing \[ inserted] }
 begin help_ptr:=2; help_line[1]:={"Where was the left brace? You said something like `\def\a]',"=} 752; help_line[0]:={"which I'm going to interpret as `\def\a[]'."=} 753; end ; error; goto found;
end

;
done: end


else scan_left_brace; {remove the compulsory left brace}

{ Scan and build the body of the token list; |goto found| when finished }
unbalance:=1;
 while true do    begin if xpand then 
{ Expand the next part of the input }
begin  while true do  begin get_next;
  if cur_cmd<=max_command then goto done2;
  if cur_cmd<>the then expand
  else  begin q:=the_toks;
    if  mem[ mem_top-3 ].hh.rh <>-{0xfffffff=}268435455   then
      begin  mem[ p].hh.rh := mem[ mem_top-3 ].hh.rh ; p:=q;
      end;
    end;
  end;
done2: x_token
end


  else get_token;
  if cur_tok<right_brace_limit then
    if cur_cmd<right_brace then incr(unbalance)
    else  begin decr(unbalance);
      if unbalance=0 then goto found;
      end
  else if cur_cmd=mac_param then
    if macro_def then 
{ Look for parameter number or \.[\#\#] }
begin s:=cur_tok;
if xpand then get_x_token else get_token;
if cur_cmd<>mac_param then
  if (cur_tok<=zero_token)or(cur_tok>t) then
    begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Illegal parameter number in definition of "=} 760); end ;
{ \xref[Illegal parameter number...] }
    sprint_cs(warning_index);
     begin help_ptr:=3; help_line[2]:={"You meant to type ## instead of #, right?"=} 761; help_line[1]:={"Or maybe a ] was forgotten somewhere earlier, and things"=} 762; help_line[0]:={"are all screwed up? I'm going to assume that you meant ##."=} 763; end ;
    back_error; cur_tok:=s;
    end
  else cur_tok:=out_param_token-{"0"=}48+cur_chr;
end

;
  begin q:=get_avail;  mem[ p].hh.rh :=q;  mem[ q].hh.lh := cur_tok; p:=q; end ;
  end

;
found: scanner_status:=normal;
if hash_brace<>0 then begin q:=get_avail;  mem[ p].hh.rh :=q;  mem[ q].hh.lh := hash_brace; p:=q; end ;
scan_toks:=p;
end;



{ 482. }

{tangle:pos tex.web:9444:1: }

{ The |read_toks| procedure constructs a token list like that for any
macro definition, and makes |cur_val| point to it. Parameter |r| points
to the control sequence that will receive this token list. } procedure read_toks( n:integer; r:halfword );
label done;
var p:halfword ; {tail of the token list}
 q:halfword ; {new node being added to the token list via |store_new_token|}
 s:integer; {saved value of |align_state|}
 m:small_number; {stream number}
begin scanner_status:=defining; warning_index:=r;
def_ref:=get_avail;   mem[  def_ref].hh.lh  :=-{0xfffffff=}268435455  ;
p:=def_ref; {the reference count}
begin q:=get_avail;  mem[ p].hh.rh :=q;  mem[ q].hh.lh := end_match_token; p:=q; end ;
if (n<0)or(n>15) then m:=16 else m:=n;
s:=align_state; align_state:=1000000; {disable tab marks, etc.}
repeat 
{ Input and store tokens from the next line of the file }
begin_file_reading; cur_input.name_field :=m+1;
if read_open[m]=closed then 
{ Input for \.[\\read] from the terminal }
if interaction>nonstop_mode then
  if n<0 then begin    ; print({""=} 335); term_input; end 
  else  begin    ;
    print_ln; sprint_cs(r); begin    ; print({"="=} 61); term_input; end ; n:=-1;
    end
else begin
  cur_input.limit_field :=0;
  fatal_error({"*** (cannot \read from terminal in nonstop modes)"=}764);
  end
{ \xref[cannot \\read] }


else if read_open[m]=just_open then 
{ Input the first line of |read_file[m]| }
if input_ln(read_file[m],false) then read_open[m]:=normal
else  begin a_close(read_file[m]); read_open[m]:=closed;
  end


else 
{ Input the next line of |read_file[m]| }
begin if not input_ln(read_file[m],true) then
  begin a_close(read_file[m]); read_open[m]:=closed;
  if align_state<>1000000 then
    begin runaway;
    begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"File ended within "=} 765); end ; print_esc({"read"=}542);
{ \xref[File ended within \\read] }
     begin help_ptr:=1; help_line[0]:={"This \read has unbalanced braces."=} 766; end ;
    align_state:=1000000; cur_input.limit_field :=0; error;
    end;
  end;
end

;
cur_input.limit_field :=last;
if  (eqtb[int_base+ end_line_char_code].int  <0)or(eqtb[int_base+ end_line_char_code].int  >255)  then decr(cur_input.limit_field )
else  buffer[cur_input.limit_field ]:=eqtb[int_base+ end_line_char_code].int  ;
first:=cur_input.limit_field +1; cur_input.loc_field :=cur_input.start_field ; cur_input.state_field :=new_line;

 while true do    begin get_token;
  if cur_tok=0 then goto done;
    {|cur_cmd=cur_chr=0| will occur at the end of the line}
  if align_state<1000000 then {unmatched `\.\]' aborts the line}
    begin repeat get_token; until cur_tok=0;
    align_state:=1000000; goto done;
    end;
  begin q:=get_avail;  mem[ p].hh.rh :=q;  mem[ q].hh.lh := cur_tok; p:=q; end ;
  end;
done: end_file_reading

;
until align_state=1000000;
cur_val:=def_ref; scanner_status:=normal; align_state:=s;
end;



{ 494. }

{tangle:pos tex.web:9653:1: }

{ Here is a procedure that ignores text until coming to an \.[\\or],
\.[\\else], or \.[\\fi] at the current level of $\.[\\if]\ldots\.[\\fi]$
nesting. After it has acted, |cur_chr| will indicate the token that
was found, but |cur_tok| will not be set (because this makes the
procedure run faster). } procedure pass_text;
label done;
var l:integer; {level of $\.[\\if]\ldots\.[\\fi]$ nesting}
 save_scanner_status:small_number; {|scanner_status| upon entry}
begin save_scanner_status:=scanner_status; scanner_status:=skipping; l:=0;
skip_line:=line;
 while true do    begin get_next;
  if cur_cmd=fi_or_else then
    begin if l=0 then goto done;
    if cur_chr=fi_code then decr(l);
    end
  else if cur_cmd=if_test then incr(l);
  end;
done: scanner_status:=save_scanner_status;
end;



{ 497. }

{tangle:pos tex.web:9693:1: }

{ Here's a procedure that changes the |if_limit| code corresponding to
a given value of |cond_ptr|. } procedure change_if_limit( l:small_number; p:halfword );
label exit;
var q:halfword ;
begin if p=cond_ptr then if_limit:=l {that's the easy case}
else  begin q:=cond_ptr;
   while true do    begin if q=-{0xfffffff=}268435455   then confusion({"if"=}767);
{ \xref[this can't happen if][\quad if] }
    if  mem[ q].hh.rh =p then
      begin  mem[ q].hh.b0 :=l;  goto exit ;
      end;
    q:= mem[ q].hh.rh ;
    end;
  end;
exit:end;



{ 498. }

{tangle:pos tex.web:9711:1: }

{ A condition is started when the |expand| procedure encounters
an |if_test| command; in that case |expand| reduces to |conditional|,
which is a recursive procedure.
\xref[recursion] } procedure conditional;
label exit,common_ending;
var b:boolean; {is the condition true?}
 r:{"<"=}60..{">"=}62; {relation to be evaluated}
 m, n:integer; {to be tested against the second operand}
 p, q:halfword ; {for traversing token lists in \.[\\ifx] tests}
 save_scanner_status:small_number; {|scanner_status| upon entry}
 save_cond_ptr:halfword ; {|cond_ptr| corresponding to this conditional}
 this_if:small_number; {type of this conditional}
begin 
{ Push the condition stack }
begin p:=get_node(if_node_size);  mem[ p].hh.rh :=cond_ptr;  mem[ p].hh.b0 :=if_limit;
 mem[ p].hh.b1 :=cur_if; mem[ p+1].int :=if_line;
cond_ptr:=p; cur_if:=cur_chr; if_limit:=if_code; if_line:=line;
end

; save_cond_ptr:=cond_ptr;this_if:=cur_chr;


{ Either process \.[\\ifcase] or set |b| to the value of a boolean condition }
case this_if of
if_char_code, if_cat_code: 
{ Test if two characters match }
begin {  } begin get_x_token; if cur_cmd=relax then if cur_chr=no_expand_flag then begin cur_cmd:=active_char; cur_chr:=cur_tok-{07777=}4095 -active_base; end; end ;
if (cur_cmd>active_char)or(cur_chr>255) then {not a character}
  begin m:=relax; n:=256;
  end
else  begin m:=cur_cmd; n:=cur_chr;
  end;
{  } begin get_x_token; if cur_cmd=relax then if cur_chr=no_expand_flag then begin cur_cmd:=active_char; cur_chr:=cur_tok-{07777=}4095 -active_base; end; end ;
if (cur_cmd>active_char)or(cur_chr>255) then
  begin cur_cmd:=relax; cur_chr:=256;
  end;
if this_if=if_char_code then b:=(n=cur_chr) else b:=(m=cur_cmd);
end

;
if_int_code, if_dim_code: 
{ Test relation between integers or dimensions }
begin if this_if=if_int_code then scan_int else scan_dimen(false,false,false) ;
n:=cur_val; 
{ Get the next non-blank non-call... }
repeat get_x_token;
until cur_cmd<>spacer

;
if (cur_tok>=other_token+{"<"=}60)and(cur_tok<=other_token+{">"=}62) then
  r:=cur_tok-other_token
else  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Missing = inserted for "=} 791); end ;
{ \xref[Missing = inserted] }
  print_cmd_chr(if_test,this_if);
   begin help_ptr:=1; help_line[0]:={"I was expecting to see `<', `=', or `>'. Didn't."=} 792; end ;
  back_error; r:={"="=}61;
  end;
if this_if=if_int_code then scan_int else scan_dimen(false,false,false) ;
case r of
{"<"=}60: b:=(n<cur_val);
{"="=}61: b:=(n=cur_val);
{">"=}62: b:=(n>cur_val);
end;
end

;
if_odd_code: 
{ Test if an integer is odd }
begin scan_int; b:=odd(cur_val);
end

;
if_vmode_code: b:=(abs(cur_list.mode_field )=vmode);
if_hmode_code: b:=(abs(cur_list.mode_field )=hmode);
if_mmode_code: b:=(abs(cur_list.mode_field )=mmode);
if_inner_code: b:=(cur_list.mode_field <0);
if_void_code, if_hbox_code, if_vbox_code: 
{ Test box register status }
begin scan_eight_bit_int; p:= eqtb[  box_base+   cur_val].hh.rh   ;
if this_if=if_void_code then b:=(p=-{0xfffffff=}268435455  )
else if p=-{0xfffffff=}268435455   then b:=false
else if this_if=if_hbox_code then b:=( mem[ p].hh.b0 =hlist_node)
else b:=( mem[ p].hh.b0 =vlist_node);
end

;
ifx_code: 
{ Test if two tokens match }
begin save_scanner_status:=scanner_status; scanner_status:=normal;
get_next; n:=cur_cs; p:=cur_cmd; q:=cur_chr;
get_next; if cur_cmd<>p then b:=false
else if cur_cmd<call then b:=(cur_chr=q)
else 
{ Test if two macro texts match }
begin p:= mem[ cur_chr].hh.rh ; q:= mem[  eqtb[   n].hh.rh  ].hh.rh ; {omit reference counts}
if p=q then b:=true
else begin while (p<>-{0xfffffff=}268435455  )and(q<>-{0xfffffff=}268435455  ) do
    if  mem[ p].hh.lh <> mem[ q].hh.lh  then p:=-{0xfffffff=}268435455  
    else  begin p:= mem[ p].hh.rh ; q:= mem[ q].hh.rh ;
      end;
  b:=((p=-{0xfffffff=}268435455  )and(q=-{0xfffffff=}268435455  ));
  end;
end

;
scanner_status:=save_scanner_status;
end

;
if_eof_code: begin scan_four_bit_int_or_18;
  if cur_val=18 then b:=not shellenabledp
  else b:=(read_open[cur_val]=closed);
  end;
if_true_code: b:=true;
if_false_code: b:=false;
if_case_code: 
{ Select the appropriate case and |return| or |goto common_ending| }
begin scan_int; n:=cur_val; {|n| is the number of cases to pass}
if eqtb[int_base+ tracing_commands_code].int  >1 then
  begin begin_diagnostic; print({"[case "=}793); print_int(n); print_char({"]"=}125);
  end_diagnostic(false);
  end;
while n<>0 do
  begin pass_text;
  if cond_ptr=save_cond_ptr then
    if cur_chr=or_code then decr(n)
    else goto common_ending
  else if cur_chr=fi_code then 
{ Pop the condition stack }
begin p:=cond_ptr; if_line:=mem[ p+1].int ;
cur_if:= mem[ p].hh.b1 ; if_limit:= mem[ p].hh.b0 ; cond_ptr:= mem[ p].hh.rh ;
free_node(p,if_node_size);
end

;
  end;
change_if_limit(or_code,save_cond_ptr);
 goto exit ; {wait for \.[\\or], \.[\\else], or \.[\\fi]}
end

;
end {there are no other cases}

;
if eqtb[int_base+ tracing_commands_code].int  >1 then 
{ Display the value of |b| }
begin begin_diagnostic;
if b then print({"[true]"=}789) else print({"[false]"=}790);
end_diagnostic(false);
end

;
if b then
  begin change_if_limit(else_code,save_cond_ptr);
   goto exit ; {wait for \.[\\else] or \.[\\fi]}
  end;

{ Skip to \.[\\else] or \.[\\fi], then |goto common_ending| }
 while true do    begin pass_text;
  if cond_ptr=save_cond_ptr then
    begin if cur_chr<>or_code then goto common_ending;
    begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Extra "=} 787); end ; print_esc({"or"=}785);
{ \xref[Extra \\or] }
     begin help_ptr:=1; help_line[0]:={"I'm ignoring this; it doesn't match any \if."=} 788; end ;
    error;
    end
  else if cur_chr=fi_code then 
{ Pop the condition stack }
begin p:=cond_ptr; if_line:=mem[ p+1].int ;
cur_if:= mem[ p].hh.b1 ; if_limit:= mem[ p].hh.b0 ; cond_ptr:= mem[ p].hh.rh ;
free_node(p,if_node_size);
end

;
  end

;
common_ending: if cur_chr=fi_code then 
{ Pop the condition stack }
begin p:=cond_ptr; if_line:=mem[ p+1].int ;
cur_if:= mem[ p].hh.b1 ; if_limit:= mem[ p].hh.b0 ; cond_ptr:= mem[ p].hh.rh ;
free_node(p,if_node_size);
end


else if_limit:=fi_code; {wait for \.[\\fi]}
exit:end;



{ 499. }

{tangle:pos tex.web:9737:1: }

{ In a construction like `\.[\\if\\iftrue abc\\else d\\fi]', the first
\.[\\else] that we come to after learning that the \.[\\if] is false is
not the \.[\\else] we're looking for. Hence the following curious
logic is needed. }

{ 511. \[29] File names }

{tangle:pos tex.web:9910:19: }

{ It's time now to fret about file names.  Besides the fact that different
operating systems treat files in different ways, we must cope with the
fact that completely different naming conventions are used by different
groups of people. The following programs show what is required for one
particular operating system; similar routines for other systems are not
difficult to devise.
\xref[fingers]
\xref[system dependencies]

\TeX\ assumes that a file name has three parts: the name proper; its
``extension''; and a ``file area'' where it is found in an external file
system.  The extension of an input file or a write file is assumed to be
`\.[.tex]' unless otherwise specified; it is `\.[.log]' on the
transcript file that records each run of \TeX; it is `\.[.tfm]' on the font
metric files that describe characters in the fonts \TeX\ uses; it is
`\.[.dvi]' on the output files that specify typesetting information; and it
is `\.[.fmt]' on the format files written by \.[INITEX] to initialize \TeX.
The file area can be arbitrary on input files, but files are usually
output to the user's current area.  If an input file cannot be
found on the specified area, \TeX\ will look for it on a special system
area; this special area is intended for commonly used input files like
\.[webmac.tex].

Simple uses of \TeX\ refer only to file names that have no explicit
extension or area. For example, a person usually says `\.[\\input] \.[paper]'
or `\.[\\font\\tenrm] \.= \.[helvetica]' instead of `\.[\\input]
\.[paper.new]' or `\.[\\font\\tenrm] \.= \.[<csd.knuth>test]'. Simple file
names are best, because they make the \TeX\ source files portable;
whenever a file name consists entirely of letters and digits, it should be
treated in the same way by all implementations of \TeX. However, users
need the ability to refer to other files in their environment, especially
when responding to error messages concerning unopenable files; therefore
we want to let them use the syntax that appears in their favorite
operating system.

The following procedures don't allow spaces to be part of
file names; but some users seem to like names that are spaced-out.
System-dependent changes to allow such things should probably
be made with reluctance, and only when an entire file name that
includes spaces is ``quoted'' somehow. }

{ 514. }

{tangle:pos tex.web:9998:1: }

{ Input files that can't be found in the user's area may appear in a standard
system area called |TEX_area|. Font metric files whose areas are not given
explicitly are assumed to appear in a standard system area called
|TEX_font_area|.  These system area names will, of course, vary from place
to place.
\xref[system dependencies]

In C, the default paths are specified separately. }

{ 515. }

{tangle:pos tex.web:10010:1: }

{ Here now is the first of the system-dependent routines for file name scanning.
\xref[system dependencies] } procedure begin_name;
begin area_delimiter:=0; ext_delimiter:=0; quoted_filename:=false;
end;



{ 516. }

{tangle:pos tex.web:10017:1: }

{ And here's the second. The string pool might change as the file name is
being scanned, since a new \.[\\csname] might be entered; therefore we keep
|area_delimiter| and |ext_delimiter| relative to the beginning of the current
string, instead of assigning an absolute address like |pool_ptr| to them.
\xref[system dependencies] } function more_name( c:ASCII_code):boolean;
begin if (c={" "=}32) and stop_at_space and (not quoted_filename) then
  more_name:=false
else  if c={""""=}34 then begin
  quoted_filename:=not quoted_filename;
  more_name:=true;
  end
else  begin  begin if pool_ptr+ 1 > pool_size then overflow({"pool size"=}257,pool_size-init_pool_ptr); { \xref[TeX capacity exceeded pool size][\quad pool size] } end ;  begin str_pool[pool_ptr]:=   c ; incr(pool_ptr); end ; {contribute |c| to the current string}
  if IS_DIR_SEP(c) then
    begin area_delimiter:= (pool_ptr - str_start[str_ptr]) ; ext_delimiter:=0;
    end
  else if c={"."=}46 then ext_delimiter:= (pool_ptr - str_start[str_ptr]) ;
  more_name:=true;
  end;
end;



{ 517. }

{tangle:pos tex.ch:1618:3: }

{ The third.
\xref[system dependencies]
If a string is already in the string pool, the function
|slow_make_string| does not create a new string but returns this string
number, thus saving string space.  Because of this new property of the
returned string number it is not possible to apply |flush_string| to
these strings. } procedure end_name;
var temp_str: str_number; {result of file name cache lookups}
 j, s, t: pool_pointer; {running indices}
 must_quote:boolean; {whether we need to quote a string}
begin if str_ptr+3>max_strings then
  overflow({"number of strings"=}258,max_strings-init_str_ptr);
{ \xref[TeX capacity exceeded number of strings][\quad number of strings] }
 begin if pool_ptr+ 6 > pool_size then overflow({"pool size"=}257,pool_size-init_pool_ptr); { \xref[TeX capacity exceeded pool size][\quad pool size] } end ; {Room for quotes, if needed.}
{add quotes if needed}
if area_delimiter<>0 then begin
  {maybe quote |cur_area|}
  must_quote:=false;
  s:=str_start[str_ptr];
  t:=str_start[str_ptr]+area_delimiter;
  j:=s;
  while (not must_quote) and (j<t) do begin
    must_quote:=str_pool[j]={" "=}32; incr(j);
    end;
  if must_quote then begin
    for j:=pool_ptr-1 downto t do str_pool[j+2]:=str_pool[j];
    str_pool[t+1]:={""""=}34;
    for j:=t-1 downto s do str_pool[j+1]:=str_pool[j];
    str_pool[s]:={""""=}34;
    if ext_delimiter<>0 then ext_delimiter:=ext_delimiter+2;
    area_delimiter:=area_delimiter+2;
    pool_ptr:=pool_ptr+2;
    end;
  end;
{maybe quote |cur_name|}
s:=str_start[str_ptr]+area_delimiter;
if ext_delimiter=0 then t:=pool_ptr else t:=str_start[str_ptr]+ext_delimiter-1;
must_quote:=false;
j:=s;
while (not must_quote) and (j<t) do begin
  must_quote:=str_pool[j]={" "=}32; incr(j);
  end;
if must_quote then begin
  for j:=pool_ptr-1 downto t do str_pool[j+2]:=str_pool[j];
  str_pool[t+1]:={""""=}34;
  for j:=t-1 downto s do str_pool[j+1]:=str_pool[j];
  str_pool[s]:={""""=}34;
  if ext_delimiter<>0 then ext_delimiter:=ext_delimiter+2;
  pool_ptr:=pool_ptr+2;
  end;
if ext_delimiter<>0 then begin
  {maybe quote |cur_ext|}
  s:=str_start[str_ptr]+ext_delimiter-1;
  t:=pool_ptr;
  must_quote:=false;
  j:=s;
  while (not must_quote) and (j<t) do begin
    must_quote:=str_pool[j]={" "=}32; incr(j);
    end;
  if must_quote then begin
    str_pool[t+1]:={""""=}34;
    for j:=t-1 downto s do str_pool[j+1]:=str_pool[j];
    str_pool[s]:={""""=}34;
    pool_ptr:=pool_ptr+2;
    end;
  end;
if area_delimiter=0 then cur_area:={""=}335
else  begin cur_area:=str_ptr;
  str_start[str_ptr+1]:=str_start[str_ptr]+area_delimiter; incr(str_ptr);
  temp_str:=search_string(cur_area);
  if temp_str>0 then
    begin cur_area:=temp_str;
    decr(str_ptr);  {no |flush_string|, |pool_ptr| will be wrong!}
    for j:=str_start[str_ptr+1] to pool_ptr-1 do
      begin str_pool[j-area_delimiter]:=str_pool[j];
      end;
    pool_ptr:=pool_ptr-area_delimiter; {update |pool_ptr|}
    end;
  end;
if ext_delimiter=0 then
  begin cur_ext:={""=}335; cur_name:=slow_make_string;
  end
else  begin cur_name:=str_ptr;
  str_start[str_ptr+1]:=str_start[str_ptr]+ext_delimiter-area_delimiter-1;
  incr(str_ptr); cur_ext:=make_string;
  decr(str_ptr); {undo extension string to look at name part}
  temp_str:=search_string(cur_name);
  if temp_str>0 then
    begin cur_name:=temp_str;
    decr(str_ptr);  {no |flush_string|, |pool_ptr| will be wrong!}
    for j:=str_start[str_ptr+1] to pool_ptr-1 do
      begin str_pool[j-ext_delimiter+area_delimiter+1]:=str_pool[j];
      end;
    pool_ptr:=pool_ptr-ext_delimiter+area_delimiter+1;  {update |pool_ptr|}
    end;
  cur_ext:=slow_make_string;  {remake extension string}
  end;
end;



{ 519. }

{tangle:pos tex.web:10064:1: }

{ Another system-dependent routine is needed to convert three internal
\TeX\ strings
into the |name_of_file| value that is used to open files. The present code
allows both lowercase and uppercase letters in the file name.
\xref[system dependencies] } procedure pack_file_name( n, a, e:str_number);
var k:integer; {number of positions filled in |name_of_file|}
 c: ASCII_code; {character being packed}
 j:pool_pointer; {index into |str_pool|}
begin k:=0;
if name_of_file then libc_free (name_of_file);
name_of_file:= xmalloc_array (ASCII_code, (str_start[ a+1]-str_start[ a]) +(str_start[ n+1]-str_start[ n]) +(str_start[ e+1]-str_start[ e]) +1);
for j:=str_start[a] to str_start[a+1]-1 do begin c:=    str_pool[  j] ; if not (c={""""=}34) then begin incr(k); if k<= maxint  then name_of_file[k]:=xchr[c]; end end ;
for j:=str_start[n] to str_start[n+1]-1 do begin c:=    str_pool[  j] ; if not (c={""""=}34) then begin incr(k); if k<= maxint  then name_of_file[k]:=xchr[c]; end end ;
for j:=str_start[e] to str_start[e+1]-1 do begin c:=    str_pool[  j] ; if not (c={""""=}34) then begin incr(k); if k<= maxint  then name_of_file[k]:=xchr[c]; end end ;
if k<= maxint  then name_length:=k else name_length:= maxint ;
name_of_file[name_length+1]:=0;
end;



{ 521. }

{tangle:pos tex.ch:1814:3: }

{ We set the name of the default format file and the length of that name
in C, instead of Pascal, since we want them to depend on the name of the
program.
\xref[TeXformats]
\xref[plain]
\xref[system dependencies] }

{ 523. }

{tangle:pos tex.web:10109:1: }

{ Here is the messy routine that was just mentioned. It sets |name_of_file|
from the first |n| characters of |TEX_format_default|, followed by
|buffer[a..b]|, followed by the last |format_ext_length| characters of
|TEX_format_default|.

We dare not give error messages here, since \TeX\ calls this routine before
the |error| routine is ready to roll. Instead, we simply drop excess characters,
since the error will be detected in another way when a strange file name
isn't found.
\xref[system dependencies] } procedure pack_buffered_name( n:small_number; a, b:integer);
var k:integer; {number of positions filled in |name_of_file|}
 c: ASCII_code; {character being packed}
 j:integer; {index into |buffer| or |TEX_format_default|}
begin if n+b-a+1+format_ext_length> maxint  then
  b:=a+ maxint -n-1-format_ext_length;
k:=0;
if name_of_file then libc_free (name_of_file);
name_of_file := xmalloc_array (ASCII_code, n+(b-a+1)+format_ext_length+1);
for j:=1 to n do begin c:= xord[ ucharcast( TEX_format_default[ j])]; if not (c={""""=}34) then begin incr(k); if k<= maxint  then name_of_file[k]:=xchr[c]; end end ;
for j:=a to b do begin c:= buffer[ j]; if not (c={""""=}34) then begin incr(k); if k<= maxint  then name_of_file[k]:=xchr[c]; end end ;
for j:=format_default_length-format_ext_length+1 to format_default_length do
  begin c:= xord[ ucharcast( TEX_format_default[ j])]; if not (c={""""=}34) then begin incr(k); if k<= maxint  then name_of_file[k]:=xchr[c]; end end ;
if k<= maxint  then name_length:=k else name_length:= maxint ;
name_of_file[name_length+1]:=0;
end;



{ 525. }

{tangle:pos tex.web:10172:1: }

{ Operating systems often make it possible to determine the exact name (and
possible version number) of a file that has been opened. The following routine,
which simply makes a \TeX\ string from the value of |name_of_file|, should
ideally be changed to deduce the full name of file~|f|, which is the file
most recently opened, if it is possible to do this in a \PASCAL\ program.
\xref[system dependencies]

This routine might be called after string memory has overflowed, hence
we dare not use `|str_room|'. } function make_name_string:str_number;
var k:1.. maxint ; {index into |name_of_file|}
save_area_delimiter, save_ext_delimiter: pool_pointer;
save_name_in_progress, save_stop_at_space: boolean;
begin if (pool_ptr+name_length>pool_size)or(str_ptr=max_strings)or
 ( (pool_ptr - str_start[str_ptr]) >0) then
  make_name_string:={"?"=}63
else  begin for k:=1 to name_length do  begin str_pool[pool_ptr]:=   xord[  name_of_file[  k]] ; incr(pool_ptr); end ;
  make_name_string:=make_string;
  {At this point we also set |cur_name|, |cur_ext|, and |cur_area| to
   match the contents of |name_of_file|.}
  save_area_delimiter:=area_delimiter; save_ext_delimiter:=ext_delimiter;
  save_name_in_progress:=name_in_progress; save_stop_at_space:=stop_at_space;
  name_in_progress:=true;
  begin_name;
  stop_at_space:=false;
  k:=1;
  while (k<=name_length)and(more_name(name_of_file[k])) do
    incr(k);
  stop_at_space:=save_stop_at_space;
  end_name;
  name_in_progress:=save_name_in_progress;
  area_delimiter:=save_area_delimiter; ext_delimiter:=save_ext_delimiter;
  end;
end;
function a_make_name_string(var f:alpha_file):str_number;
begin a_make_name_string:=make_name_string;
end;
function b_make_name_string(var f:byte_file):str_number;
begin b_make_name_string:=make_name_string;
end;
function w_make_name_string(var f:word_file):str_number;
begin w_make_name_string:=make_name_string;
end;



{ 526. }

{tangle:pos tex.web:10201:1: }

{ Now let's consider the ``driver''
routines by which \TeX\ deals with file names
in a system-independent manner.  First comes a procedure that looks for a
file name in the input by calling |get_x_token| for the information. } procedure scan_file_name;
label done;
var
   save_warning_index: halfword ;
begin
  save_warning_index := warning_index;
  warning_index := cur_cs; {store |cur_cs| here to remember until later}
  
{ Get the next non-blank non-relax non-call... }
repeat get_x_token;
until (cur_cmd<>spacer)and(cur_cmd<>relax)

; {here the program expands
    tokens and removes spaces and \.[\\relax]es from the input. The \.[\\relax]
    removal follows LuaTeX''s implementation, and other cases of
    balanced text scanning.}
  back_input; {return the last token to be read by either code path}
  if cur_cmd=left_brace then
    scan_file_name_braced
  else
begin name_in_progress:=true; begin_name;

{ Get the next non-blank non-call... }
repeat get_x_token;
until cur_cmd<>spacer

;
 while true do  begin if (cur_cmd>other_char)or(cur_chr>255) then {not a character}
    begin back_input; goto done;
    end;
  {If |cur_chr| is a space and we're not scanning a token list, check
   whether we're at the end of the buffer. Otherwise we end up adding
   spurious spaces to file names in some cases.}
  if (cur_chr={" "=}32) and (cur_input.state_field <>token_list) and (cur_input.loc_field >cur_input.limit_field ) then goto done;
  if not more_name(cur_chr) then goto done;
  get_x_token;
  end;
end;
done: end_name; name_in_progress:=false;
warning_index := save_warning_index; {restore |warning_index|}
end;



{ 529. }

{tangle:pos tex.web:10243:1: }

{ Here is a routine that manufactures the output file names, assuming that
|job_name<>0|. It ignores and changes the current settings of |cur_area|
and |cur_ext|. } procedure pack_job_name( s:str_number); {|s = ".log"|, |".dvi"|, or
  |format_extension|}
begin cur_area:={""=}335; cur_ext:=s;
cur_name:=job_name; pack_file_name(cur_name,cur_area,cur_ext) ;
end;



{ 530. }

{tangle:pos tex.web:10255:1: }

{ If some trouble arises when \TeX\ tries to open a file, the following
routine calls upon the user to supply another file name. Parameter~|s|
is used in the error message to identify the type of file; parameter~|e|
is the default extension if none is given. Upon exit from the routine,
variables |cur_name|, |cur_area|, |cur_ext|, and |name_of_file| are
ready for another attempt at file opening. } procedure prompt_file_name( s, e:str_number);
label done;
var k:0..buf_size; {index into |buffer|}
 saved_cur_name:str_number; {to catch empty terminal input}
 saved_cur_ext:str_number; {to catch empty terminal input}
 saved_cur_area:str_number; {to catch empty terminal input}
begin if interaction=scroll_mode then    ;
if s={"input file name"=}795 then begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"I can't find file `"=} 796); end 
{ \xref[I can't find file x] }
else begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"I can't write on file `"=} 797); end ;
{ \xref[I can't write on file x] }
print_file_name(cur_name,cur_area,cur_ext); print({"'."=}798);
if (e={".tex"=}799) or (e={""=}335) then show_context;
print_ln; print_c_string(prompt_file_name_help_msg);
if (e<>{""=}335) then
  begin
    print({"; default file extension is `"=}800); print(e); print({"'"=}39);
  end;
print({")"=}41); print_ln;
print_nl({"Please type another "=}801); print(s);
{ \xref[Please type...] }
if interaction<scroll_mode then
  fatal_error({"*** (job aborted, file error in nonstop mode)"=}802);
{ \xref[job aborted, file error...] }
saved_cur_name:=cur_name;
saved_cur_ext:=cur_ext;
saved_cur_area:=cur_area;
   ; begin    ; print({": "=} 576); term_input; end ; 
{ Scan file name in the buffer }
begin begin_name; k:=first;
while (buffer[k]={" "=}32)and(k<last) do incr(k);
 while true do    begin if k=last then goto done;
  if not more_name(buffer[k]) then goto done;
  incr(k);
  end;
done:end_name;
end

;
if ((str_start[ cur_name+1]-str_start[ cur_name]) =0) and (cur_ext={""=}335) and (cur_area={""=}335) then
  begin
    cur_name:=saved_cur_name;
    cur_ext:=saved_cur_ext;
    cur_area:=saved_cur_area;
  end
else
  if cur_ext={""=}335 then cur_ext:=e;
pack_file_name(cur_name,cur_area,cur_ext) ;
end;



{ 534. }

{tangle:pos tex.web:10310:1: }

{ The |open_log_file| routine is used to open the transcript file and to help
it catch up to what has previously been printed on the terminal. } procedure open_log_file;
var old_setting:0..max_selector; {previous |selector| setting}
 k:0..buf_size; {index into |months| and |buffer|}
 l:0..buf_size; {end of first input line}
 months:const_cstring;
begin old_setting:=selector;
if job_name=0 then job_name:=get_job_name({"texput"=}805);
{ \xref[texput] }
pack_job_name({".fls"=}806);
recorder_change_filename(stringcast(name_of_file+1));
pack_job_name({".log"=}807);
while not a_open_out(log_file) do 
{ Try to get a different log file name }
begin selector:=term_only;
prompt_file_name({"transcript file name"=}809,{".log"=}807);
end

;
 texmf_log_name :=a_make_name_string(log_file);
selector:=log_only; log_opened:=true;

{ Print the banner line, including the date and time }
begin
if src_specials_p or file_line_error_style_p or parse_first_line_p
then
  write(log_file, 'This is GoTeXk, Version 3.141592653 (gotex v0.0-prerelease)'  ) 
else
  write(log_file, 'This is GoTeX, Version 3.141592653 (gotex v0.0-prerelease)'  ) ;
write(log_file, version_string) ;
slow_print(format_ident); print({"  "=}810);
print_int(sys_day); print_char({" "=}32);
months := ' JANFEBMARAPRMAYJUNJULAUGSEPOCTNOVDEC';
for k:=3*sys_month-2 to 3*sys_month do write(log_file, months[ k]) ;
print_char({" "=}32); print_int(sys_year); print_char({" "=}32);
print_two(sys_time div 60); print_char({":"=}58); print_two(sys_time mod 60);
if shellenabledp then begin
  writeln( log_file)  ;
  write(log_file,' ') ;
  if restrictedshell then begin
    write(log_file,'restricted ') ;
  end;
  write(log_file,'\write18 enabled.') 
  end;
if src_specials_p then begin
  writeln( log_file)  ;
  write(log_file,' Source specials enabled.') 
  end;
if file_line_error_style_p then begin
  writeln( log_file)  ;
  write(log_file,' file:line:error style messages enabled.') 
  end;
if parse_first_line_p then begin
  writeln( log_file)  ;
  write(log_file,' %&-line parsing enabled.') ;
  end;
if translate_filename then begin
  writeln( log_file)  ;
  write(log_file,' (') ;
  fputs(translate_filename, log_file);
  write(log_file,')') ;
  end;
end

;
if mltex_enabled_p then
  begin writeln( log_file)  ; write(log_file,'MLTeX v2.2 enabled') ;
  end;
input_stack[input_ptr]:=cur_input; {make sure bottom level is in memory}
print_nl({"**"=}808);
{ \xref[**] }
l:=input_stack[0].limit_field; {last position of first line}
if buffer[l]=eqtb[int_base+ end_line_char_code].int   then decr(l);
for k:=1 to l do print(buffer[k]);
print_ln; {now the transcript file contains the first line of input}
selector:=old_setting+2; {|log_only| or |term_and_log|}
end;



{ 537. }

{tangle:pos tex.web:10365:1: }

{ Let's turn now to the procedure that is used to initiate file reading
when an `\.[\\input]' command is being processed.
Beware: For historic reasons, this code foolishly conserves a tiny bit
of string pool space; but that can confuse the interactive `\.E' option.
\xref[system dependencies] } procedure start_input; {\TeX\ will \.[\\input] something}
label done;
var temp_str: str_number;
begin scan_file_name; {set |cur_name| to desired file name}
pack_file_name(cur_name,cur_area,cur_ext) ;
 while true do  begin
  begin_file_reading; {set up |cur_file| and new level of input}
  tex_input_type := 1; {Tell |open_input| we are \.[\\input].}
  {Kpathsea tries all the various ways to get the file.}
  if kpse_in_name_ok(stringcast(name_of_file+1))
     and a_open_in(input_file[cur_input.index_field ] , kpse_tex_format) then
    goto done;
  end_file_reading; {remove the level that didn't work}
  prompt_file_name({"input file name"=}795,{""=}335);
  end;
done: cur_input.name_field :=a_make_name_string(input_file[cur_input.index_field ] );
source_filename_stack[in_open]:=cur_input.name_field ;
full_source_filename_stack[in_open]:=make_full_name_string;
if cur_input.name_field =str_ptr-1 then {we can try to conserve string pool space now}
  begin temp_str:=search_string(cur_input.name_field );
  if temp_str>0 then
    begin cur_input.name_field :=temp_str; begin decr(str_ptr); pool_ptr:=str_start[str_ptr]; end ;
    end;
  end;
if job_name=0 then
  begin job_name:=get_job_name(cur_name); open_log_file;
  end; {|open_log_file| doesn't |show_context|, so |limit|
    and |loc| needn't be set to meaningful values yet}
if term_offset+(str_start[ full_source_filename_stack[ in_open]+1]-str_start[ full_source_filename_stack[ in_open]]) >max_print_line-2
then print_ln
else if (term_offset>0)or(file_offset>0) then print_char({" "=}32);
print_char({"("=}40); incr(open_parens);
slow_print(full_source_filename_stack[in_open]);  fflush (stdout ) ;
cur_input.state_field :=new_line;


{ Read the first line of the new file }
begin line:=1;
if input_ln(input_file[cur_input.index_field ] ,false) then  ;
firm_up_the_line;
if  (eqtb[int_base+ end_line_char_code].int  <0)or(eqtb[int_base+ end_line_char_code].int  >255)  then decr(cur_input.limit_field )
else  buffer[cur_input.limit_field ]:=eqtb[int_base+ end_line_char_code].int  ;
first:=cur_input.limit_field +1; cur_input.loc_field :=cur_input.start_field ;
end

;
end;



{ 540. }

{tangle:pos tex.web:10433:1: }

{ The first 24 bytes (6 words) of a \.[TFM] file contain twelve 16-bit
integers that give the lengths of the various subsequent portions
of the file. These twelve integers are, in order:
$$\vbox[\halign[\hfil#&$\null=\null$#\hfil\cr
|lf|&length of the entire file, in words;\cr
|lh|&length of the header data, in words;\cr
|bc|&smallest character code in the font;\cr
|ec|&largest character code in the font;\cr
|nw|&number of words in the width table;\cr
|nh|&number of words in the height table;\cr
|nd|&number of words in the depth table;\cr
|ni|&number of words in the italic correction table;\cr
|nl|&number of words in the lig/kern table;\cr
|nk|&number of words in the kern table;\cr
|ne|&number of words in the extensible character table;\cr
|np|&number of font parameter words.\cr]]$$
They are all nonnegative and less than $2^[15]$. We must have |bc-1<=ec<=255|,
and
$$\hbox[|lf=6+lh+(ec-bc+1)+nw+nh+nd+ni+nl+nk+ne+np|.]$$
Note that a font may contain as many as 256 characters (if |bc=0| and |ec=255|),
and as few as 0 characters (if |bc=ec+1|).

Incidentally, when two or more 8-bit bytes are combined to form an integer of
16 or more bits, the most significant bytes appear first in the file.
This is called BigEndian order.
 \xref[BigEndian order] }

{ 541. }

{tangle:pos tex.web:10460:1: }

{ The rest of the \.[TFM] file may be regarded as a sequence of ten data
arrays having the informal specification
$$\def\arr$[#1]#2$[\&[array] $[#1]$ \&[of] #2]
\vbox[\halign[\hfil\\[#]&$\,:\,$\arr#\hfil\cr
header&|[0..lh-1]\\[stuff]|\cr
char\_info&|[bc..ec]char_info_word|\cr
width&|[0..nw-1]fix_word|\cr
height&|[0..nh-1]fix_word|\cr
depth&|[0..nd-1]fix_word|\cr
italic&|[0..ni-1]fix_word|\cr
lig\_kern&|[0..nl-1]lig_kern_command|\cr
kern&|[0..nk-1]fix_word|\cr
exten&|[0..ne-1]extensible_recipe|\cr
param&|[1..np]fix_word|\cr]]$$
The most important data type used here is a | fix_word|, which is
a 32-bit representation of a binary fraction. A |fix_word| is a signed
quantity, with the two's complement of the entire word used to represent
negation. Of the 32 bits in a |fix_word|, exactly 12 are to the left of the
binary point; thus, the largest |fix_word| value is $2048-2^[-20]$, and
the smallest is $-2048$. We will see below, however, that all but two of
the |fix_word| values must lie between $-16$ and $+16$. }

{ 542. }

{tangle:pos tex.web:10482:1: }

{ The first data array is a block of header information, which contains
general facts about the font. The header must contain at least two words,
|header[0]| and |header[1]|, whose meaning is explained below.
Additional header information of use to other software routines might
also be included, but \TeX82 does not need to know about such details.
For example, 16 more words of header information are in use at the Xerox
Palo Alto Research Center; the first ten specify the character coding
scheme used (e.g., `\.[XEROX text]' or `\.[TeX math symbols]'), the next five
give the font identifier (e.g., `\.[HELVETICA]' or `\.[CMSY]'), and the
last gives the ``face byte.'' The program that converts \.[DVI] files
to Xerox printing format gets this information by looking at the \.[TFM]
file, which it needs to read anyway because of other information that
is not explicitly repeated in \.[DVI]~format.

\yskip\hang|header[0]| is a 32-bit check sum that \TeX\ will copy into
the \.[DVI] output file. Later on when the \.[DVI] file is printed,
possibly on another computer, the actual font that gets used is supposed
to have a check sum that agrees with the one in the \.[TFM] file used by
\TeX. In this way, users will be warned about potential incompatibilities.
(However, if the check sum is zero in either the font file or the \.[TFM]
file, no check is made.)  The actual relation between this check sum and
the rest of the \.[TFM] file is not important; the check sum is simply an
identification number with the property that incompatible fonts almost
always have distinct check sums.
\xref[check sum]

\yskip\hang|header[1]| is a |fix_word| containing the design size of
the font, in units of \TeX\ points. This number must be at least 1.0; it is
fairly arbitrary, but usually the design size is 10.0 for a ``10 point''
font, i.e., a font that was designed to look best at a 10-point size,
whatever that really means. When a \TeX\ user asks for a font
`\.[at] $\delta$ \.[pt]', the effect is to override the design size
and replace it by $\delta$, and to multiply the $x$ and~$y$ coordinates
of the points in the font image by a factor of $\delta$ divided by the
design size.  [\sl All other dimensions in the\/ \.[TFM] file are
|fix_word|\kern-1pt\ numbers in design-size units], with the exception of
|param[1]| (which denotes the slant ratio). Thus, for example, the value
of |param[6]|, which defines the \.[em] unit, is often the |fix_word| value
$2^[20]=1.0$, since many fonts have a design size equal to one em.
The other dimensions must be less than 16 design-size units in absolute
value; thus, |header[1]| and |param[1]| are the only |fix_word|
entries in the whole \.[TFM] file whose first byte might be something
besides 0 or 255. }

{ 543. }

{tangle:pos tex.web:10526:1: }

{ Next comes the |char_info| array, which contains one | char_info_word|
per character. Each word in this part of the file contains six fields
packed into four bytes as follows.

\yskip\hang first byte: | width_index| (8 bits)\par
\hang second byte: | height_index| (4 bits) times 16, plus | depth_index|
  (4~bits)\par
\hang third byte: | italic_index| (6 bits) times 4, plus | tag|
  (2~bits)\par
\hang fourth byte: | remainder| (8 bits)\par
\yskip\noindent
The actual width of a character is \\[width]|[width_index]|, in design-size
units; this is a device for compressing information, since many characters
have the same width. Since it is quite common for many characters
to have the same height, depth, or italic correction, the \.[TFM] format
imposes a limit of 16 different heights, 16 different depths, and
64 different italic corrections.

 \xref[italic correction]
The italic correction of a character has two different uses.
(a)~In ordinary text, the italic correction is added to the width only if
the \TeX\ user specifies `\.[\\/]' after the character.
(b)~In math formulas, the italic correction is always added to the width,
except with respect to the positioning of subscripts.

Incidentally, the relation $\\[width][0]=\\[height][0]=\\[depth][0]=
\\[italic][0]=0$ should always hold, so that an index of zero implies a
value of zero.  The |width_index| should never be zero unless the
character does not exist in the font, since a character is valid if and
only if it lies between |bc| and |ec| and has a nonzero |width_index|. }

{ 544. }

{tangle:pos tex.web:10557:1: }

{ The |tag| field in a |char_info_word| has four values that explain how to
interpret the |remainder| field.

\yskip\hangg|tag=0| (|no_tag|) means that |remainder| is unused.\par
\hangg|tag=1| (|lig_tag|) means that this character has a ligature/kerning
program starting at position |remainder| in the |lig_kern| array.\par
\hangg|tag=2| (|list_tag|) means that this character is part of a chain of
characters of ascending sizes, and not the largest in the chain.  The
|remainder| field gives the character code of the next larger character.\par
\hangg|tag=3| (|ext_tag|) means that this character code represents an
extensible character, i.e., a character that is built up of smaller pieces
so that it can be made arbitrarily large. The pieces are specified in
| exten[remainder]|.\par
\yskip\noindent
Characters with |tag=2| and |tag=3| are treated as characters with |tag=0|
unless they are used in special circumstances in math formulas. For example,
the \.[\\sum] operation looks for a |list_tag|, and the \.[\\left]
operation looks for both |list_tag| and |ext_tag|. }

{ 545. }

{tangle:pos tex.web:10581:1: }

{ The |lig_kern| array contains instructions in a simple programming language
that explains what to do for special letter pairs. Each word in this array is a
| lig_kern_command| of four bytes.

\yskip\hang first byte: |skip_byte|, indicates that this is the final program
  step if the byte is 128 or more, otherwise the next step is obtained by
  skipping this number of intervening steps.\par
\hang second byte: |next_char|, ``if |next_char| follows the current character,
  then perform the operation and stop, otherwise continue.''\par
\hang third byte: |op_byte|, indicates a ligature step if less than~128,
  a kern step otherwise.\par
\hang fourth byte: |remainder|.\par
\yskip\noindent
In a kern step, an
additional space equal to |kern[256*(op_byte-128)+remainder]| is inserted
between the current character and |next_char|. This amount is
often negative, so that the characters are brought closer together
by kerning; but it might be positive.

There are eight kinds of ligature steps, having |op_byte| codes $4a+2b+c$ where
$0\le a\le b+c$ and $0\le b,c\le1$. The character whose code is
|remainder| is inserted between the current character and |next_char|;
then the current character is deleted if $b=0$, and |next_char| is
deleted if $c=0$; then we pass over $a$~characters to reach the next
current character (which may have a ligature/kerning program of its own).

If the very first instruction of the |lig_kern| array has |skip_byte=255|,
the |next_char| byte is the so-called boundary character of this font;
the value of |next_char| need not lie between |bc| and~|ec|.
If the very last instruction of the |lig_kern| array has |skip_byte=255|,
there is a special ligature/kerning program for a boundary character at the
left, beginning at location |256*op_byte+remainder|.
The interpretation is that \TeX\ puts implicit boundary characters
before and after each consecutive string of characters from the same font.
These implicit characters do not appear in the output, but they can affect
ligatures and kerning.

If the very first instruction of a character's |lig_kern| program has
|skip_byte>128|, the program actually begins in location
|256*op_byte+remainder|. This feature allows access to large |lig_kern|
arrays, because the first instruction must otherwise
appear in a location |<=255|.

Any instruction with |skip_byte>128| in the |lig_kern| array must satisfy
the condition
$$\hbox[|256*op_byte+remainder<nl|.]$$
If such an instruction is encountered during
normal program execution, it denotes an unconditional halt; no ligature
or kerning command is performed. }

{ 546. }

{tangle:pos tex.web:10638:1: }

{ Extensible characters are specified by an | extensible_recipe|, which
consists of four bytes called | top|, | mid|, | bot|, and | rep| (in this
order). These bytes are the character codes of individual pieces used to
build up a large symbol.  If |top|, |mid|, or |bot| are zero, they are not
present in the built-up result. For example, an extensible vertical line is
like an extensible bracket, except that the top and bottom pieces are missing.

Let $T$, $M$, $B$, and $R$ denote the respective pieces, or an empty box
if the piece isn't present. Then the extensible characters have the form
$TR^kMR^kB$ from top to bottom, for some |k>=0|, unless $M$ is absent;
in the latter case we can have $TR^kB$ for both even and odd values of~|k|.
The width of the extensible character is the width of $R$; and the
height-plus-depth is the sum of the individual height-plus-depths of the
components used, since the pieces are butted together in a vertical list. }

{ 547. }

{tangle:pos tex.web:10658:1: }

{ The final portion of a \.[TFM] file is the |param| array, which is another
sequence of |fix_word| values.

\yskip\hang|param[1]=slant| is the amount of italic slant, which is used
to help position accents. For example, |slant=.25| means that when you go
up one unit, you also go .25 units to the right. The |slant| is a pure
number; it's the only |fix_word| other than the design size itself that is
not scaled by the design size.

\hang|param[2]=space| is the normal spacing between words in text.
Note that character |" "| in the font need not have anything to do with
blank spaces.

\hang|param[3]=space_stretch| is the amount of glue stretching between words.

\hang|param[4]=space_shrink| is the amount of glue shrinking between words.

\hang|param[5]=x_height| is the size of one ex in the font; it is also
the height of letters for which accents don't have to be raised or lowered.

\hang|param[6]=quad| is the size of one em in the font.

\hang|param[7]=extra_space| is the amount added to |param[2]| at the
ends of sentences.

\yskip\noindent
If fewer than seven parameters are present, \TeX\ sets the missing parameters
to zero. Fonts used for math symbols are required to have
additional parameter information, which is explained later. }

{ 554. }

{tangle:pos tex.web:10804:1: }

{ Of course we want to define macros that suppress the detail of how font
information is actually packed, so that we don't have to write things like
$$\hbox[|font_info[width_base[f]+font_info[char_base[f]+c].qqqq.b0].sc|]$$
too often. The \.[WEB] definitions here make |char_info(f)(c)| the
|four_quarters| word of font information corresponding to character
|c| of font |f|. If |q| is such a word, |char_width(f)(q)| will be
the character's width; hence the long formula above is at least
abbreviated to
$$\hbox[|char_width(f)(char_info(f)(c))|.]$$
Usually, of course, we will fetch |q| first and look at several of its
fields at the same time.

The italic correction of a character will be denoted by
|char_italic(f)(q)|, so it is analogous to |char_width|.  But we will get
at the height and depth in a slightly different way, since we usually want
to compute both height and depth if we want either one.  The value of
|height_depth(q)| will be the 8-bit quantity
$$b=|height_index|\times16+|depth_index|,$$ and if |b| is such a byte we
will write |char_height(f)(b)| and |char_depth(f)(b)| for the height and
depth of the character |c| for which |q=char_info(f)(c)|. Got that?

The tag field will be called |char_tag(q)|; the remainder byte will be
called |rem_byte(q)|, using a macro that we have already defined above.

Access to a character's |width|, |height|, |depth|, and |tag| fields is
part of \TeX's inner loop, so we want these macros to produce code that is
as fast as possible under the circumstances.
\xref[inner loop]

ML\TeX[] will assume that a character |c| exists iff either exists in
the current font or a character substitution definition for this
character was defined using \.[\\charsubdef].  To avoid the
distinction between these two cases, ML\TeX[] introduces the notion
``effective character'' of an input character |c|.  If |c| exists in
the current font, the effective character of |c| is the character |c|
itself.  If it doesn't exist but a character substitution is defined,
the effective character of |c| is the base character defined in the
character substitution.  If there is an effective character for a
non-existing character |c|, the ``virtual character'' |c| will get
appended to the horizontal lists.

The effective character is used within |char_info| to access
appropriate character descriptions in the font.  For example, when
calculating the width of a box, ML\TeX[] will use the metrics of the
effective characters.  For the case of a substitution, ML\TeX[] uses
the metrics of the base character, ignoring the metrics of the accent
character.

If character substitutions are changed, it will be possible that a
character |c| neither exists in a font nor there is a valid character
substitution for |c|.  To handle these cases |effective_char| should
be called with its first argument set to |true| to ensure that it
will still return an existing character in the font.  If neither |c|
nor the substituted base character in the current character
substitution exists, |effective_char| will output a warning and
return the character |font_bc[f]| (which is incorrect, but can not be
changed within the current framework).

Sometimes character substitutions are unwanted, therefore the
original definition of |char_info| can be used using the macro
|orig_char_info|.  Operations in which character substitutions should
be avoided are, for example, loading a new font and checking the font
metric information in this font, and character accesses in math mode. }

{ 557. }

{tangle:pos tex.web:10858:1: }

{ Here are some macros that help process ligatures and kerns.
We write |char_kern(f)(j)| to find the amount of kerning specified by
kerning command~|j| in font~|f|. If |j| is the |char_info| for a character
with a ligature/kern program, the first instruction of that program is either
|i=font_info[lig_kern_start(f)(j)]| or |font_info[lig_kern_restart(f)(i)]|,
depending on whether or not |skip_byte(i)<=stop_flag|.

The constant |kern_base_offset| should be simplified, for \PASCAL\ compilers
that do not do local optimization.
\xref[system dependencies] }

{ 560. }

{tangle:pos tex.web:10892:1: }

{ \TeX\ checks the information of a \.[TFM] file for validity as the
file is being read in, so that no further checks will be needed when
typesetting is going on. The somewhat tedious subroutine that does this
is called |read_font_info|. It has four parameters: the user font
identifier~|u|, the file name and area strings |nom| and |aire|, and the
``at'' size~|s|. If |s|~is negative, it's the negative of a scale factor
to be applied to the design size; |s=-1000| is the normal case.
Otherwise |s| will be substituted for the design size; in this
case, |s| must be positive and less than $2048\rm\,pt$
(i.e., it must be less than $2^[27]$ when considered as an integer).

The subroutine opens and closes a global file variable called |tfm_file|.
It returns the value of the internal font number that was just loaded.
If an error is detected, an error message is issued and no font
information is stored; |null_font| is returned in this case. } { \4 }
{ Declare additional functions for ML\TeX }
function effective_char( err_p:boolean;
                         f:internal_font_number; c:quarterword):integer;
label found;
var base_c: integer; {or |eightbits|: replacement base character}
 result: integer; {or |quarterword|}
begin result:=c;  {return |c| unless it does not exist in the font}
if not mltex_enabled_p then goto found;
if font_ec[f]>= c  then if font_bc[f]<= c  then
  if ( font_info[char_base[  f]+  c].qqqq .b0>min_quarterword)  then  {N.B.: not |char_info|(f)(c)}
    goto found;
if  c >=eqtb[int_base+ char_sub_def_min_code].int   then if  c <=eqtb[int_base+ char_sub_def_max_code].int   then
  if ( eqtb[  char_sub_code_base+         c ].hh.rh   > 0 )  then
    begin base_c:=(  eqtb[  char_sub_code_base+           c ].hh.rh     mod 256) ;
    result:= base_c ;  {return |base_c|}
    if not err_p then goto found;
    if font_ec[f]>=base_c then if font_bc[f]<=base_c then
      if ( font_info[char_base[  f]+     base_c ].qqqq .b0>min_quarterword)  then goto found;
    end;
if err_p then  {print error and return existing character?}
  begin begin_diagnostic;
  print_nl({"Missing character: There is no "=}837); print({"substitution for "=}1330);
{ \xref[Missing character] }
   print ( c ); print({" in font "=}838);
  slow_print(font_name[f]); print_char({"!"=}33); end_diagnostic(false);
  result:= font_bc[ f] ; {N.B.: not non-existing character |c|!}
  end;
found: effective_char:=result;
end;



function effective_char_info( f:internal_font_number;
                              c:quarterword):four_quarters;
label exit;
var ci:four_quarters; {character information bytes for |c|}
 base_c:integer; {or |eightbits|: replacement base character}
begin if not mltex_enabled_p then
  begin effective_char_info:=font_info[char_base[ f]+ c].qqqq ;  goto exit ;
  end;
if font_ec[f]>= c  then if font_bc[f]<= c  then
  begin ci:=font_info[char_base[ f]+ c].qqqq ;  {N.B.: not |char_info|(f)(c)}
  if ( ci.b0>min_quarterword)  then
    begin effective_char_info:=ci;  goto exit ;
    end;
  end;
if  c >=eqtb[int_base+ char_sub_def_min_code].int   then if  c <=eqtb[int_base+ char_sub_def_max_code].int   then
  if ( eqtb[  char_sub_code_base+         c ].hh.rh   > 0 )  then
    begin {|effective_char_info:=char_info(f)(qi(char_list_char(qo(c))));|}
    base_c:=(  eqtb[  char_sub_code_base+           c ].hh.rh     mod 256) ;
    if font_ec[f]>=base_c then if font_bc[f]<=base_c then
      begin ci:=font_info[char_base[ f]+   base_c ].qqqq ;  {N.B.: not |char_info|(f)(c)}
      if ( ci.b0>min_quarterword)  then
        begin effective_char_info:=ci;  goto exit ;
        end;
      end;
    end;
effective_char_info:=null_character;
exit:end;





function read_font_info( u:halfword ; nom, aire:str_number;
   s:scaled):internal_font_number; {input a \.[TFM] file}
label done,bad_tfm,not_found;
var k:font_index; {index into |font_info|}
 name_too_long:boolean; {|nom| or |aire| exceeds 255 bytes?}
 file_opened:boolean; {was |tfm_file| successfully opened?}
 lf, lh, bc, ec, nw, nh, nd, ni, nl, nk, ne, np:halfword;
  {sizes of subfiles}
 f:internal_font_number; {the new font's number}
 g:internal_font_number; {the number to return}
 a, b, c, d:eight_bits; {byte variables}
 qw:four_quarters; sw:scaled; {accumulators}
 bch_label:integer; {left boundary start location, or infinity}
 bchar:0..256; {boundary character, or 256}
 z:scaled; {the design size or the ``at'' size}
 alpha:integer; beta:1..16;
  {auxiliary quantities used in fixed-point multiplication}
begin g:=font_base ;


{ Read and check the font data; |abort| if the \.[TFM] file is malformed; if there's no room for this font, say so and |goto done|; otherwise |incr(font_ptr)| and |goto done| }

{ Open |tfm_file| for input }
file_opened:=false;
name_too_long:=((str_start[ nom+1]-str_start[ nom]) >255)or((str_start[ aire+1]-str_start[ aire]) >255);
if name_too_long then goto bad_tfm ;
{|kpse_find_file| will append the |".tfm"|, and avoid searching the disk
 before the font alias files as well.}
pack_file_name(nom,aire,{""=}335);
if not b_open_in(tfm_file) then goto bad_tfm ;
file_opened:=true

;

{ Read the [\.[TFM]] size fields }
begin begin  lf:=tfm_temp ; if  lf>127 then goto bad_tfm ; tfm_temp:=getc(tfm_file) ;  lf:= lf*{0400=}256+tfm_temp ; end ;
tfm_temp:=getc(tfm_file) ; begin  lh:=tfm_temp ; if  lh>127 then goto bad_tfm ; tfm_temp:=getc(tfm_file) ;  lh:= lh*{0400=}256+tfm_temp ; end ;
tfm_temp:=getc(tfm_file) ; begin  bc:=tfm_temp ; if  bc>127 then goto bad_tfm ; tfm_temp:=getc(tfm_file) ;  bc:= bc*{0400=}256+tfm_temp ; end ;
tfm_temp:=getc(tfm_file) ; begin  ec:=tfm_temp ; if  ec>127 then goto bad_tfm ; tfm_temp:=getc(tfm_file) ;  ec:= ec*{0400=}256+tfm_temp ; end ;
if (bc>ec+1)or(ec>255) then goto bad_tfm ;
if bc>255 then {|bc=256| and |ec=255|}
  begin bc:=1; ec:=0;
  end;
tfm_temp:=getc(tfm_file) ; begin  nw:=tfm_temp ; if  nw>127 then goto bad_tfm ; tfm_temp:=getc(tfm_file) ;  nw:= nw*{0400=}256+tfm_temp ; end ;
tfm_temp:=getc(tfm_file) ; begin  nh:=tfm_temp ; if  nh>127 then goto bad_tfm ; tfm_temp:=getc(tfm_file) ;  nh:= nh*{0400=}256+tfm_temp ; end ;
tfm_temp:=getc(tfm_file) ; begin  nd:=tfm_temp ; if  nd>127 then goto bad_tfm ; tfm_temp:=getc(tfm_file) ;  nd:= nd*{0400=}256+tfm_temp ; end ;
tfm_temp:=getc(tfm_file) ; begin  ni:=tfm_temp ; if  ni>127 then goto bad_tfm ; tfm_temp:=getc(tfm_file) ;  ni:= ni*{0400=}256+tfm_temp ; end ;
tfm_temp:=getc(tfm_file) ; begin  nl:=tfm_temp ; if  nl>127 then goto bad_tfm ; tfm_temp:=getc(tfm_file) ;  nl:= nl*{0400=}256+tfm_temp ; end ;
tfm_temp:=getc(tfm_file) ; begin  nk:=tfm_temp ; if  nk>127 then goto bad_tfm ; tfm_temp:=getc(tfm_file) ;  nk:= nk*{0400=}256+tfm_temp ; end ;
tfm_temp:=getc(tfm_file) ; begin  ne:=tfm_temp ; if  ne>127 then goto bad_tfm ; tfm_temp:=getc(tfm_file) ;  ne:= ne*{0400=}256+tfm_temp ; end ;
tfm_temp:=getc(tfm_file) ; begin  np:=tfm_temp ; if  np>127 then goto bad_tfm ; tfm_temp:=getc(tfm_file) ;  np:= np*{0400=}256+tfm_temp ; end ;
if lf<>6+lh+(ec-bc+1)+nw+nh+nd+ni+nl+nk+ne+np then goto bad_tfm ;
if (nw=0)or(nh=0)or(nd=0)or(ni=0) then goto bad_tfm ;
end

;

{ Use size fields to allocate font information }
lf:=lf-6-lh; {|lf| words should be loaded into |font_info|}
if np<7 then lf:=lf+7-np; {at least seven parameters will appear}
if (font_ptr=font_max)or(fmem_ptr+lf>font_mem_size) then
  
{ Apologize for not loading the font, |goto done| }
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Font "=} 812); end ; sprint_cs(u); print_char({"="=}61); print_file_name(nom,aire,{""=}335); if s>=0 then begin print({" at "=}751); print_scaled(s); print({"pt"=}402); end else if s<>-1000 then begin print({" scaled "=}813); print_int(-s); end ;
print({" not loaded: Not enough room left"=}822);
{ \xref[Font x=xx not loaded...] }
 begin help_ptr:=4; help_line[3]:={"I'm afraid I won't be able to make use of this font,"=} 823; help_line[2]:={"because my memory for character-size data is too small."=} 824; help_line[1]:={"If you're really stuck, ask a wizard to enlarge me."=} 825; help_line[0]:={"Or maybe try `I\font<same font id>=<name of loaded font>'."=} 826; end ;
error; goto done;
end

;
f:=font_ptr+1;
char_base[f]:=fmem_ptr-bc;
width_base[f]:=char_base[f]+ec+1;
height_base[f]:=width_base[f]+nw;
depth_base[f]:=height_base[f]+nh;
italic_base[f]:=depth_base[f]+nd;
lig_kern_base[f]:=italic_base[f]+ni;
kern_base[f]:=lig_kern_base[f]+nl-256*(128+min_quarterword) ;
exten_base[f]:=kern_base[f]+256*(128+min_quarterword) +nk;
param_base[f]:=exten_base[f]+ne

;

{ Read the [\.[TFM]] header }
begin if lh<2 then goto bad_tfm ;
begin tfm_temp:=getc(tfm_file) ; a:=tfm_temp ; qw.b0:= a ; tfm_temp:=getc(tfm_file) ; b:=tfm_temp ; qw.b1:= b ; tfm_temp:=getc(tfm_file) ; c:=tfm_temp ; qw.b2:= c ; tfm_temp:=getc(tfm_file) ; d:=tfm_temp ; qw.b3:= d ;  font_check[ f]:=qw; end ;
tfm_temp:=getc(tfm_file) ; begin  z:=tfm_temp ; if  z>127 then goto bad_tfm ; tfm_temp:=getc(tfm_file) ;  z:= z*{0400=}256+tfm_temp ; end ; {this rejects a negative design size}
tfm_temp:=getc(tfm_file) ; z:=z*{0400=}256+tfm_temp ; tfm_temp:=getc(tfm_file) ; z:=(z*{020=}16)+(tfm_temp  div{020=}16);
if z< {0200000=}65536  then goto bad_tfm ;
while lh>2 do
  begin tfm_temp:=getc(tfm_file) ;tfm_temp:=getc(tfm_file) ;tfm_temp:=getc(tfm_file) ;tfm_temp:=getc(tfm_file) ;decr(lh); {ignore the rest of the header}
  end;
font_dsize[f]:=z;
if s<>-1000 then
  if s>=0 then z:=s
  else begin
    save_arith_error:=arith_error;
    sw:=z; z:=xn_over_d(z,-s,1000);
    if arith_error or z>={01000000000=}134217728 then begin
       begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Font "=} 812); end ; sprint_cs(u); print_char({"="=}61); print_file_name(nom,aire,{""=}335); if s>=0 then begin print({" at "=}751); print_scaled(s); print({"pt"=}402); end else if s<>-1000 then begin print({" scaled "=}813); print_int(-s); end ; print({" scaled to 2048pt or higher"=}827);
        begin help_ptr:=1; help_line[0]:={"I will ignore the scaling factor."=} 828; end ; error; z:=sw;
       end;
    arith_error:=save_arith_error;
  end;
font_size[f]:=z;
end

;

{ Read character data }
for k:=fmem_ptr to width_base[f]-1 do
  begin begin tfm_temp:=getc(tfm_file) ; a:=tfm_temp ; qw.b0:= a ; tfm_temp:=getc(tfm_file) ; b:=tfm_temp ; qw.b1:= b ; tfm_temp:=getc(tfm_file) ; c:=tfm_temp ; qw.b2:= c ; tfm_temp:=getc(tfm_file) ; d:=tfm_temp ; qw.b3:= d ;  font_info[ k]. qqqq:=qw; end ;
  if (a>=nw)or(b div {020=}16>=nh)or(b mod {020=}16>=nd)or
    (c div 4>=ni) then goto bad_tfm ;
  case c mod 4 of
  lig_tag: if d>=nl then goto bad_tfm ;
  ext_tag: if d>=ne then goto bad_tfm ;
  list_tag: 
{ Check for charlist cycle }
begin begin if ( d<bc)or( d>ec) then goto bad_tfm  end ;
while d<k+bc-fmem_ptr  do
  begin qw:=font_info[char_base[ f]+ d].qqqq ;
  {N.B.: not |qi(d)|, since |char_base[f]| hasn't been adjusted yet}
  if ((  qw. b2 ) mod 4) <>list_tag then goto not_found;
  d:=   qw.b3  ; {next character on the list}
  end;
if d=k+bc-fmem_ptr  then goto bad_tfm ; {yes, there's a cycle}
not_found:end

;
   else    {|no_tag|}
   end ;
  end

;

{ Read box dimensions }
begin 
{ Replace |z| by $|z|^\prime$ and compute $\alpha,\beta$ }
begin alpha:=16;
while z>={040000000=}8388608 do
  begin z:=z div 2; alpha:=alpha+alpha;
  end;
beta:=256 div alpha; alpha:=alpha*z;
end

;
for k:=width_base[f] to lig_kern_base[f]-1 do
  begin tfm_temp:=getc(tfm_file) ; a:=tfm_temp ; tfm_temp:=getc(tfm_file) ; b:=tfm_temp ; tfm_temp:=getc(tfm_file) ; c:=tfm_temp ; tfm_temp:=getc(tfm_file) ; d:=tfm_temp ; sw:=(((((d*z)div{0400=}256)+(c*z))div{0400=}256)+(b*z))div beta; if a=0 then  font_info[ k]. int :=sw else if a=255 then  font_info[ k]. int :=sw-alpha else goto bad_tfm ; end ;
if font_info[width_base[f]].int <>0 then goto bad_tfm ; {\\[width][0] must be zero}
if font_info[height_base[f]].int <>0 then goto bad_tfm ; {\\[height][0] must be zero}
if font_info[depth_base[f]].int <>0 then goto bad_tfm ; {\\[depth][0] must be zero}
if font_info[italic_base[f]].int <>0 then goto bad_tfm ; {\\[italic][0] must be zero}
end

;

{ Read ligature/kern program }
bch_label:={077777=}32767; bchar:=256;
if nl>0 then
  begin for k:=lig_kern_base[f] to kern_base[f]+256*(128+min_quarterword) -1 do
    begin begin tfm_temp:=getc(tfm_file) ; a:=tfm_temp ; qw.b0:= a ; tfm_temp:=getc(tfm_file) ; b:=tfm_temp ; qw.b1:= b ; tfm_temp:=getc(tfm_file) ; c:=tfm_temp ; qw.b2:= c ; tfm_temp:=getc(tfm_file) ; d:=tfm_temp ; qw.b3:= d ;  font_info[ k]. qqqq:=qw; end ;
    if a>128 then
      begin if 256*c+d>=nl then goto bad_tfm ;
      if a=255 then if k=lig_kern_base[f] then bchar:=b;
      end
    else begin if b<>bchar then {  } begin begin if (  b<bc)or(  b>ec) then goto bad_tfm  end ; qw:=font_info[char_base[ f]+  b].qqqq ; if not ( qw.b0>min_quarterword)  then goto bad_tfm ; end ;
      if c<128 then {  } begin begin if (  d<bc)or(  d>ec) then goto bad_tfm  end ; qw:=font_info[char_base[ f]+  d].qqqq ; if not ( qw.b0>min_quarterword)  then goto bad_tfm ; end  {check ligature}
      else if 256*(c-128)+d>=nk then goto bad_tfm ; {check kern}
      if a<128 then if k-lig_kern_base[f]+a+1>=nl then goto bad_tfm ;
      end;
    end;
  if a=255 then bch_label:=256*c+d;
  end;
for k:=kern_base[f]+256*(128+min_quarterword)  to exten_base[f]-1 do
  begin tfm_temp:=getc(tfm_file) ; a:=tfm_temp ; tfm_temp:=getc(tfm_file) ; b:=tfm_temp ; tfm_temp:=getc(tfm_file) ; c:=tfm_temp ; tfm_temp:=getc(tfm_file) ; d:=tfm_temp ; sw:=(((((d*z)div{0400=}256)+(c*z))div{0400=}256)+(b*z))div beta; if a=0 then  font_info[ k]. int :=sw else if a=255 then  font_info[ k]. int :=sw-alpha else goto bad_tfm ; end ;

;

{ Read extensible character recipes }
for k:=exten_base[f] to param_base[f]-1 do
  begin begin tfm_temp:=getc(tfm_file) ; a:=tfm_temp ; qw.b0:= a ; tfm_temp:=getc(tfm_file) ; b:=tfm_temp ; qw.b1:= b ; tfm_temp:=getc(tfm_file) ; c:=tfm_temp ; qw.b2:= c ; tfm_temp:=getc(tfm_file) ; d:=tfm_temp ; qw.b3:= d ;  font_info[ k]. qqqq:=qw; end ;
  if a<>0 then {  } begin begin if (  a<bc)or(  a>ec) then goto bad_tfm  end ; qw:=font_info[char_base[ f]+  a].qqqq ; if not ( qw.b0>min_quarterword)  then goto bad_tfm ; end ;
  if b<>0 then {  } begin begin if (  b<bc)or(  b>ec) then goto bad_tfm  end ; qw:=font_info[char_base[ f]+  b].qqqq ; if not ( qw.b0>min_quarterword)  then goto bad_tfm ; end ;
  if c<>0 then {  } begin begin if (  c<bc)or(  c>ec) then goto bad_tfm  end ; qw:=font_info[char_base[ f]+  c].qqqq ; if not ( qw.b0>min_quarterword)  then goto bad_tfm ; end ;
  {  } begin begin if (  d<bc)or(  d>ec) then goto bad_tfm  end ; qw:=font_info[char_base[ f]+  d].qqqq ; if not ( qw.b0>min_quarterword)  then goto bad_tfm ; end ;
  end

;

{ Read font parameters }
begin for k:=1 to np do
  if k=1 then {the |slant| parameter is a pure number}
    begin tfm_temp:=getc(tfm_file) ; sw:=tfm_temp ; if sw>127 then sw:=sw-256;
    tfm_temp:=getc(tfm_file) ; sw:=sw*{0400=}256+tfm_temp ; tfm_temp:=getc(tfm_file) ; sw:=sw*{0400=}256+tfm_temp ;
    tfm_temp:=getc(tfm_file) ; font_info[param_base[f]].int :=
      (sw*{020=}16)+(tfm_temp  div{020=}16);
    end
  else begin tfm_temp:=getc(tfm_file) ; a:=tfm_temp ; tfm_temp:=getc(tfm_file) ; b:=tfm_temp ; tfm_temp:=getc(tfm_file) ; c:=tfm_temp ; tfm_temp:=getc(tfm_file) ; d:=tfm_temp ; sw:=(((((d*z)div{0400=}256)+(c*z))div{0400=}256)+(b*z))div beta; if a=0 then  font_info[ param_base[ f]+ k- 1]. int :=sw else if a=255 then  font_info[ param_base[ f]+ k- 1]. int :=sw-alpha else goto bad_tfm ; end ;
if feof(tfm_file) then goto bad_tfm ;
for k:=np+1 to 7 do font_info[param_base[f]+k-1].int :=0;
end

;

{ Make final adjustments and |goto done| }
if np>=7 then font_params[f]:=np else font_params[f]:=7;
hyphen_char[f]:=eqtb[int_base+ default_hyphen_char_code].int  ; skew_char[f]:=eqtb[int_base+ default_skew_char_code].int  ;
if bch_label<nl then bchar_label[f]:=bch_label+lig_kern_base[f]
else bchar_label[f]:=non_address;
font_bchar[f]:= bchar ;
font_false_bchar[f]:= bchar ;
if bchar<=ec then if bchar>=bc then
  begin qw:=font_info[char_base[ f]+ bchar].qqqq ; {N.B.: not |qi(bchar)|}
  if ( qw.b0>min_quarterword)  then font_false_bchar[f]:= 256  ;
  end;
font_name[f]:=nom;
font_area[f]:=aire;
font_bc[f]:=bc; font_ec[f]:=ec; font_glue[f]:=-{0xfffffff=}268435455  ;
 char_base[f]:=  char_base[ f]  ;  width_base[f]:=  width_base[ f]  ;  lig_kern_base[f]:=  lig_kern_base[ f]  ;
 kern_base[f]:=  kern_base[ f]  ;  exten_base[f]:=  exten_base[ f]  ;
decr(param_base[f]);
fmem_ptr:=fmem_ptr+lf; font_ptr:=f; g:=f; goto done



;
bad_tfm: 
{ Report that the font won't be loaded }
begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Font "=} 812); end ; sprint_cs(u); print_char({"="=}61); print_file_name(nom,aire,{""=}335); if s>=0 then begin print({" at "=}751); print_scaled(s); print({"pt"=}402); end else if s<>-1000 then begin print({" scaled "=}813); print_int(-s); end ;
{ \xref[Font x=xx not loadable...] }
if file_opened then print({" not loadable: Bad metric (TFM) file"=}814)
else if name_too_long then print({" not loadable: Metric (TFM) file name too long"=}815)
else print({" not loadable: Metric (TFM) file not found"=}816);
 begin help_ptr:=5; help_line[4]:={"I wasn't able to read the size data for this font,"=} 817; help_line[3]:={"so I will ignore the font specification."=} 818; help_line[2]:={"[Wizards can fix TFM files using TFtoPL/PLtoTF.]"=} 819; help_line[1]:={"You might try inserting a different font spec;"=} 820; help_line[0]:={"e.g., type `I\font<same font id>=<substitute font name>'."=} 821; end ;
error

;
done: if file_opened then b_close(tfm_file);
read_font_info:=g;
end;



{ 564. }

{tangle:pos tex.web:10983:1: }

{ Note: A malformed \.[TFM] file might be shorter than it claims to be;
thus |eof(tfm_file)| might be true when |read_font_info| refers to
|tfm_file^| or when it says |get(tfm_file)|. If such circumstances
cause system error messages, you will have to defeat them somehow,
for example by defining |fget| to be `\ignorespaces|begin get(tfm_file);|
|if eof(tfm_file) then abort; end|\unskip'.
\xref[system dependencies] }

{ 581. }

{tangle:pos tex.web:11293:1: }

{ When \TeX\ wants to typeset a character that doesn't exist, the
character node is not created; thus the output routine can assume
that characters exist when it sees them. The following procedure
prints a warning message unless the user has suppressed it. } procedure char_warning( f:internal_font_number; c:eight_bits);
begin if eqtb[int_base+ tracing_lost_chars_code].int  >0 then
  begin begin_diagnostic;
  print_nl({"Missing character: There is no "=}837);
{ \xref[Missing character] }
   print (c); print({" in font "=}838);
  slow_print(font_name[f]); print_char({"!"=}33); end_diagnostic(false);
  end;
end;



{ 582. }

{tangle:pos tex.web:11308:1: }

{ Here is a function that returns a pointer to a character node for a
given character in a given font. If that character doesn't exist,
|null| is returned instead.

This allows a character node to be used if there is an equivalent
in the |char_sub_code| list. } function new_character( f:internal_font_number; c:eight_bits):halfword ;
label exit;
var p:halfword ; {newly allocated node}
 ec:quarterword;  {effective character of |c|}
begin ec:=effective_char(false,f, c );
if font_bc[f]<= ec  then if font_ec[f]>= ec  then
  if ( font_info[char_base[  f]+  ec].qqqq .b0>min_quarterword)  then  {N.B.: not |char_info|}
    begin p:=get_avail;   mem[ p].hh.b0 :=f;   mem[ p].hh.b1 := c ;
    new_character:=p;  goto exit ;
    end;
char_warning(f,c);
new_character:=-{0xfffffff=}268435455  ;
exit:end;



{ 583. \[31] Device-independent file format }

{tangle:pos tex.web:11324:39: }

{ The most important output produced by a run of \TeX\ is the ``device
independent'' (\.[DVI]) file that specifies where characters and rules
are to appear on printed pages. The form of these files was designed by
David R. Fuchs in 1979. Almost any reasonable typesetting device can be
\xref[Fuchs, David Raymond]
\xref[DVI_files][\.[DVI] files]
driven by a program that takes \.[DVI] files as input, and dozens of such
\.[DVI]-to-whatever programs have been written. Thus, it is possible to
print the output of \TeX\ on many different kinds of equipment, using \TeX\
as a device-independent ``front end.''

A \.[DVI] file is a stream of 8-bit bytes, which may be regarded as a
series of commands in a machine-like language. The first byte of each command
is the operation code, and this code is followed by zero or more bytes
that provide parameters to the command. The parameters themselves may consist
of several consecutive bytes; for example, the `|set_rule|' command has two
parameters, each of which is four bytes long. Parameters are usually
regarded as nonnegative integers; but four-byte-long parameters,
and shorter parameters that denote distances, can be
either positive or negative. Such parameters are given in two's complement
notation. For example, a two-byte-long distance parameter has a value between
$-2^[15]$ and $2^[15]-1$. As in \.[TFM] files, numbers that occupy
more than one byte position appear in BigEndian order.

A \.[DVI] file consists of a ``preamble,'' followed by a sequence of one
or more ``pages,'' followed by a ``postamble.'' The preamble is simply a
|pre| command, with its parameters that define the dimensions used in the
file; this must come first.  Each ``page'' consists of a |bop| command,
followed by any number of other commands that tell where characters are to
be placed on a physical page, followed by an |eop| command. The pages
appear in the order that \TeX\ generated them. If we ignore |nop| commands
and \\[fnt\_def] commands (which are allowed between any two commands in
the file), each |eop| command is immediately followed by a |bop| command,
or by a |post| command; in the latter case, there are no more pages in the
file, and the remaining bytes form the postamble.  Further details about
the postamble will be explained later.

Some parameters in \.[DVI] commands are ``pointers.'' These are four-byte
quantities that give the location number of some other byte in the file;
the first byte is number~0, then comes number~1, and so on. For example,
one of the parameters of a |bop| command points to the previous |bop|;
this makes it feasible to read the pages in backwards order, in case the
results are being directed to a device that stacks its output face up.
Suppose the preamble of a \.[DVI] file occupies bytes 0 to 99. Now if the
first page occupies bytes 100 to 999, say, and if the second
page occupies bytes 1000 to 1999, then the |bop| that starts in byte 1000
points to 100 and the |bop| that starts in byte 2000 points to 1000. (The
very first |bop|, i.e., the one starting in byte 100, has a pointer of~$-1$.) }

{ 584. }

{tangle:pos tex.web:11374:1: }

{ The \.[DVI] format is intended to be both compact and easily interpreted
by a machine. Compactness is achieved by making most of the information
implicit instead of explicit. When a \.[DVI]-reading program reads the
commands for a page, it keeps track of several quantities: (a)~The current
font |f| is an integer; this value is changed only
by \\[fnt] and \\[fnt\_num] commands. (b)~The current position on the page
is given by two numbers called the horizontal and vertical coordinates,
|h| and |v|. Both coordinates are zero at the upper left corner of the page;
moving to the right corresponds to increasing the horizontal coordinate, and
moving down corresponds to increasing the vertical coordinate. Thus, the
coordinates are essentially Cartesian, except that vertical directions are
flipped; the Cartesian version of |(h,v)| would be |(h,-v)|.  (c)~The
current spacing amounts are given by four numbers |w|, |x|, |y|, and |z|,
where |w| and~|x| are used for horizontal spacing and where |y| and~|z|
are used for vertical spacing. (d)~There is a stack containing
|(h,v,w,x,y,z)| values; the \.[DVI] commands |push| and |pop| are used to
change the current level of operation. Note that the current font~|f| is
not pushed and popped; the stack contains only information about
positioning.

The values of |h|, |v|, |w|, |x|, |y|, and |z| are signed integers having up
to 32 bits, including the sign. Since they represent physical distances,
there is a small unit of measurement such that increasing |h| by~1 means
moving a certain tiny distance to the right. The actual unit of
measurement is variable, as explained below; \TeX\ sets things up so that
its \.[DVI] output is in sp units, i.e., scaled points, in agreement with
all the |scaled| dimensions in \TeX's data structures. }

{ 585. }

{tangle:pos tex.web:11402:1: }

{ Here is a list of all the commands that may appear in a \.[DVI] file. Each
command is specified by its symbolic name (e.g., |bop|), its opcode byte
(e.g., 139), and its parameters (if any). The parameters are followed
by a bracketed number telling how many bytes they occupy; for example,
`|p[4]|' means that parameter |p| is four bytes long.

\yskip\hang|set_char_0| 0. Typeset character number~0 from font~|f|
such that the reference point of the character is at |(h,v)|. Then
increase |h| by the width of that character. Note that a character may
have zero or negative width, so one cannot be sure that |h| will advance
after this command; but |h| usually does increase.

\yskip\hang\\[set\_char\_1] through \\[set\_char\_127] (opcodes 1 to 127).
Do the operations of |set_char_0|; but use the character whose number
matches the opcode, instead of character~0.

\yskip\hang|set1| 128 |c[1]|. Same as |set_char_0|, except that character
number~|c| is typeset. \TeX82 uses this command for characters in the
range |128<=c<256|.

\yskip\hang| set2| 129 |c[2]|. Same as |set1|, except that |c|~is two
bytes long, so it is in the range |0<=c<65536|. \TeX82 never uses this
command, but it should come in handy for extensions of \TeX\ that deal
with oriental languages.
\xref[oriental characters]\xref[Chinese characters]\xref[Japanese characters]

\yskip\hang| set3| 130 |c[3]|. Same as |set1|, except that |c|~is three
bytes long, so it can be as large as $2^[24]-1$. Not even the Chinese
language has this many characters, but this command might prove useful
in some yet unforeseen extension.

\yskip\hang| set4| 131 |c[4]|. Same as |set1|, except that |c|~is four
bytes long. Imagine that.

\yskip\hang|set_rule| 132 |a[4]| |b[4]|. Typeset a solid black rectangle
of height~|a| and width~|b|, with its bottom left corner at |(h,v)|. Then
set |h:=h+b|. If either |a<=0| or |b<=0|, nothing should be typeset. Note
that if |b<0|, the value of |h| will decrease even though nothing else happens.
See below for details about how to typeset rules so that consistency with
\MF\ is guaranteed.

\yskip\hang| put1| 133 |c[1]|. Typeset character number~|c| from font~|f|
such that the reference point of the character is at |(h,v)|. (The `put'
commands are exactly like the `set' commands, except that they simply put out a
character or a rule without moving the reference point afterwards.)

\yskip\hang| put2| 134 |c[2]|. Same as |set2|, except that |h| is not changed.

\yskip\hang| put3| 135 |c[3]|. Same as |set3|, except that |h| is not changed.

\yskip\hang| put4| 136 |c[4]|. Same as |set4|, except that |h| is not changed.

\yskip\hang|put_rule| 137 |a[4]| |b[4]|. Same as |set_rule|, except that
|h| is not changed.

\yskip\hang|nop| 138. No operation, do nothing. Any number of |nop|'s
may occur between \.[DVI] commands, but a |nop| cannot be inserted between
a command and its parameters or between two parameters.

\yskip\hang|bop| 139 $c_0[4]$ $c_1[4]$ $\ldots$ $c_9[4]$ $p[4]$. Beginning
of a page: Set |(h,v,w,x,y,z):=(0,0,0,0,0,0)| and set the stack empty. Set
the current font |f| to an undefined value.  The ten $c_i$ parameters hold
the values of \.[\\count0] $\ldots$ \.[\\count9] in \TeX\ at the time
\.[\\shipout] was invoked for this page; they can be used to identify
pages, if a user wants to print only part of a \.[DVI] file. The parameter
|p| points to the previous |bop| in the file; the first
|bop| has $p=-1$.

\yskip\hang|eop| 140.  End of page: Print what you have read since the
previous |bop|. At this point the stack should be empty. (The \.[DVI]-reading
programs that drive most output devices will have kept a buffer of the
material that appears on the page that has just ended. This material is
largely, but not entirely, in order by |v| coordinate and (for fixed |v|) by
|h|~coordinate; so it usually needs to be sorted into some order that is
appropriate for the device in question.)

\yskip\hang|push| 141. Push the current values of |(h,v,w,x,y,z)| onto the
top of the stack; do not change any of these values. Note that |f| is
not pushed.

\yskip\hang|pop| 142. Pop the top six values off of the stack and assign
them respectively to |(h,v,w,x,y,z)|. The number of pops should never
exceed the number of pushes, since it would be highly embarrassing if the
stack were empty at the time of a |pop| command.

\yskip\hang|right1| 143 |b[1]|. Set |h:=h+b|, i.e., move right |b| units.
The parameter is a signed number in two's complement notation, |-128<=b<128|;
if |b<0|, the reference point moves left.

\yskip\hang| right2| 144 |b[2]|. Same as |right1|, except that |b| is a
two-byte quantity in the range |-32768<=b<32768|.

\yskip\hang| right3| 145 |b[3]|. Same as |right1|, except that |b| is a
three-byte quantity in the range |$-2^[23]$<=b<$2^[23]$|.

\yskip\hang| right4| 146 |b[4]|. Same as |right1|, except that |b| is a
four-byte quantity in the range |$-2^[31]$<=b<$2^[31]$|.

\yskip\hang|w0| 147. Set |h:=h+w|; i.e., move right |w| units. With luck,
this parameterless command will usually suffice, because the same kind of motion
will occur several times in succession; the following commands explain how
|w| gets particular values.

\yskip\hang|w1| 148 |b[1]|. Set |w:=b| and |h:=h+b|. The value of |b| is a
signed quantity in two's complement notation, |-128<=b<128|. This command
changes the current |w|~spacing and moves right by |b|.

\yskip\hang| w2| 149 |b[2]|. Same as |w1|, but |b| is two bytes long,
|-32768<=b<32768|.

\yskip\hang| w3| 150 |b[3]|. Same as |w1|, but |b| is three bytes long,
|$-2^[23]$<=b<$2^[23]$|.

\yskip\hang| w4| 151 |b[4]|. Same as |w1|, but |b| is four bytes long,
|$-2^[31]$<=b<$2^[31]$|.

\yskip\hang|x0| 152. Set |h:=h+x|; i.e., move right |x| units. The `|x|'
commands are like the `|w|' commands except that they involve |x| instead
of |w|.

\yskip\hang|x1| 153 |b[1]|. Set |x:=b| and |h:=h+b|. The value of |b| is a
signed quantity in two's complement notation, |-128<=b<128|. This command
changes the current |x|~spacing and moves right by |b|.

\yskip\hang| x2| 154 |b[2]|. Same as |x1|, but |b| is two bytes long,
|-32768<=b<32768|.

\yskip\hang| x3| 155 |b[3]|. Same as |x1|, but |b| is three bytes long,
|$-2^[23]$<=b<$2^[23]$|.

\yskip\hang| x4| 156 |b[4]|. Same as |x1|, but |b| is four bytes long,
|$-2^[31]$<=b<$2^[31]$|.

\yskip\hang|down1| 157 |a[1]|. Set |v:=v+a|, i.e., move down |a| units.
The parameter is a signed number in two's complement notation, |-128<=a<128|;
if |a<0|, the reference point moves up.

\yskip\hang| down2| 158 |a[2]|. Same as |down1|, except that |a| is a
two-byte quantity in the range |-32768<=a<32768|.

\yskip\hang| down3| 159 |a[3]|. Same as |down1|, except that |a| is a
three-byte quantity in the range |$-2^[23]$<=a<$2^[23]$|.

\yskip\hang| down4| 160 |a[4]|. Same as |down1|, except that |a| is a
four-byte quantity in the range |$-2^[31]$<=a<$2^[31]$|.

\yskip\hang|y0| 161. Set |v:=v+y|; i.e., move down |y| units. With luck,
this parameterless command will usually suffice, because the same kind of motion
will occur several times in succession; the following commands explain how
|y| gets particular values.

\yskip\hang|y1| 162 |a[1]|. Set |y:=a| and |v:=v+a|. The value of |a| is a
signed quantity in two's complement notation, |-128<=a<128|. This command
changes the current |y|~spacing and moves down by |a|.

\yskip\hang| y2| 163 |a[2]|. Same as |y1|, but |a| is two bytes long,
|-32768<=a<32768|.

\yskip\hang| y3| 164 |a[3]|. Same as |y1|, but |a| is three bytes long,
|$-2^[23]$<=a<$2^[23]$|.

\yskip\hang| y4| 165 |a[4]|. Same as |y1|, but |a| is four bytes long,
|$-2^[31]$<=a<$2^[31]$|.

\yskip\hang|z0| 166. Set |v:=v+z|; i.e., move down |z| units. The `|z|' commands
are like the `|y|' commands except that they involve |z| instead of |y|.

\yskip\hang|z1| 167 |a[1]|. Set |z:=a| and |v:=v+a|. The value of |a| is a
signed quantity in two's complement notation, |-128<=a<128|. This command
changes the current |z|~spacing and moves down by |a|.

\yskip\hang| z2| 168 |a[2]|. Same as |z1|, but |a| is two bytes long,
|-32768<=a<32768|.

\yskip\hang| z3| 169 |a[3]|. Same as |z1|, but |a| is three bytes long,
|$-2^[23]$<=a<$2^[23]$|.

\yskip\hang| z4| 170 |a[4]|. Same as |z1|, but |a| is four bytes long,
|$-2^[31]$<=a<$2^[31]$|.

\yskip\hang|fnt_num_0| 171. Set |f:=0|. Font 0 must previously have been
defined by a \\[fnt\_def] instruction, as explained below.

\yskip\hang\\[fnt\_num\_1] through \\[fnt\_num\_63] (opcodes 172 to 234). Set
|f:=1|, \dots, \hbox[|f:=63|], respectively.

\yskip\hang|fnt1| 235 |k[1]|. Set |f:=k|. \TeX82 uses this command for font
numbers in the range |64<=k<256|.

\yskip\hang| fnt2| 236 |k[2]|. Same as |fnt1|, except that |k|~is two
bytes long, so it is in the range |0<=k<65536|. \TeX82 never generates this
command, but large font numbers may prove useful for specifications of
color or texture, or they may be used for special fonts that have fixed
numbers in some external coding scheme.

\yskip\hang| fnt3| 237 |k[3]|. Same as |fnt1|, except that |k|~is three
bytes long, so it can be as large as $2^[24]-1$.

\yskip\hang| fnt4| 238 |k[4]|. Same as |fnt1|, except that |k|~is four
bytes long; this is for the really big font numbers (and for the negative ones).

\yskip\hang|xxx1| 239 |k[1]| |x[k]|. This command is undefined in
general; it functions as a $(k+2)$-byte |nop| unless special \.[DVI]-reading
programs are being used. \TeX82 generates |xxx1| when a short enough
\.[\\special] appears, setting |k| to the number of bytes being sent. It
is recommended that |x| be a string having the form of a keyword followed
by possible parameters relevant to that keyword.

\yskip\hang| xxx2| 240 |k[2]| |x[k]|. Like |xxx1|, but |0<=k<65536|.

\yskip\hang| xxx3| 241 |k[3]| |x[k]|. Like |xxx1|, but |0<=k<$2^[24]$|.

\yskip\hang|xxx4| 242 |k[4]| |x[k]|. Like |xxx1|, but |k| can be ridiculously
large. \TeX82 uses |xxx4| when sending a string of length 256 or more.

\yskip\hang|fnt_def1| 243 |k[1]| |c[4]| |s[4]| |d[4]| |a[1]| |l[1]| |n[a+l]|.
Define font |k|, where |0<=k<256|; font definitions will be explained shortly.

\yskip\hang| fnt_def2| 244 |k[2]| |c[4]| |s[4]| |d[4]| |a[1]| |l[1]| |n[a+l]|.
Define font |k|, where |0<=k<65536|.

\yskip\hang| fnt_def3| 245 |k[3]| |c[4]| |s[4]| |d[4]| |a[1]| |l[1]| |n[a+l]|.
Define font |k|, where |0<=k<$2^[24]$|.

\yskip\hang| fnt_def4| 246 |k[4]| |c[4]| |s[4]| |d[4]| |a[1]| |l[1]| |n[a+l]|.
Define font |k|, where |$-2^[31]$<=k<$2^[31]$|.

\yskip\hang|pre| 247 |i[1]| |num[4]| |den[4]| |mag[4]| |k[1]| |x[k]|.
Beginning of the preamble; this must come at the very beginning of the
file. Parameters |i|, |num|, |den|, |mag|, |k|, and |x| are explained below.

\yskip\hang|post| 248. Beginning of the postamble, see below.

\yskip\hang|post_post| 249. Ending of the postamble, see below.

\yskip\noindent Commands 250--255 are undefined at the present time. }

{ 586. }

{ 587. }

{tangle:pos tex.web:11667:1: }

{ The preamble contains basic information about the file as a whole. As
stated above, there are six parameters:
$$\hbox[| i[1]| | num[4]| | den[4]| | mag[4]| | k[1]| | x[k]|.]$$
The |i| byte identifies \.[DVI] format; currently this byte is always set
to~2. (The value |i=3| is currently used for an extended format that
allows a mixture of right-to-left and left-to-right typesetting.
Some day we will set |i=4|, when \.[DVI] format makes another
incompatible change---perhaps in the year 2048.)

The next two parameters, |num| and |den|, are positive integers that define
the units of measurement; they are the numerator and denominator of a
fraction by which all dimensions in the \.[DVI] file could be multiplied
in order to get lengths in units of $10^[-7]$ meters. Since $\rm 7227[pt] =
254[cm]$, and since \TeX\ works with scaled points where there are $2^[16]$
sp in a point, \TeX\ sets
$|num|/|den|=(254\cdot10^5)/(7227\cdot2^[16])=25400000/473628672$.
\xref[sp]

The |mag| parameter is what \TeX\ calls \.[\\mag], i.e., 1000 times the
desired magnification. The actual fraction by which dimensions are
multiplied is therefore $|mag|\cdot|num|/1000|den|$. Note that if a \TeX\
source document does not call for any `\.[true]' dimensions, and if you
change it only by specifying a different \.[\\mag] setting, the \.[DVI]
file that \TeX\ creates will be completely unchanged except for the value
of |mag| in the preamble and postamble. (Fancy \.[DVI]-reading programs allow
users to override the |mag|~setting when a \.[DVI] file is being printed.)

Finally, |k| and |x| allow the \.[DVI] writer to include a comment, which is not
interpreted further. The length of comment |x| is |k|, where |0<=k<256|. }

{ 588. }

{tangle:pos tex.web:11699:1: }

{ Font definitions for a given font number |k| contain further parameters
$$\hbox[|c[4]| |s[4]| |d[4]| |a[1]| |l[1]| |n[a+l]|.]$$
The four-byte value |c| is the check sum that \TeX\ found in the \.[TFM]
file for this font; |c| should match the check sum of the font found by
programs that read this \.[DVI] file.
\xref[check sum]

Parameter |s| contains a fixed-point scale factor that is applied to
the character widths in font |k|; font dimensions in \.[TFM] files and
other font files are relative to this quantity, which is called the
``at size'' elsewhere in this documentation. The value of |s| is
always positive and less than $2^[27]$. It is given in the same units
as the other \.[DVI] dimensions, i.e., in sp when \TeX82 has made the
file.  Parameter |d| is similar to |s|; it is the ``design size,'' and
(like~|s|) it is given in \.[DVI] units. Thus, font |k| is to be used
at $|mag|\cdot s/1000d$ times its normal size.

The remaining part of a font definition gives the external name of the font,
which is an ASCII string of length |a+l|. The number |a| is the length
of the ``area'' or directory, and |l| is the length of the font name itself;
the standard local system font area is supposed to be used when |a=0|.
The |n| field contains the area in its first |a| bytes.

Font definitions must appear before the first use of a particular font number.
Once font |k| is defined, it must not be defined again; however, we
shall see below that font definitions appear in the postamble as well as
in the pages, so in this sense each font number is defined exactly twice,
if at all. Like |nop| commands, font definitions can
appear before the first |bop|, or between an |eop| and a |bop|. }

{ 589. }

{tangle:pos tex.web:11729:1: }

{ Sometimes it is desirable to make horizontal or vertical rules line up
precisely with certain features in characters of a font. It is possible to
guarantee the correct matching between \.[DVI] output and the characters
generated by \MF\ by adhering to the following principles: (1)~The \MF\
characters should be positioned so that a bottom edge or left edge that is
supposed to line up with the bottom or left edge of a rule appears at the
reference point, i.e., in row~0 and column~0 of the \MF\ raster. This
ensures that the position of the rule will not be rounded differently when
the pixel size is not a perfect multiple of the units of measurement in
the \.[DVI] file. (2)~A typeset rule of height $a>0$ and width $b>0$
should be equivalent to a \MF-generated character having black pixels in
precisely those raster positions whose \MF\ coordinates satisfy
|0<=x<$\alpha$b| and |0<=y<$\alpha$a|, where $\alpha$ is the number
of pixels per \.[DVI] unit.
\xref[METAFONT][\MF]
\xref[alignment of rules with characters]
\xref[rules aligning with characters] }

{ 590. }

{tangle:pos tex.web:11747:1: }

{ The last page in a \.[DVI] file is followed by `|post|'; this command
introduces the postamble, which summarizes important facts that \TeX\ has
accumulated about the file, making it possible to print subsets of the data
with reasonable efficiency. The postamble has the form
$$\vbox[\halign[\hbox[#\hfil]\cr
  |post| |p[4]| |num[4]| |den[4]| |mag[4]| |l[4]| |u[4]| |s[2]| |t[2]|\cr
  $\langle\,$font definitions$\,\rangle$\cr
  |post_post| |q[4]| |i[1]| 223's$[[\G]4]$\cr]]$$
Here |p| is a pointer to the final |bop| in the file. The next three
parameters, |num|, |den|, and |mag|, are duplicates of the quantities that
appeared in the preamble.

Parameters |l| and |u| give respectively the height-plus-depth of the tallest
page and the width of the widest page, in the same units as other dimensions
of the file. These numbers might be used by a \.[DVI]-reading program to
position individual ``pages'' on large sheets of film or paper; however,
the standard convention for output on normal size paper is to position each
page so that the upper left-hand corner is exactly one inch from the left
and the top. Experience has shown that it is unwise to design \.[DVI]-to-printer
software that attempts cleverly to center the output; a fixed position of
the upper left corner is easiest for users to understand and to work with.
Therefore |l| and~|u| are often ignored.

Parameter |s| is the maximum stack depth (i.e., the largest excess of
|push| commands over |pop| commands) needed to process this file. Then
comes |t|, the total number of pages (|bop| commands) present.

The postamble continues with font definitions, which are any number of
\\[fnt\_def] commands as described above, possibly interspersed with |nop|
commands. Each font number that is used in the \.[DVI] file must be defined
exactly twice: Once before it is first selected by a \\[fnt] command, and once
in the postamble. }

{ 591. }

{tangle:pos tex.web:11780:1: }

{ The last part of the postamble, following the |post_post| byte that
signifies the end of the font definitions, contains |q|, a pointer to the
|post| command that started the postamble.  An identification byte, |i|,
comes next; this currently equals~2, as in the preamble.

The |i| byte is followed by four or more bytes that are all equal to
the decimal number 223 (i.e., @'337 in octal). \TeX\ puts out four to seven of
these trailing bytes, until the total length of the file is a multiple of
four bytes, since this works out best on machines that pack four bytes per
word; but any number of 223's is allowed, as long as there are at least four
of them. In effect, 223 is a sort of signature that is added at the very end.
\xref[Fuchs, David Raymond]

This curious way to finish off a \.[DVI] file makes it feasible for
\.[DVI]-reading programs to find the postamble first, on most computers,
even though \TeX\ wants to write the postamble last. Most operating
systems permit random access to individual words or bytes of a file, so
the \.[DVI] reader can start at the end and skip backwards over the 223's
until finding the identification byte. Then it can back up four bytes, read
|q|, and move to byte |q| of the file. This byte should, of course,
contain the value 248 (|post|); now the postamble can be read, so the
\.[DVI] reader can discover all the information needed for typesetting the
pages. Note that it is also possible to skip through the \.[DVI] file at
reasonably high speed to locate a particular page, if that proves
desirable. This saves a lot of time, since \.[DVI] files used in production
jobs tend to be large.

Unfortunately, however, standard \PASCAL\ does not include the ability to
\xref[system dependencies]
access a random position in a file, or even to determine the length of a file.
Almost all systems nowadays provide the necessary capabilities, so \.[DVI]
format has been designed to work most efficiently with modern operating systems.
But if \.[DVI] files have to be processed under the restrictions of standard
\PASCAL, one can simply read them from front to back, since the necessary
header information is present in the preamble and in the font definitions.
(The |l| and |u| and |s| and |t| parameters, which appear only in the
postamble, are ``frills'' that are handy but not absolutely necessary.) }

{ 597. }

{tangle:pos tex.web:11911:1: }

{ The actual output of |dvi_buf[a..b]| to |dvi_file| is performed by calling
|write_dvi(a,b)|. For best results, this procedure should be optimized to
run as fast as possible on each particular system, since it is part of
\TeX's inner loop. It is safe to assume that |a| and |b+1| will both be
multiples of 4 when |write_dvi(a,b)| is called; therefore it is possible on
many machines to use efficient methods to pack four bytes per word and to
output an array of words with one system call.
\xref[system dependencies]
\xref[inner loop]
\xref[defecation]

In C, we use a macro to call |fwrite| or |write| directly, writing all
the bytes in one shot.  Much better even than writing four
bytes at a time. }

{ 598. }

{tangle:pos tex.web:11927:1: }

{ To put a byte in the buffer without paying the cost of invoking a procedure
each time, we use the macro |dvi_out|.

The length of |dvi_file| should not exceed |@"7FFFFFFF|; we set |cur_s:=-2|
to prevent further \.[DVI] output causing infinite recursion. } procedure dvi_swap; {outputs half of the buffer}
begin if dvi_ptr>({0x7fffffff=}2147483647-dvi_offset) then
  begin cur_s:=-2;
  fatal_error({"dvi length exceeds ""7FFFFFFF"=}839);
{ \xref[dvi length exceeds...] }
  end;
if dvi_limit=dvi_buf_size then
  begin write_dvi(0,half_buf-1); dvi_limit:=half_buf;
  dvi_offset:=dvi_offset+dvi_buf_size; dvi_ptr:=0;
  end
else  begin write_dvi(half_buf,dvi_buf_size-1); dvi_limit:=dvi_buf_size;
  end;
dvi_gone:=dvi_gone+half_buf;
end;



{ 600. }

{tangle:pos tex.web:11951:1: }

{ The |dvi_four| procedure outputs four bytes in two's complement notation,
without risking arithmetic overflow. } procedure dvi_four( x:integer);
begin if x>=0 then  begin dvi_buf[dvi_ptr]:= x  div {0100000000=} 16777216; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end 
else  begin x:=x+{010000000000=}1073741824;
  x:=x+{010000000000=}1073741824;
   begin dvi_buf[dvi_ptr]:=( x  div {0100000000=} 16777216) +  128; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
  end;
x:=x mod {0100000000=}16777216;  begin dvi_buf[dvi_ptr]:= x  div {0200000=} 65536; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
x:=x mod {0200000=}65536;  begin dvi_buf[dvi_ptr]:= x  div {0400=} 256; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
 begin dvi_buf[dvi_ptr]:= x  mod {0400=} 256; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
end;



{ 601. }

{tangle:pos tex.web:11965:1: }

{ A mild optimization of the output is performed by the |dvi_pop|
routine, which issues a |pop| unless it is possible to cancel a
`|push| |pop|' pair. The parameter to |dvi_pop| is the byte address
following the old |push| that matches the new |pop|. } procedure dvi_pop( l:integer);
begin if (l=dvi_offset+dvi_ptr)and(dvi_ptr>0) then decr(dvi_ptr)
else  begin dvi_buf[dvi_ptr]:= pop; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
end;



{ 602. }

{tangle:pos tex.web:11975:1: }

{ Here's a procedure that outputs a font definition. Since \TeX82 uses at
most 256 different fonts per job, |fnt_def1| is always used as the command code. } procedure dvi_font_def( f:internal_font_number);
var k:pool_pointer; {index into |str_pool|}
begin if f<=256+font_base then
  begin  begin dvi_buf[dvi_ptr]:= fnt_def1; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
   begin dvi_buf[dvi_ptr]:= f- font_base- 1; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
  end
else begin  begin dvi_buf[dvi_ptr]:= fnt_def1+ 1; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
   begin dvi_buf[dvi_ptr]:=( f- font_base- 1)  div {0400=} 256; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
   begin dvi_buf[dvi_ptr]:=( f- font_base- 1)  mod {0400=} 256; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
  end;
 begin dvi_buf[dvi_ptr]:=   font_check[  f].  b0 ; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
 begin dvi_buf[dvi_ptr]:=   font_check[  f].  b1 ; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
 begin dvi_buf[dvi_ptr]:=   font_check[  f].  b2 ; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
 begin dvi_buf[dvi_ptr]:=   font_check[  f].  b3 ; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;

dvi_four(font_size[f]);
dvi_four(font_dsize[f]);

 begin dvi_buf[dvi_ptr]:= (str_start[  font_area[  f]+1]-str_start[  font_area[  f]]) ; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
 begin dvi_buf[dvi_ptr]:= (str_start[  font_name[  f]+1]-str_start[  font_name[  f]]) ; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;

{ Output the font name whose internal number is |f| }
for k:=str_start[font_area[f]] to str_start[font_area[f]+1]-1 do
   begin dvi_buf[dvi_ptr]:=    str_pool[  k] ; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
for k:=str_start[font_name[f]] to str_start[font_name[f]+1]-1 do
   begin dvi_buf[dvi_ptr]:=    str_pool[  k] ; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end 

;
end;



{ 604. }

{tangle:pos tex.web:11999:1: }

{ Versions of \TeX\ intended for small computers might well choose to omit
the ideas in the next few parts of this program, since it is not really
necessary to optimize the \.[DVI] code by making use of the |w0|, |x0|,
|y0|, and |z0| commands. Furthermore, the algorithm that we are about to
describe does not pretend to give an optimum reduction in the length
of the \.[DVI] code; after all, speed is more important than compactness.
But the method is surprisingly effective, and it takes comparatively little
time.

We can best understand the basic idea by first considering a simpler problem
that has the same essential characteristics. Given a sequence of digits,
say $3\,1\,4\,1\,5\,9\,2\,6\,5\,3\,5\,8\,9$, we want to assign subscripts
$d$, $y$, or $z$ to each digit so as to maximize the number of ``$y$-hits''
and ``$z$-hits''; a $y$-hit is an instance of two appearances of the same
digit with the subscript $y$, where no $y$'s intervene between the two
appearances, and a $z$-hit is defined similarly. For example, the sequence
above could be decorated with subscripts as follows:
$$3_z\,1_y\,4_d\,1_y\,5_y\,9_d\,2_d\,6_d\,5_y\,3_z\,5_y\,8_d\,9_d.$$
There are three $y$-hits ($1_y\ldots1_y$ and $5_y\ldots5_y\ldots5_y$) and
one $z$-hit ($3_z\ldots3_z$); there are no $d$-hits, since the two appearances
of $9_d$ have $d$'s between them, but we don't count $d$-hits so it doesn't
matter how many there are. These subscripts are analogous to the \.[DVI]
commands called \\[down], $y$, and $z$, and the digits are analogous to
different amounts of vertical motion; a $y$-hit or $z$-hit corresponds to
the opportunity to use the one-byte commands |y0| or |z0| in a \.[DVI] file.

\TeX's method of assigning subscripts works like this: Append a new digit,
say $\delta$, to the right of the sequence. Now look back through the
sequence until one of the following things happens: (a)~You see
$\delta_y$ or $\delta_z$, and this was the first time you encountered a
$y$ or $z$ subscript, respectively.  Then assign $y$ or $z$ to the new
$\delta$; you have scored a hit. (b)~You see $\delta_d$, and no $y$
subscripts have been encountered so far during this search.  Then change
the previous $\delta_d$ to $\delta_y$ (this corresponds to changing a
command in the output buffer), and assign $y$ to the new $\delta$; it's
another hit.  (c)~You see $\delta_d$, and a $y$ subscript has been seen
but not a $z$.  Change the previous $\delta_d$ to $\delta_z$ and assign
$z$ to the new $\delta$. (d)~You encounter both $y$ and $z$ subscripts
before encountering a suitable $\delta$, or you scan all the way to the
front of the sequence. Assign $d$ to the new $\delta$; this assignment may
be changed later.

The subscripts $3_z\,1_y\,4_d\ldots\,$ in the example above were, in fact,
produced by this procedure, as the reader can verify. (Go ahead and try it.) }

{ 607. }

{tangle:pos tex.web:12066:1: }

{ Here is a subroutine that produces a \.[DVI] command for some specified
downward or rightward motion. It has two parameters: |w| is the amount
of motion, and |o| is either |down1| or |right1|. We use the fact that
the command codes have convenient arithmetic properties: |y1-down1=w1-right1|
and |z1-down1=x1-right1|. } procedure movement( w:scaled; o:eight_bits);
label exit,found,not_found,2,1;
var mstate:small_number; {have we seen a |y| or |z|?}
 p, q:halfword ; {current and top nodes on the stack}
 k:integer; {index into |dvi_buf|, modulo |dvi_buf_size|}
begin q:=get_node(movement_node_size); {new node for the top of the stack}
 mem[ q+width_offset].int  :=w; mem[ q+2].int :=dvi_offset+dvi_ptr;
if o=down1 then
  begin  mem[ q].hh.rh :=down_ptr; down_ptr:=q;
  end
else  begin  mem[ q].hh.rh :=right_ptr; right_ptr:=q;
  end;

{ Look at the other stack entries until deciding what sort of \.[DVI] command to generate; |goto found| if node |p| is a ``hit'' }
p:= mem[ q].hh.rh ; mstate:=none_seen;
while p<>-{0xfffffff=}268435455   do
  begin if  mem[ p+width_offset].int  =w then 
{ Consider a node with matching width; |goto found| if it's a hit }
case mstate+ mem[ p].hh.lh  of
none_seen+yz_OK,none_seen+y_OK,z_seen+yz_OK,z_seen+y_OK:{  } 

  if mem[ p+2].int <dvi_gone then goto not_found
  else 
{ Change buffered instruction to |y| or |w| and |goto found| }
begin k:=mem[ p+2].int -dvi_offset;
if k<0 then k:=k+dvi_buf_size;
dvi_buf[k]:=dvi_buf[k]+y1-down1;
 mem[ p].hh.lh :=y_here; goto found;
end

;
none_seen+z_OK,y_seen+yz_OK,y_seen+z_OK:{  } 

  if mem[ p+2].int <dvi_gone then goto not_found
  else 
{ Change buffered instruction to |z| or |x| and |goto found| }
begin k:=mem[ p+2].int -dvi_offset;
if k<0 then k:=k+dvi_buf_size;
dvi_buf[k]:=dvi_buf[k]+z1-down1;
 mem[ p].hh.lh :=z_here; goto found;
end

;
none_seen+y_here,none_seen+z_here,y_seen+z_here,z_seen+y_here: goto found;
 else   
 end 


  else  case mstate+ mem[ p].hh.lh  of
    none_seen+y_here: mstate:=y_seen;
    none_seen+z_here: mstate:=z_seen;
    y_seen+z_here,z_seen+y_here: goto not_found;
     else   
     end ;
  p:= mem[ p].hh.rh ;
  end;
not_found:

;

{ Generate a |down| or |right| command for |w| and |return| }
 mem[ q].hh.lh :=yz_OK;
if abs(w)>={040000000=}8388608 then
  begin  begin dvi_buf[dvi_ptr]:= o+ 3; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ; {|down4| or |right4|}
  dvi_four(w);  goto exit ;
  end;
if abs(w)>={0100000=}32768 then
  begin  begin dvi_buf[dvi_ptr]:= o+ 2; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ; {|down3| or |right3|}
  if w<0 then w:=w+{0100000000=}16777216;
   begin dvi_buf[dvi_ptr]:= w  div {0200000=} 65536; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ; w:=w mod {0200000=}65536; goto 2;
  end;
if abs(w)>={0200=}128 then
  begin  begin dvi_buf[dvi_ptr]:= o+ 1; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ; {|down2| or |right2|}
  if w<0 then w:=w+{0200000=}65536;
  goto 2;
  end;
 begin dvi_buf[dvi_ptr]:= o; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ; {|down1| or |right1|}
if w<0 then w:=w+{0400=}256;
goto 1;
2:  begin dvi_buf[dvi_ptr]:= w  div {0400=} 256; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
1:  begin dvi_buf[dvi_ptr]:= w  mod {0400=} 256; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;  goto exit 

;
found: 
{ Generate a |y0| or |z0| command in order to reuse a previous appearance of~|w| }
 mem[ q].hh.lh := mem[ p].hh.lh ;
if  mem[ q].hh.lh =y_here then
  begin  begin dvi_buf[dvi_ptr]:= o+ y0- down1; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ; {|y0| or |w0|}
  while  mem[ q].hh.rh <>p do
    begin q:= mem[ q].hh.rh ;
    case  mem[ q].hh.lh  of
    yz_OK:  mem[ q].hh.lh :=z_OK;
    y_OK:  mem[ q].hh.lh :=d_fixed;
     else   
     end ;
    end;
  end
else  begin  begin dvi_buf[dvi_ptr]:= o+ z0- down1; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ; {|z0| or |x0|}
  while  mem[ q].hh.rh <>p do
    begin q:= mem[ q].hh.rh ;
    case  mem[ q].hh.lh  of
    yz_OK:  mem[ q].hh.lh :=y_OK;
    z_OK:  mem[ q].hh.lh :=d_fixed;
     else   
     end ;
    end;
  end

;
exit:end;



{ 608. }

{tangle:pos tex.web:12091:1: }

{ The |info| fields in the entries of the down stack or the right stack
have six possible settings: |y_here| or |z_here| mean that the \.[DVI]
command refers to |y| or |z|, respectively (or to |w| or |x|, in the
case of horizontal motion); |yz_OK| means that the \.[DVI] command is
\\[down] (or \\[right]) but can be changed to either |y| or |z| (or
to either |w| or |x|); |y_OK| means that it is \\[down] and can be changed
to |y| but not |z|; |z_OK| is similar; and |d_fixed| means it must stay
\\[down].

The four settings |yz_OK|, |y_OK|, |z_OK|, |d_fixed| would not need to
be distinguished from each other if we were simply solving the
digit-subscripting problem mentioned above. But in \TeX's case there is
a complication because of the nested structure of |push| and |pop|
commands. Suppose we add parentheses to the digit-subscripting problem,
redefining hits so that $\delta_y\ldots \delta_y$ is a hit if all $y$'s between
the $\delta$'s are enclosed in properly nested parentheses, and if the
parenthesis level of the right-hand $\delta_y$ is deeper than or equal to
that of the left-hand one. Thus, `(' and `)' correspond to `|push|'
and `|pop|'. Now if we want to assign a subscript to the final 1 in the
sequence
$$2_y\,7_d\,1_d\,(\,8_z\,2_y\,8_z\,)\,1$$
we cannot change the previous $1_d$ to $1_y$, since that would invalidate
the $2_y\ldots2_y$ hit. But we can change it to $1_z$, scoring a hit
since the intervening $8_z$'s are enclosed in parentheses.

The program below removes movement nodes that are introduced after a |push|,
before it outputs the corresponding |pop|. }

{ 615. }

{tangle:pos tex.web:12233:1: }

{ In case you are wondering when all the movement nodes are removed from
\TeX's memory, the answer is that they are recycled just before
|hlist_out| and |vlist_out| finish outputting a box. This restores the
down and right stacks to the state they were in before the box was output,
except that some |info|'s may have become more restrictive. } procedure prune_movements( l:integer);
  {delete movement nodes with |location>=l|}
label done,exit;
var p:halfword ; {node being deleted}
begin while down_ptr<>-{0xfffffff=}268435455   do
  begin if mem[ down_ptr+2].int <l then goto done;
  p:=down_ptr; down_ptr:= mem[ p].hh.rh ; free_node(p,movement_node_size);
  end;
done: while right_ptr<>-{0xfffffff=}268435455   do
  begin if mem[ right_ptr+2].int <l then  goto exit ;
  p:=right_ptr; right_ptr:= mem[ p].hh.rh ; free_node(p,movement_node_size);
  end;
exit:end;



{ 618. }

{tangle:pos tex.web:12307:1: }

{ When |hlist_out| is called, its duty is to output the box represented
by the |hlist_node| pointed to by |temp_ptr|. The reference point of that
box has coordinates |(cur_h,cur_v)|.

Similarly, when |vlist_out| is called, its duty is to output the box represented
by the |vlist_node| pointed to by |temp_ptr|. The reference point of that
box has coordinates |(cur_h,cur_v)|.
\xref[recursion] } procedure vlist_out; forward; {|hlist_out| and |vlist_out| are mutually
  recursive}



{ 619. }

{tangle:pos tex.web:12319:1: }

{ The recursive procedures |hlist_out| and |vlist_out| each have local variables
|save_h| and |save_v| to hold the values of |dvi_h| and |dvi_v| just before
entering a new level of recursion.  In effect, the values of |save_h| and
|save_v| on \TeX's run-time stack correspond to the values of |h| and |v|
that a \.[DVI]-reading program will push onto its coordinate stack. } { \4 }
{ Declare procedures needed in |hlist_out|, |vlist_out| }
procedure special_out( p:halfword );
var old_setting:0..max_selector; {holds print |selector|}
 k:pool_pointer; {index into |str_pool|}
begin if cur_h<>dvi_h then begin movement(cur_h-dvi_h,right1); dvi_h:=cur_h; end ; if cur_v<>dvi_v then begin movement(cur_v-dvi_v,down1); dvi_v:=cur_v; end ;

old_setting:=selector; selector:=new_string;
show_token_list( mem[   mem[   p+ 1].hh.rh  ].hh.rh ,-{0xfffffff=}268435455  ,pool_size-pool_ptr);
selector:=old_setting;
 begin if pool_ptr+ 1 > pool_size then overflow({"pool size"=}257,pool_size-init_pool_ptr); { \xref[TeX capacity exceeded pool size][\quad pool size] } end ;
if  (pool_ptr - str_start[str_ptr]) <256 then
  begin  begin dvi_buf[dvi_ptr]:= xxx1; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;  begin dvi_buf[dvi_ptr]:=  (pool_ptr - str_start[str_ptr]) ; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
  end
else  begin  begin dvi_buf[dvi_ptr]:= xxx4; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ; dvi_four( (pool_ptr - str_start[str_ptr]) );
  end;
for k:=str_start[str_ptr] to pool_ptr-1 do  begin dvi_buf[dvi_ptr]:=    str_pool[  k] ; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
pool_ptr:=str_start[str_ptr]; {erase the string}
end;


procedure write_out( p:halfword );
var old_setting:0..max_selector; {holds print |selector|}
 old_mode:integer; {saved |mode|}
 j:small_number; {write stream number}
 q, r:halfword ; {temporary variables for list manipulation}
 d:integer; {number of characters in incomplete current string}
 clobbered:boolean; {system string is ok?}
 runsystem_ret:integer; {return value from |runsystem|}
begin 
{ Expand macros in the token list and make |link(def_ref)| point to the result }
q:=get_avail;  mem[ q].hh.lh :=right_brace_token+{"]"=}125;

r:=get_avail;  mem[ q].hh.rh :=r;  mem[ r].hh.lh :={07777=}4095 +end_write ; begin_token_list( q,inserted) ;

begin_token_list(  mem[  p+ 1].hh.rh  ,write_text);

q:=get_avail;  mem[ q].hh.lh :=left_brace_token+{"["=}123; begin_token_list( q,inserted) ;
{now we're ready to scan
  `\.\[$\langle\,$token list$\,\rangle$\.[\] \\endwrite]'}
old_mode:=cur_list.mode_field ; cur_list.mode_field :=0;
  {disable \.[\\prevdepth], \.[\\spacefactor], \.[\\lastskip], \.[\\prevgraf]}
cur_cs:=write_loc; q:=scan_toks(false,true); {expand macros, etc.}
get_token; if cur_tok<>{07777=}4095 +end_write  then
  
{ Recover from an unbalanced write command }
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Unbalanced write command"=} 1324); end ;
{ \xref[Unbalanced write...] }
 begin help_ptr:=2; help_line[1]:={"On this page there's a \write with fewer real ['s than ]'s."=} 1325; help_line[0]:={"I can't handle that very well; good luck."=} 1026; end ; error;
repeat get_token;
until cur_tok={07777=}4095 +end_write ;
end

;
cur_list.mode_field :=old_mode;
end_token_list {conserve stack space}

;
old_setting:=selector; j:=  mem[  p+ 1].hh.lh  ;
if j=18 then selector := new_string
else if write_open[j] then selector:=j
else  begin {write to the terminal if file isn't open}
  if (j=17)and(selector=term_and_log) then selector:=log_only;
  print_nl({""=}335);
  end;
token_show(def_ref); print_ln;
flush_list(def_ref);
if j=18 then
  begin if (eqtb[int_base+ tracing_online_code].int  <=0) then
    selector:=log_only  {Show what we're doing in the log file.}
  else selector:=term_and_log;  {Show what we're doing.}
  {If the log file isn't open yet, we can only send output to the terminal.
   Calling |open_log_file| from here seems to result in bad data in the log.}
  if not log_opened then selector:=term_only;
  print_nl({"runsystem("=}1316);
  for d:=0 to  (pool_ptr - str_start[str_ptr]) -1 do
    begin {|print| gives up if passed |str_ptr|, so do it by hand.}
    print(  str_pool[ str_start[ str_ptr]+ d] ); {N.B.: not |print_char|}
    end;
  print({")..."=}1317);
  if shellenabledp then begin
     begin if pool_ptr+ 1 > pool_size then overflow({"pool size"=}257,pool_size-init_pool_ptr); { \xref[TeX capacity exceeded pool size][\quad pool size] } end ;  begin str_pool[pool_ptr]:=   0 ; incr(pool_ptr); end ; {Append a null byte to the expansion.}
    clobbered:=false;
    for d:=0 to  (pool_ptr - str_start[str_ptr]) -1 do {Convert to external character set.}
      begin
        str_pool[str_start[str_ptr]+d]:=xchr[str_pool[str_start[str_ptr]+d]];
        if (str_pool[str_start[str_ptr]+d]=null_code)
           and (d< (pool_ptr - str_start[str_ptr]) -1) then clobbered:=true;
        {minimal checking: NUL not allowed in argument string of |system|()}
      end;
    if clobbered then print({"clobbered"=}1318)
    else begin {We have the command.  See if we're allowed to execute it,
         and report in the log.  We don't check the actual exit status of
         the command, or do anything with the output.}
      runsystem_ret := runsystem(conststringcast(addressof(
                                              str_pool[str_start[str_ptr]])));
      if runsystem_ret = -1 then print({"quotation error in system command"=}1319)
      else if runsystem_ret = 0 then print({"disabled (restricted)"=}1320)
      else if runsystem_ret = 1 then print({"executed"=}1321)
      else if runsystem_ret = 2 then print({"executed safely (allowed)"=}1322)
    end;
  end else begin
    print({"disabled"=}1323); {|shellenabledp| false}
  end;
  print_char({"."=}46); print_nl({""=}335); print_ln;
  pool_ptr:=str_start[str_ptr];  {erase the string}
end;
selector:=old_setting;
end;


procedure out_what( p:halfword );
var j:small_number; {write stream number}
     old_setting:0..max_selector;
begin case  mem[ p].hh.b1  of
open_node,write_node,close_node:
{ Do some work that has been queued up for \.[\\write] }
if not doing_leaders then
  begin j:=  mem[  p+ 1].hh.lh  ;
  if  mem[ p].hh.b1 =write_node then write_out(p)
  else  begin if write_open[j] then begin a_close(write_file[j]);
                                          write_open[j]:=false; end;
    if  mem[ p].hh.b1 =close_node then   {already closed}
    else if j<16 then
      begin cur_name:=  mem[  p+ 1].hh.rh  ; cur_area:=  mem[  p+ 2].hh.lh  ;
      cur_ext:=  mem[  p+ 2].hh.rh  ;
      if cur_ext={""=}335 then cur_ext:={".tex"=}799;
      pack_file_name(cur_name,cur_area,cur_ext) ;
      while not kpse_out_name_ok(stringcast(name_of_file+1))
            or not a_open_out(write_file[j]) do
        prompt_file_name({"output file name"=}1327,{".tex"=}799);
      write_open[j]:=true;
      {If on first line of input, log file is not ready yet, so don't log.}
      if log_opened and texmf_yesno('log_openout') then begin
        old_setting:=selector;
        if (eqtb[int_base+ tracing_online_code].int  <=0) then
          selector:=log_only  {Show what we're doing in the log file.}
        else selector:=term_and_log;  {Show what we're doing.}
        print_nl({"\openout"=}1328);
        print_int(j);
        print({" = `"=}1329);
        print_file_name(cur_name,cur_area,cur_ext);
        print({"'."=}798); print_nl({""=}335); print_ln;
        selector:=old_setting;
      end;
      end;
    end;
  end

;
special_node:special_out(p);
language_node: ;
 else  confusion({"ext4"=}1326)
{ \xref[this can't happen ext4][\quad ext4] }
 end ;
end;

{  }

procedure hlist_out; {output an |hlist_node| box}
label reswitch, move_past, fin_rule, next_p, continue, found;
var base_line: scaled; {the baseline coordinate for this box}
 left_edge: scaled; {the left coordinate for this box}
 save_h, save_v: scaled; {what |dvi_h| and |dvi_v| should pop to}
 this_box: halfword ; {pointer to containing box}
 g_order: glue_ord; {applicable order of infinity for glue}
 g_sign: normal..shrinking; {selects type of glue}
 p:halfword ; {current position in the hlist}
 save_loc:integer; {\.[DVI] byte location upon entry}
 leader_box:halfword ; {the leader box being replicated}
 leader_wd:scaled; {width of leader box being replicated}
 lx:scaled; {extra space between leader boxes}
 outer_doing_leaders:boolean; {were we doing leaders?}
 edge:scaled; {left edge of sub-box, or right edge of leader space}
 glue_temp:real; {glue value before rounding}
 cur_glue:real; {glue seen so far}
 cur_g:scaled; {rounded equivalent of |cur_glue| times the glue ratio}
begin cur_g:=0; cur_glue:=  0.0 ;
this_box:=temp_ptr; g_order:=  mem[  this_box+ list_offset].hh.b1  ;
g_sign:=  mem[  this_box+ list_offset].hh.b0  ; p:=  mem[  this_box+ list_offset].hh.rh  ;
incr(cur_s);
if cur_s>0 then  begin dvi_buf[dvi_ptr]:= push; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
if cur_s>max_push then max_push:=cur_s;
save_loc:=dvi_offset+dvi_ptr; base_line:=cur_v; left_edge:=cur_h;
while p<>-{0xfffffff=}268435455   do 
{ Output node |p| for |hlist_out| and move to the next node, maintaining the condition |cur_v=base_line| }
reswitch: if  ( p>=hi_mem_min)  then
  begin if cur_h<>dvi_h then begin movement(cur_h-dvi_h,right1); dvi_h:=cur_h; end ; if cur_v<>dvi_v then begin movement(cur_v-dvi_v,down1); dvi_v:=cur_v; end ;
  repeat f:=  mem[ p].hh.b0 ; c:=  mem[ p].hh.b1 ;
  if f<>dvi_f then 
{ Change font |dvi_f| to |f| }
begin if not font_used[f] then
  begin dvi_font_def(f); font_used[f]:=true;
  end;
if f<=64+font_base then  begin dvi_buf[dvi_ptr]:= f- font_base- 1+ fnt_num_0; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end 
else if f<=256+font_base then
  begin  begin dvi_buf[dvi_ptr]:= fnt1; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;  begin dvi_buf[dvi_ptr]:= f- font_base- 1; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
  end
else begin  begin dvi_buf[dvi_ptr]:= fnt1+ 1; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
   begin dvi_buf[dvi_ptr]:=( f- font_base- 1)  div {0400=} 256; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
   begin dvi_buf[dvi_ptr]:=( f- font_base- 1)  mod {0400=} 256; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
  end;
dvi_f:=f;
end

;
  if font_ec[f]>= c  then if font_bc[f]<= c  then
    if ( font_info[char_base[  f]+  c].qqqq .b0>min_quarterword)  then  {N.B.: not |char_info|}
      begin if c>= 128  then  begin dvi_buf[dvi_ptr]:= set1; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
       begin dvi_buf[dvi_ptr]:=   c ; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;

      cur_h:=cur_h+font_info[width_base[ f]+ font_info[char_base[  f]+  c].qqqq .b0].int  ;
      goto continue;
      end;
  if mltex_enabled_p then
    
{ Output a substitution, |goto continue| if not possible }
  begin
  
{ Get substitution information, check it, goto |found| if all is ok, otherwise goto |continue| }
  if  c >=eqtb[int_base+ char_sub_def_min_code].int   then if  c <=eqtb[int_base+ char_sub_def_max_code].int   then
    if ( eqtb[  char_sub_code_base+         c ].hh.rh   > 0 )  then
      begin  base_c:=(  eqtb[  char_sub_code_base+           c ].hh.rh     mod 256) ;
      accent_c:=(  eqtb[  char_sub_code_base+           c ].hh.rh     div 256) ;
      if (font_ec[f]>=base_c) then if (font_bc[f]<=base_c) then
        if (font_ec[f]>=accent_c) then if (font_bc[f]<=accent_c) then
          begin ia_c:= font_info[char_base[ f]+effective_char(true, f,    accent_c )].qqqq ;
          ib_c:= font_info[char_base[ f]+effective_char(true, f,    base_c )].qqqq ;
          if ( ib_c.b0>min_quarterword)  then
            if ( ia_c.b0>min_quarterword)  then goto found;
          end;
      begin_diagnostic;
      print_nl({"Missing character: Incomplete substitution "=}1331);
{ \xref[Missing character] }
       print ( c ); print({" = "=}1215);  print (accent_c);
      print({" "=}32);  print (base_c); print({" in font "=}838);
      slow_print(font_name[f]); print_char({"!"=}33); end_diagnostic(false);
      goto continue;
      end;
  begin_diagnostic;
  print_nl({"Missing character: There is no "=}837); print({"substitution for "=}1330);
{ \xref[Missing character] }
   print ( c ); print({" in font "=}838);
  slow_print(font_name[f]); print_char({"!"=}33); end_diagnostic(false);
  goto continue


;
found: 
{ Print character substitution tracing log }
 if eqtb[int_base+ tracing_lost_chars_code].int  >99 then
   begin begin_diagnostic;
   print_nl({"Using character substitution: "=}1332);
    print ( c ); print({" = "=}1215);
    print (accent_c); print({" "=}32);  print (base_c);
   print({" in font "=}838); slow_print(font_name[f]); print_char({"."=}46);
   end_diagnostic(false);
   end


;
  
{ Rebuild character using substitution information }
  base_x_height:=font_info[ x_height_code+param_base[ f]].int  ;
  base_slant:=font_info[ slant_code+param_base[ f]].int  /  65536.0 ;
{ \xref[real division] }
  accent_slant:=base_slant; {slant of accent character font}
  base_width:=font_info[width_base[ f]+ ib_c.b0].int  ;
  base_height:=font_info[height_base[ f]+(    ib_c. b1  ) div 16].int  ;
  accent_width:=font_info[width_base[ f]+ ia_c.b0].int  ;
  accent_height:=font_info[height_base[ f]+(    ia_c. b1  ) div 16].int  ;
  
{compute necessary horizontal shift (don't forget slant)}

  delta:=round((base_width-accent_width)/  2.0 +
            base_height*base_slant-base_x_height*accent_slant);
{ \xref[real multiplication] }
{ \xref[real addition] }
  dvi_h:=cur_h;  {update |dvi_h|, similar to the last statement in module 620}
  
{1. For centering/horizontal shifting insert a kern node.}

  cur_h:=cur_h+delta; if cur_h<>dvi_h then begin movement(cur_h-dvi_h,right1); dvi_h:=cur_h; end ;
  
{2. Then insert the accent character possibly shifted up or down.}

  if ((base_height<>base_x_height) and (accent_height>0)) then
    begin {the accent must be shifted up or down}
    cur_v:=base_line+(base_x_height-base_height); if cur_v<>dvi_v then begin movement(cur_v-dvi_v,down1); dvi_v:=cur_v; end ;
    if accent_c>=128 then  begin dvi_buf[dvi_ptr]:= set1; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
     begin dvi_buf[dvi_ptr]:= accent_c; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;

    cur_v:=base_line;
    end
  else begin if cur_v<>dvi_v then begin movement(cur_v-dvi_v,down1); dvi_v:=cur_v; end ;
    if accent_c>=128 then  begin dvi_buf[dvi_ptr]:= set1; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
     begin dvi_buf[dvi_ptr]:= accent_c; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;

    end;
  cur_h:=cur_h+accent_width; dvi_h:=cur_h;
  
{3. For centering/horizontal shifting insert another kern node.}

  cur_h:=cur_h+(-accent_width-delta);
  
{4. Output the base character.}

  if cur_h<>dvi_h then begin movement(cur_h-dvi_h,right1); dvi_h:=cur_h; end ; if cur_v<>dvi_v then begin movement(cur_v-dvi_v,down1); dvi_v:=cur_v; end ;
  if base_c>=128 then  begin dvi_buf[dvi_ptr]:= set1; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
   begin dvi_buf[dvi_ptr]:= base_c; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;

  cur_h:=cur_h+base_width;
  dvi_h:=cur_h {update of |dvi_h| is unnecessary, will be set in module 620}

;
  end


;
continue:
  p:= mem[ p].hh.rh ;
  until not  ( p>=hi_mem_min) ;
  dvi_h:=cur_h;
  end
else 
{ Output the non-|char_node| |p| for |hlist_out| and move to the next node }
begin case  mem[ p].hh.b0  of
hlist_node,vlist_node:
{ Output a box in an hlist }
if   mem[  p+ list_offset].hh.rh  =-{0xfffffff=}268435455   then cur_h:=cur_h+ mem[ p+width_offset].int  
else  begin save_h:=dvi_h; save_v:=dvi_v;
  cur_v:=base_line+ mem[ p+4].int  ; {shift the box down}
  temp_ptr:=p; edge:=cur_h;
  if  mem[ p].hh.b0 =vlist_node then vlist_out else hlist_out;
  dvi_h:=save_h; dvi_v:=save_v;
  cur_h:=edge+ mem[ p+width_offset].int  ; cur_v:=base_line;
  end

;
rule_node: begin rule_ht:= mem[ p+height_offset].int  ; rule_dp:= mem[ p+depth_offset].int  ; rule_wd:= mem[ p+width_offset].int  ;
  goto fin_rule;
  end;
whatsit_node: 
{ Output the whatsit node |p| in an hlist }
out_what(p)

;
glue_node: 
{ Move right or output leaders }
begin g:=  mem[  p+ 1].hh.lh  ; rule_wd:= mem[ g+width_offset].int  -cur_g;
if g_sign<>normal then
  begin if g_sign=stretching then
    begin if   mem[ g].hh.b0 =g_order then
      begin cur_glue:=cur_glue+ mem[ g+2].int  ;
       glue_temp:=     mem[   this_box+glue_offset].gr  * cur_glue; if glue_temp>  1000000000.0   then glue_temp:=  1000000000.0   else if glue_temp<-  1000000000.0   then glue_temp:=-  1000000000.0   ;
{ \xref[real multiplication] }
      cur_g:=round(glue_temp);
      end;
    end
  else if   mem[ g].hh.b1 =g_order then
      begin cur_glue:=cur_glue- mem[ g+3].int  ;
       glue_temp:=     mem[   this_box+glue_offset].gr  * cur_glue; if glue_temp>  1000000000.0   then glue_temp:=  1000000000.0   else if glue_temp<-  1000000000.0   then glue_temp:=-  1000000000.0   ;
      cur_g:=round(glue_temp);
      end;
  end;
rule_wd:=rule_wd+cur_g;
if  mem[ p].hh.b1 >=a_leaders then
  
{ Output leaders in an hlist, |goto fin_rule| if a rule or to |next_p| if done }
begin leader_box:=  mem[  p+ 1].hh.rh  ;
if  mem[ leader_box].hh.b0 =rule_node then
  begin rule_ht:= mem[ leader_box+height_offset].int  ; rule_dp:= mem[ leader_box+depth_offset].int  ;
  goto fin_rule;
  end;
leader_wd:= mem[ leader_box+width_offset].int  ;
if (leader_wd>0)and(rule_wd>0) then
  begin rule_wd:=rule_wd+10; {compensate for floating-point rounding}
  edge:=cur_h+rule_wd; lx:=0;
  
{ Let |cur_h| be the position of the first box, and set |leader_wd+lx| to the spacing between corresponding parts of boxes }
if  mem[ p].hh.b1 =a_leaders then
  begin save_h:=cur_h;
  cur_h:=left_edge+leader_wd*((cur_h-left_edge) div leader_wd);
  if cur_h<save_h then cur_h:=cur_h+leader_wd;
  end
else  begin lq:=rule_wd div leader_wd; {the number of box copies}
  lr:=rule_wd mod leader_wd; {the remaining space}
  if  mem[ p].hh.b1 =c_leaders then cur_h:=cur_h+(lr div 2)
  else  begin lx:=lr div (lq+1);
    cur_h:=cur_h+((lr-(lq-1)*lx) div 2);
    end;
  end

;
  while cur_h+leader_wd<=edge do
    
{ Output a leader box at |cur_h|, then advance |cur_h| by |leader_wd+lx| }
begin cur_v:=base_line+ mem[ leader_box+4].int  ; if cur_v<>dvi_v then begin movement(cur_v-dvi_v,down1); dvi_v:=cur_v; end ; save_v:=dvi_v;

if cur_h<>dvi_h then begin movement(cur_h-dvi_h,right1); dvi_h:=cur_h; end ; save_h:=dvi_h; temp_ptr:=leader_box;
outer_doing_leaders:=doing_leaders; doing_leaders:=true;
if  mem[ leader_box].hh.b0 =vlist_node then vlist_out else hlist_out;
doing_leaders:=outer_doing_leaders;
dvi_v:=save_v; dvi_h:=save_h; cur_v:=base_line;
cur_h:=save_h+leader_wd+lx;
end

;
  cur_h:=edge-10; goto next_p;
  end;
end

;
goto move_past;
end

;
kern_node,math_node:cur_h:=cur_h+ mem[ p+width_offset].int  ;
ligature_node: 
{ Make node |p| look like a |char_node| and |goto reswitch| }
begin mem[mem_top-12 ]:=mem[ p+1 ];  mem[ mem_top-12 ].hh.rh := mem[ p].hh.rh ;
p:=mem_top-12 ; goto reswitch;
end

;
 else   
 end ;

goto next_p;
fin_rule: 
{ Output a rule in an hlist }
if  ( rule_ht=-{010000000000=}1073741824 )  then rule_ht:= mem[ this_box+height_offset].int  ;
if  ( rule_dp=-{010000000000=}1073741824 )  then rule_dp:= mem[ this_box+depth_offset].int  ;
rule_ht:=rule_ht+rule_dp; {this is the rule thickness}
if (rule_ht>0)and(rule_wd>0) then {we don't output empty rules}
  begin if cur_h<>dvi_h then begin movement(cur_h-dvi_h,right1); dvi_h:=cur_h; end ; cur_v:=base_line+rule_dp; if cur_v<>dvi_v then begin movement(cur_v-dvi_v,down1); dvi_v:=cur_v; end ;
   begin dvi_buf[dvi_ptr]:= set_rule; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ; dvi_four(rule_ht); dvi_four(rule_wd);
  cur_v:=base_line; dvi_h:=dvi_h+rule_wd;
  end

;
move_past: cur_h:=cur_h+rule_wd;
next_p:p:= mem[ p].hh.rh ;
end



;
prune_movements(save_loc);
if cur_s>0 then dvi_pop(save_loc);
decr(cur_s);
end;



{ 629. }

{tangle:pos tex.web:12524:1: }

{ The |vlist_out| routine is similar to |hlist_out|, but a bit simpler. } procedure vlist_out; {output a |vlist_node| box}
label move_past, fin_rule, next_p;
var left_edge: scaled; {the left coordinate for this box}
 top_edge: scaled; {the top coordinate for this box}
 save_h, save_v: scaled; {what |dvi_h| and |dvi_v| should pop to}
 this_box: halfword ; {pointer to containing box}
 g_order: glue_ord; {applicable order of infinity for glue}
 g_sign: normal..shrinking; {selects type of glue}
 p:halfword ; {current position in the vlist}
 save_loc:integer; {\.[DVI] byte location upon entry}
 leader_box:halfword ; {the leader box being replicated}
 leader_ht:scaled; {height of leader box being replicated}
 lx:scaled; {extra space between leader boxes}
 outer_doing_leaders:boolean; {were we doing leaders?}
 edge:scaled; {bottom boundary of leader space}
 glue_temp:real; {glue value before rounding}
 cur_glue:real; {glue seen so far}
 cur_g:scaled; {rounded equivalent of |cur_glue| times the glue ratio}
begin cur_g:=0; cur_glue:=  0.0 ;
this_box:=temp_ptr; g_order:=  mem[  this_box+ list_offset].hh.b1  ;
g_sign:=  mem[  this_box+ list_offset].hh.b0  ; p:=  mem[  this_box+ list_offset].hh.rh  ;
incr(cur_s);
if cur_s>0 then  begin dvi_buf[dvi_ptr]:= push; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
if cur_s>max_push then max_push:=cur_s;
save_loc:=dvi_offset+dvi_ptr; left_edge:=cur_h; cur_v:=cur_v- mem[ this_box+height_offset].int  ;
top_edge:=cur_v;
while p<>-{0xfffffff=}268435455   do 
{ Output node |p| for |vlist_out| and move to the next node, maintaining the condition |cur_h=left_edge| }
begin if  ( p>=hi_mem_min)  then confusion({"vlistout"=}841)
{ \xref[this can't happen vlistout][\quad vlistout] }
else 
{ Output the non-|char_node| |p| for |vlist_out| }
begin case  mem[ p].hh.b0  of
hlist_node,vlist_node:
{ Output a box in a vlist }
if   mem[  p+ list_offset].hh.rh  =-{0xfffffff=}268435455   then cur_v:=cur_v+ mem[ p+height_offset].int  + mem[ p+depth_offset].int  
else  begin cur_v:=cur_v+ mem[ p+height_offset].int  ; if cur_v<>dvi_v then begin movement(cur_v-dvi_v,down1); dvi_v:=cur_v; end ;
  save_h:=dvi_h; save_v:=dvi_v;
  cur_h:=left_edge+ mem[ p+4].int  ; {shift the box right}
  temp_ptr:=p;
  if  mem[ p].hh.b0 =vlist_node then vlist_out else hlist_out;
  dvi_h:=save_h; dvi_v:=save_v;
  cur_v:=save_v+ mem[ p+depth_offset].int  ; cur_h:=left_edge;
  end

;
rule_node: begin rule_ht:= mem[ p+height_offset].int  ; rule_dp:= mem[ p+depth_offset].int  ; rule_wd:= mem[ p+width_offset].int  ;
  goto fin_rule;
  end;
whatsit_node: 
{ Output the whatsit node |p| in a vlist }
out_what(p)

;
glue_node: 
{ Move down or output leaders }
begin g:=  mem[  p+ 1].hh.lh  ; rule_ht:= mem[ g+width_offset].int  -cur_g;
if g_sign<>normal then
  begin if g_sign=stretching then
    begin if   mem[ g].hh.b0 =g_order then
      begin cur_glue:=cur_glue+ mem[ g+2].int  ;
       glue_temp:=     mem[   this_box+glue_offset].gr  * cur_glue; if glue_temp>  1000000000.0   then glue_temp:=  1000000000.0   else if glue_temp<-  1000000000.0   then glue_temp:=-  1000000000.0   ;
{ \xref[real multiplication] }
      cur_g:=round(glue_temp);
      end;
    end
  else if   mem[ g].hh.b1 =g_order then
      begin cur_glue:=cur_glue- mem[ g+3].int  ;
       glue_temp:=     mem[   this_box+glue_offset].gr  * cur_glue; if glue_temp>  1000000000.0   then glue_temp:=  1000000000.0   else if glue_temp<-  1000000000.0   then glue_temp:=-  1000000000.0   ;
      cur_g:=round(glue_temp);
      end;
  end;
rule_ht:=rule_ht+cur_g;
if  mem[ p].hh.b1 >=a_leaders then
  
{ Output leaders in a vlist, |goto fin_rule| if a rule or to |next_p| if done }
begin leader_box:=  mem[  p+ 1].hh.rh  ;
if  mem[ leader_box].hh.b0 =rule_node then
  begin rule_wd:= mem[ leader_box+width_offset].int  ; rule_dp:=0;
  goto fin_rule;
  end;
leader_ht:= mem[ leader_box+height_offset].int  + mem[ leader_box+depth_offset].int  ;
if (leader_ht>0)and(rule_ht>0) then
  begin rule_ht:=rule_ht+10; {compensate for floating-point rounding}
  edge:=cur_v+rule_ht; lx:=0;
  
{ Let |cur_v| be the position of the first box, and set |leader_ht+lx| to the spacing between corresponding parts of boxes }
if  mem[ p].hh.b1 =a_leaders then
  begin save_v:=cur_v;
  cur_v:=top_edge+leader_ht*((cur_v-top_edge) div leader_ht);
  if cur_v<save_v then cur_v:=cur_v+leader_ht;
  end
else  begin lq:=rule_ht div leader_ht; {the number of box copies}
  lr:=rule_ht mod leader_ht; {the remaining space}
  if  mem[ p].hh.b1 =c_leaders then cur_v:=cur_v+(lr div 2)
  else  begin lx:=lr div (lq+1);
    cur_v:=cur_v+((lr-(lq-1)*lx) div 2);
    end;
  end

;
  while cur_v+leader_ht<=edge do
    
{ Output a leader box at |cur_v|, then advance |cur_v| by |leader_ht+lx| }
begin cur_h:=left_edge+ mem[ leader_box+4].int  ; if cur_h<>dvi_h then begin movement(cur_h-dvi_h,right1); dvi_h:=cur_h; end ; save_h:=dvi_h;

cur_v:=cur_v+ mem[ leader_box+height_offset].int  ; if cur_v<>dvi_v then begin movement(cur_v-dvi_v,down1); dvi_v:=cur_v; end ; save_v:=dvi_v;
temp_ptr:=leader_box;
outer_doing_leaders:=doing_leaders; doing_leaders:=true;
if  mem[ leader_box].hh.b0 =vlist_node then vlist_out else hlist_out;
doing_leaders:=outer_doing_leaders;
dvi_v:=save_v; dvi_h:=save_h; cur_h:=left_edge;
cur_v:=save_v- mem[ leader_box+height_offset].int  +leader_ht+lx;
end

;
  cur_v:=edge-10; goto next_p;
  end;
end

;
goto move_past;
end

;
kern_node:cur_v:=cur_v+ mem[ p+width_offset].int  ;
 else   
 end ;

goto next_p;
fin_rule: 
{ Output a rule in a vlist, |goto next_p| }
if  ( rule_wd=-{010000000000=}1073741824 )  then rule_wd:= mem[ this_box+width_offset].int  ;
rule_ht:=rule_ht+rule_dp; {this is the rule thickness}
cur_v:=cur_v+rule_ht;
if (rule_ht>0)and(rule_wd>0) then {we don't output empty rules}
  begin if cur_h<>dvi_h then begin movement(cur_h-dvi_h,right1); dvi_h:=cur_h; end ; if cur_v<>dvi_v then begin movement(cur_v-dvi_v,down1); dvi_v:=cur_v; end ;
   begin dvi_buf[dvi_ptr]:= put_rule; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ; dvi_four(rule_ht); dvi_four(rule_wd);
  end;
goto next_p

;
move_past: cur_v:=cur_v+rule_ht;
end

;
next_p:p:= mem[ p].hh.rh ;
end

;
prune_movements(save_loc);
if cur_s>0 then dvi_pop(save_loc);
decr(cur_s);
end;



{ 638. }

{tangle:pos tex.web:12678:1: }

{ The |hlist_out| and |vlist_out| procedures are now complete, so we are
ready for the |ship_out| routine that gets them started in the first place. } procedure ship_out( p:halfword ); {output the box |p|}
label done;
var page_loc:integer; {location of the current |bop|}
 j, k:0..9; {indices to first ten count registers}
 s:pool_pointer; {index into |str_pool|}
 old_setting:0..max_selector; {saved |selector| setting}
begin if eqtb[int_base+ tracing_output_code].int  >0 then
  begin print_nl({""=}335); print_ln;
  print({"Completed box being shipped out"=}842);
{ \xref[Completed box...] }
  end;
if term_offset>max_print_line-9 then print_ln
else if (term_offset>0)or(file_offset>0) then print_char({" "=}32);
print_char({"["=}91); j:=9;
while (eqtb[count_base+ j].int =0)and(j>0) do decr(j);
for k:=0 to j do
  begin print_int(eqtb[count_base+ k].int );
  if k<j then print_char({"."=}46);
  end;
 fflush (stdout ) ;
if eqtb[int_base+ tracing_output_code].int  >0 then
  begin print_char({"]"=}93);
  begin_diagnostic; show_box(p); end_diagnostic(true);
  end;

{ Ship box |p| out }

{ Update the values of |max_h| and |max_v|; but if the page is too large, |goto done| }
if ( mem[ p+height_offset].int  >{07777777777=}1073741823 )or ( mem[ p+depth_offset].int  >{07777777777=}1073741823 )or 
   ( mem[ p+height_offset].int  + mem[ p+depth_offset].int  +eqtb[dimen_base+ v_offset_code].int   >{07777777777=}1073741823 )or 
   ( mem[ p+width_offset].int  +eqtb[dimen_base+ h_offset_code].int   >{07777777777=}1073741823 ) then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Huge page cannot be shipped out"=} 846); end ;
{ \xref[Huge page...] }
   begin help_ptr:=2; help_line[1]:={"The page just created is more than 18 feet tall or"=} 847; help_line[0]:={"more than 18 feet wide, so I suspect something went wrong."=} 848; end ;
  error;
  if eqtb[int_base+ tracing_output_code].int  <=0 then
    begin begin_diagnostic;
    print_nl({"The following box has been deleted:"=}849);
{ \xref[The following...deleted] }
    show_box(p);
    end_diagnostic(true);
    end;
  goto done;
  end;
if  mem[ p+height_offset].int  + mem[ p+depth_offset].int  +eqtb[dimen_base+ v_offset_code].int   >max_v then max_v:= mem[ p+height_offset].int  + mem[ p+depth_offset].int  +eqtb[dimen_base+ v_offset_code].int   ;
if  mem[ p+width_offset].int  +eqtb[dimen_base+ h_offset_code].int   >max_h then max_h:= mem[ p+width_offset].int  +eqtb[dimen_base+ h_offset_code].int   

;

{ Initialize variables as |ship_out| begins }
dvi_h:=0; dvi_v:=0; cur_h:=eqtb[dimen_base+ h_offset_code].int   ; dvi_f:=font_base ;
if output_file_name=0 then begin if job_name=0 then open_log_file; pack_job_name({".dvi"=}803); while not b_open_out(dvi_file) do prompt_file_name({"file name for output"=}804,{".dvi"=}803); output_file_name:=b_make_name_string(dvi_file); end ;
if total_pages=0 then
  begin  begin dvi_buf[dvi_ptr]:= pre; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;  begin dvi_buf[dvi_ptr]:= id_byte; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ; {output the preamble}
{ \xref[preamble of \.[DVI] file] }
  dvi_four(25400000); dvi_four(473628672); {conversion ratio for sp}
  prepare_mag; dvi_four(eqtb[int_base+ mag_code].int  ); {magnification factor is frozen}
  if output_comment then
  begin l:=strlen(output_comment);  begin dvi_buf[dvi_ptr]:= l; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
  for s:=0 to l-1 do  begin dvi_buf[dvi_ptr]:= output_comment[ s]; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
  end
else begin {the default code is unchanged}
  old_setting:=selector; selector:=new_string;
  print({" TeX output "=}840); print_int(eqtb[int_base+ year_code].int  ); print_char({"."=}46);
  print_two(eqtb[int_base+ month_code].int  ); print_char({"."=}46); print_two(eqtb[int_base+ day_code].int  );
  print_char({":"=}58); print_two(eqtb[int_base+ time_code].int   div 60);
  print_two(eqtb[int_base+ time_code].int   mod 60);
  selector:=old_setting;  begin dvi_buf[dvi_ptr]:=  (pool_ptr - str_start[str_ptr]) ; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
  for s:=str_start[str_ptr] to pool_ptr-1 do  begin dvi_buf[dvi_ptr]:=    str_pool[  s] ; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
  pool_ptr:=str_start[str_ptr]; {flush the current string}
  end;
  end

;
page_loc:=dvi_offset+dvi_ptr;
 begin dvi_buf[dvi_ptr]:= bop; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;
for k:=0 to 9 do dvi_four(eqtb[count_base+ k].int );
dvi_four(last_bop); last_bop:=page_loc;
cur_v:= mem[ p+height_offset].int  +eqtb[dimen_base+ v_offset_code].int   ; temp_ptr:=p;
if  mem[ p].hh.b0 =vlist_node then vlist_out else hlist_out;
 begin dvi_buf[dvi_ptr]:= eop; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ; incr(total_pages); cur_s:=-1;
ifdef ('IPC')
if ipc_on>0 then
  begin if dvi_limit=half_buf then
    begin write_dvi(half_buf, dvi_buf_size-1);
    flush_dvi;
    dvi_gone:=dvi_gone+half_buf;
    end;
  if dvi_ptr>({0x7fffffff=}2147483647-dvi_offset) then
    begin cur_s:=-2;
    fatal_error({"dvi length exceeds ""7FFFFFFF"=}839);
{ \xref[dvi length exceeds...] }
    end;
  if dvi_ptr>0 then
    begin write_dvi(0, dvi_ptr-1);
    flush_dvi;
    dvi_offset:=dvi_offset+dvi_ptr; dvi_gone:=dvi_gone+dvi_ptr;
    end;
  dvi_ptr:=0; dvi_limit:=dvi_buf_size;
  ipc_page(dvi_gone);
  end;
endif ('IPC');
done:

;
if eqtb[int_base+ tracing_output_code].int  <=0 then print_char({"]"=}93);
dead_cycles:=0;
 fflush (stdout ) ; {progress report}

{ Flush the box from memory, showing statistics if requested }
 ifdef('STAT')  if eqtb[int_base+ tracing_stats_code].int  >1 then
  begin print_nl({"Memory usage before: "=}843);
{ \xref[Memory usage...] }
  print_int(var_used); print_char({"&"=}38);
  print_int(dyn_used); print_char({";"=}59);
  end;
endif('STAT') 

flush_node_list(p);
 ifdef('STAT')  if eqtb[int_base+ tracing_stats_code].int  >1 then
  begin print({" after: "=}844);
  print_int(var_used); print_char({"&"=}38);
  print_int(dyn_used); print({"; still untouched: "=}845);
  print_int(hi_mem_min-lo_mem_max-1); print_ln;
  end;
endif('STAT') 

;
end;



{ 644. \[33] Packaging }

{tangle:pos tex.web:12811:18: }

{ We're essentially done with the parts of \TeX\ that are concerned with
the input (|get_next|) and the output (|ship_out|). So it's time to
get heavily into the remaining part, which does the real work of typesetting.

After lists are constructed, \TeX\ wraps them up and puts them into boxes.
Two major subroutines are given the responsibility for this task: |hpack|
applies to horizontal lists (hlists) and |vpack| applies to vertical lists
(vlists). The main duty of |hpack| and |vpack| is to compute the dimensions
of the resulting boxes, and to adjust the glue if one of those dimensions
is pre-specified. The computed sizes normally enclose all of the material
inside the new box; but some items may stick out if negative glue is used,
if the box is overfull, or if a \.[\\vbox] includes other boxes that have
been shifted left.

The subroutine call |hpack(p,w,m)| returns a pointer to an |hlist_node|
for a box containing the hlist that starts at |p|. Parameter |w| specifies
a width; and parameter |m| is either `|exactly|' or `|additional|'.  Thus,
|hpack(p,w,exactly)| produces a box whose width is exactly |w|, while
|hpack(p,w,additional)| yields a box whose width is the natural width plus
|w|.  It is convenient to define a macro called `|natural|' to cover the
most common case, so that we can say |hpack(p,natural)| to get a box that
has the natural width of list |p|.

Similarly, |vpack(p,w,m)| returns a pointer to a |vlist_node| for a
box containing the vlist that starts at |p|. In this case |w| represents
a height instead of a width; the parameter |m| is interpreted as in |hpack|. }

{ 645. }

{tangle:pos tex.web:12843:1: }

{ The parameters to |hpack| and |vpack| correspond to \TeX's primitives
like `\.[\\hbox] \.[to] \.[300pt]', `\.[\\hbox] \.[spread] \.[10pt]'; note
that `\.[\\hbox]' with no dimension following it is equivalent to
`\.[\\hbox] \.[spread] \.[0pt]'.  The |scan_spec| subroutine scans such
constructions in the user's input, including the mandatory left brace that
follows them, and it puts the specification onto |save_stack| so that the
desired box can later be obtained by executing the following code:
$$\vbox[\halign[#\hfil\cr
|save_ptr:=save_ptr-2;|\cr
|hpack(p,saved(1),saved(0)).|\cr]]$$
Special care is necessary to ensure that the special |save_stack| codes
are placed just below the new group code, because scanning can change
|save_stack| when \.[\\csname] appears. } procedure scan_spec( c:group_code; three_codes:boolean);
  {scans a box specification and left brace}
label found;
var  s:integer; {temporarily saved value}
 spec_code:exactly..additional;
begin if three_codes then s:=save_stack[save_ptr+ 0].int ;
if scan_keyword({"to"=}856) then spec_code:=exactly
{ \xref[to] }
else if scan_keyword({"spread"=}857) then spec_code:=additional
{ \xref[spread] }
else  begin spec_code:=additional; cur_val:=0;
  goto found;
  end;
scan_dimen(false,false,false) ;
found: if three_codes then
  begin save_stack[save_ptr+ 0].int :=s; incr(save_ptr);
  end;
save_stack[save_ptr+ 0].int :=spec_code; save_stack[save_ptr+ 1].int :=cur_val; save_ptr:=save_ptr+2;
new_save_level(c); scan_left_brace;
end;



{ 649. }

{tangle:pos tex.web:12910:1: }

{ Here now is |hpack|, which contains few if any surprises. } function hpack( p:halfword ; w:scaled; m:small_number):halfword ;
label reswitch, common_ending, exit;
var r:halfword ; {the box node that will be returned}
 q:halfword ; {trails behind |p|}
 h, d, x:scaled; {height, depth, and natural width}
 s:scaled; {shift amount}
 g:halfword ; {points to a glue specification}
 o:glue_ord; {order of infinity}
 f:internal_font_number; {the font in a |char_node|}
 i:four_quarters; {font information about a |char_node|}
 hd:eight_bits; {height and depth indices for a character}
begin last_badness:=0; r:=get_node(box_node_size);  mem[ r].hh.b0 :=hlist_node;
 mem[ r].hh.b1 :=min_quarterword;  mem[ r+4].int  :=0;
q:=r+list_offset;  mem[ q].hh.rh :=p;

h:=0; 
{ Clear dimensions to zero }
d:=0; x:=0;
total_stretch[normal]:=0; total_shrink[normal]:=0;
total_stretch[fil]:=0; total_shrink[fil]:=0;
total_stretch[fill]:=0; total_shrink[fill]:=0;
total_stretch[filll]:=0; total_shrink[filll]:=0

;
while p<>-{0xfffffff=}268435455   do 
{ Examine node |p| in the hlist, taking account of its effect on the dimensions of the new box, or moving it to the adjustment list; then advance |p| to the next node }
{ \xref[inner loop] }
begin reswitch: while  ( p>=hi_mem_min)  do
  
{ Incorporate character dimensions into the dimensions of the hbox that will contain~it, then move to the next node }
begin f:=  mem[ p].hh.b0 ; i:= font_info[char_base[ f]+effective_char(true, f,    mem[  p].hh.b1 )].qqqq ; hd:=  i. b1  ;
x:=x+font_info[width_base[ f]+ i.b0].int  ;

s:=font_info[height_base[ f]+( hd) div 16].int  ; if s>h then h:=s;
s:=font_info[depth_base[ f]+( hd) mod 16].int  ; if s>d then d:=s;
p:= mem[ p].hh.rh ;
end

;
if p<>-{0xfffffff=}268435455   then
  begin case  mem[ p].hh.b0  of
  hlist_node,vlist_node,rule_node,unset_node:
    
{ Incorporate box dimensions into the dimensions of the hbox that will contain~it }
begin x:=x+ mem[ p+width_offset].int  ;
if  mem[ p].hh.b0 >=rule_node then s:=0  else s:= mem[ p+4].int  ;
if  mem[ p+height_offset].int  -s>h then h:= mem[ p+height_offset].int  -s;
if  mem[ p+depth_offset].int  +s>d then d:= mem[ p+depth_offset].int  +s;
end

;
  ins_node,mark_node,adjust_node: if adjust_tail<>-{0xfffffff=}268435455   then
    
{ Transfer node |p| to the adjustment list }
begin while  mem[ q].hh.rh <>p do q:= mem[ q].hh.rh ;
if  mem[ p].hh.b0 =adjust_node then
  begin  mem[ adjust_tail].hh.rh :=mem[ p+1].int ;
  while  mem[ adjust_tail].hh.rh <>-{0xfffffff=}268435455   do adjust_tail:= mem[ adjust_tail].hh.rh ;
  p:= mem[ p].hh.rh ; free_node( mem[ q].hh.rh ,small_node_size);
  end
else  begin  mem[ adjust_tail].hh.rh :=p; adjust_tail:=p; p:= mem[ p].hh.rh ;
  end;
 mem[ q].hh.rh :=p; p:=q;
end

;
  whatsit_node:
{ Incorporate a whatsit node into an hbox } 

;
  glue_node:
{ Incorporate glue into the horizontal totals }
begin g:=  mem[  p+ 1].hh.lh  ; x:=x+ mem[ g+width_offset].int  ;

o:=  mem[ g].hh.b0 ; total_stretch[o]:=total_stretch[o]+ mem[ g+2].int  ;
o:=  mem[ g].hh.b1 ; total_shrink[o]:=total_shrink[o]+ mem[ g+3].int  ;
if  mem[ p].hh.b1 >=a_leaders then
  begin g:=  mem[  p+ 1].hh.rh  ;
  if  mem[ g+height_offset].int  >h then h:= mem[ g+height_offset].int  ;
  if  mem[ g+depth_offset].int  >d then d:= mem[ g+depth_offset].int  ;
  end;
end

;
  kern_node,math_node: x:=x+ mem[ p+width_offset].int  ;
  ligature_node: 
{ Make node |p| look like a |char_node| and |goto reswitch| }
begin mem[mem_top-12 ]:=mem[ p+1 ];  mem[ mem_top-12 ].hh.rh := mem[ p].hh.rh ;
p:=mem_top-12 ; goto reswitch;
end

;
   else   
   end ;

  p:= mem[ p].hh.rh ;
  end;
end


;
if adjust_tail<>-{0xfffffff=}268435455   then  mem[ adjust_tail].hh.rh :=-{0xfffffff=}268435455  ;
 mem[ r+height_offset].int  :=h;  mem[ r+depth_offset].int  :=d;


{ Determine the value of |width(r)| and the appropriate glue setting; then |return| or |goto common_ending| }
if m=additional then w:=x+w;
 mem[ r+width_offset].int  :=w; x:=w-x; {now |x| is the excess to be made up}
if x=0 then
  begin   mem[  r+ list_offset].hh.b0  :=normal;   mem[  r+ list_offset].hh.b1  :=normal;
     mem[  r+glue_offset].gr :=0.0 ;
   goto exit ;
  end
else if x>0 then 
{ Determine horizontal glue stretch setting, then |return| or \hbox[|goto common_ending|] }
begin 
{ Determine the stretch order }
if total_stretch[filll]<>0 then o:=filll
else if total_stretch[fill]<>0 then o:=fill
else if total_stretch[fil]<>0 then o:=fil
else o:=normal

;
  mem[  r+ list_offset].hh.b1  :=o;   mem[  r+ list_offset].hh.b0  :=stretching;
if total_stretch[o]<>0 then  mem[ r+glue_offset].gr :=  x/ total_stretch[ o] 
{ \xref[real division] }
else  begin   mem[  r+ list_offset].hh.b0  :=normal;
     mem[  r+glue_offset].gr :=0.0 ; {there's nothing to stretch}
  end;
if o=normal then if   mem[  r+ list_offset].hh.rh  <>-{0xfffffff=}268435455   then
  
{ Report an underfull hbox and |goto common_ending|, if this box is sufficiently bad }
begin last_badness:=badness(x,total_stretch[normal]);
if last_badness>eqtb[int_base+ hbadness_code].int   then
  begin print_ln;
  if last_badness>100 then print_nl({"Underfull"=}858) else print_nl({"Loose"=}859);
  print({" \hbox (badness "=}860); print_int(last_badness);
{ \xref[Underfull \\hbox...] }
{ \xref[Loose \\hbox...] }
  goto common_ending;
  end;
end

;
 goto exit ;
end


else 
{ Determine horizontal glue shrink setting, then |return| or \hbox[|goto common_ending|] }
begin 
{ Determine the shrink order }
if total_shrink[filll]<>0 then o:=filll
else if total_shrink[fill]<>0 then o:=fill
else if total_shrink[fil]<>0 then o:=fil
else o:=normal

;
  mem[  r+ list_offset].hh.b1  :=o;   mem[  r+ list_offset].hh.b0  :=shrinking;
if total_shrink[o]<>0 then  mem[ r+glue_offset].gr := (- x)/ total_shrink[ o] 
{ \xref[real division] }
else  begin   mem[  r+ list_offset].hh.b0  :=normal;
     mem[  r+glue_offset].gr :=0.0 ; {there's nothing to shrink}
  end;
if (total_shrink[o]<-x)and(o=normal)and(  mem[  r+ list_offset].hh.rh  <>-{0xfffffff=}268435455  ) then
  begin last_badness:=1000000;
     mem[  r+glue_offset].gr :=1.0 ; {use the maximum shrinkage}
  
{ Report an overfull hbox and |goto common_ending|, if this box is sufficiently bad }
if (-x-total_shrink[normal]>eqtb[dimen_base+ hfuzz_code].int   )or(eqtb[int_base+ hbadness_code].int  <100) then
  begin if (eqtb[dimen_base+ overfull_rule_code].int   >0)and(-x-total_shrink[normal]>eqtb[dimen_base+ hfuzz_code].int   ) then
    begin while  mem[ q].hh.rh <>-{0xfffffff=}268435455   do q:= mem[ q].hh.rh ;
     mem[ q].hh.rh :=new_rule;
     mem[  mem[  q].hh.rh +width_offset].int  :=eqtb[dimen_base+ overfull_rule_code].int   ;
    end;
  print_ln; print_nl({"Overfull \hbox ("=}866);
{ \xref[Overfull \\hbox...] }
  print_scaled(-x-total_shrink[normal]); print({"pt too wide"=}867);
  goto common_ending;
  end

;
  end
else if o=normal then if   mem[  r+ list_offset].hh.rh  <>-{0xfffffff=}268435455   then
  
{ Report a tight hbox and |goto common_ending|, if this box is sufficiently bad }
begin last_badness:=badness(-x,total_shrink[normal]);
if last_badness>eqtb[int_base+ hbadness_code].int   then
  begin print_ln; print_nl({"Tight \hbox (badness "=}868); print_int(last_badness);
{ \xref[Tight \\hbox...] }
  goto common_ending;
  end;
end

;
 goto exit ;
end



;
common_ending: 
{ Finish issuing a diagnostic message for an overfull or underfull hbox }
if output_active then print({") has occurred while \output is active"=}861)
else  begin if pack_begin_line<>0 then
    begin if pack_begin_line>0 then print({") in paragraph at lines "=}862)
    else print({") in alignment at lines "=}863);
    print_int(abs(pack_begin_line));
    print({"--"=}864);
    end
  else print({") detected at line "=}865);
  print_int(line);
  end;
print_ln;

font_in_short_display:=font_base ; short_display(  mem[  r+ list_offset].hh.rh  ); print_ln;

begin_diagnostic; show_box(r); end_diagnostic(true)

;
exit: hpack:=r;
end;



{ 668. }

{tangle:pos tex.web:13152:1: }

{ The |vpack| subroutine is actually a special case of a slightly more
general routine called |vpackage|, which has four parameters. The fourth
parameter, which is |max_dimen| in the case of |vpack|, specifies the
maximum depth of the page box that is constructed. The depth is first
computed by the normal rules; if it exceeds this limit, the reference
point is simply moved down until the limiting depth is attained. } function vpackage( p:halfword ; h:scaled; m:small_number; l:scaled):
  halfword ;
label common_ending, exit;
var r:halfword ; {the box node that will be returned}
 w, d, x:scaled; {width, depth, and natural height}
 s:scaled; {shift amount}
 g:halfword ; {points to a glue specification}
 o:glue_ord; {order of infinity}
begin last_badness:=0; r:=get_node(box_node_size);  mem[ r].hh.b0 :=vlist_node;
 mem[ r].hh.b1 :=min_quarterword;  mem[ r+4].int  :=0;
  mem[  r+ list_offset].hh.rh  :=p;

w:=0; 
{ Clear dimensions to zero }
d:=0; x:=0;
total_stretch[normal]:=0; total_shrink[normal]:=0;
total_stretch[fil]:=0; total_shrink[fil]:=0;
total_stretch[fill]:=0; total_shrink[fill]:=0;
total_stretch[filll]:=0; total_shrink[filll]:=0

;
while p<>-{0xfffffff=}268435455   do 
{ Examine node |p| in the vlist, taking account of its effect on the dimensions of the new box; then advance |p| to the next node }
begin if  ( p>=hi_mem_min)  then confusion({"vpack"=}869)
{ \xref[this can't happen vpack][\quad vpack] }
else  case  mem[ p].hh.b0  of
  hlist_node,vlist_node,rule_node,unset_node:
    
{ Incorporate box dimensions into the dimensions of the vbox that will contain~it }
begin x:=x+d+ mem[ p+height_offset].int  ; d:= mem[ p+depth_offset].int  ;
if  mem[ p].hh.b0 >=rule_node then s:=0  else s:= mem[ p+4].int  ;
if  mem[ p+width_offset].int  +s>w then w:= mem[ p+width_offset].int  +s;
end

;
  whatsit_node:
{ Incorporate a whatsit node into a vbox } 

;
  glue_node: 
{ Incorporate glue into the vertical totals }
begin x:=x+d; d:=0;

g:=  mem[  p+ 1].hh.lh  ; x:=x+ mem[ g+width_offset].int  ;

o:=  mem[ g].hh.b0 ; total_stretch[o]:=total_stretch[o]+ mem[ g+2].int  ;
o:=  mem[ g].hh.b1 ; total_shrink[o]:=total_shrink[o]+ mem[ g+3].int  ;
if  mem[ p].hh.b1 >=a_leaders then
  begin g:=  mem[  p+ 1].hh.rh  ;
  if  mem[ g+width_offset].int  >w then w:= mem[ g+width_offset].int  ;
  end;
end

;
  kern_node: begin x:=x+d+ mem[ p+width_offset].int  ; d:=0;
    end;
   else   
   end ;
p:= mem[ p].hh.rh ;
end

;
 mem[ r+width_offset].int  :=w;
if d>l then
  begin x:=x+d-l;  mem[ r+depth_offset].int  :=l;
  end
else  mem[ r+depth_offset].int  :=d;

{ Determine the value of |height(r)| and the appropriate glue setting; then |return| or |goto common_ending| }
if m=additional then h:=x+h;
 mem[ r+height_offset].int  :=h; x:=h-x; {now |x| is the excess to be made up}
if x=0 then
  begin   mem[  r+ list_offset].hh.b0  :=normal;   mem[  r+ list_offset].hh.b1  :=normal;
     mem[  r+glue_offset].gr :=0.0 ;
   goto exit ;
  end
else if x>0 then 
{ Determine vertical glue stretch setting, then |return| or \hbox[|goto common_ending|] }
begin 
{ Determine the stretch order }
if total_stretch[filll]<>0 then o:=filll
else if total_stretch[fill]<>0 then o:=fill
else if total_stretch[fil]<>0 then o:=fil
else o:=normal

;
  mem[  r+ list_offset].hh.b1  :=o;   mem[  r+ list_offset].hh.b0  :=stretching;
if total_stretch[o]<>0 then  mem[ r+glue_offset].gr :=  x/ total_stretch[ o] 
{ \xref[real division] }
else  begin   mem[  r+ list_offset].hh.b0  :=normal;
     mem[  r+glue_offset].gr :=0.0 ; {there's nothing to stretch}
  end;
if o=normal then if   mem[  r+ list_offset].hh.rh  <>-{0xfffffff=}268435455   then
  
{ Report an underfull vbox and |goto common_ending|, if this box is sufficiently bad }
begin last_badness:=badness(x,total_stretch[normal]);
if last_badness>eqtb[int_base+ vbadness_code].int   then
  begin print_ln;
  if last_badness>100 then print_nl({"Underfull"=}858) else print_nl({"Loose"=}859);
  print({" \vbox (badness "=}870); print_int(last_badness);
{ \xref[Underfull \\vbox...] }
{ \xref[Loose \\vbox...] }
  goto common_ending;
  end;
end

;
 goto exit ;
end


else 
{ Determine vertical glue shrink setting, then |return| or \hbox[|goto common_ending|] }
begin 
{ Determine the shrink order }
if total_shrink[filll]<>0 then o:=filll
else if total_shrink[fill]<>0 then o:=fill
else if total_shrink[fil]<>0 then o:=fil
else o:=normal

;
  mem[  r+ list_offset].hh.b1  :=o;   mem[  r+ list_offset].hh.b0  :=shrinking;
if total_shrink[o]<>0 then  mem[ r+glue_offset].gr := (- x)/ total_shrink[ o] 
{ \xref[real division] }
else  begin   mem[  r+ list_offset].hh.b0  :=normal;
     mem[  r+glue_offset].gr :=0.0 ; {there's nothing to shrink}
  end;
if (total_shrink[o]<-x)and(o=normal)and(  mem[  r+ list_offset].hh.rh  <>-{0xfffffff=}268435455  ) then
  begin last_badness:=1000000;
     mem[  r+glue_offset].gr :=1.0 ; {use the maximum shrinkage}
  
{ Report an overfull vbox and |goto common_ending|, if this box is sufficiently bad }
if (-x-total_shrink[normal]>eqtb[dimen_base+ vfuzz_code].int   )or(eqtb[int_base+ vbadness_code].int  <100) then
  begin print_ln; print_nl({"Overfull \vbox ("=}871);
{ \xref[Overfull \\vbox...] }
  print_scaled(-x-total_shrink[normal]); print({"pt too high"=}872);
  goto common_ending;
  end

;
  end
else if o=normal then if   mem[  r+ list_offset].hh.rh  <>-{0xfffffff=}268435455   then
  
{ Report a tight vbox and |goto common_ending|, if this box is sufficiently bad }
begin last_badness:=badness(-x,total_shrink[normal]);
if last_badness>eqtb[int_base+ vbadness_code].int   then
  begin print_ln; print_nl({"Tight \vbox (badness "=}873); print_int(last_badness);
{ \xref[Tight \\vbox...] }
  goto common_ending;
  end;
end

;
 goto exit ;
end



;
common_ending: 
{ Finish issuing a diagnostic message for an overfull or underfull vbox }
if output_active then print({") has occurred while \output is active"=}861)
else  begin if pack_begin_line<>0 then {it's actually negative}
    begin print({") in alignment at lines "=}863);
    print_int(abs(pack_begin_line));
    print({"--"=}864);
    end
  else print({") detected at line "=}865);
  print_int(line);
  print_ln;

  end;
begin_diagnostic; show_box(r); end_diagnostic(true)

;
exit: vpackage:=r;
end;



{ 679. }

{tangle:pos tex.web:13312:1: }

{ When a box is being appended to the current vertical list, the
baselineskip calculation is handled by the |append_to_vlist| routine. } procedure append_to_vlist( b:halfword );
var d:scaled; {deficiency of space between baselines}
 p:halfword ; {a new glue node}
begin if cur_list.aux_field .int  >-65536000  then
  begin d:= mem[  eqtb[  glue_base+   baseline_skip_code].hh.rh    +width_offset].int  -cur_list.aux_field .int  - mem[ b+height_offset].int  ;
  if d<eqtb[dimen_base+ line_skip_limit_code].int    then p:=new_param_glue(line_skip_code)
  else  begin p:=new_skip_param(baseline_skip_code);
     mem[ temp_ptr+width_offset].int  :=d; {|temp_ptr=glue_ptr(p)|}
    end;
   mem[ cur_list.tail_field ].hh.rh :=p; cur_list.tail_field :=p;
  end;
 mem[ cur_list.tail_field ].hh.rh :=b; cur_list.tail_field :=b; cur_list.aux_field .int  := mem[ b+depth_offset].int  ;
end;



{ 680. \[34] Data structures for math mode }

{tangle:pos tex.web:13329:38: }

{ When \TeX\ reads a formula that is enclosed between \.\$'s, it constructs an
[\sl mlist], which is essentially a tree structure representing that
formula.  An mlist is a linear sequence of items, but we can regard it as
a tree structure because mlists can appear within mlists. For example, many
of the entries can be subscripted or superscripted, and such ``scripts''
are mlists in their own right.

An entire formula is parsed into such a tree before any of the actual
typesetting is done, because the current style of type is usually not
known until the formula has been fully scanned. For example, when the
formula `\.[\$a+b \\over c+d\$]' is being read, there is no way to tell
that `\.[a+b]' will be in script size until `\.[\\over]' has appeared.

During the scanning process, each element of the mlist being built is
classified as a relation, a binary operator, an open parenthesis, etc.,
or as a construct like `\.[\\sqrt]' that must be built up. This classification
appears in the mlist data structure.

After a formula has been fully scanned, the mlist is converted to an hlist
so that it can be incorporated into the surrounding text. This conversion is
controlled by a recursive procedure that decides all of the appropriate
styles by a ``top-down'' process starting at the outermost level and working
in towards the subformulas. The formula is ultimately pasted together using
combinations of horizontal and vertical boxes, with glue and penalty nodes
inserted as necessary.

An mlist is represented internally as a linked list consisting chiefly
of ``noads'' (pronounced ``no-adds''), to distinguish them from the somewhat
similar ``nodes'' in hlists and vlists. Certain kinds of ordinary nodes are
allowed to appear in mlists together with the noads; \TeX\ tells the difference
by means of the |type| field, since a noad's |type| is always greater than
that of a node. An mlist does not contain character nodes, hlist nodes, vlist
nodes, math nodes, ligature nodes,
or unset nodes; in particular, each mlist item appears in the
variable-size part of |mem|, so the |type| field is always present. }

{ 681. }

{tangle:pos tex.web:13366:1: }

{ Each noad is four or more words long. The first word contains the |type|
and |subtype| and |link| fields that are already so familiar to us; the
second, third, and fourth words are called the noad's |nucleus|, |subscr|,
and |supscr| fields.

Consider, for example, the simple formula `\.[\$x\^2\$]', which would be
parsed into an mlist containing a single element called an |ord_noad|.
The |nucleus| of this noad is a representation of `\.x', the |subscr| is
empty, and the |supscr| is a representation of `\.2'.

The |nucleus|, |subscr|, and |supscr| fields are further broken into
subfields. If |p| points to a noad, and if |q| is one of its principal
fields (e.g., |q=subscr(p)|), there are several possibilities for the
subfields, depending on the |math_type| of |q|.

\yskip\hang|math_type(q)=math_char| means that |fam(q)| refers to one of
the sixteen font families, and |character(q)| is the number of a character
within a font of that family, as in a character node.

\yskip\hang|math_type(q)=math_text_char| is similar, but the character is
unsubscripted and unsuperscripted and it is followed immediately by another
character from the same font. (This |math_type| setting appears only
briefly during the processing; it is used to suppress unwanted italic
corrections.)

\yskip\hang|math_type(q)=empty| indicates a field with no value (the
corresponding attribute of noad |p| is not present).

\yskip\hang|math_type(q)=sub_box| means that |info(q)| points to a box
node (either an |hlist_node| or a |vlist_node|) that should be used as the
value of the field.  The |shift_amount| in the subsidiary box node is the
amount by which that box will be shifted downward.

\yskip\hang|math_type(q)=sub_mlist| means that |info(q)| points to
an mlist; the mlist must be converted to an hlist in order to obtain
the value of this field.

\yskip\noindent In the latter case, we might have |info(q)=null|. This
is not the same as |math_type(q)=empty|; for example, `\.[\$P\_\[\]\$]'
and `\.[\$P\$]' produce different results (the former will not have the
``italic correction'' added to the width of |P|, but the ``script skip''
will be added).

The definitions of subfields given here are evidently wasteful of space,
since a halfword is being used for the |math_type| although only three
bits would be needed. However, there are hardly ever many noads present at
once, since they are soon converted to nodes that take up even more space,
so we can afford to represent them in whatever way simplifies the
programming. }

{ 682. }

{tangle:pos tex.web:13427:1: }

{ Each portion of a formula is classified as Ord, Op, Bin, Rel, Open,
Close, Punct, or Inner, for purposes of spacing and line breaking. An
|ord_noad|, |op_noad|, |bin_noad|, |rel_noad|, |open_noad|, |close_noad|,
|punct_noad|, or |inner_noad| is used to represent portions of the various
types. For example, an `\.=' sign in a formula leads to the creation of a
|rel_noad| whose |nucleus| field is a representation of an equals sign
(usually |fam=0|, |character=@'75|).  A formula preceded by \.[\\mathrel]
also results in a |rel_noad|.  When a |rel_noad| is followed by an
|op_noad|, say, and possibly separated by one or more ordinary nodes (not
noads), \TeX\ will insert a penalty node (with the current |rel_penalty|)
just after the formula that corresponds to the |rel_noad|, unless there
already was a penalty immediately following; and a ``thick space'' will be
inserted just before the formula that corresponds to the |op_noad|.

A noad of type |ord_noad|, |op_noad|, \dots, |inner_noad| usually
has a |subtype=normal|. The only exception is that an |op_noad| might
have |subtype=limits| or |no_limits|, if the normal positioning of
limits has been overridden for this operator. }

{ 683. }

{tangle:pos tex.web:13457:1: }

{ A |radical_noad| is five words long; the fifth word is the |left_delimiter|
field, which usually represents a square root sign.

A |fraction_noad| is six words long; it has a |right_delimiter| field
as well as a |left_delimiter|.

Delimiter fields are of type |four_quarters|, and they have four subfields
called |small_fam|, |small_char|, |large_fam|, |large_char|. These subfields
represent variable-size delimiters by giving the ``small'' and ``large''
starting characters, as explained in Chapter~17 of [\sl The \TeX book].
\xref[TeXbook][\sl The \TeX book]

A |fraction_noad| is actually quite different from all other noads. Not
only does it have six words, it has |thickness|, |denominator|, and
|numerator| fields instead of |nucleus|, |subscr|, and |supscr|. The
|thickness| is a scaled value that tells how thick to make a fraction
rule; however, the special value |default_code| is used to stand for the
|default_rule_thickness| of the current size. The |numerator| and
|denominator| point to mlists that define a fraction; we always have
$$\hbox[|math_type(numerator)=math_type(denominator)=sub_mlist|].$$ The
|left_delimiter| and |right_delimiter| fields specify delimiters that will
be placed at the left and right of the fraction. In this way, a
|fraction_noad| is able to represent all of \TeX's operators \.[\\over],
\.[\\atop], \.[\\above], \.[\\overwithdelims], \.[\\atopwithdelims], and
 \.[\\abovewithdelims]. }

{ 686. }

{tangle:pos tex.web:13511:1: }

{ The |new_noad| function creates an |ord_noad| that is completely null. } function new_noad:halfword ;
var p:halfword ;
begin p:=get_node(noad_size);
 mem[ p].hh.b0 :=ord_noad;  mem[ p].hh.b1 :=normal;
mem[ p+1 ].hh:=empty_field;
mem[ p+3 ].hh:=empty_field;
mem[ p+2 ].hh:=empty_field;
new_noad:=p;
end;



{ 687. }

{tangle:pos tex.web:13523:1: }

{ A few more kinds of noads will complete the set: An |under_noad| has its
nucleus underlined; an |over_noad| has it overlined. An |accent_noad| places
an accent over its nucleus; the accent character appears as
|fam(accent_chr(p))| and |character(accent_chr(p))|. A |vcenter_noad|
centers its nucleus vertically with respect to the axis of the formula;
in such noads we always have |math_type(nucleus(p))=sub_box|.

And finally, we have |left_noad| and |right_noad| types, to implement
\TeX's \.[\\left] and \.[\\right]. The |nucleus| of such noads is
replaced by a |delimiter| field; thus, for example, `\.[\\left(]' produces
a |left_noad| such that |delimiter(p)| holds the family and character
codes for all left parentheses. A |left_noad| never appears in an mlist
except as the first element, and a |right_noad| never appears in an mlist
except as the last element; furthermore, we either have both a |left_noad|
and a |right_noad|, or neither one is present. The |subscr| and |supscr|
fields are always |empty| in a |left_noad| and a |right_noad|. }

{ 688. }

{tangle:pos tex.web:13551:1: }

{ Math formulas can also contain instructions like \.[\\textstyle] that
override \TeX's normal style rules. A |style_node| is inserted into the
data structure to record such instructions; it is three words long, so it
is considered a node instead of a noad. The |subtype| is either |display_style|
or |text_style| or |script_style| or |script_script_style|. The
second and third words of a |style_node| are not used, but they are
present because a |choice_node| is converted to a |style_node|.

\TeX\ uses even numbers 0, 2, 4, 6 to encode the basic styles
|display_style|, \dots, |script_script_style|, and adds~1 to get the
``cramped'' versions of these styles. This gives a numerical order that
is backwards from the convention of Appendix~G in [\sl The \TeX book\/];
i.e., a smaller style has a larger numerical value.
\xref[TeXbook][\sl The \TeX book] } function new_style( s:small_number):halfword ; {create a style node}
var p:halfword ; {the new node}
begin p:=get_node(style_node_size);  mem[ p].hh.b0 :=style_node;
 mem[ p].hh.b1 :=s;  mem[ p+width_offset].int  :=0;  mem[ p+depth_offset].int  :=0; {the |width| and |depth| are not used}
new_style:=p;
end;



{ 689. }

{tangle:pos tex.web:13581:1: }

{ Finally, the \.[\\mathchoice] primitive creates a |choice_node|, which
has special subfields |display_mlist|, |text_mlist|, |script_mlist|,
and |script_script_mlist| pointing to the mlists for each style. } function new_choice:halfword ; {create a choice node}
var p:halfword ; {the new node}
begin p:=get_node(style_node_size);  mem[ p].hh.b0 :=choice_node;
 mem[ p].hh.b1 :=0; {the |subtype| is not used}
 mem[  p+ 1].hh.lh  :=-{0xfffffff=}268435455  ;  mem[  p+ 1].hh.rh  :=-{0xfffffff=}268435455  ;  mem[  p+ 2].hh.lh  :=-{0xfffffff=}268435455  ;
 mem[  p+ 2].hh.rh  :=-{0xfffffff=}268435455  ;
new_choice:=p;
end;



{ 693. }

{tangle:pos tex.web:13669:1: }

{ The inelegant introduction of |show_info| in the code above seems better
than the alternative of using \PASCAL's strange |forward| declaration for a
procedure with parameters. The \PASCAL\ convention about dropping parameters
from a post-|forward| procedure is, frankly, so intolerable to the author
of \TeX\ that he would rather stoop to communication via a global temporary
variable. (A similar stoopidity occurred with respect to |hlist_out| and
|vlist_out| above, and it will occur with respect to |mlist_to_hlist| below.)
\xref[Knuth, Donald Ervin]
\xref[PASCAL][\PASCAL] } procedure show_info; {the reader will kindly forgive this}
begin show_node_list( mem[ temp_ptr].hh.lh );
end;



{ 700. }

{tangle:pos tex.web:13804:1: }

{ Before an mlist is converted to an hlist, \TeX\ makes sure that
the fonts in family~2 have enough parameters to be math-symbol
fonts, and that the fonts in family~3 have enough parameters to be
math-extension fonts. The math-symbol parameters are referred to by using the
following macros, which take a size code as their parameter; for example,
|num1(cur_size)| gives the value of the |num1| parameter for the current size.
\xref[parameters for symbols]
\xref[font parameters] }

{ 701. }

{tangle:pos tex.web:13835:1: }

{ The math-extension parameters have similar macros, but the size code is
omitted (since it is always |cur_size| when we refer to such parameters).
\xref[parameters for symbols]
\xref[font parameters] }

{ 702. }

{tangle:pos tex.web:13849:1: }

{ We also need to compute the change in style between mlists and their
subsidiaries. The following macros define the subsidiary style for
an overlined nucleus (|cramped_style|), for a subscript or a superscript
(|sub_style| or |sup_style|), or for a numerator or denominator (|num_style|
or |denom_style|). }

{ 704. }

{tangle:pos tex.web:13870:1: }

{ Here is a function that returns a pointer to a rule node having a given
thickness |t|. The rule will extend horizontally to the boundary of the vlist
that eventually contains it. } function fraction_rule( t:scaled):halfword ;
  {construct the bar for a fraction}
var p:halfword ; {the new node}
begin p:=new_rule;  mem[ p+height_offset].int  :=t;  mem[ p+depth_offset].int  :=0; fraction_rule:=p;
end;



{ 705. }

{tangle:pos tex.web:13880:1: }

{ The |overbar| function returns a pointer to a vlist box that consists of
a given box |b|, above which has been placed a kern of height |k| under a
fraction rule of thickness |t| under additional space of height |t|. } function overbar( b:halfword ; k, t:scaled):halfword ;
var p, q:halfword ; {nodes being constructed}
begin p:=new_kern(k);  mem[ p].hh.rh :=b; q:=fraction_rule(t);  mem[ q].hh.rh :=p;
p:=new_kern(t);  mem[ p].hh.rh :=q; overbar:=vpackage( p, 0,additional ,{07777777777=}1073741823 ) ;
end;



{ 706. }

{tangle:pos tex.web:13890:1: }

{ The |var_delimiter| function, which finds or constructs a sufficiently
large delimiter, is the most interesting of the auxiliary functions that
currently concern us. Given a pointer |d| to a delimiter field in some noad,
together with a size code |s| and a vertical distance |v|, this function
returns a pointer to a box that contains the smallest variant of |d| whose
height plus depth is |v| or more. (And if no variant is large enough, it
returns the largest available variant.) In particular, this routine will
construct arbitrarily large delimiters from extensible components, if
|d| leads to such characters.

The value returned is a box whose |shift_amount| has been set so that
the box is vertically centered with respect to the axis in the given size.
If a built-up symbol is returned, the height of the box before shifting
will be the height of its topmost component. }{ \4 }
{ Declare subprocedures for |var_delimiter| }
function char_box( f:internal_font_number; c:quarterword):halfword ;
var q:four_quarters;
 hd:eight_bits; {|height_depth| byte}
 b, p:halfword ; {the new box and its character node}
begin q:= font_info[char_base[ f]+effective_char(true, f,  c)].qqqq ; hd:=  q. b1  ;
b:=new_null_box;  mem[ b+width_offset].int  :=font_info[width_base[ f]+ q.b0].int  +font_info[italic_base[ f]+(  q. b2 ) div 4].int  ;
 mem[ b+height_offset].int  :=font_info[height_base[ f]+( hd) div 16].int  ;  mem[ b+depth_offset].int  :=font_info[depth_base[ f]+( hd) mod 16].int  ;
p:=get_avail;   mem[ p].hh.b1 :=c;   mem[ p].hh.b0 :=f;   mem[  b+ list_offset].hh.rh  :=p; char_box:=b;
end;


procedure stack_into_box( b:halfword ; f:internal_font_number;
   c:quarterword);
var p:halfword ; {new node placed into |b|}
begin p:=char_box(f,c);  mem[ p].hh.rh :=  mem[  b+ list_offset].hh.rh  ;   mem[  b+ list_offset].hh.rh  :=p;
 mem[ b+height_offset].int  := mem[ p+height_offset].int  ;
end;


function height_plus_depth( f:internal_font_number; c:quarterword):scaled;
var q:four_quarters;
 hd:eight_bits; {|height_depth| byte}
begin q:= font_info[char_base[ f]+effective_char(true, f,  c)].qqqq ; hd:=  q. b1  ;
height_plus_depth:=font_info[height_base[ f]+( hd) div 16].int  +font_info[depth_base[ f]+( hd) mod 16].int  ;
end;


function var_delimiter( d:halfword ; s:small_number; v:scaled):halfword ;
label found,continue;
var b:halfword ; {the box that will be constructed}
 f, g: internal_font_number; {best-so-far and tentative font codes}
 c, x, y: quarterword; {best-so-far and tentative character codes}
 m, n: integer; {the number of extensible pieces}
 u: scaled; {height-plus-depth of a tentative character}
 w: scaled; {largest height-plus-depth so far}
 q: four_quarters; {character info}
 hd: eight_bits; {height-depth byte}
 r: four_quarters; {extensible pieces}
 z: small_number; {runs through font family members}
 large_attempt: boolean; {are we trying the ``large'' variant?}
begin f:=font_base ; w:=0; large_attempt:=false;
z:=mem[ d].qqqq.b0 ; x:=mem[ d].qqqq.b1 ;
 while true do    begin 
{ Look at the variants of |(z,x)|; set |f| and |c| whenever a better character is found; |goto found| as soon as a large enough variant is encountered }
if (z<>0)or(x<>min_quarterword) then
  begin z:=z+s+16;
  repeat z:=z-16; g:= eqtb[  math_font_base+   z].hh.rh   ;
  if g<>font_base  then
    
{ Look at the list of characters starting with |x| in font |g|; set |f| and |c| whenever a better character is found; |goto found| as soon as a large enough variant is encountered }
begin y:=x;
if ( y >=font_bc[g])and( y <=font_ec[g]) then
  begin continue: q:=font_info[char_base[ g]+ y].qqqq ;
  if ( q.b0>min_quarterword)  then
    begin if ((  q. b2 ) mod 4) =ext_tag then
      begin f:=g; c:=y; goto found;
      end;
    hd:=  q. b1  ;
    u:=font_info[height_base[ g]+( hd) div 16].int  +font_info[depth_base[ g]+( hd) mod 16].int  ;
    if u>w then
      begin f:=g; c:=y; w:=u;
      if u>=v then goto found;
      end;
    if ((  q. b2 ) mod 4) =list_tag then
      begin y:= q.b3 ; goto continue;
      end;
    end;
  end;
end

;
  until z<16;
  end

;
  if large_attempt then goto found; {there were none large enough}
  large_attempt:=true; z:=mem[ d].qqqq.b2 ; x:=mem[ d].qqqq.b3 ;
  end;
found: if f<>font_base  then
  
{ Make variable |b| point to a box for |(f,c)| }
if ((  q. b2 ) mod 4) =ext_tag then
  
{ Construct an extensible character in a new box |b|, using recipe |rem_byte(q)| and font |f| }
begin b:=new_null_box;
 mem[ b].hh.b0 :=vlist_node;
r:=font_info[exten_base[f]+ q.b3 ].qqqq;


{ Compute the minimum suitable height, |w|, and the corresponding number of extension steps, |n|; also set |width(b)| }
c:= r.b3 ; u:=height_plus_depth(f,c);
w:=0; q:= font_info[char_base[ f]+effective_char(true, f,  c)].qqqq ;  mem[ b+width_offset].int  :=font_info[width_base[ f]+ q.b0].int  +font_info[italic_base[ f]+(  q. b2 ) div 4].int  ;

c:= r.b2 ; if c<>min_quarterword then w:=w+height_plus_depth(f,c);
c:= r.b1 ; if c<>min_quarterword then w:=w+height_plus_depth(f,c);
c:= r.b0 ; if c<>min_quarterword then w:=w+height_plus_depth(f,c);
n:=0;
if u>0 then while w<v do
  begin w:=w+u; incr(n);
  if  r.b1 <>min_quarterword then w:=w+u;
  end

;
c:= r.b2 ;
if c<>min_quarterword then stack_into_box(b,f,c);
c:= r.b3 ;
for m:=1 to n do stack_into_box(b,f,c);
c:= r.b1 ;
if c<>min_quarterword then
  begin stack_into_box(b,f,c); c:= r.b3 ;
  for m:=1 to n do stack_into_box(b,f,c);
  end;
c:= r.b0 ;
if c<>min_quarterword then stack_into_box(b,f,c);
 mem[ b+depth_offset].int  :=w- mem[ b+height_offset].int  ;
end


else b:=char_box(f,c)


else  begin b:=new_null_box;
   mem[ b+width_offset].int  :=eqtb[dimen_base+ null_delimiter_space_code].int   ; {use this width if no delimiter was found}
  end;
 mem[ b+4].int  :=half( mem[ b+height_offset].int  - mem[ b+depth_offset].int  ) - font_info[ 22+param_base[ eqtb[  math_font_base+   2+    s].hh.rh   ]].int  ;
var_delimiter:=b;
end;



{ 715. }

{tangle:pos tex.web:14059:1: }

{ The next subroutine is much simpler; it is used for numerators and
denominators of fractions as well as for displayed operators and
their limits above and below. It takes a given box~|b| and
changes it so that the new box is centered in a box of width~|w|.
The centering is done by putting \.[\\hss] glue at the left and right
of the list inside |b|, then packaging the new box; thus, the
actual box might not really be centered, if it already contains
infinite glue.

The given box might contain a single character whose italic correction
has been added to the width of the box; in this case a compensating
kern is inserted. } function rebox( b:halfword ; w:scaled):halfword ;
var p:halfword ; {temporary register for list manipulation}
 f:internal_font_number; {font in a one-character box}
 v:scaled; {width of a character without italic correction}
begin if ( mem[ b+width_offset].int  <>w)and(  mem[  b+ list_offset].hh.rh  <>-{0xfffffff=}268435455  ) then
  begin if  mem[ b].hh.b0 =vlist_node then b:=hpack(b,0,additional );
  p:=  mem[  b+ list_offset].hh.rh  ;
  if ( ( p>=hi_mem_min) )and( mem[ p].hh.rh =-{0xfffffff=}268435455  ) then
    begin f:=  mem[ p].hh.b0 ; v:=font_info[width_base[ f]+  font_info[char_base[  f]+effective_char(true,  f,     mem[   p].hh.b1 )].qqqq .b0].int  ;
    if v<> mem[ b+width_offset].int   then  mem[ p].hh.rh :=new_kern( mem[ b+width_offset].int  -v);
    end;
  free_node(b,box_node_size);
  b:=new_glue(mem_bot +glue_spec_size +glue_spec_size +glue_spec_size );  mem[ b].hh.rh :=p;
  while  mem[ p].hh.rh <>-{0xfffffff=}268435455   do p:= mem[ p].hh.rh ;
   mem[ p].hh.rh :=new_glue(mem_bot +glue_spec_size +glue_spec_size +glue_spec_size );
  rebox:=hpack(b,w,exactly);
  end
else  begin  mem[ b+width_offset].int  :=w; rebox:=b;
  end;
end;



{ 716. }

{tangle:pos tex.web:14093:1: }

{ Here is a subroutine that creates a new glue specification from another
one that is expressed in `\.[mu]', given the value of the math unit. } function math_glue( g:halfword ; m:scaled):halfword ;
var p:halfword ; {the new glue specification}
 n:integer; {integer part of |m|}
 f:scaled; {fraction part of |m|}
begin n:=x_over_n(m,{0200000=}65536); f:=tex_remainder ;

if f<0 then
  begin decr(n); f:=f+{0200000=}65536;
  end;
p:=get_node(glue_spec_size);
 mem[ p+width_offset].int  :=mult_and_add( n,   mem[   g+width_offset].int  , xn_over_d(   mem[   g+width_offset].int  , f,{0200000=} 65536),{07777777777=}1073741823)  ; {convert \.[mu] to \.[pt]}
  mem[ p].hh.b0 :=  mem[ g].hh.b0 ;
if   mem[ p].hh.b0 =normal then  mem[ p+2].int  :=mult_and_add( n,   mem[   g+2].int  , xn_over_d(   mem[   g+2].int  , f,{0200000=} 65536),{07777777777=}1073741823)  
else  mem[ p+2].int  := mem[ g+2].int  ;
  mem[ p].hh.b1 :=  mem[ g].hh.b1 ;
if   mem[ p].hh.b1 =normal then  mem[ p+3].int  :=mult_and_add( n,   mem[   g+3].int  , xn_over_d(   mem[   g+3].int  , f,{0200000=} 65536),{07777777777=}1073741823)  
else  mem[ p+3].int  := mem[ g+3].int  ;
math_glue:=p;
end;



{ 717. }

{tangle:pos tex.web:14117:1: }

{ The |math_kern| subroutine removes |mu_glue| from a kern node, given
the value of the math unit. } procedure math_kern( p:halfword ; m:scaled);
var  n:integer; {integer part of |m|}
 f:scaled; {fraction part of |m|}
begin if  mem[ p].hh.b1 =mu_glue then
  begin n:=x_over_n(m,{0200000=}65536); f:=tex_remainder ;

  if f<0 then
    begin decr(n); f:=f+{0200000=}65536;
    end;
   mem[ p+width_offset].int  :=mult_and_add( n,   mem[   p+width_offset].int  , xn_over_d(   mem[   p+width_offset].int  , f,{0200000=} 65536),{07777777777=}1073741823)  ;  mem[ p].hh.b1 :=explicit;
  end;
end;



{ 718. }

{tangle:pos tex.web:14132:1: }

{ Sometimes it is necessary to destroy an mlist. The following
subroutine empties the current list, assuming that |abs(mode)=mmode|. } procedure flush_math;
begin flush_node_list( mem[ cur_list.head_field ].hh.rh ); flush_node_list(cur_list.aux_field .int );
 mem[ cur_list.head_field ].hh.rh :=-{0xfffffff=}268435455  ; cur_list.tail_field :=cur_list.head_field ; cur_list.aux_field .int :=-{0xfffffff=}268435455  ;
end;



{ 720. }

{tangle:pos tex.web:14163:1: }

{ The recursion in |mlist_to_hlist| is due primarily to a subroutine
called |clean_box| that puts a given noad field into a box using a given
math style; |mlist_to_hlist| can call |clean_box|, which can call
|mlist_to_hlist|.
\xref[recursion]

The box returned by |clean_box| is ``clean'' in the
sense that its |shift_amount| is zero. } procedure mlist_to_hlist; forward;{ \2 }

function clean_box( p:halfword ; s:small_number):halfword ;
label found;
var q:halfword ; {beginning of a list to be boxed}
 save_style:small_number; {|cur_style| to be restored}
 x:halfword ; {box to be returned}
 r:halfword ; {temporary pointer}
begin case  mem[ p].hh.rh  of
math_char: begin cur_mlist:=new_noad; mem[ cur_mlist+1 ]:=mem[p];
  end;
sub_box: begin q:= mem[ p].hh.lh ; goto found;
  end;
sub_mlist: cur_mlist:= mem[ p].hh.lh ;
 else  begin q:=new_null_box; goto found;
  end
 end ;

save_style:=cur_style; cur_style:=s; mlist_penalties:=false;

mlist_to_hlist; q:= mem[ mem_top-3 ].hh.rh ; {recursive call}
cur_style:=save_style; {restore the style}

{ Set up the values of |cur_size| and |cur_mu|, based on |cur_style| }
begin if cur_style<script_style then cur_size:=text_size
else cur_size:=16*((cur_style-text_style) div 2);
cur_mu:=x_over_n(font_info[ 6+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  ,18);
end

;
found: if  ( q>=hi_mem_min) or(q=-{0xfffffff=}268435455  ) then x:=hpack(q,0,additional )
  else if ( mem[ q].hh.rh =-{0xfffffff=}268435455  )and( mem[ q].hh.b0 <=vlist_node)and( mem[ q+4].int  =0) then
    x:=q {it's already clean}
  else x:=hpack(q,0,additional );

{ Simplify a trivial box }
q:=  mem[  x+ list_offset].hh.rh  ;
if  ( q>=hi_mem_min)  then
  begin r:= mem[ q].hh.rh ;
  if r<>-{0xfffffff=}268435455   then if  mem[ r].hh.rh =-{0xfffffff=}268435455   then if not  ( r>=hi_mem_min)  then
   if  mem[ r].hh.b0 =kern_node then {unneeded italic correction}
    begin free_node(r,small_node_size);  mem[ q].hh.rh :=-{0xfffffff=}268435455  ;
    end;
  end

;
clean_box:=x;
end;



{ 722. }

{tangle:pos tex.web:14212:1: }

{ It is convenient to have a procedure that converts a |math_char|
field to an ``unpacked'' form. The |fetch| routine sets |cur_f|, |cur_c|,
and |cur_i| to the font code, character code, and character information bytes of
a given noad field. It also takes care of issuing error messages for
nonexistent characters; in such cases, |char_exists(cur_i)| will be |false|
after |fetch| has acted, and the field will also have been reset to |empty|. } procedure fetch( a:halfword ); {unpack the |math_char| field |a|}
begin cur_c:=  mem[ a].hh.b1 ; cur_f:= eqtb[  math_font_base+     mem[    a].hh.b0 +   cur_size].hh.rh   ;
if cur_f=font_base  then
  
{ Complain about an undefined family and set |cur_i| null }
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({""=} 335); end ; print_size(cur_size); print_char({" "=}32);
print_int(  mem[ a].hh.b0 ); print({" is undefined (character "=}898);
 print ( cur_c ); print_char({")"=}41);
 begin help_ptr:=4; help_line[3]:={"Somewhere in the math formula just ended, you used the"=} 899; help_line[2]:={"stated character from an undefined font family. For example,"=} 900; help_line[1]:={"plain TeX doesn't allow \it or \sl in subscripts. Proceed,"=} 901; help_line[0]:={"and I'll try to forget that I needed that character."=} 902; end ;
error; cur_i:=null_character;  mem[ a].hh.rh :=empty;
end


else  begin if ( cur_c >=font_bc[cur_f])and( cur_c <=font_ec[cur_f]) then
    cur_i:=font_info[char_base[ cur_f]+ cur_c].qqqq 
  else cur_i:=null_character;
  if not(( cur_i.b0>min_quarterword) ) then
    begin char_warning(cur_f, cur_c );
     mem[ a].hh.rh :=empty; cur_i:=null_character;
    end;
  end;
end;



{ 725. }

{tangle:pos tex.web:14252:1: }

{ We need to do a lot of different things, so |mlist_to_hlist| makes two
passes over the given mlist.

The first pass does most of the processing: It removes ``mu'' spacing from
glue, it recursively evaluates all subsidiary mlists so that only the
top-level mlist remains to be handled, it puts fractions and square roots
and such things into boxes, it attaches subscripts and superscripts, and
it computes the overall height and depth of the top-level mlist so that
the size of delimiters for a |left_noad| and a |right_noad| will be known.
The hlist resulting from each noad is recorded in that noad's |new_hlist|
field, an integer field that replaces the |nucleus| or |thickness|.
\xref[recursion]

The second pass eliminates all noads and inserts the correct glue and
penalties between nodes. }

{ 726. }

{tangle:pos tex.web:14270:1: }

{ Here is the overall plan of |mlist_to_hlist|, and the list of its
local variables. }{ \4 }
{ Declare math construction procedures }
procedure make_over( q:halfword );
begin  mem[   q+1 ].hh.lh := 
  overbar(clean_box( q+1 ,2*( cur_style div 2)+cramped ), 
  3*font_info[ 8+param_base[ eqtb[  math_font_base+   3+   cur_size].hh.rh   ]].int   ,font_info[ 8+param_base[ eqtb[  math_font_base+   3+   cur_size].hh.rh   ]].int   );
 mem[   q+1 ].hh.rh :=sub_box;
end;


procedure make_under( q:halfword );
var p, x, y: halfword ; {temporary registers for box construction}
 delta:scaled; {overall height plus depth}
begin x:=clean_box( q+1 ,cur_style);
p:=new_kern(3*font_info[ 8+param_base[ eqtb[  math_font_base+   3+   cur_size].hh.rh   ]].int   );  mem[ x].hh.rh :=p;
 mem[ p].hh.rh :=fraction_rule(font_info[ 8+param_base[ eqtb[  math_font_base+   3+   cur_size].hh.rh   ]].int   );
y:=vpackage( x, 0,additional ,{07777777777=}1073741823 ) ;
delta:= mem[ y+height_offset].int  + mem[ y+depth_offset].int  +font_info[ 8+param_base[ eqtb[  math_font_base+   3+   cur_size].hh.rh   ]].int   ;
 mem[ y+height_offset].int  := mem[ x+height_offset].int  ;  mem[ y+depth_offset].int  :=delta- mem[ y+height_offset].int  ;
 mem[   q+1 ].hh.lh :=y;  mem[   q+1 ].hh.rh :=sub_box;
end;


procedure make_vcenter( q:halfword );
var v:halfword ; {the box that should be centered vertically}
 delta:scaled; {its height plus depth}
begin v:= mem[   q+1 ].hh.lh ;
if  mem[ v].hh.b0 <>vlist_node then confusion({"vcenter"=}547);
{ \xref[this can't happen vcenter][\quad vcenter] }
delta:= mem[ v+height_offset].int  + mem[ v+depth_offset].int  ;
 mem[ v+height_offset].int  :=font_info[ 22+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  +half(delta);
 mem[ v+depth_offset].int  :=delta- mem[ v+height_offset].int  ;
end;


procedure make_radical( q:halfword );
var x, y:halfword ; {temporary registers for box construction}
 delta, clr:scaled; {dimensions involved in the calculation}
begin x:=clean_box( q+1 ,2*( cur_style div 2)+cramped );
if cur_style<text_style then {display style}
  clr:=font_info[ 8+param_base[ eqtb[  math_font_base+   3+   cur_size].hh.rh   ]].int   +(abs(font_info[ 5+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  ) div 4)
else  begin clr:=font_info[ 8+param_base[ eqtb[  math_font_base+   3+   cur_size].hh.rh   ]].int   ; clr:=clr + (abs(clr) div 4);
  end;
y:=var_delimiter( q+4 ,cur_size, mem[ x+height_offset].int  + mem[ x+depth_offset].int  +clr+
  font_info[ 8+param_base[ eqtb[  math_font_base+   3+   cur_size].hh.rh   ]].int   );
delta:= mem[ y+depth_offset].int  -( mem[ x+height_offset].int  + mem[ x+depth_offset].int  +clr);
if delta>0 then clr:=clr+half(delta); {increase the actual clearance}
 mem[ y+4].int  :=-( mem[ x+height_offset].int  +clr);
 mem[ y].hh.rh :=overbar(x,clr, mem[ y+height_offset].int  );
 mem[   q+1 ].hh.lh :=hpack(y,0,additional );  mem[   q+1 ].hh.rh :=sub_box;
end;


procedure make_math_accent( q:halfword );
label done,done1;
var p, x, y:halfword ; {temporary registers for box construction}
 a:integer; {address of lig/kern instruction}
 c:quarterword; {accent character}
 f:internal_font_number; {its font}
 i:four_quarters; {its |char_info|}
 s:scaled; {amount to skew the accent to the right}
 h:scaled; {height of character being accented}
 delta:scaled; {space to remove between accent and accentee}
 w:scaled; {width of the accentee, not including sub/superscripts}
begin fetch( q+4 );
if ( cur_i.b0>min_quarterword)  then
  begin i:=cur_i; c:=cur_c; f:=cur_f;

  
{ Compute the amount of skew }
s:=0;
if  mem[   q+1 ].hh.rh =math_char then
  begin fetch( q+1 );
  if ((  cur_i. b2 ) mod 4) =lig_tag then
    begin a:=lig_kern_base[ cur_f]+ cur_i.b3 ;
    cur_i:=font_info[a].qqqq;
    if  cur_i.b0 > 128   then
      begin a:=lig_kern_base[ cur_f]+256*  cur_i.b2 +  cur_i.b3 +32768-256*(128+min_quarterword)  ;
      cur_i:=font_info[a].qqqq;
      end;
     while true do   begin if    cur_i.b1  =skew_char[cur_f] then
        begin if  cur_i.b2 >= 128   then
          if  cur_i.b0 <= 128   then s:=font_info[kern_base[ cur_f]+256*  cur_i.b2 +  cur_i.b3 ].int  ;
        goto done1;
        end;
      if  cur_i.b0 >= 128   then goto done1;
      a:=a+   cur_i.b0  +1;
      cur_i:=font_info[a].qqqq;
      end;
    end;
  end;
done1:

;
  x:=clean_box( q+1 ,2*( cur_style div 2)+cramped ); w:= mem[ x+width_offset].int  ; h:= mem[ x+height_offset].int  ;
  
{ Switch to a larger accent if available and appropriate }
 while true do    begin if ((  i. b2 ) mod 4) <>list_tag then goto done;
  y:= i.b3 ;
  i:=font_info[char_base[ f]+ y].qqqq ;
  if not ( i.b0>min_quarterword)  then goto done;
  if font_info[width_base[ f]+ i.b0].int  >w then goto done;
  c:=y;
  end;
done:

;
  if h<font_info[ x_height_code+param_base[ f]].int   then delta:=h else delta:=font_info[ x_height_code+param_base[ f]].int  ;
  if ( mem[   q+2 ].hh.rh <>empty)or( mem[   q+3 ].hh.rh <>empty) then
    if  mem[   q+1 ].hh.rh =math_char then
      
{ Swap the subscript and superscript into box |x| }
begin flush_node_list(x); x:=new_noad;
mem[ x+1 ]:=mem[ q+1 ];
mem[ x+2 ]:=mem[ q+2 ];
mem[ x+3 ]:=mem[ q+3 ];

mem[ q+2 ].hh:=empty_field;
mem[ q+3 ].hh:=empty_field;

 mem[   q+1 ].hh.rh :=sub_mlist;  mem[   q+1 ].hh.lh :=x;
x:=clean_box( q+1 ,cur_style); delta:=delta+ mem[ x+height_offset].int  -h; h:= mem[ x+height_offset].int  ;
end

;
  y:=char_box(f,c);
   mem[ y+4].int  :=s+half(w- mem[ y+width_offset].int  );
   mem[ y+width_offset].int  :=0; p:=new_kern(-delta);  mem[ p].hh.rh :=x;  mem[ y].hh.rh :=p;
  y:=vpackage( y, 0,additional ,{07777777777=}1073741823 ) ;  mem[ y+width_offset].int  := mem[ x+width_offset].int  ;
  if  mem[ y+height_offset].int  <h then 
{ Make the height of box |y| equal to |h| }
begin p:=new_kern(h- mem[ y+height_offset].int  );  mem[ p].hh.rh :=  mem[  y+ list_offset].hh.rh  ;   mem[  y+ list_offset].hh.rh  :=p;
 mem[ y+height_offset].int  :=h;
end

;
   mem[   q+1 ].hh.lh :=y;
   mem[   q+1 ].hh.rh :=sub_box;
  end;
end;


procedure make_fraction( q:halfword );
var p, v, x, y, z:halfword ; {temporary registers for box construction}
 delta, delta1, delta2, shift_up, shift_down, clr:scaled;
  {dimensions for box calculations}
begin if  mem[ q+width_offset].int  ={010000000000=}1073741824  then  mem[ q+width_offset].int  :=font_info[ 8+param_base[ eqtb[  math_font_base+   3+   cur_size].hh.rh   ]].int   ;

{ Create equal-width boxes |x| and |z| for the numerator and denominator, and compute the default amounts |shift_up| and |shift_down| by which they are displaced from the baseline }
x:=clean_box( q+2 , cur_style+2-2*( cur_style div 6) );
z:=clean_box( q+3 ,2*( cur_style div 2)+cramped+2-2*( cur_style div 6) );
if  mem[ x+width_offset].int  < mem[ z+width_offset].int   then x:=rebox(x, mem[ z+width_offset].int  )
else z:=rebox(z, mem[ x+width_offset].int  );
if cur_style<text_style then {display style}
  begin shift_up:=font_info[ 8+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  ; shift_down:=font_info[ 11+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  ;
  end
else  begin shift_down:=font_info[ 12+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  ;
  if  mem[ q+width_offset].int  <>0 then shift_up:=font_info[ 9+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  
  else shift_up:=font_info[ 10+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  ;
  end

;
if  mem[ q+width_offset].int  =0 then 
{ Adjust \(s)|shift_up| and |shift_down| for the case of no fraction line }
begin if cur_style<text_style then clr:=7*font_info[ 8+param_base[ eqtb[  math_font_base+   3+   cur_size].hh.rh   ]].int   
else clr:=3*font_info[ 8+param_base[ eqtb[  math_font_base+   3+   cur_size].hh.rh   ]].int   ;
delta:=half(clr-((shift_up- mem[ x+depth_offset].int  )-( mem[ z+height_offset].int  -shift_down)));
if delta>0 then
  begin shift_up:=shift_up+delta;
  shift_down:=shift_down+delta;
  end;
end


else 
{ Adjust \(s)|shift_up| and |shift_down| for the case of a fraction line }
begin if cur_style<text_style then clr:=3* mem[ q+width_offset].int  
else clr:= mem[ q+width_offset].int  ;
delta:=half( mem[ q+width_offset].int  );
delta1:=clr-((shift_up- mem[ x+depth_offset].int  )-(font_info[ 22+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  +delta));
delta2:=clr-((font_info[ 22+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  -delta)-( mem[ z+height_offset].int  -shift_down));
if delta1>0 then shift_up:=shift_up+delta1;
if delta2>0 then shift_down:=shift_down+delta2;
end

;

{ Construct a vlist box for the fraction, according to |shift_up| and |shift_down| }
v:=new_null_box;  mem[ v].hh.b0 :=vlist_node;
 mem[ v+height_offset].int  :=shift_up+ mem[ x+height_offset].int  ;  mem[ v+depth_offset].int  := mem[ z+depth_offset].int  +shift_down;
 mem[ v+width_offset].int  := mem[ x+width_offset].int  ; {this also equals |width(z)|}
if  mem[ q+width_offset].int  =0 then
  begin p:=new_kern((shift_up- mem[ x+depth_offset].int  )-( mem[ z+height_offset].int  -shift_down));
   mem[ p].hh.rh :=z;
  end
else  begin y:=fraction_rule( mem[ q+width_offset].int  );

  p:=new_kern((font_info[ 22+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  -delta)- ( mem[ z+height_offset].int  -shift_down));

   mem[ y].hh.rh :=p;  mem[ p].hh.rh :=z;

  p:=new_kern((shift_up- mem[ x+depth_offset].int  )-(font_info[ 22+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  +delta));
   mem[ p].hh.rh :=y;
  end;
 mem[ x].hh.rh :=p;   mem[  v+ list_offset].hh.rh  :=x

;

{ Put the \(f)fraction into a box with its delimiters, and make |new_hlist(q)| point to it }
if cur_style<text_style then delta:=font_info[ 20+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  
else delta:=font_info[ 21+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  ;
x:=var_delimiter( q+4 , cur_size, delta);  mem[ x].hh.rh :=v;

z:=var_delimiter( q+5 , cur_size, delta);  mem[ v].hh.rh :=z;

mem[  q+1 ].int :=hpack(x,0,additional )

;
end;


function make_op( q:halfword ):scaled;
var delta:scaled; {offset between subscript and superscript}
 p, v, x, y, z:halfword ; {temporary registers for box construction}
 c:quarterword;  i:four_quarters; {registers for character examination}
 shift_up, shift_down:scaled; {dimensions for box calculation}
begin if ( mem[ q].hh.b1 =normal)and(cur_style<text_style) then
   mem[ q].hh.b1 :=limits;
if  mem[   q+1 ].hh.rh =math_char then
  begin fetch( q+1 );
  if (cur_style<text_style)and(((  cur_i. b2 ) mod 4) =list_tag) then {make it larger}
    begin c:= cur_i.b3 ; i:=font_info[char_base[ cur_f]+ c].qqqq ;
    if ( i.b0>min_quarterword)  then
      begin cur_c:=c; cur_i:=i;   mem[   q+1 ].hh.b1 :=c;
      end;
    end;
  delta:=font_info[italic_base[ cur_f]+(  cur_i. b2 ) div 4].int  ; x:=clean_box( q+1 ,cur_style);
  if ( mem[   q+3 ].hh.rh <>empty)and( mem[ q].hh.b1 <>limits) then
     mem[ x+width_offset].int  := mem[ x+width_offset].int  -delta; {remove italic correction}
   mem[ x+4].int  :=half( mem[ x+height_offset].int  - mem[ x+depth_offset].int  ) - font_info[ 22+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  ;
    {center vertically}
   mem[   q+1 ].hh.rh :=sub_box;  mem[   q+1 ].hh.lh :=x;
  end
else delta:=0;
if  mem[ q].hh.b1 =limits then
  
{ Construct a box with limits above and below it, skewed by |delta| }
begin x:=clean_box( q+2 ,2*( cur_style div 4)+script_style+( cur_style mod 2) );
y:=clean_box( q+1 ,cur_style);
z:=clean_box( q+3 ,2*( cur_style div 4)+script_style+cramped );
v:=new_null_box;  mem[ v].hh.b0 :=vlist_node;  mem[ v+width_offset].int  := mem[ y+width_offset].int  ;
if  mem[ x+width_offset].int  > mem[ v+width_offset].int   then  mem[ v+width_offset].int  := mem[ x+width_offset].int  ;
if  mem[ z+width_offset].int  > mem[ v+width_offset].int   then  mem[ v+width_offset].int  := mem[ z+width_offset].int  ;
x:=rebox(x, mem[ v+width_offset].int  ); y:=rebox(y, mem[ v+width_offset].int  ); z:=rebox(z, mem[ v+width_offset].int  );

 mem[ x+4].int  :=half(delta);  mem[ z+4].int  :=- mem[ x+4].int  ;
 mem[ v+height_offset].int  := mem[ y+height_offset].int  ;  mem[ v+depth_offset].int  := mem[ y+depth_offset].int  ;

{ Attach the limits to |y| and adjust |height(v)|, |depth(v)| to account for their presence }
if  mem[   q+2 ].hh.rh =empty then
  begin free_node(x,box_node_size);   mem[  v+ list_offset].hh.rh  :=y;
  end
else  begin shift_up:=font_info[ 11+param_base[ eqtb[  math_font_base+   3+   cur_size].hh.rh   ]].int   - mem[ x+depth_offset].int  ;
  if shift_up<font_info[ 9+param_base[ eqtb[  math_font_base+   3+   cur_size].hh.rh   ]].int    then shift_up:=font_info[ 9+param_base[ eqtb[  math_font_base+   3+   cur_size].hh.rh   ]].int   ;
  p:=new_kern(shift_up);  mem[ p].hh.rh :=y;  mem[ x].hh.rh :=p;

  p:=new_kern(font_info[ 13+param_base[ eqtb[  math_font_base+   3+   cur_size].hh.rh   ]].int   );  mem[ p].hh.rh :=x;   mem[  v+ list_offset].hh.rh  :=p;
   mem[ v+height_offset].int  := mem[ v+height_offset].int  +font_info[ 13+param_base[ eqtb[  math_font_base+   3+   cur_size].hh.rh   ]].int   + mem[ x+height_offset].int  + mem[ x+depth_offset].int  +shift_up;
  end;
if  mem[   q+3 ].hh.rh =empty then free_node(z,box_node_size)
else  begin shift_down:=font_info[ 12+param_base[ eqtb[  math_font_base+   3+   cur_size].hh.rh   ]].int   - mem[ z+height_offset].int  ;
  if shift_down<font_info[ 10+param_base[ eqtb[  math_font_base+   3+   cur_size].hh.rh   ]].int    then shift_down:=font_info[ 10+param_base[ eqtb[  math_font_base+   3+   cur_size].hh.rh   ]].int   ;
  p:=new_kern(shift_down);  mem[ y].hh.rh :=p;  mem[ p].hh.rh :=z;

  p:=new_kern(font_info[ 13+param_base[ eqtb[  math_font_base+   3+   cur_size].hh.rh   ]].int   );  mem[ z].hh.rh :=p;
   mem[ v+depth_offset].int  := mem[ v+depth_offset].int  +font_info[ 13+param_base[ eqtb[  math_font_base+   3+   cur_size].hh.rh   ]].int   + mem[ z+height_offset].int  + mem[ z+depth_offset].int  +shift_down;
  end

;
mem[  q+1 ].int :=v;
end

;
make_op:=delta;
end;


procedure make_ord( q:halfword );
label restart,exit;
var a:integer; {address of lig/kern instruction}
 p, r:halfword ; {temporary registers for list manipulation}
begin restart:{  } 

if  mem[   q+3 ].hh.rh =empty then if  mem[   q+2 ].hh.rh =empty then
 if  mem[   q+1 ].hh.rh =math_char then
  begin p:= mem[ q].hh.rh ;
  if p<>-{0xfffffff=}268435455   then if ( mem[ p].hh.b0 >=ord_noad)and( mem[ p].hh.b0 <=punct_noad) then
    if  mem[   p+1 ].hh.rh =math_char then
    if   mem[   p+1 ].hh.b0 =  mem[   q+1 ].hh.b0  then
      begin  mem[   q+1 ].hh.rh :=math_text_char;
      fetch( q+1 );
      if ((  cur_i. b2 ) mod 4) =lig_tag then
        begin a:=lig_kern_base[ cur_f]+ cur_i.b3 ;
        cur_c:=  mem[   p+1 ].hh.b1 ;
        cur_i:=font_info[a].qqqq;
        if  cur_i.b0 > 128   then
          begin a:=lig_kern_base[ cur_f]+256*  cur_i.b2 +  cur_i.b3 +32768-256*(128+min_quarterword)  ;
          cur_i:=font_info[a].qqqq;
          end;
         while true do   begin 
{ If instruction |cur_i| is a kern with |cur_c|, attach the kern after~|q|; or if it is a ligature with |cur_c|, combine noads |q| and~|p| appropriately; then |return| if the cursor has moved past a noad, or |goto restart| }
if  cur_i.b1 =cur_c then if  cur_i.b0 <= 128   then
  if  cur_i.b2 >= 128   then
    begin p:=new_kern(font_info[kern_base[ cur_f]+256*  cur_i.b2 +  cur_i.b3 ].int  );
     mem[ p].hh.rh := mem[ q].hh.rh ;  mem[ q].hh.rh :=p;  goto exit ;
    end
  else  begin begin if interrupt<>0 then pause_for_instructions; end ; {allow a way out of infinite ligature loop}
    case  cur_i.b2  of
   1 , 5 :   mem[   q+1 ].hh.b1 := cur_i.b3 ; {\.[=:\?], \.[=:\?>]}
   2 , 6 :   mem[   p+1 ].hh.b1 := cur_i.b3 ; {\.[\?=:], \.[\?=:>]}
   3 , 7 , 11 :begin r:=new_noad; {\.[\?=:\?], \.[\?=:\?>], \.[\?=:\?>>]}
        mem[   r+1 ].hh.b1 := cur_i.b3 ;
        mem[   r+1 ].hh.b0 :=  mem[   q+1 ].hh.b0 ;

       mem[ q].hh.rh :=r;  mem[ r].hh.rh :=p;
      if  cur_i.b2 < 11  then  mem[   r+1 ].hh.rh :=math_char
      else  mem[   r+1 ].hh.rh :=math_text_char; {prevent combination}
      end;
     else  begin  mem[ q].hh.rh := mem[ p].hh.rh ;
        mem[   q+1 ].hh.b1 := cur_i.b3 ; {\.[=:]}
      mem[ q+3 ]:=mem[ p+3 ]; mem[ q+2 ]:=mem[ p+2 ];

      free_node(p,noad_size);
      end
     end ;
    if  cur_i.b2 > 3  then  goto exit ;
     mem[   q+1 ].hh.rh :=math_char; goto restart;
    end

;
          if  cur_i.b0 >= 128   then  goto exit ;
          a:=a+   cur_i.b0  +1;
          cur_i:=font_info[a].qqqq;
          end;
        end;
      end;
  end;
exit:end;


procedure make_scripts( q:halfword ; delta:scaled);
var p, x, y, z:halfword ; {temporary registers for box construction}
 shift_up, shift_down, clr:scaled; {dimensions in the calculation}
 t:small_number; {subsidiary size code}
begin p:=mem[  q+1 ].int ;
if  ( p>=hi_mem_min)  then
  begin shift_up:=0; shift_down:=0;
  end
else  begin z:=hpack(p,0,additional );
  if cur_style<script_style then t:=script_size else t:=script_script_size;
  shift_up:= mem[ z+height_offset].int  -font_info[ 18+param_base[ eqtb[  math_font_base+   2+    t].hh.rh   ]].int  ;
  shift_down:= mem[ z+depth_offset].int  +font_info[ 19+param_base[ eqtb[  math_font_base+   2+    t].hh.rh   ]].int  ;
  free_node(z,box_node_size);
  end;
if  mem[   q+2 ].hh.rh =empty then
  
{ Construct a subscript box |x| when there is no superscript }
begin x:=clean_box( q+3 ,2*( cur_style div 4)+script_style+cramped );
 mem[ x+width_offset].int  := mem[ x+width_offset].int  +eqtb[dimen_base+ script_space_code].int   ;
if shift_down<font_info[ 16+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int   then shift_down:=font_info[ 16+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  ;
clr:= mem[ x+height_offset].int  -(abs(font_info[ 5+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  *4) div 5);
if shift_down<clr then shift_down:=clr;
 mem[ x+4].int  :=shift_down;
end


else  begin 
{ Construct a superscript box |x| }
begin x:=clean_box( q+2 ,2*( cur_style div 4)+script_style+( cur_style mod 2) );
 mem[ x+width_offset].int  := mem[ x+width_offset].int  +eqtb[dimen_base+ script_space_code].int   ;
if odd(cur_style) then clr:=font_info[ 15+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  
else if cur_style<text_style then clr:=font_info[ 13+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  
else clr:=font_info[ 14+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  ;
if shift_up<clr then shift_up:=clr;
clr:= mem[ x+depth_offset].int  +(abs(font_info[ 5+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  ) div 4);
if shift_up<clr then shift_up:=clr;
end

;
  if  mem[   q+3 ].hh.rh =empty then  mem[ x+4].int  :=-shift_up
  else 
{ Construct a sub/superscript combination box |x|, with the superscript offset by |delta| }
begin y:=clean_box( q+3 ,2*( cur_style div 4)+script_style+cramped );
 mem[ y+width_offset].int  := mem[ y+width_offset].int  +eqtb[dimen_base+ script_space_code].int   ;
if shift_down<font_info[ 17+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int   then shift_down:=font_info[ 17+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  ;
clr:=4*font_info[ 8+param_base[ eqtb[  math_font_base+   3+   cur_size].hh.rh   ]].int   -
  ((shift_up- mem[ x+depth_offset].int  )-( mem[ y+height_offset].int  -shift_down));
if clr>0 then
  begin shift_down:=shift_down+clr;
  clr:=(abs(font_info[ 5+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  *4) div 5)-(shift_up- mem[ x+depth_offset].int  );
  if clr>0 then
    begin shift_up:=shift_up+clr;
    shift_down:=shift_down-clr;
    end;
  end;
 mem[ x+4].int  :=delta; {superscript is |delta| to the right of the subscript}
p:=new_kern((shift_up- mem[ x+depth_offset].int  )-( mem[ y+height_offset].int  -shift_down));  mem[ x].hh.rh :=p;  mem[ p].hh.rh :=y;
x:=vpackage( x, 0,additional ,{07777777777=}1073741823 ) ;  mem[ x+4].int  :=shift_down;
end

;
  end;
if mem[  q+1 ].int =-{0xfffffff=}268435455   then mem[  q+1 ].int :=x
else  begin p:=mem[  q+1 ].int ;
  while  mem[ p].hh.rh <>-{0xfffffff=}268435455   do p:= mem[ p].hh.rh ;
   mem[ p].hh.rh :=x;
  end;
end;


function make_left_right( q:halfword ; style:small_number;
   max_d, max_h:scaled):small_number;
var delta, delta1, delta2:scaled; {dimensions used in the calculation}
begin if style<script_style then cur_size:=text_size
else cur_size:=16*((style-text_style) div 2);
delta2:=max_d+font_info[ 22+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  ;
delta1:=max_h+max_d-delta2;
if delta2>delta1 then delta1:=delta2; {|delta1| is max distance from axis}
delta:=(delta1 div 500)*eqtb[int_base+ delimiter_factor_code].int  ;
delta2:=delta1+delta1-eqtb[dimen_base+ delimiter_shortfall_code].int   ;
if delta<delta2 then delta:=delta2;
mem[  q+1 ].int :=var_delimiter( q+1 ,cur_size,delta);
make_left_right:= mem[ q].hh.b0 -(left_noad-open_noad); {|open_noad| or |close_noad|}
end;


procedure mlist_to_hlist;
label reswitch, check_dimensions, done_with_noad, done_with_node, delete_q,
  done;
var mlist:halfword ; {beginning of the given list}
 penalties:boolean; {should penalty nodes be inserted?}
 style:small_number; {the given style}
 save_style:small_number; {holds |cur_style| during recursion}
 q:halfword ; {runs through the mlist}
 r:halfword ; {the most recent noad preceding |q|}
 r_type:small_number; {the |type| of noad |r|, or |op_noad| if |r=null|}
 t:small_number; {the effective |type| of noad |q| during the second pass}
 p, x, y, z: halfword ; {temporary registers for list construction}
 pen:integer; {a penalty to be inserted}
 s:small_number; {the size of a noad to be deleted}
 max_h, max_d:scaled; {maximum height and depth of the list translated so far}
 delta:scaled; {offset between subscript and superscript}
begin mlist:=cur_mlist; penalties:=mlist_penalties;
style:=cur_style; {tuck global parameters away as local variables}
q:=mlist; r:=-{0xfffffff=}268435455  ; r_type:=op_noad; max_h:=0; max_d:=0;

{ Set up the values of |cur_size| and |cur_mu|, based on |cur_style| }
begin if cur_style<script_style then cur_size:=text_size
else cur_size:=16*((cur_style-text_style) div 2);
cur_mu:=x_over_n(font_info[ 6+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  ,18);
end

;
while q<>-{0xfffffff=}268435455   do 
{ Process node-or-noad |q| as much as possible in preparation for the second pass of |mlist_to_hlist|, then move to the next item in the mlist }
begin 
{ Do first-pass processing based on |type(q)|; |goto done_with_noad| if a noad has been fully processed, |goto check_dimensions| if it has been translated into |new_hlist(q)|, or |goto done_with_node| if a node has been fully processed }
reswitch: delta:=0;
case  mem[ q].hh.b0  of
bin_noad: case r_type of
  bin_noad,op_noad,rel_noad,open_noad,punct_noad,left_noad:
    begin  mem[ q].hh.b0 :=ord_noad; goto reswitch;
    end;
   else   
   end ;
rel_noad,close_noad,punct_noad,right_noad: begin{  } 

  
{ Convert \(a)a final |bin_noad| to an |ord_noad| }
if r_type=bin_noad then  mem[ r].hh.b0 :=ord_noad

;
  if  mem[ q].hh.b0 =right_noad then goto done_with_noad;
  end;
{ \4 }
{ Cases for noads that can follow a |bin_noad| }
left_noad: goto done_with_noad;
fraction_noad: begin make_fraction(q); goto check_dimensions;
  end;
op_noad: begin delta:=make_op(q);
  if  mem[ q].hh.b1 =limits then goto check_dimensions;
  end;
ord_noad: make_ord(q);
open_noad,inner_noad:  ;
radical_noad: make_radical(q);
over_noad: make_over(q);
under_noad: make_under(q);
accent_noad: make_math_accent(q);
vcenter_noad: make_vcenter(q);

 
{ \4 }
{ Cases for nodes that can appear in an mlist, after which we |goto done_with_node| }
style_node: begin cur_style:= mem[ q].hh.b1 ;
  
{ Set up the values of |cur_size| and |cur_mu|, based on |cur_style| }
begin if cur_style<script_style then cur_size:=text_size
else cur_size:=16*((cur_style-text_style) div 2);
cur_mu:=x_over_n(font_info[ 6+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  ,18);
end

;
  goto done_with_node;
  end;
choice_node: 
{ Change this node to a style node followed by the correct choice, then |goto done_with_node| }
begin case cur_style div 2 of
0: begin p:=  mem[  q+ 1].hh.lh  ;   mem[  q+ 1].hh.lh  :=-{0xfffffff=}268435455  ; end ; {|display_style=0|}
1: begin p:=  mem[  q+ 1].hh.rh  ;   mem[  q+ 1].hh.rh  :=-{0xfffffff=}268435455  ; end ; {|text_style=2|}
2: begin p:=  mem[  q+ 2].hh.lh  ;   mem[  q+ 2].hh.lh  :=-{0xfffffff=}268435455  ; end ; {|script_style=4|}
3: begin p:=  mem[  q+ 2].hh.rh  ;   mem[  q+ 2].hh.rh  :=-{0xfffffff=}268435455  ; end ; {|script_script_style=6|}
end; {there are no other cases}
flush_node_list( mem[  q+ 1].hh.lh  );
flush_node_list( mem[  q+ 1].hh.rh  );
flush_node_list( mem[  q+ 2].hh.lh  );
flush_node_list( mem[  q+ 2].hh.rh  );

 mem[ q].hh.b0 :=style_node;  mem[ q].hh.b1 :=cur_style;  mem[ q+width_offset].int  :=0;  mem[ q+depth_offset].int  :=0;
if p<>-{0xfffffff=}268435455   then
  begin z:= mem[ q].hh.rh ;  mem[ q].hh.rh :=p;
  while  mem[ p].hh.rh <>-{0xfffffff=}268435455   do p:= mem[ p].hh.rh ;
   mem[ p].hh.rh :=z;
  end;
goto done_with_node;
end

;
ins_node,mark_node,adjust_node,
  whatsit_node,penalty_node,disc_node: goto done_with_node;
rule_node: begin if  mem[ q+height_offset].int  >max_h then max_h:= mem[ q+height_offset].int  ;
  if  mem[ q+depth_offset].int  >max_d then max_d:= mem[ q+depth_offset].int  ; goto done_with_node;
  end;
glue_node: begin 
{ Convert \(m)math glue to ordinary glue }
if  mem[ q].hh.b1 =mu_glue then
  begin x:=  mem[  q+ 1].hh.lh  ;
  y:=math_glue(x,cur_mu); delete_glue_ref(x);   mem[  q+ 1].hh.lh  :=y;
   mem[ q].hh.b1 :=normal;
  end
else if (cur_size<>text_size)and( mem[ q].hh.b1 =cond_math_glue) then
  begin p:= mem[ q].hh.rh ;
  if p<>-{0xfffffff=}268435455   then if ( mem[ p].hh.b0 =glue_node)or( mem[ p].hh.b0 =kern_node) then
    begin  mem[ q].hh.rh := mem[ p].hh.rh ;  mem[ p].hh.rh :=-{0xfffffff=}268435455  ; flush_node_list(p);
    end;
  end

;
  goto done_with_node;
  end;
kern_node: begin math_kern(q,cur_mu); goto done_with_node;
  end;

 
 else  confusion({"mlist1"=}903)
{ \xref[this can't happen mlist1][\quad mlist1] }
 end ;


{ Convert \(n)|nucleus(q)| to an hlist and attach the sub/superscripts }
case  mem[   q+1 ].hh.rh  of
math_char, math_text_char:
  
{ Create a character node |p| for |nucleus(q)|, possibly followed by a kern node for the italic correction, and set |delta| to the italic correction if a subscript is present }
begin fetch( q+1 );
if ( cur_i.b0>min_quarterword)  then
  begin delta:=font_info[italic_base[ cur_f]+(  cur_i. b2 ) div 4].int  ; p:=new_character(cur_f, cur_c );
  if ( mem[   q+1 ].hh.rh =math_text_char)and(font_info[ space_code+param_base[ cur_f]].int  <>0) then
    delta:=0; {no italic correction in mid-word of text font}
  if ( mem[   q+3 ].hh.rh =empty)and(delta<>0) then
    begin  mem[ p].hh.rh :=new_kern(delta); delta:=0;
    end;
  end
else p:=-{0xfffffff=}268435455  ;
end

;
empty: p:=-{0xfffffff=}268435455  ;
sub_box: p:= mem[   q+1 ].hh.lh ;
sub_mlist: begin cur_mlist:= mem[   q+1 ].hh.lh ; save_style:=cur_style;
  mlist_penalties:=false; mlist_to_hlist; {recursive call}
{ \xref[recursion] }
  cur_style:=save_style; 
{ Set up the values... }
begin if cur_style<script_style then cur_size:=text_size
else cur_size:=16*((cur_style-text_style) div 2);
cur_mu:=x_over_n(font_info[ 6+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  ,18);
end

;
  p:=hpack( mem[ mem_top-3 ].hh.rh ,0,additional );
  end;
 else  confusion({"mlist2"=}904)
{ \xref[this can't happen mlist2][\quad mlist2] }
 end ;

mem[  q+1 ].int :=p;
if ( mem[   q+3 ].hh.rh =empty)and( mem[   q+2 ].hh.rh =empty) then
  goto check_dimensions;
make_scripts(q,delta)



;
check_dimensions: z:=hpack(mem[  q+1 ].int ,0,additional );
if  mem[ z+height_offset].int  >max_h then max_h:= mem[ z+height_offset].int  ;
if  mem[ z+depth_offset].int  >max_d then max_d:= mem[ z+depth_offset].int  ;
free_node(z,box_node_size);
done_with_noad: r:=q; r_type:= mem[ r].hh.b0 ;
done_with_node: q:= mem[ q].hh.rh ;
end

;

{ Convert \(a)a final |bin_noad| to an |ord_noad| }
if r_type=bin_noad then  mem[ r].hh.b0 :=ord_noad

;

{ Make a second pass over the mlist, removing all noads and inserting the proper spacing and penalties }
p:=mem_top-3 ;  mem[ p].hh.rh :=-{0xfffffff=}268435455  ; q:=mlist; r_type:=0; cur_style:=style;

{ Set up the values of |cur_size| and |cur_mu|, based on |cur_style| }
begin if cur_style<script_style then cur_size:=text_size
else cur_size:=16*((cur_style-text_style) div 2);
cur_mu:=x_over_n(font_info[ 6+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  ,18);
end

;
while q<>-{0xfffffff=}268435455   do
  begin 
{ If node |q| is a style node, change the style and |goto delete_q|; otherwise if it is not a noad, put it into the hlist, advance |q|, and |goto done|; otherwise set |s| to the size of noad |q|, set |t| to the associated type (|ord_noad.. inner_noad|), and set |pen| to the associated penalty }
t:=ord_noad; s:=noad_size; pen:=inf_penalty;
case  mem[ q].hh.b0  of
op_noad,open_noad,close_noad,punct_noad,inner_noad: t:= mem[ q].hh.b0 ;
bin_noad: begin t:=bin_noad; pen:=eqtb[int_base+ bin_op_penalty_code].int  ;
  end;
rel_noad: begin t:=rel_noad; pen:=eqtb[int_base+ rel_penalty_code].int  ;
  end;
ord_noad,vcenter_noad,over_noad,under_noad:  ;
radical_noad: s:=radical_noad_size;
accent_noad: s:=accent_noad_size;
fraction_noad: s:=fraction_noad_size;
left_noad,right_noad: t:=make_left_right(q,style,max_d,max_h);
style_node: 
{ Change the current style and |goto delete_q| }
begin cur_style:= mem[ q].hh.b1 ; s:=style_node_size;

{ Set up the values of |cur_size| and |cur_mu|, based on |cur_style| }
begin if cur_style<script_style then cur_size:=text_size
else cur_size:=16*((cur_style-text_style) div 2);
cur_mu:=x_over_n(font_info[ 6+param_base[ eqtb[  math_font_base+   2+    cur_size].hh.rh   ]].int  ,18);
end

;
goto delete_q;
end

;
whatsit_node,penalty_node,rule_node,disc_node,adjust_node,ins_node,mark_node,
 glue_node,kern_node:{  } 

  begin  mem[ p].hh.rh :=q; p:=q; q:= mem[ q].hh.rh ;  mem[ p].hh.rh :=-{0xfffffff=}268435455  ; goto done;
  end;
 else  confusion({"mlist3"=}905)
{ \xref[this can't happen mlist3][\quad mlist3] }
 end 

;
  
{ Append inter-element spacing based on |r_type| and |t| }
if r_type>0 then {not the first noad}
  begin case   str_pool[ r_type* 8+ t+ magic_offset]  of
  {"0"=}48: x:=0;
  {"1"=}49: if cur_style<script_style then x:=thin_mu_skip_code else x:=0;
  {"2"=}50: x:=thin_mu_skip_code;
  {"3"=}51: if cur_style<script_style then x:=med_mu_skip_code else x:=0;
  {"4"=}52: if cur_style<script_style then x:=thick_mu_skip_code else x:=0;
   else  confusion({"mlist4"=}907)
{ \xref[this can't happen mlist4][\quad mlist4] }
   end ;
  if x<>0 then
    begin y:=math_glue( eqtb[  glue_base+   x].hh.rh   ,cur_mu);
    z:=new_glue(y);   mem[  y].hh.rh  :=-{0xfffffff=}268435455  ;  mem[ p].hh.rh :=z; p:=z;

     mem[ z].hh.b1 :=x+1; {store a symbolic subtype}
    end;
  end

;
  
{ Append any |new_hlist| entries for |q|, and any appropriate penalties }
if mem[  q+1 ].int <>-{0xfffffff=}268435455   then
  begin  mem[ p].hh.rh :=mem[  q+1 ].int ;
  repeat p:= mem[ p].hh.rh ;
  until  mem[ p].hh.rh =-{0xfffffff=}268435455  ;
  end;
if penalties then if  mem[ q].hh.rh <>-{0xfffffff=}268435455   then if pen<inf_penalty then
  begin r_type:= mem[  mem[  q].hh.rh ].hh.b0 ;
  if r_type<>penalty_node then if r_type<>rel_noad then
    begin z:=new_penalty(pen);  mem[ p].hh.rh :=z; p:=z;
    end;
  end

;
  r_type:=t;
  delete_q: r:=q; q:= mem[ q].hh.rh ; free_node(r,s);
  done: end

;
end;



{ 768. \[37] Alignment }

{tangle:pos tex.web:15108:18: }

{ It's sort of a miracle whenever \.[\\halign] and \.[\\valign] work, because
they cut across so many of the control structures of \TeX.

Therefore the
present page is probably not the best place for a beginner to start reading
this program; it is better to master everything else first.

Let us focus our thoughts on an example of what the input might be, in order
to get some idea about how the alignment miracle happens. The example doesn't
do anything useful, but it is sufficiently general to indicate all of the
special cases that must be dealt with; please do not be disturbed by its
apparent complexity and meaninglessness.
$$\vbox[\halign[\.[#]\hfil\cr
[]\\tabskip 2pt plus 3pt\cr
[]\\halign to 300pt\[u1\#v1\&\cr
\hskip 50pt\\tabskip 1pt plus 1fil u2\#v2\&\cr
\hskip 50pt u3\#v3\\cr\cr
\hskip 25pt a1\&\\omit a2\&\\vrule\\cr\cr
\hskip 25pt \\noalign\[\\vskip 3pt\]\cr
\hskip 25pt b1\\span b2\\cr\cr
\hskip 25pt \\omit\&c2\\span\\omit\\cr\]\cr]]$$
Here's what happens:

\yskip
(0) When `\.[\\halign to 300pt\[]' is scanned, the |scan_spec| routine
places the 300pt dimension onto the |save_stack|, and an |align_group|
code is placed above it. This will make it possible to complete the alignment
when the matching `\.\]' is found.

(1) The preamble is scanned next. Macros in the preamble are not expanded,
\xref[preamble]
except as part of a tabskip specification. For example, if \.[u2] had been
a macro in the preamble above, it would have been expanded, since \TeX\
must look for `\.[minus...]' as part of the tabskip glue. A ``preamble list''
is constructed based on the user's preamble; in our case it contains the
following seven items:
$$\vbox[\halign[\.[#]\hfil\qquad&(#)\hfil\cr
[]\\glue 2pt plus 3pt&the tabskip preceding column 1\cr
[]\\alignrecord, width $-\infty$&preamble info for column 1\cr
[]\\glue 2pt plus 3pt&the tabskip between columns 1 and 2\cr
[]\\alignrecord, width $-\infty$&preamble info for column 2\cr
[]\\glue 1pt plus 1fil&the tabskip between columns 2 and 3\cr
[]\\alignrecord, width $-\infty$&preamble info for column 3\cr
[]\\glue 1pt plus 1fil&the tabskip following column 3\cr]]$$
These ``alignrecord'' entries have the same size as an |unset_node|,
since they will later be converted into such nodes. However, at the
moment they have no |type| or |subtype| fields; they have |info| fields
instead, and these |info| fields are initially set to the value |end_span|,
for reasons explained below. Furthermore, the alignrecord nodes have no
|height| or |depth| fields; these are renamed |u_part| and |v_part|,
and they point to token lists for the templates of the alignment.
For example, the |u_part| field in the first alignrecord points to the
token list `\.[u1]', i.e., the template preceding the `\.\#' for column~1.

(2) \TeX\ now looks at what follows the \.[\\cr] that ended the preamble.
It is not `\.[\\noalign]' or `\.[\\omit]', so this input is put back to
be read again, and the template `\.[u1]' is fed to the scanner. Just
before reading `\.[u1]', \TeX\ goes into restricted horizontal mode.
Just after reading `\.[u1]', \TeX\ will see `\.[a1]', and then (when the
[\.\&] is sensed) \TeX\ will see `\.[v1]'. Then \TeX\ scans an |endv|
token, indicating the end of a column. At this point an |unset_node| is
created, containing the contents of the current hlist (i.e., `\.[u1a1v1]').
The natural width of this unset node replaces the |width| field of the
alignrecord for column~1; in general, the alignrecords will record the
maximum natural width that has occurred so far in a given column.

(3) Since `\.[\\omit]' follows the `\.\&', the templates for column~2
are now bypassed. Again \TeX\ goes into restricted horizontal mode and
makes an |unset_node| from the resulting hlist; but this time the
hlist contains simply `\.[a2]'. The natural width of the new unset box
is remembered in the |width| field of the alignrecord for column~2.

(4) A third |unset_node| is created for column 3, using essentially the
mechanism that worked for column~1; this unset box contains `\.[u3\\vrule
v3]'. The vertical rule in this case has running dimensions that will later
extend to the height and depth of the whole first row, since each |unset_node|
in a row will eventually inherit the height and depth of its enclosing box.

(5) The first row has now ended; it is made into a single unset box
comprising the following seven items:
$$\vbox[\halign[\hbox to 325pt[\qquad\.[#]\hfil]\cr
[]\\glue 2pt plus 3pt\cr
[]\\unsetbox for 1 column: u1a1v1\cr
[]\\glue 2pt plus 3pt\cr
[]\\unsetbox for 1 column: a2\cr
[]\\glue 1pt plus 1fil\cr
[]\\unsetbox for 1 column: u3\\vrule v3\cr
[]\\glue 1pt plus 1fil\cr]]$$
The width of this unset row is unimportant, but it has the correct height
and depth, so the correct baselineskip glue will be computed as the row
is inserted into a vertical list.

(6) Since `\.[\\noalign]' follows the current \.[\\cr], \TeX\ appends
additional material (in this case \.[\\vskip 3pt]) to the vertical list.
While processing this material, \TeX\ will be in internal vertical
mode, and |no_align_group| will be on |save_stack|.

(7) The next row produces an unset box that looks like this:
$$\vbox[\halign[\hbox to 325pt[\qquad\.[#]\hfil]\cr
[]\\glue 2pt plus 3pt\cr
[]\\unsetbox for 2 columns: u1b1v1u2b2v2\cr
[]\\glue 1pt plus 1fil\cr
[]\\unsetbox for 1 column: [\rm(empty)]\cr
[]\\glue 1pt plus 1fil\cr]]$$
The natural width of the unset box that spans columns 1~and~2 is stored
in a ``span node,'' which we will explain later; the |info| field of the
alignrecord for column~1 now points to the new span node, and the |info|
of the span node points to |end_span|.

(8) The final row produces the unset box
$$\vbox[\halign[\hbox to 325pt[\qquad\.[#]\hfil]\cr
[]\\glue 2pt plus 3pt\cr
[]\\unsetbox for 1 column: [\rm(empty)]\cr
[]\\glue 2pt plus 3pt\cr
[]\\unsetbox for 2 columns: u2c2v2\cr
[]\\glue 1pt plus 1fil\cr]]$$
A new span node is attached to the alignrecord for column 2.

(9) The last step is to compute the true column widths and to change all the
unset boxes to hboxes, appending the whole works to the vertical list that
encloses the \.[\\halign]. The rules for deciding on the final widths of
each unset column box will be explained below.

\yskip\noindent
Note that as \.[\\halign] is being processed, we fearlessly give up control
to the rest of \TeX. At critical junctures, an alignment routine is
called upon to step in and do some little action, but most of the time
these routines just lurk in the background. It's something like
post-hypnotic suggestion. }

{ 769. }

{tangle:pos tex.web:15239:1: }

{ We have mentioned that alignrecords contain no |height| or |depth| fields.
Their |glue_sign| and |glue_order| are pre-empted as well, since it
is necessary to store information about what to do when a template ends.
This information is called the |extra_info| field. }

{ 772. }

{tangle:pos tex.web:15281:1: }

{ Alignment stack maintenance is handled by a pair of trivial routines
called |push_alignment| and |pop_alignment|. } procedure push_alignment;
var p:halfword ; {the new alignment stack node}
begin p:=get_node(align_stack_node_size);
 mem[ p].hh.rh :=align_ptr;  mem[ p].hh.lh :=cur_align;
  mem[  p+ 1].hh.lh  := mem[ mem_top-8 ].hh.rh  ;   mem[  p+ 1].hh.rh  :=cur_span;
mem[p+2].int:=cur_loop; mem[p+3].int:=align_state;
 mem[ p+ 4].hh.lh :=cur_head;  mem[ p+ 4].hh.rh :=cur_tail;
align_ptr:=p;
cur_head:=get_avail;
end;


procedure pop_alignment;
var p:halfword ; {the top alignment stack node}
begin  begin  mem[  cur_head].hh.rh :=avail; avail:= cur_head; ifdef('STAT')  decr(dyn_used); endif('STAT')  end ;
p:=align_ptr;
cur_tail:= mem[ p+ 4].hh.rh ; cur_head:= mem[ p+ 4].hh.lh ;
align_state:=mem[p+3].int; cur_loop:=mem[p+2].int;
cur_span:=  mem[  p+ 1].hh.rh  ;  mem[ mem_top-8 ].hh.rh  :=  mem[  p+ 1].hh.lh  ;
cur_align:= mem[ p].hh.lh ; align_ptr:= mem[ p].hh.rh ;
free_node(p,align_stack_node_size);
end;



{ 773. }

{tangle:pos tex.web:15306:1: }

{ \TeX\ has eight procedures that govern alignments: |init_align| and
|fin_align| are used at the very beginning and the very end; |init_row| and
|fin_row| are used at the beginning and end of individual rows; |init_span|
is used at the beginning of a sequence of spanned columns (possibly involving
only one column); |init_col| and |fin_col| are used at the beginning and
end of individual columns; and |align_peek| is used after \.[\\cr] to see
whether the next item is \.[\\noalign].

We shall consider these routines in the order they are first used during
the course of a complete \.[\\halign], namely |init_align|, |align_peek|,
|init_row|, |init_span|, |init_col|, |fin_col|, |fin_row|, |fin_align|. }

{ 774. }

{tangle:pos tex.web:15318:1: }

{ When \.[\\halign] or \.[\\valign] has been scanned in an appropriate
mode, \TeX\ calls |init_align|, whose task is to get everything off to a
good start. This mostly involves scanning the preamble and putting its
information into the preamble list.
\xref[preamble] } { \4 }
{ Declare the procedure called |get_preamble_token| }
procedure get_preamble_token;
label restart;
begin restart: get_token;
while (cur_chr=span_code)and(cur_cmd=tab_mark) do
  begin get_token; {this token will be expanded once}
  if cur_cmd>max_command then
    begin expand; get_token;
    end;
  end;
if cur_cmd=endv then
  fatal_error({"(interwoven alignment preambles are not allowed)"=}602);
{ \xref[interwoven alignment preambles...] }
if (cur_cmd=assign_glue)and(cur_chr=glue_base+tab_skip_code) then
  begin scan_optional_equals; scan_glue(glue_val);
  if eqtb[int_base+ global_defs_code].int  >0 then geq_define(glue_base+tab_skip_code,glue_ref,cur_val)
  else eq_define(glue_base+tab_skip_code,glue_ref,cur_val);
  goto restart;
  end;
end;

{  }

procedure align_peek; forward;{ \2 }

procedure normal_paragraph; forward;{ \2 }

procedure init_align;
label done, done1, done2, continue;
var save_cs_ptr:halfword ; {|warning_index| value for error messages}
 p:halfword ; {for short-term temporary use}
begin save_cs_ptr:=cur_cs; {\.[\\halign] or \.[\\valign], usually}
push_alignment; align_state:=-1000000; {enter a new alignment level}

{ Check for improper alignment in displayed math }
if (cur_list.mode_field =mmode)and((cur_list.tail_field <>cur_list.head_field )or(cur_list.aux_field .int <>-{0xfffffff=}268435455  )) then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Improper "=} 690); end ; print_esc({"halign"=}528); print({" inside $$'s"=}908);
{ \xref[Improper \\halign...] }
   begin help_ptr:=3; help_line[2]:={"Displays can use special alignments (like \eqalignno)"=} 909; help_line[1]:={"only if nothing but the alignment itself is between $$'s."=} 910; help_line[0]:={"So I've deleted the formulas that preceded this alignment."=} 911; end ;
  error; flush_math;
  end

;
push_nest; {enter a new semantic level}

{ Change current mode to |-vmode| for \.[\\halign], |-hmode| for \.[\\valign] }
if cur_list.mode_field =mmode then
  begin cur_list.mode_field :=-vmode; cur_list.aux_field .int  :=nest[nest_ptr-2].aux_field.int ;
  end
else if cur_list.mode_field >0 then   cur_list.mode_field :=- cur_list.mode_field  

;
scan_spec(align_group,false);


{ Scan the preamble and record it in the |preamble| list }
 mem[ mem_top-8 ].hh.rh  :=-{0xfffffff=}268435455  ; cur_align:=mem_top-8 ; cur_loop:=-{0xfffffff=}268435455  ; scanner_status:=aligning;
warning_index:=save_cs_ptr; align_state:=-1000000;
  {at this point, |cur_cmd=left_brace|}
 while true do    begin 
{ Append the current tabskip glue to the preamble list }
 mem[ cur_align].hh.rh :=new_param_glue(tab_skip_code);
cur_align:= mem[ cur_align].hh.rh 

;
  if cur_cmd=car_ret then goto done; {\.[\\cr] ends the preamble}
  
{ Scan preamble text until |cur_cmd| is |tab_mark| or |car_ret|, looking for changes in the tabskip glue; append an alignrecord to the preamble list }

{ Scan the template \<u_j>, putting the resulting token list in |hold_head| }
p:=mem_top-4 ;  mem[ p].hh.rh :=-{0xfffffff=}268435455  ;
 while true do    begin get_preamble_token;
  if cur_cmd=mac_param then goto done1;
  if (cur_cmd<=car_ret)and(cur_cmd>=tab_mark)and(align_state=-1000000) then
   if (p=mem_top-4 )and(cur_loop=-{0xfffffff=}268435455  )and(cur_cmd=tab_mark)
    then cur_loop:=cur_align
   else  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Missing # inserted in alignment preamble"=} 917); end ;
{ \xref[Missing \# inserted...] }
     begin help_ptr:=3; help_line[2]:={"There should be exactly one # between &'s, when an"=} 918; help_line[1]:={"\halign or \valign is being set up. In this case you had"=} 919; help_line[0]:={"none, so I've put one in; maybe that will work."=} 920; end ;
    back_error; goto done1;
    end
  else if (cur_cmd<>spacer)or(p<>mem_top-4 ) then
    begin  mem[ p].hh.rh :=get_avail; p:= mem[ p].hh.rh ;  mem[ p].hh.lh :=cur_tok;
    end;
  end;
done1:

;
 mem[ cur_align].hh.rh :=new_null_box; cur_align:= mem[ cur_align].hh.rh ; {a new alignrecord}
 mem[ cur_align].hh.lh :=mem_top-9 ;  mem[ cur_align+width_offset].int  :=-{010000000000=}1073741824 ;
mem[ cur_align+height_offset].int := mem[ mem_top-4 ].hh.rh ;

{ Scan the template \<v_j>, putting the resulting token list in |hold_head| }
p:=mem_top-4 ;  mem[ p].hh.rh :=-{0xfffffff=}268435455  ;
 while true do    begin continue: get_preamble_token;
  if (cur_cmd<=car_ret)and(cur_cmd>=tab_mark)and(align_state=-1000000) then
    goto done2;
  if cur_cmd=mac_param then
    begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Only one # is allowed per tab"=} 921); end ;
{ \xref[Only one \# is allowed...] }
     begin help_ptr:=3; help_line[2]:={"There should be exactly one # between &'s, when an"=} 918; help_line[1]:={"\halign or \valign is being set up. In this case you had"=} 919; help_line[0]:={"more than one, so I'm ignoring all but the first."=} 922; end ;
    error; goto continue;
    end;
   mem[ p].hh.rh :=get_avail; p:= mem[ p].hh.rh ;  mem[ p].hh.lh :=cur_tok;
  end;
done2:  mem[ p].hh.rh :=get_avail; p:= mem[ p].hh.rh ;
 mem[ p].hh.lh :={07777=}4095 +frozen_end_template  {put \.[\\endtemplate] at the end}

;
mem[ cur_align+depth_offset].int := mem[ mem_top-4 ].hh.rh 

;
  end;
done: scanner_status:=normal

;
new_save_level(align_group);
if  eqtb[  every_cr_loc].hh.rh   <>-{0xfffffff=}268435455   then begin_token_list( eqtb[  every_cr_loc].hh.rh   ,every_cr_text);
align_peek; {look for \.[\\noalign] or \.[\\omit]}
end;



{ 786. }

{tangle:pos tex.web:15525:1: }

{ To start a row (i.e., a `row' that rhymes with `dough' but not with `bough'),
we enter a new semantic level, copy the first tabskip glue, and change
from internal vertical mode to restricted horizontal mode or vice versa.
The |space_factor| and |prev_depth| are not used on this semantic level,
but we clear them to zero just to be tidy. } { \4 }
{ Declare the procedure called |init_span| }
procedure init_span( p:halfword );
begin push_nest;
if cur_list.mode_field =-hmode then cur_list.aux_field .hh.lh :=1000
else  begin cur_list.aux_field .int  :=-65536000 ; normal_paragraph;
  end;
cur_span:=p;
end;

{  }

procedure init_row;
begin push_nest; cur_list.mode_field :=(-hmode-vmode)-cur_list.mode_field ;
if cur_list.mode_field =-hmode then cur_list.aux_field .hh.lh :=0  else cur_list.aux_field .int  :=0;
begin  mem[ cur_list.tail_field ].hh.rh := new_glue(   mem[    mem[ mem_top-8 ].hh.rh  + 1].hh.lh  ); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
 mem[ cur_list.tail_field ].hh.b1 :=tab_skip_code+1;

cur_align:= mem[  mem[ mem_top-8 ].hh.rh  ].hh.rh ; cur_tail:=cur_head; init_span(cur_align);
end;



{ 788. }

{tangle:pos tex.web:15553:1: }

{ When a column begins, we assume that |cur_cmd| is either |omit| or else
the current token should be put back into the input until the \<u_j>
template has been scanned.  (Note that |cur_cmd| might be |tab_mark| or
|car_ret|.)  We also assume that |align_state| is approximately 1000000 at
this time.  We remain in the same mode, and start the template if it is
called for. } procedure init_col;
begin  mem[  cur_align+ list_offset].hh.lh  :=cur_cmd;
if cur_cmd=omit then align_state:=0
else  begin back_input; begin_token_list(mem[ cur_align+height_offset].int ,u_template);
  end; {now |align_state=1000000|}
end;



{ 791. }

{tangle:pos tex.web:15592:1: }

{ When the |endv| command at the end of a \<v_j> template comes through the
scanner, things really start to happen; and it is the |fin_col| routine
that makes them happen. This routine returns |true| if a row as well as a
column has been finished. } function fin_col:boolean;
label exit;
var p:halfword ; {the alignrecord after the current one}
 q, r:halfword ; {temporary pointers for list manipulation}
 s:halfword ; {a new span node}
 u:halfword ; {a new unset box}
 w:scaled; {natural width}
 o:glue_ord; {order of infinity}
 n:halfword; {span counter}
begin if cur_align=-{0xfffffff=}268435455   then confusion({"endv"=}923);
q:= mem[ cur_align].hh.rh ; if q=-{0xfffffff=}268435455   then confusion({"endv"=}923);
{ \xref[this can't happen endv][\quad endv] }
if align_state<500000 then
  fatal_error({"(interwoven alignment preambles are not allowed)"=}602);
{ \xref[interwoven alignment preambles...] }
p:= mem[ q].hh.rh ;

{ If the preamble list has been traversed, check that the row has ended }
if (p=-{0xfffffff=}268435455  )and( mem[  cur_align+ list_offset].hh.lh  <cr_code) then
 if cur_loop<>-{0xfffffff=}268435455   then 
{ Lengthen the preamble periodically }
begin  mem[ q].hh.rh :=new_null_box; p:= mem[ q].hh.rh ; {a new alignrecord}
 mem[ p].hh.lh :=mem_top-9 ;  mem[ p+width_offset].int  :=-{010000000000=}1073741824 ; cur_loop:= mem[ cur_loop].hh.rh ;

{ Copy the templates from node |cur_loop| into node |p| }
q:=mem_top-4 ; r:=mem[ cur_loop+height_offset].int ;
while r<>-{0xfffffff=}268435455   do
  begin  mem[ q].hh.rh :=get_avail; q:= mem[ q].hh.rh ;  mem[ q].hh.lh := mem[ r].hh.lh ; r:= mem[ r].hh.rh ;
  end;
 mem[ q].hh.rh :=-{0xfffffff=}268435455  ; mem[ p+height_offset].int := mem[ mem_top-4 ].hh.rh ;
q:=mem_top-4 ; r:=mem[ cur_loop+depth_offset].int ;
while r<>-{0xfffffff=}268435455   do
  begin  mem[ q].hh.rh :=get_avail; q:= mem[ q].hh.rh ;  mem[ q].hh.lh := mem[ r].hh.lh ; r:= mem[ r].hh.rh ;
  end;
 mem[ q].hh.rh :=-{0xfffffff=}268435455  ; mem[ p+depth_offset].int := mem[ mem_top-4 ].hh.rh 

;
cur_loop:= mem[ cur_loop].hh.rh ;
 mem[ p].hh.rh :=new_glue(  mem[  cur_loop+ 1].hh.lh  );
 mem[  mem[  p].hh.rh ].hh.b1 :=tab_skip_code+1;
end


 else  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Extra alignment tab has been changed to "=} 924); end ;
{ \xref[Extra alignment tab...] }
  print_esc({"cr"=}913);
   begin help_ptr:=3; help_line[2]:={"You have given more \span or & marks than there were"=} 925; help_line[1]:={"in the preamble to the \halign or \valign now in progress."=} 926; help_line[0]:={"So I'll assume that you meant to type \cr instead."=} 927; end ;
   mem[  cur_align+ list_offset].hh.lh  :=cr_code; error;
  end

;
if  mem[  cur_align+ list_offset].hh.lh  <>span_code then
  begin unsave; new_save_level(align_group);

  
{ Package an unset box for the current column and record its width }
begin if cur_list.mode_field =-hmode then
  begin adjust_tail:=cur_tail; u:=hpack( mem[ cur_list.head_field ].hh.rh ,0,additional ); w:= mem[ u+width_offset].int  ;
  cur_tail:=adjust_tail; adjust_tail:=-{0xfffffff=}268435455  ;
  end
else  begin u:=vpackage( mem[ cur_list.head_field ].hh.rh ,0,additional ,0); w:= mem[ u+height_offset].int  ;
  end;
n:=min_quarterword; {this represents a span count of 1}
if cur_span<>cur_align then 
{ Update width entry for spanned columns }
begin q:=cur_span;
repeat incr(n); q:= mem[  mem[  q].hh.rh ].hh.rh ;
until q=cur_align;
if n>max_quarterword then confusion({"256 spans"=}928); {this can happen, but won't}
{ \xref[system dependencies] }
{ \xref[this can't happen 256 spans][\quad 256 spans] }
q:=cur_span; while  mem[  mem[  q].hh.lh ].hh.rh <n do q:= mem[ q].hh.lh ;
if  mem[  mem[  q].hh.lh ].hh.rh >n then
  begin s:=get_node(span_node_size);  mem[ s].hh.lh := mem[ q].hh.lh ;  mem[ s].hh.rh :=n;
   mem[ q].hh.lh :=s;  mem[ s+width_offset].int  :=w;
  end
else if  mem[  mem[  q].hh.lh +width_offset].int  <w then  mem[  mem[  q].hh.lh +width_offset].int  :=w;
end


else if w> mem[ cur_align+width_offset].int   then  mem[ cur_align+width_offset].int  :=w;
 mem[ u].hh.b0 :=unset_node;  mem[ u].hh.b1 :=n;


{ Determine the stretch order }
if total_stretch[filll]<>0 then o:=filll
else if total_stretch[fill]<>0 then o:=fill
else if total_stretch[fil]<>0 then o:=fil
else o:=normal

;
  mem[  u+ list_offset].hh.b1  :=o; mem[ u+glue_offset].int  :=total_stretch[o];


{ Determine the shrink order }
if total_shrink[filll]<>0 then o:=filll
else if total_shrink[fill]<>0 then o:=fill
else if total_shrink[fil]<>0 then o:=fil
else o:=normal

;
  mem[  u+ list_offset].hh.b0  :=o;  mem[ u+4].int  :=total_shrink[o];

pop_nest;  mem[ cur_list.tail_field ].hh.rh :=u; cur_list.tail_field :=u;
end

;
  
{ Copy the tabskip glue between columns }
begin  mem[ cur_list.tail_field ].hh.rh := new_glue(   mem[    mem[    cur_align].hh.rh + 1].hh.lh  ); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
 mem[ cur_list.tail_field ].hh.b1 :=tab_skip_code+1

;
  if  mem[  cur_align+ list_offset].hh.lh  >=cr_code then
    begin fin_col:=true;  goto exit ;
    end;
  init_span(p);
  end;
align_state:=1000000; 
{ Get the next non-blank non-call token }
repeat get_x_token;
until cur_cmd<>spacer

;
cur_align:=p;
init_col; fin_col:=false;
exit: end;



{ 799. }

{tangle:pos tex.web:15714:1: }

{ At the end of a row, we append an unset box to the current vlist (for
\.[\\halign]) or the current hlist (for \.[\\valign]). This unset box
contains the unset boxes for the columns, separated by the tabskip glue.
Everything will be set later. } procedure fin_row;
var p:halfword ; {the new unset box}
begin if cur_list.mode_field =-hmode then
  begin p:=hpack( mem[ cur_list.head_field ].hh.rh ,0,additional );
  pop_nest; append_to_vlist(p);
  if cur_head<>cur_tail then
    begin  mem[ cur_list.tail_field ].hh.rh := mem[ cur_head].hh.rh ; cur_list.tail_field :=cur_tail;
    end;
  end
else  begin p:=vpackage(  mem[  cur_list.head_field ].hh.rh , 0,additional ,{07777777777=}1073741823 ) ; pop_nest;
   mem[ cur_list.tail_field ].hh.rh :=p; cur_list.tail_field :=p; cur_list.aux_field .hh.lh :=1000;
  end;
 mem[ p].hh.b0 :=unset_node; mem[ p+glue_offset].int  :=0;
if  eqtb[  every_cr_loc].hh.rh   <>-{0xfffffff=}268435455   then begin_token_list( eqtb[  every_cr_loc].hh.rh   ,every_cr_text);
align_peek;
end; {note that |glue_shrink(p)=0| since |glue_shrink==shift_amount|}



{ 800. }

{tangle:pos tex.web:15736:1: }

{ Finally, we will reach the end of the alignment, and we can breathe a
sigh of relief that memory hasn't overflowed. All the unset boxes will now be
set so that the columns line up, taking due account of spanned columns. } procedure do_assignments; forward;{ \2 }

procedure resume_after_display; forward;{ \2 }

procedure build_page; forward;{ \2 }

procedure fin_align;
var  p, q, r, s, u, v: halfword ; {registers for the list operations}
 t, w:scaled; {width of column}
 o:scaled; {shift offset for unset boxes}
 n:halfword; {matching span amount}
 rule_save:scaled; {temporary storage for |overfull_rule|}
 aux_save:memory_word; {temporary storage for |aux|}
begin if cur_group<>align_group then confusion({"align1"=}929);
{ \xref[this can't happen align][\quad align] }
unsave; {that |align_group| was for individual entries}
if cur_group<>align_group then confusion({"align0"=}930);
unsave; {that |align_group| was for the whole alignment}
if nest[nest_ptr-1].mode_field=mmode then o:=eqtb[dimen_base+ display_indent_code].int   
  else o:=0;

{ Go through the preamble list, determining the column widths and changing the alignrecords to dummy unset boxes }
q:= mem[  mem[ mem_top-8 ].hh.rh  ].hh.rh ;
repeat flush_list(mem[ q+height_offset].int ); flush_list(mem[ q+depth_offset].int );
p:= mem[  mem[  q].hh.rh ].hh.rh ;
if  mem[ q+width_offset].int  =-{010000000000=}1073741824  then
  
{ Nullify |width(q)| and the tabskip glue following this column }
begin  mem[ q+width_offset].int  :=0; r:= mem[ q].hh.rh ; s:=  mem[  r+ 1].hh.lh  ;
if s<>mem_bot  then
  begin incr(  mem[   mem_bot ].hh.rh  ) ; delete_glue_ref(s);
    mem[  r+ 1].hh.lh  :=mem_bot ;
  end;
end

;
if  mem[ q].hh.lh <>mem_top-9  then
  
{ Merge the widths in the span nodes of |q| with those of |p|, destroying the span nodes of |q| }
begin t:= mem[ q+width_offset].int  + mem[   mem[    mem[    q].hh.rh + 1].hh.lh  +width_offset].int  ;
r:= mem[ q].hh.lh ; s:=mem_top-9 ;  mem[ s].hh.lh :=p; n:=min_quarterword+1;
repeat  mem[ r+width_offset].int  := mem[ r+width_offset].int  -t; u:= mem[ r].hh.lh ;
while  mem[ r].hh.rh >n do
  begin s:= mem[ s].hh.lh ; n:= mem[  mem[  s].hh.lh ].hh.rh +1;
  end;
if  mem[ r].hh.rh <n then
  begin  mem[ r].hh.lh := mem[ s].hh.lh ;  mem[ s].hh.lh :=r; decr( mem[ r].hh.rh ); s:=r;
  end
else  begin if  mem[ r+width_offset].int  > mem[  mem[  s].hh.lh +width_offset].int   then  mem[  mem[  s].hh.lh +width_offset].int  := mem[ r+width_offset].int  ;
  free_node(r,span_node_size);
  end;
r:=u;
until r=mem_top-9 ;
end

;
 mem[ q].hh.b0 :=unset_node;  mem[ q].hh.b1 :=min_quarterword;  mem[ q+height_offset].int  :=0;
 mem[ q+depth_offset].int  :=0;   mem[  q+ list_offset].hh.b1  :=normal;   mem[  q+ list_offset].hh.b0  :=normal;
mem[ q+glue_offset].int  :=0;  mem[ q+4].int  :=0; q:=p;
until q=-{0xfffffff=}268435455  

;

{ Package the preamble list, to determine the actual tabskip glue amounts, and let |p| point to this prototype box }
save_ptr:=save_ptr-2; pack_begin_line:=-cur_list.ml_field ;
if cur_list.mode_field =-vmode then
  begin rule_save:=eqtb[dimen_base+ overfull_rule_code].int   ;
  eqtb[dimen_base+ overfull_rule_code].int   :=0; {prevent rule from being packaged}
  p:=hpack( mem[ mem_top-8 ].hh.rh  ,save_stack[save_ptr+ 1].int ,save_stack[save_ptr+ 0].int ); eqtb[dimen_base+ overfull_rule_code].int   :=rule_save;
  end
else  begin q:= mem[  mem[ mem_top-8 ].hh.rh  ].hh.rh ;
  repeat  mem[ q+height_offset].int  := mem[ q+width_offset].int  ;  mem[ q+width_offset].int  :=0; q:= mem[  mem[  q].hh.rh ].hh.rh ;
  until q=-{0xfffffff=}268435455  ;
  p:=vpackage(  mem[ mem_top-8 ].hh.rh  , save_stack[save_ptr+  1].int , save_stack[save_ptr+  0].int ,{07777777777=}1073741823 ) ;
  q:= mem[  mem[ mem_top-8 ].hh.rh  ].hh.rh ;
  repeat  mem[ q+width_offset].int  := mem[ q+height_offset].int  ;  mem[ q+height_offset].int  :=0; q:= mem[  mem[  q].hh.rh ].hh.rh ;
  until q=-{0xfffffff=}268435455  ;
  end;
pack_begin_line:=0

;

{ Set the glue in all the unset boxes of the current list }
q:= mem[ cur_list.head_field ].hh.rh ; s:=cur_list.head_field ;
while q<>-{0xfffffff=}268435455   do
  begin if not  ( q>=hi_mem_min)  then
    if  mem[ q].hh.b0 =unset_node then
      
{ Set the unset box |q| and the unset boxes in it }
begin if cur_list.mode_field =-vmode then
  begin  mem[ q].hh.b0 :=hlist_node;  mem[ q+width_offset].int  := mem[ p+width_offset].int  ;
  end
else  begin  mem[ q].hh.b0 :=vlist_node;  mem[ q+height_offset].int  := mem[ p+height_offset].int  ;
  end;
  mem[  q+ list_offset].hh.b1  :=  mem[  p+ list_offset].hh.b1  ;   mem[  q+ list_offset].hh.b0  :=  mem[  p+ list_offset].hh.b0  ;
 mem[ q+glue_offset].gr := mem[ p+glue_offset].gr ;  mem[ q+4].int  :=o;
r:= mem[   mem[   q+ list_offset].hh.rh  ].hh.rh ; s:= mem[   mem[   p+ list_offset].hh.rh  ].hh.rh ;
repeat 
{ Set the glue in node |r| and change it from an unset node }
n:= mem[ r].hh.b1 ; t:= mem[ s+width_offset].int  ; w:=t; u:=mem_top-4 ;
while n>min_quarterword do
  begin decr(n);
  
{ Append tabskip glue and an empty box to list |u|, and update |s| and |t| as the prototype nodes are passed }
s:= mem[ s].hh.rh ; v:=  mem[  s+ 1].hh.lh  ;  mem[ u].hh.rh :=new_glue(v); u:= mem[ u].hh.rh ;
 mem[ u].hh.b1 :=tab_skip_code+1; t:=t+ mem[ v+width_offset].int  ;
if   mem[  p+ list_offset].hh.b0  =stretching then
  begin if   mem[ v].hh.b0 =  mem[  p+ list_offset].hh.b1   then
    t:=t+round(   mem[  p+glue_offset].gr  * mem[ v+2].int  );
{ \xref[real multiplication] }
  end
else if   mem[  p+ list_offset].hh.b0  =shrinking then
  begin if   mem[ v].hh.b1 =  mem[  p+ list_offset].hh.b1   then
    t:=t-round(   mem[  p+glue_offset].gr  * mem[ v+3].int  );
  end;
s:= mem[ s].hh.rh ;  mem[ u].hh.rh :=new_null_box; u:= mem[ u].hh.rh ; t:=t+ mem[ s+width_offset].int  ;
if cur_list.mode_field =-vmode then  mem[ u+width_offset].int  := mem[ s+width_offset].int   else
  begin  mem[ u].hh.b0 :=vlist_node;  mem[ u+height_offset].int  := mem[ s+width_offset].int  ;
  end

;
  end;
if cur_list.mode_field =-vmode then
  
{ Make the unset node |r| into an |hlist_node| of width |w|, setting the glue as if the width were |t| }
begin  mem[ r+height_offset].int  := mem[ q+height_offset].int  ;  mem[ r+depth_offset].int  := mem[ q+depth_offset].int  ;
if t= mem[ r+width_offset].int   then
  begin   mem[  r+ list_offset].hh.b0  :=normal;   mem[  r+ list_offset].hh.b1  :=normal;
     mem[  r+glue_offset].gr :=0.0 ;
  end
else if t> mem[ r+width_offset].int   then
  begin   mem[  r+ list_offset].hh.b0  :=stretching;
  if mem[ r+glue_offset].int  =0 then    mem[  r+glue_offset].gr :=0.0 
  else  mem[ r+glue_offset].gr := ( t-  mem[  r+width_offset].int  )/ mem[  r+glue_offset].int   ;
{ \xref[real division] }
  end
else  begin   mem[  r+ list_offset].hh.b1  :=  mem[  r+ list_offset].hh.b0  ;   mem[  r+ list_offset].hh.b0  :=shrinking;
  if  mem[ r+4].int  =0 then    mem[  r+glue_offset].gr :=0.0 
  else if (  mem[  r+ list_offset].hh.b1  =normal)and( mem[ r+width_offset].int  -t> mem[ r+4].int  ) then
       mem[  r+glue_offset].gr :=1.0 
  else  mem[ r+glue_offset].gr := (  mem[  r+width_offset].int  - t)/  mem[  r+4].int   ;
  end;
 mem[ r+width_offset].int  :=w;  mem[ r].hh.b0 :=hlist_node;
end


else 
{ Make the unset node |r| into a |vlist_node| of height |w|, setting the glue as if the height were |t| }
begin  mem[ r+width_offset].int  := mem[ q+width_offset].int  ;
if t= mem[ r+height_offset].int   then
  begin   mem[  r+ list_offset].hh.b0  :=normal;   mem[  r+ list_offset].hh.b1  :=normal;
     mem[  r+glue_offset].gr :=0.0 ;
  end
else if t> mem[ r+height_offset].int   then
  begin   mem[  r+ list_offset].hh.b0  :=stretching;
  if mem[ r+glue_offset].int  =0 then    mem[  r+glue_offset].gr :=0.0 
  else  mem[ r+glue_offset].gr := ( t-  mem[  r+height_offset].int  )/ mem[  r+glue_offset].int   ;
{ \xref[real division] }
  end
else  begin   mem[  r+ list_offset].hh.b1  :=  mem[  r+ list_offset].hh.b0  ;   mem[  r+ list_offset].hh.b0  :=shrinking;
  if  mem[ r+4].int  =0 then    mem[  r+glue_offset].gr :=0.0 
  else if (  mem[  r+ list_offset].hh.b1  =normal)and( mem[ r+height_offset].int  -t> mem[ r+4].int  ) then
       mem[  r+glue_offset].gr :=1.0 
  else  mem[ r+glue_offset].gr := (  mem[  r+height_offset].int  - t)/  mem[  r+4].int   ;
  end;
 mem[ r+height_offset].int  :=w;  mem[ r].hh.b0 :=vlist_node;
end

;
 mem[ r+4].int  :=0;
if u<>mem_top-4  then {append blank boxes to account for spanned nodes}
  begin  mem[ u].hh.rh := mem[ r].hh.rh ;  mem[ r].hh.rh := mem[ mem_top-4 ].hh.rh ; r:=u;
  end

;
r:= mem[  mem[  r].hh.rh ].hh.rh ; s:= mem[  mem[  s].hh.rh ].hh.rh ;
until r=-{0xfffffff=}268435455  ;
end


    else if  mem[ q].hh.b0 =rule_node then
      
{ Make the running dimensions in rule |q| extend to the boundaries of the alignment }
begin if  (  mem[  q+width_offset].int  =-{010000000000=}1073741824 )  then  mem[ q+width_offset].int  := mem[ p+width_offset].int  ;
if  (  mem[  q+height_offset].int  =-{010000000000=}1073741824 )  then  mem[ q+height_offset].int  := mem[ p+height_offset].int  ;
if  (  mem[  q+depth_offset].int  =-{010000000000=}1073741824 )  then  mem[ q+depth_offset].int  := mem[ p+depth_offset].int  ;
if o<>0 then
  begin r:= mem[ q].hh.rh ;  mem[ q].hh.rh :=-{0xfffffff=}268435455  ; q:=hpack(q,0,additional );
   mem[ q+4].int  :=o;  mem[ q].hh.rh :=r;  mem[ s].hh.rh :=q;
  end;
end

;
  s:=q; q:= mem[ q].hh.rh ;
  end

;
flush_node_list(p); pop_alignment;

{ Insert the \(c)current list into its environment }
aux_save:=cur_list.aux_field ; p:= mem[ cur_list.head_field ].hh.rh ; q:=cur_list.tail_field ; pop_nest;
if cur_list.mode_field =mmode then 
{ Finish an alignment in a display }
begin do_assignments;
if cur_cmd<>math_shift then 
{ Pontificate about improper alignment in display }
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Missing $$ inserted"=} 1183); end ;
{ \xref[Missing [\$\$] inserted] }
 begin help_ptr:=2; help_line[1]:={"Displays can use special alignments (like \eqalignno)"=} 909; help_line[0]:={"only if nothing but the alignment itself is between $$'s."=} 910; end ;
back_error;
end


else 
{ Check that another \.\$ follows }
begin get_x_token;
if cur_cmd<>math_shift then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Display math should end with $$"=} 1179); end ;
{ \xref[Display math...with \$\$] }
   begin help_ptr:=2; help_line[1]:={"The `$' that I just saw supposedly matches a previous `$$'."=} 1180; help_line[0]:={"So I shall assume that you typed `$$' both times."=} 1181; end ;
  back_error;
  end;
end

;
pop_nest;
begin  mem[ cur_list.tail_field ].hh.rh := new_penalty( eqtb[int_base+ pre_display_penalty_code].int  ); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
begin  mem[ cur_list.tail_field ].hh.rh := new_param_glue( above_display_skip_code); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
 mem[ cur_list.tail_field ].hh.rh :=p;
if p<>-{0xfffffff=}268435455   then cur_list.tail_field :=q;
begin  mem[ cur_list.tail_field ].hh.rh := new_penalty( eqtb[int_base+ post_display_penalty_code].int  ); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
begin  mem[ cur_list.tail_field ].hh.rh := new_param_glue( below_display_skip_code); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
cur_list.aux_field .int  :=aux_save.int ; resume_after_display;
end


else  begin cur_list.aux_field :=aux_save;  mem[ cur_list.tail_field ].hh.rh :=p;
  if p<>-{0xfffffff=}268435455   then cur_list.tail_field :=q;
  if cur_list.mode_field =vmode then build_page;
  end

;
end;

{ \4 }
{ Declare the procedure called |align_peek| }
procedure align_peek;
label restart;
begin restart: align_state:=1000000; 
{ Get the next non-blank non-call token }
repeat get_x_token;
until cur_cmd<>spacer

;
if cur_cmd=no_align then
  begin scan_left_brace; new_save_level(no_align_group);
  if cur_list.mode_field =-vmode then normal_paragraph;
  end
else if cur_cmd=right_brace then fin_align
else if (cur_cmd=car_ret)and(cur_chr=cr_cr_code) then
  goto restart {ignore \.[\\crcr]}
else  begin init_row; {start a new row}
  init_col; {start a new column and replace what we peeked at}
  end;
end;





{ 813. \[38] Breaking paragraphs into lines }

{tangle:pos tex.web:15997:39: }

{ We come now to what is probably the most interesting algorithm of \TeX:
the mechanism for choosing the ``best possible'' breakpoints that yield
the individual lines of a paragraph. \TeX's line-breaking algorithm takes
a given horizontal list and converts it to a sequence of boxes that are
appended to the current vertical list. In the course of doing this, it
creates a special data structure containing three kinds of records that are
not used elsewhere in \TeX. Such nodes are created while a paragraph is
being processed, and they are destroyed afterwards; thus, the other parts
of \TeX\ do not need to know anything about how line-breaking is done.

The method used here is based on an approach devised by Michael F. Plass and
\xref[Plass, Michael Frederick]
\xref[Knuth, Donald Ervin]
the author in 1977, subsequently generalized and improved by the same two
people in 1980. A detailed discussion appears in [\sl Software---Practice
and Experience \bf11] (1981), 1119--1184, where it is shown that the
line-breaking problem can be regarded as a special case of the problem of
computing the shortest path in an acyclic network. The cited paper includes
numerous examples and describes the history of line breaking as it has been
practiced by printers through the ages. The present implementation adds two
new ideas to the algorithm of 1980: Memory space requirements are considerably
reduced by using smaller records for inactive nodes than for active ones,
and arithmetic overflow is avoided by using ``delta distances'' instead of
keeping track of the total distance from the beginning of the paragraph to the
current point. }

{ 815. }

{tangle:pos tex.web:16048:1: }

{ Since |line_break| is a rather lengthy procedure---sort of a small world unto
itself---we must build it up little by little, somewhat more cautiously
than we have done with the simpler procedures of \TeX. Here is the
general outline. }{ \4 }
{ Declare subprocedures for |line_break| }
function finite_shrink( p:halfword ):halfword ; {recovers from infinite shrinkage}
var q:halfword ; {new glue specification}
begin if no_shrink_error_yet then
  begin no_shrink_error_yet:=false;
   ifdef('STAT')  if eqtb[int_base+ tracing_paragraphs_code].int  >0 then end_diagnostic(true); endif('STAT')  
  begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Infinite glue shrinkage found in a paragraph"=} 931); end ;
{ \xref[Infinite glue shrinkage...] }
   begin help_ptr:=5; help_line[4]:={"The paragraph just ended includes some glue that has"=} 932; help_line[3]:={"infinite shrinkability, e.g., `\hskip 0pt minus 1fil'."=} 933; help_line[2]:={"Such glue doesn't belong there---it allows a paragraph"=} 934; help_line[1]:={"of any length to fit on one line. But it's safe to proceed,"=} 935; help_line[0]:={"since the offensive shrinkability has been made finite."=} 936; end ;
  error;
   ifdef('STAT')  if eqtb[int_base+ tracing_paragraphs_code].int  >0 then begin_diagnostic; endif('STAT')  
  end;
q:=new_spec(p);   mem[ q].hh.b1 :=normal;
delete_glue_ref(p); finite_shrink:=q;
end;


procedure try_break( pi:integer; break_type:small_number);
label exit,done,done1,continue,deactivate;
var r:halfword ; {runs through the active list}
 prev_r:halfword ; {stays a step behind |r|}
 old_l:halfword; {maximum line number in current equivalence class of lines}
 no_break_yet:boolean; {have we found a feasible break at |cur_p|?}

{ Other local variables for |try_break| }
 prev_prev_r:halfword ; {a step behind |prev_r|, if |type(prev_r)=delta_node|}
 s:halfword ; {runs through nodes ahead of |cur_p|}
 q:halfword ; {points to a new node being created}
 v:halfword ; {points to a glue specification or a node ahead of |cur_p|}
 t:integer; {node count, if |cur_p| is a discretionary node}
 f:internal_font_number; {used in character width calculation}
 l:halfword; {line number of current active node}
 node_r_stays_active:boolean; {should node |r| remain in the active list?}
 line_width:scaled; {the current line will be justified to this width}
 fit_class:very_loose_fit..tight_fit; {possible fitness class of test line}
 b:halfword; {badness of test line}
 d:integer; {demerits of test line}
 artificial_demerits:boolean; {has |d| been forced to zero?}
 shortfall:scaled; {used in badness calculations}

 
begin 
{ Make sure that |pi| is in the proper range }
if abs(pi)>=inf_penalty then
  if pi>0 then  goto exit  {this breakpoint is inhibited by infinite penalty}
  else pi:=eject_penalty {this breakpoint will be forced}

;
no_break_yet:=true; prev_r:=mem_top-7 ; old_l:=0;
 cur_active_width[ 1]:=active_width[ 1] ; cur_active_width[ 2]:=active_width[ 2] ; cur_active_width[ 3]:=active_width[ 3] ; cur_active_width[ 4]:=active_width[ 4] ; cur_active_width[ 5]:=active_width[ 5] ; cur_active_width[ 6]:=active_width[ 6]  ;
 while true do    begin continue: r:= mem[ prev_r].hh.rh ;
  
{ If node |r| is of type |delta_node|, update |cur_active_width|, set |prev_r| and |prev_prev_r|, then |goto continue| }
{ \xref[inner loop] }
if  mem[ r].hh.b0 =delta_node then
  begin   cur_active_width[ 1]:=cur_active_width[ 1]+mem[r+ 1].int  ;  cur_active_width[ 2]:=cur_active_width[ 2]+mem[r+ 2].int  ;  cur_active_width[ 3]:=cur_active_width[ 3]+mem[r+ 3].int  ;  cur_active_width[ 4]:=cur_active_width[ 4]+mem[r+ 4].int  ;  cur_active_width[ 5]:=cur_active_width[ 5]+mem[r+ 5].int  ;  cur_active_width[ 6]:=cur_active_width[ 6]+mem[r+ 6].int   ;
  prev_prev_r:=prev_r; prev_r:=r; goto continue;
  end

;
  
{ If a line number class has ended, create new active nodes for the best feasible breaks in that class; then |return| if |r=last_active|, otherwise compute the new |line_width| }
begin l:=  mem[  r+ 1].hh.lh  ;
if l>old_l then
  begin {now we are no longer in the inner loop}
  if (minimum_demerits<{07777777777=}1073741823 )and 
      ((old_l<>easy_line)or(r=mem_top-7  )) then
    
{ Create new active nodes for the best feasible breaks just found }
begin if no_break_yet then 
{ Compute the values of |break_width| }
begin no_break_yet:=false;  break_width[ 1]:=background[ 1] ; break_width[ 2]:=background[ 2] ; break_width[ 3]:=background[ 3] ; break_width[ 4]:=background[ 4] ; break_width[ 5]:=background[ 5] ; break_width[ 6]:=background[ 6]  ;
s:=cur_p;
if break_type>unhyphenated then if cur_p<>-{0xfffffff=}268435455   then
  
{ Compute the discretionary |break_width| values }
begin t:= mem[ cur_p].hh.b1 ; v:=cur_p; s:=  mem[  cur_p+ 1].hh.rh  ;
while t>0 do
  begin decr(t); v:= mem[ v].hh.rh ;
  
{ Subtract the width of node |v| from |break_width| }
if  ( v>=hi_mem_min)  then
  begin f:=  mem[ v].hh.b0 ;
  break_width[1]:=break_width[1]-font_info[width_base[ f]+  font_info[char_base[  f]+effective_char(true,  f,     mem[   v].hh.b1 )].qqqq .b0].int  ;
  end
else  case  mem[ v].hh.b0  of
  ligature_node: begin f:=  mem[   v+1 ].hh.b0 ;

    break_width[1]:= break_width[1]-
      font_info[width_base[ f]+  font_info[char_base[  f]+effective_char(true,  f,     mem[       v+1 ].hh.b1 )].qqqq .b0].int  ;
    end;
  hlist_node,vlist_node,rule_node,kern_node:
    break_width[1]:=break_width[1]- mem[ v+width_offset].int  ;
   else  confusion({"disc1"=}937)
{ \xref[this can't happen disc1][\quad disc1] }
   end 

;
  end;
while s<>-{0xfffffff=}268435455   do
  begin 
{ Add the width of node |s| to |break_width| }
if  ( s>=hi_mem_min)  then
  begin f:=  mem[ s].hh.b0 ;
  break_width[1]:= break_width[1]+font_info[width_base[ f]+  font_info[char_base[  f]+effective_char(true,  f,     mem[   s].hh.b1 )].qqqq .b0].int  ;
  end
else  case  mem[ s].hh.b0  of
  ligature_node: begin f:=  mem[   s+1 ].hh.b0 ;
    break_width[1]:=break_width[1]+
      font_info[width_base[ f]+  font_info[char_base[  f]+effective_char(true,  f,     mem[       s+1 ].hh.b1 )].qqqq .b0].int  ;
    end;
  hlist_node,vlist_node,rule_node,kern_node:
    break_width[1]:=break_width[1]+ mem[ s+width_offset].int  ;
   else  confusion({"disc2"=}938)
{ \xref[this can't happen disc2][\quad disc2] }
   end 

;
  s:= mem[ s].hh.rh ;
  end;
break_width[1]:=break_width[1]+disc_width;
if   mem[  cur_p+ 1].hh.rh  =-{0xfffffff=}268435455   then s:= mem[ v].hh.rh ;
          {nodes may be discardable after the break}
end

;
while s<>-{0xfffffff=}268435455   do
  begin if  ( s>=hi_mem_min)  then goto done;
  case  mem[ s].hh.b0  of
  glue_node:
{ Subtract glue from |break_width| }
begin v:=  mem[  s+ 1].hh.lh  ; break_width[1]:=break_width[1]- mem[ v+width_offset].int  ;
break_width[2+  mem[ v].hh.b0 ]:=break_width[2+  mem[ v].hh.b0 ]- mem[ v+2].int  ;
break_width[6]:=break_width[6]- mem[ v+3].int  ;
end

;
  penalty_node:  ;
  math_node: break_width[1]:=break_width[1]- mem[ s+width_offset].int  ;
  kern_node: if  mem[ s].hh.b1 <>explicit then goto done
    else break_width[1]:=break_width[1]- mem[ s+width_offset].int  ;
   else  goto done
   end ;

  s:= mem[ s].hh.rh ;
  end;
done: end

;

{ Insert a delta node to prepare for breaks at |cur_p| }
if  mem[ prev_r].hh.b0 =delta_node then {modify an existing delta node}
  begin   mem[prev_r+ 1].int := { \hskip10pt }mem[prev_r+ 1].int  -cur_active_width[ 1]+break_width[ 1] ;  mem[prev_r+ 2].int := { \hskip10pt }mem[prev_r+ 2].int  -cur_active_width[ 2]+break_width[ 2] ;  mem[prev_r+ 3].int := { \hskip10pt }mem[prev_r+ 3].int  -cur_active_width[ 3]+break_width[ 3] ;  mem[prev_r+ 4].int := { \hskip10pt }mem[prev_r+ 4].int  -cur_active_width[ 4]+break_width[ 4] ;  mem[prev_r+ 5].int := { \hskip10pt }mem[prev_r+ 5].int  -cur_active_width[ 5]+break_width[ 5] ;  mem[prev_r+ 6].int := { \hskip10pt }mem[prev_r+ 6].int  -cur_active_width[ 6]+break_width[ 6]  ;
  end
else if prev_r=mem_top-7  then {no delta node needed at the beginning}
  begin  active_width[ 1]:=break_width[ 1] ; active_width[ 2]:=break_width[ 2] ; active_width[ 3]:=break_width[ 3] ; active_width[ 4]:=break_width[ 4] ; active_width[ 5]:=break_width[ 5] ; active_width[ 6]:=break_width[ 6]  ;
  end
else  begin q:=get_node(delta_node_size);  mem[ q].hh.rh :=r;  mem[ q].hh.b0 :=delta_node;

   mem[ q].hh.b1 :=0; {the |subtype| is not used}
    mem[q+ 1].int :=break_width[ 1]-cur_active_width[ 1] ;  mem[q+ 2].int :=break_width[ 2]-cur_active_width[ 2] ;  mem[q+ 3].int :=break_width[ 3]-cur_active_width[ 3] ;  mem[q+ 4].int :=break_width[ 4]-cur_active_width[ 4] ;  mem[q+ 5].int :=break_width[ 5]-cur_active_width[ 5] ;  mem[q+ 6].int :=break_width[ 6]-cur_active_width[ 6]  ;
   mem[ prev_r].hh.rh :=q; prev_prev_r:=prev_r; prev_r:=q;
  end

;
if abs(eqtb[int_base+ adj_demerits_code].int  )>={07777777777=}1073741823 -minimum_demerits then
  minimum_demerits:={07777777777=}1073741823 -1
else minimum_demerits:=minimum_demerits+abs(eqtb[int_base+ adj_demerits_code].int  );
for fit_class:=very_loose_fit to tight_fit do
  begin if minimal_demerits[fit_class]<=minimum_demerits then
    
{ Insert a new active node from |best_place[fit_class]| to |cur_p| }
begin q:=get_node(passive_node_size);
 mem[ q].hh.rh :=passive; passive:=q;   mem[  q+ 1].hh.rh  :=cur_p;
 ifdef('STAT')  incr(pass_number);  mem[ q].hh.lh :=pass_number; endif('STAT')  

  mem[  q+ 1].hh.lh  :=best_place[fit_class];

q:=get_node(active_node_size);   mem[  q+ 1].hh.rh  :=passive;
  mem[  q+ 1].hh.lh  :=best_pl_line[fit_class]+1;
 mem[ q].hh.b1 :=fit_class;  mem[ q].hh.b0 :=break_type;
mem[ q+2].int :=minimal_demerits[fit_class];
 mem[ q].hh.rh :=r;  mem[ prev_r].hh.rh :=q; prev_r:=q;
 ifdef('STAT')  if eqtb[int_base+ tracing_paragraphs_code].int  >0 then
  
{ Print a symbolic description of the new break node }
begin print_nl({"@@"=}939); print_int( mem[ passive].hh.lh );
{ \xref[\AT!\AT!] }
print({": line "=}940); print_int(  mem[  q+ 1].hh.lh  -1);
print_char({"."=}46); print_int(fit_class);
if break_type=hyphenated then print_char({"-"=}45);
print({" t="=}941); print_int(mem[ q+2].int );
print({" -> @@"=}942);
if   mem[  passive+ 1].hh.lh  =-{0xfffffff=}268435455   then print_char({"0"=}48)
else print_int( mem[   mem[   passive+ 1].hh.lh  ].hh.lh );
end

;
endif('STAT')  

end

;
  minimal_demerits[fit_class]:={07777777777=}1073741823 ;
  end;
minimum_demerits:={07777777777=}1073741823 ;

{ Insert a delta node to prepare for the next active node }
if r<>mem_top-7   then
  begin q:=get_node(delta_node_size);  mem[ q].hh.rh :=r;  mem[ q].hh.b0 :=delta_node;

   mem[ q].hh.b1 :=0; {the |subtype| is not used}
    mem[q+ 1].int := cur_active_width[ 1]-break_width[ 1] ;  mem[q+ 2].int := cur_active_width[ 2]-break_width[ 2] ;  mem[q+ 3].int := cur_active_width[ 3]-break_width[ 3] ;  mem[q+ 4].int := cur_active_width[ 4]-break_width[ 4] ;  mem[q+ 5].int := cur_active_width[ 5]-break_width[ 5] ;  mem[q+ 6].int := cur_active_width[ 6]-break_width[ 6]  ;
   mem[ prev_r].hh.rh :=q; prev_prev_r:=prev_r; prev_r:=q;
  end

;
end

;
  if r=mem_top-7   then  goto exit ;
  
{ Compute the new line width }
if l>easy_line then
  begin line_width:=second_width; old_l:={0xfffffff=}268435455 -1;
  end
else  begin old_l:=l;
  if l>last_special_line then line_width:=second_width
  else if  eqtb[  par_shape_loc].hh.rh   =-{0xfffffff=}268435455   then line_width:=first_width
  else line_width:=mem[ eqtb[  par_shape_loc].hh.rh   +2*l ].int ;
  end

;
  end;
end

;
  
{ Consider the demerits for a line from |r| to |cur_p|; deactivate node |r| if it should no longer be active; then |goto continue| if a line from |r| to |cur_p| is infeasible, otherwise record a new feasible break }
begin artificial_demerits:=false;

{ \xref[inner loop] }
shortfall:=line_width-cur_active_width[1]; {we're this much too short}
if shortfall>0 then
  
{ Set the value of |b| to the badness for stretching the line, and compute the corresponding |fit_class| }
if (cur_active_width[3]<>0)or(cur_active_width[4]<>0)or 
  (cur_active_width[5]<>0) then
  begin b:=0; fit_class:=decent_fit; {infinite stretch}
  end
else  begin if shortfall>7230584 then if cur_active_width[2]<1663497 then
    begin b:=inf_bad; fit_class:=very_loose_fit; goto done1;
    end;
  b:=badness(shortfall,cur_active_width[2]);
  if b>12 then
    if b>99 then fit_class:=very_loose_fit
    else fit_class:=loose_fit
  else fit_class:=decent_fit;
  done1:
  end


else 
{ Set the value of |b| to the badness for shrinking the line, and compute the corresponding |fit_class| }
begin if -shortfall>cur_active_width[6] then b:=inf_bad+1
else b:=badness(-shortfall,cur_active_width[6]);
if b>12 then fit_class:=tight_fit else fit_class:=decent_fit;
end

;
if (b>inf_bad)or(pi=eject_penalty) then
  
{ Prepare to deactivate node~|r|, and |goto deactivate| unless there is a reason to consider lines of text from |r| to |cur_p| }
begin if final_pass and (minimum_demerits={07777777777=}1073741823 ) and 
   ( mem[ r].hh.rh =mem_top-7  ) and
   (prev_r=mem_top-7 ) then
  artificial_demerits:=true {set demerits zero, this break is forced}
else if b>threshold then goto deactivate;
node_r_stays_active:=false;
end


else  begin prev_r:=r;
  if b>threshold then goto continue;
  node_r_stays_active:=true;
  end;

{ Record a new feasible break }
if artificial_demerits then d:=0
else 
{ Compute the demerits, |d|, from |r| to |cur_p| }
begin d:=eqtb[int_base+ line_penalty_code].int  +b;
if abs(d)>=10000 then d:=100000000 else d:=d*d;
if pi<>0 then
  if pi>0 then d:=d+pi*pi
  else if pi>eject_penalty then d:=d-pi*pi;
if (break_type=hyphenated)and( mem[ r].hh.b0 =hyphenated) then
  if cur_p<>-{0xfffffff=}268435455   then d:=d+eqtb[int_base+ double_hyphen_demerits_code].int  
  else d:=d+eqtb[int_base+ final_hyphen_demerits_code].int  ;
if abs(fit_class- mem[ r].hh.b1 )>1 then d:=d+eqtb[int_base+ adj_demerits_code].int  ;
end

;
 ifdef('STAT')  if eqtb[int_base+ tracing_paragraphs_code].int  >0 then
  
{ Print a symbolic description of this feasible break }
begin if printed_node<>cur_p then
  
{ Print the list between |printed_node| and |cur_p|, then set |printed_node:=cur_p| }
begin print_nl({""=}335);
if cur_p=-{0xfffffff=}268435455   then short_display( mem[ printed_node].hh.rh )
else  begin save_link:= mem[ cur_p].hh.rh ;
   mem[ cur_p].hh.rh :=-{0xfffffff=}268435455  ; print_nl({""=}335); short_display( mem[ printed_node].hh.rh );
   mem[ cur_p].hh.rh :=save_link;
  end;
printed_node:=cur_p;
end

;
print_nl({"@"=}64);
{ \xref[\AT!] }
if cur_p=-{0xfffffff=}268435455   then print_esc({"par"=}604)
else if  mem[ cur_p].hh.b0 <>glue_node then
  begin if  mem[ cur_p].hh.b0 =penalty_node then print_esc({"penalty"=}539)
  else if  mem[ cur_p].hh.b0 =disc_node then print_esc({"discretionary"=}346)
  else if  mem[ cur_p].hh.b0 =kern_node then print_esc({"kern"=}337)
  else print_esc({"math"=}340);
  end;
print({" via @@"=}943);
if   mem[  r+ 1].hh.rh  =-{0xfffffff=}268435455   then print_char({"0"=}48)
else print_int( mem[   mem[   r+ 1].hh.rh  ].hh.lh );
print({" b="=}944);
if b>inf_bad then print_char({"*"=}42) else print_int(b);
{ \xref[*\relax] }
print({" p="=}945); print_int(pi); print({" d="=}946);
if artificial_demerits then print_char({"*"=}42) else print_int(d);
end

;
endif('STAT')  

d:=d+mem[ r+2].int ; {this is the minimum total demerits
  from the beginning to |cur_p| via |r|}
if d<=minimal_demerits[fit_class] then
  begin minimal_demerits[fit_class]:=d;
  best_place[fit_class]:=  mem[  r+ 1].hh.rh  ; best_pl_line[fit_class]:=l;
  if d<minimum_demerits then minimum_demerits:=d;
  end

;
if node_r_stays_active then goto continue; {|prev_r| has been set to |r|}
deactivate: 
{ Deactivate node |r| }
 mem[ prev_r].hh.rh := mem[ r].hh.rh ; free_node(r,active_node_size);
if prev_r=mem_top-7  then 
{ Update the active widths, since the first active node has been deleted }
begin r:= mem[ mem_top-7 ].hh.rh ;
if  mem[ r].hh.b0 =delta_node then
  begin  active_width[ 1]:=active_width[ 1]+mem[r+ 1].int  ; active_width[ 2]:=active_width[ 2]+mem[r+ 2].int  ; active_width[ 3]:=active_width[ 3]+mem[r+ 3].int  ; active_width[ 4]:=active_width[ 4]+mem[r+ 4].int  ; active_width[ 5]:=active_width[ 5]+mem[r+ 5].int  ; active_width[ 6]:=active_width[ 6]+mem[r+ 6].int   ;
   cur_active_width[ 1]:=active_width[ 1] ; cur_active_width[ 2]:=active_width[ 2] ; cur_active_width[ 3]:=active_width[ 3] ; cur_active_width[ 4]:=active_width[ 4] ; cur_active_width[ 5]:=active_width[ 5] ; cur_active_width[ 6]:=active_width[ 6]  ;
   mem[ mem_top-7 ].hh.rh := mem[ r].hh.rh ; free_node(r,delta_node_size);
  end;
end


else if  mem[ prev_r].hh.b0 =delta_node then
  begin r:= mem[ prev_r].hh.rh ;
  if r=mem_top-7   then
    begin   cur_active_width[ 1]:=cur_active_width[ 1]- mem[prev_r+ 1].int  ;  cur_active_width[ 2]:=cur_active_width[ 2]- mem[prev_r+ 2].int  ;  cur_active_width[ 3]:=cur_active_width[ 3]- mem[prev_r+ 3].int  ;  cur_active_width[ 4]:=cur_active_width[ 4]- mem[prev_r+ 4].int  ;  cur_active_width[ 5]:=cur_active_width[ 5]- mem[prev_r+ 5].int  ;  cur_active_width[ 6]:=cur_active_width[ 6]- mem[prev_r+ 6].int   ;
     mem[ prev_prev_r].hh.rh :=mem_top-7  ;
    free_node(prev_r,delta_node_size); prev_r:=prev_prev_r;
    end
  else if  mem[ r].hh.b0 =delta_node then
    begin   cur_active_width[ 1]:=cur_active_width[ 1]+mem[r+ 1].int  ;  cur_active_width[ 2]:=cur_active_width[ 2]+mem[r+ 2].int  ;  cur_active_width[ 3]:=cur_active_width[ 3]+mem[r+ 3].int  ;  cur_active_width[ 4]:=cur_active_width[ 4]+mem[r+ 4].int  ;  cur_active_width[ 5]:=cur_active_width[ 5]+mem[r+ 5].int  ;  cur_active_width[ 6]:=cur_active_width[ 6]+mem[r+ 6].int   ;
      mem[prev_r+ 1].int :=mem[prev_r+ 1].int +mem[r+ 1].int  ;  mem[prev_r+ 2].int :=mem[prev_r+ 2].int +mem[r+ 2].int  ;  mem[prev_r+ 3].int :=mem[prev_r+ 3].int +mem[r+ 3].int  ;  mem[prev_r+ 4].int :=mem[prev_r+ 4].int +mem[r+ 4].int  ;  mem[prev_r+ 5].int :=mem[prev_r+ 5].int +mem[r+ 5].int  ;  mem[prev_r+ 6].int :=mem[prev_r+ 6].int +mem[r+ 6].int   ;
     mem[ prev_r].hh.rh := mem[ r].hh.rh ; free_node(r,delta_node_size);
    end;
  end

;
end

;
  end;
exit:  ifdef('STAT')  
{ Update the value of |printed_node| for symbolic displays }
if cur_p=printed_node then if cur_p<>-{0xfffffff=}268435455   then if  mem[ cur_p].hh.b0 =disc_node then
  begin t:= mem[ cur_p].hh.b1 ;
  while t>0 do
    begin decr(t); printed_node:= mem[ printed_node].hh.rh ;
    end;
  end

 endif('STAT')  
end;


procedure post_line_break( final_widow_penalty:integer);
label done,done1;
var q, r, s:halfword ; {temporary registers for list manipulation}
 disc_break:boolean; {was the current break at a discretionary node?}
 post_disc_break:boolean; {and did it have a nonempty post-break part?}
 cur_width:scaled; {width of line number |cur_line|}
 cur_indent:scaled; {left margin of line number |cur_line|}
 t:quarterword; {used for replacement counts in discretionary nodes}
 pen:integer; {use when calculating penalties between lines}
 cur_line: halfword; {the current line number being justified}
begin 
{ Reverse the links of the relevant passive nodes, setting |cur_p| to the first breakpoint }
q:=  mem[  best_bet+ 1].hh.rh  ; cur_p:=-{0xfffffff=}268435455  ;
repeat r:=q; q:=  mem[  q+ 1].hh.lh  ;   mem[  r+ 1].hh.lh  :=cur_p; cur_p:=r;
until q=-{0xfffffff=}268435455  

;
cur_line:=cur_list.pg_field +1;
repeat 
{ Justify the line ending at breakpoint |cur_p|, and append it to the current vertical list, together with associated penalties and other insertions }

{ Modify the end of the line to reflect the nature of the break and to include \.[\\rightskip]; also set the proper value of |disc_break| }
q:=  mem[  cur_p+ 1].hh.rh  ; disc_break:=false; post_disc_break:=false;
if q<>-{0xfffffff=}268435455   then {|q| cannot be a |char_node|}
  if  mem[ q].hh.b0 =glue_node then
    begin delete_glue_ref(  mem[  q+ 1].hh.lh  );
      mem[  q+ 1].hh.lh  := eqtb[  glue_base+   right_skip_code].hh.rh    ;
     mem[ q].hh.b1 :=right_skip_code+1; incr(  mem[    eqtb[  glue_base+   right_skip_code].hh.rh    ].hh.rh  ) ;
    goto done;
    end
  else  begin if  mem[ q].hh.b0 =disc_node then
      
{ Change discretionary to compulsory and set |disc_break:=true| }
begin t:= mem[ q].hh.b1 ;

{ Destroy the |t| nodes following |q|, and make |r| point to the following node }
if t=0 then r:= mem[ q].hh.rh 
else  begin r:=q;
  while t>1 do
    begin r:= mem[ r].hh.rh ; decr(t);
    end;
  s:= mem[ r].hh.rh ;
  r:= mem[ s].hh.rh ;  mem[ s].hh.rh :=-{0xfffffff=}268435455  ;
  flush_node_list( mem[ q].hh.rh );  mem[ q].hh.b1 :=0;
  end

;
if   mem[  q+ 1].hh.rh  <>-{0xfffffff=}268435455   then 
{ Transplant the post-break list }
begin s:=  mem[  q+ 1].hh.rh  ;
while  mem[ s].hh.rh <>-{0xfffffff=}268435455   do s:= mem[ s].hh.rh ;
 mem[ s].hh.rh :=r; r:=  mem[  q+ 1].hh.rh  ;   mem[  q+ 1].hh.rh  :=-{0xfffffff=}268435455  ; post_disc_break:=true;
end

;
if   mem[  q+ 1].hh.lh  <>-{0xfffffff=}268435455   then 
{ Transplant the pre-break list }
begin s:=  mem[  q+ 1].hh.lh  ;  mem[ q].hh.rh :=s;
while  mem[ s].hh.rh <>-{0xfffffff=}268435455   do s:= mem[ s].hh.rh ;
  mem[  q+ 1].hh.lh  :=-{0xfffffff=}268435455  ; q:=s;
end

;
 mem[ q].hh.rh :=r; disc_break:=true;
end


    else if ( mem[ q].hh.b0 =math_node)or( mem[ q].hh.b0 =kern_node) then  mem[ q+width_offset].int  :=0;
    end
else  begin q:=mem_top-3 ;
  while  mem[ q].hh.rh <>-{0xfffffff=}268435455   do q:= mem[ q].hh.rh ;
  end;

{ Put the \(r)\.[\\rightskip] glue after node |q| }
r:=new_param_glue(right_skip_code);  mem[ r].hh.rh := mem[ q].hh.rh ;  mem[ q].hh.rh :=r; q:=r

;
done:

;

{ Put the \(l)\.[\\leftskip] glue at the left and detach this line }
r:= mem[ q].hh.rh ;  mem[ q].hh.rh :=-{0xfffffff=}268435455  ; q:= mem[ mem_top-3 ].hh.rh ;  mem[ mem_top-3 ].hh.rh :=r;
if  eqtb[  glue_base+   left_skip_code].hh.rh    <>mem_bot  then
  begin r:=new_param_glue(left_skip_code);
   mem[ r].hh.rh :=q; q:=r;
  end

;

{ Call the packaging subroutine, setting |just_box| to the justified box }
if cur_line>last_special_line then
  begin cur_width:=second_width; cur_indent:=second_indent;
  end
else if  eqtb[  par_shape_loc].hh.rh   =-{0xfffffff=}268435455   then
  begin cur_width:=first_width; cur_indent:=first_indent;
  end
else  begin cur_width:=mem[ eqtb[  par_shape_loc].hh.rh   +2*cur_line].int ;
  cur_indent:=mem[ eqtb[  par_shape_loc].hh.rh   +2*cur_line-1].int ;
  end;
adjust_tail:=mem_top-5 ; just_box:=hpack(q,cur_width,exactly);
 mem[ just_box+4].int  :=cur_indent

;

{ Append the new box to the current vertical list, followed by the list of special nodes taken out of the box by the packager }
append_to_vlist(just_box);
if mem_top-5 <>adjust_tail then
  begin  mem[ cur_list.tail_field ].hh.rh := mem[ mem_top-5 ].hh.rh ; cur_list.tail_field :=adjust_tail;
   end;
adjust_tail:=-{0xfffffff=}268435455  

;

{ Append a penalty node, if a nonzero penalty is appropriate }
if cur_line+1<>best_line then
  begin pen:=eqtb[int_base+ inter_line_penalty_code].int  ;
  if cur_line=cur_list.pg_field +1 then pen:=pen+eqtb[int_base+ club_penalty_code].int  ;
  if cur_line+2=best_line then pen:=pen+final_widow_penalty;
  if disc_break then pen:=pen+eqtb[int_base+ broken_penalty_code].int  ;
  if pen<>0 then
    begin r:=new_penalty(pen);
     mem[ cur_list.tail_field ].hh.rh :=r; cur_list.tail_field :=r;
    end;
  end



;
incr(cur_line); cur_p:=  mem[  cur_p+ 1].hh.lh  ;
if cur_p<>-{0xfffffff=}268435455   then if not post_disc_break then
  
{ Prune unwanted nodes at the beginning of the next line }
begin r:=mem_top-3 ;
 while true do    begin q:= mem[ r].hh.rh ;
  if q=  mem[  cur_p+ 1].hh.rh   then goto done1;
    {|cur_break(cur_p)| is the next breakpoint}
  {now |q| cannot be |null|}
  if  ( q>=hi_mem_min)  then goto done1;
  if ( mem[  q].hh.b0 <math_node)  then goto done1;
  if  mem[ q].hh.b0 =kern_node then if  mem[ q].hh.b1 <>explicit then goto done1;
  r:=q; {now |type(q)=glue_node|, |kern_node|, |math_node|, or |penalty_node|}
  end;
done1: if r<>mem_top-3  then
  begin  mem[ r].hh.rh :=-{0xfffffff=}268435455  ; flush_node_list( mem[ mem_top-3 ].hh.rh );
   mem[ mem_top-3 ].hh.rh :=q;
  end;
end

;
until cur_p=-{0xfffffff=}268435455  ;
if (cur_line<>best_line)or( mem[ mem_top-3 ].hh.rh <>-{0xfffffff=}268435455  ) then
  confusion({"line breaking"=}953);
{ \xref[this can't happen line breaking][\quad line breaking] }
cur_list.pg_field :=best_line-1;
end;


{ \4 }
{ Declare the function called |reconstitute| }
function reconstitute( j, n:small_number; bchar, hchar:halfword):
  small_number;
label continue,done;
var  p:halfword ; {temporary register for list manipulation}
 t:halfword ; {a node being appended to}
 q:four_quarters; {character information or a lig/kern instruction}
 cur_rh:halfword; {hyphen character for ligature testing}
 test_char:halfword; {hyphen or other character for ligature testing}
 w:scaled; {amount of kerning}
 k:font_index; {position of current lig/kern instruction}
begin hyphen_passed:=0; t:=mem_top-4 ; w:=0;  mem[ mem_top-4 ].hh.rh :=-{0xfffffff=}268435455  ;
 {at this point |ligature_present=lft_hit=rt_hit=false|}

{ Set up data structures with the cursor following position |j| }
cur_l:= hu[ j] ; cur_q:=t;
if j=0 then
  begin ligature_present:=init_lig; p:=init_list;
  if ligature_present then lft_hit:=init_lft;
  while p>-{0xfffffff=}268435455   do
    begin  begin  mem[ t].hh.rh :=get_avail; t:= mem[ t].hh.rh ;   mem[ t].hh.b0 :=hf;   mem[ t].hh.b1 :=   mem[  p].hh.b1 ; end ; p:= mem[ p].hh.rh ;
    end;
  end
else if cur_l< 256   then  begin  mem[ t].hh.rh :=get_avail; t:= mem[ t].hh.rh ;   mem[ t].hh.b0 :=hf;   mem[ t].hh.b1 := cur_l; end ;
lig_stack:=-{0xfffffff=}268435455  ; begin if j<n then cur_r:= hu[ j+ 1]  else cur_r:=bchar; if odd(hyf[j]) then cur_rh:=hchar else cur_rh:= 256  ; end 

;
continue:
{ If there's a ligature or kern at the cursor position, update the data structures, possibly advancing~|j|; continue until the cursor moves }
if cur_l= 256   then
  begin k:=bchar_label[hf];
  if k=non_address then goto done else q:=font_info[k].qqqq;
  end
else begin q:= font_info[char_base[ hf]+effective_char(true, hf,  cur_l)].qqqq ;
  if ((  q. b2 ) mod 4) <>lig_tag then goto done;
  k:=lig_kern_base[ hf]+ q.b3 ; q:=font_info[k].qqqq;
  if  q.b0 > 128   then
    begin k:=lig_kern_base[ hf]+256*  q.b2 +  q.b3 +32768-256*(128+min_quarterword)  ; q:=font_info[k].qqqq;
    end;
  end; {now |k| is the starting address of the lig/kern program}
if cur_rh< 256   then test_char:=cur_rh else test_char:=cur_r;
 while true do  begin if  q.b1 =test_char then if  q.b0 <= 128   then
    if cur_rh< 256   then
      begin hyphen_passed:=j; hchar:= 256  ; cur_rh:= 256  ;
      goto continue;
      end
    else begin if hchar< 256   then if odd(hyf[j]) then
        begin hyphen_passed:=j; hchar:= 256  ;
        end;
      if  q.b2 < 128   then
      
{ Carry out a ligature replacement, updating the cursor structure and possibly advancing~|j|; |goto continue| if the cursor doesn't advance, otherwise |goto done| }
begin if cur_l= 256   then lft_hit:=true;
if j=n then if lig_stack=-{0xfffffff=}268435455   then rt_hit:=true;
begin if interrupt<>0 then pause_for_instructions; end ; {allow a way out in case there's an infinite ligature loop}
case  q.b2  of
 1 , 5 :begin cur_l:= q.b3 ; {\.[=:\?], \.[=:\?>]}
  ligature_present:=true;
  end;
 2 , 6 :begin cur_r:= q.b3 ; {\.[\?=:], \.[\?=:>]}
  if lig_stack>-{0xfffffff=}268435455   then   mem[ lig_stack].hh.b1 :=cur_r
  else begin lig_stack:=new_lig_item(cur_r);
    if j=n then bchar:= 256  
    else begin p:=get_avail;  mem[    lig_stack+1 ].hh.rh  :=p;
        mem[ p].hh.b1 := hu[ j+ 1] ;   mem[ p].hh.b0 :=hf;
      end;
    end;
  end;
 3 :begin cur_r:= q.b3 ; {\.[\?=:\?]}
  p:=lig_stack; lig_stack:=new_lig_item(cur_r);  mem[ lig_stack].hh.rh :=p;
  end;
 7 , 11 :begin if ligature_present then begin p:=new_ligature(hf,cur_l, mem[ cur_q].hh.rh ); if lft_hit then begin  mem[ p].hh.b1 :=2; lft_hit:=false; end; if  false then if lig_stack=-{0xfffffff=}268435455   then begin incr( mem[ p].hh.b1 ); rt_hit:=false; end;  mem[ cur_q].hh.rh :=p; t:=p; ligature_present:=false; end ; {\.[\?=:\?>], \.[\?=:\?>>]}
  cur_q:=t; cur_l:= q.b3 ; ligature_present:=true;
  end;
 else  begin cur_l:= q.b3 ; ligature_present:=true; {\.[=:]}
  if lig_stack>-{0xfffffff=}268435455   then begin if  mem[    lig_stack+1 ].hh.rh  >-{0xfffffff=}268435455   then begin  mem[ t].hh.rh := mem[    lig_stack+1 ].hh.rh  ; t:= mem[ t].hh.rh ; incr(j); end; p:=lig_stack; lig_stack:= mem[ p].hh.rh ; free_node(p,small_node_size); if lig_stack=-{0xfffffff=}268435455   then begin if j<n then cur_r:= hu[ j+ 1]  else cur_r:=bchar; if odd(hyf[j]) then cur_rh:=hchar else cur_rh:= 256  ; end  else cur_r:=  mem[ lig_stack].hh.b1 ; end 
  else if j=n then goto done
  else begin  begin  mem[ t].hh.rh :=get_avail; t:= mem[ t].hh.rh ;   mem[ t].hh.b0 :=hf;   mem[ t].hh.b1 := cur_r; end ; incr(j); begin if j<n then cur_r:= hu[ j+ 1]  else cur_r:=bchar; if odd(hyf[j]) then cur_rh:=hchar else cur_rh:= 256  ; end ;
    end;
  end
 end ;
if  q.b2 > 4  then if  q.b2 <> 7  then goto done;
goto continue;
end

;
      w:=font_info[kern_base[ hf]+256*  q.b2 +  q.b3 ].int  ; goto done; {this kern will be inserted below}
     end;
  if  q.b0 >= 128   then
    if cur_rh= 256   then goto done
    else begin cur_rh:= 256  ; goto continue;
      end;
  k:=k+   q.b0  +1; q:=font_info[k].qqqq;
  end;
done:

;

{ Append a ligature and/or kern to the translation; |goto continue| if the stack of inserted ligatures is nonempty }
if ligature_present then begin p:=new_ligature(hf,cur_l, mem[ cur_q].hh.rh ); if lft_hit then begin  mem[ p].hh.b1 :=2; lft_hit:=false; end; if  rt_hit then if lig_stack=-{0xfffffff=}268435455   then begin incr( mem[ p].hh.b1 ); rt_hit:=false; end;  mem[ cur_q].hh.rh :=p; t:=p; ligature_present:=false; end ;
if w<>0 then
  begin  mem[ t].hh.rh :=new_kern(w); t:= mem[ t].hh.rh ; w:=0;
  end;
if lig_stack>-{0xfffffff=}268435455   then
  begin cur_q:=t; cur_l:=  mem[ lig_stack].hh.b1 ; ligature_present:=true;
  begin if  mem[    lig_stack+1 ].hh.rh  >-{0xfffffff=}268435455   then begin  mem[ t].hh.rh := mem[    lig_stack+1 ].hh.rh  ; t:= mem[ t].hh.rh ; incr(j); end; p:=lig_stack; lig_stack:= mem[ p].hh.rh ; free_node(p,small_node_size); if lig_stack=-{0xfffffff=}268435455   then begin if j<n then cur_r:= hu[ j+ 1]  else cur_r:=bchar; if odd(hyf[j]) then cur_rh:=hchar else cur_rh:= 256  ; end  else cur_r:=  mem[ lig_stack].hh.b1 ; end ; goto continue;
  end

;
reconstitute:=j;
end;


procedure hyphenate;
label common_ending,done,found,found1,found2,not_found,exit;
var 
{ Local variables for hyphenation }
 i, j, l:0..65; {indices into |hc| or |hu|}
 q, r, s:halfword ; {temporary registers for list manipulation}
 bchar:halfword; {boundary character of hyphenated word, or |non_char|}


 major_tail, minor_tail:halfword ; {the end of lists in the main and
  discretionary branches being reconstructed}
 c:ASCII_code; {character temporarily replaced by a hyphen}
 c_loc:0..63; {where that character came from}
 r_count:integer; {replacement count for discretionary}
 hyf_node:halfword ; {the hyphen, if it exists}


 z:trie_pointer; {an index into |trie|}
 v:integer; {an index into |hyf_distance|, etc.}


 h:hyph_pointer; {an index into |hyph_word| and |hyph_list|}
 k:str_number; {an index into |str_start|}
 u:pool_pointer; {an index into |str_pool|}

 
begin 
{ Find hyphen locations for the word in |hc|, or |return| }
for j:=0 to hn do hyf[j]:=0;

{ Look for the word |hc[1..hn]| in the exception table, and |goto found| (with |hyf| containing the hyphens) if an entry is found }
h:=hc[1]; incr(hn); hc[hn]:=cur_lang;
for j:=2 to hn do h:=(h+h+hc[j]) mod hyph_prime;
 while true do    begin 
{ If the string |hyph_word[h]| is less than \(hc)|hc[1..hn]|, |goto not_found|; but if the two strings are equal, set |hyf| to the hyphen positions and |goto found| }
{This is now a simple hash list, not an ordered one, so
the module title is no longer descriptive.}
k:=hyph_word[h]; if k=0 then goto not_found;
if (str_start[ k+1]-str_start[ k]) =hn then
  begin j:=1; u:=str_start[k];
  repeat
  if   str_pool[ u] <>hc[j] then goto done;
  incr(j); incr(u);
  until j>hn;
  
{ Insert hyphens as specified in |hyph_list[h]| }
s:=hyph_list[h];
while s<>-{0xfffffff=}268435455   do
  begin hyf[ mem[ s].hh.lh ]:=1; s:= mem[ s].hh.rh ;
  end

;
  decr(hn); goto found;
  end;
done:

;
  h:=hyph_link[h]; if h=0 then goto not_found;
  decr(h);
  end;
not_found: decr(hn)

;
if trie_trc[ cur_lang+ 1] <> cur_lang  then  goto exit ; {no patterns for |cur_lang|}
hc[0]:=0; hc[hn+1]:=0; hc[hn+2]:=256; {insert delimiters}
for j:=0 to hn-r_hyf+1 do
  begin z:=trie_trl[ cur_lang+ 1] +hc[j]; l:=j;
  while hc[l]= trie_trc[  z]   do
    begin if trie_tro[ z] <>min_trie_op then
      
{ Store \(m)maximum values in the |hyf| table }
begin v:=trie_tro[ z] ;
repeat v:=v+op_start[cur_lang]; i:=l-hyf_distance[v];
if hyf_num[v]>hyf[i] then hyf[i]:=hyf_num[v];
v:=hyf_next[v];
until v=min_trie_op;
end

;
    incr(l); z:=trie_trl[ z] +hc[l];
    end;
  end;
found: for j:=0 to l_hyf-1 do hyf[j]:=0;
for j:=0 to r_hyf-1 do hyf[hn-j]:=0

;

{ If no hyphens were found, |return| }
for j:=l_hyf to hn-r_hyf do if odd(hyf[j]) then goto found1;
 goto exit ;
found1:

;

{ Replace nodes |ha..hb| by a sequence of nodes that includes the discretionary hyphens }
q:= mem[ hb].hh.rh ;  mem[ hb].hh.rh :=-{0xfffffff=}268435455  ; r:= mem[ ha].hh.rh ;  mem[ ha].hh.rh :=-{0xfffffff=}268435455  ; bchar:=hyf_bchar;
if  ( ha>=hi_mem_min)  then
  if   mem[ ha].hh.b0 <>hf then goto found2
  else begin init_list:=ha; init_lig:=false; hu[0]:=   mem[  ha].hh.b1  ;
    end
else if  mem[ ha].hh.b0 =ligature_node then
  if   mem[   ha+1 ].hh.b0 <>hf then goto found2
  else begin init_list:= mem[    ha+1 ].hh.rh  ; init_lig:=true; init_lft:=( mem[ ha].hh.b1 >1);
    hu[0]:=   mem[     ha+1 ].hh.b1  ;
    if init_list=-{0xfffffff=}268435455   then if init_lft then
      begin hu[0]:=256; init_lig:=false;
      end; {in this case a ligature will be reconstructed from scratch}
    free_node(ha,small_node_size);
    end
else begin {no punctuation found; look for left boundary}
  if not  ( r>=hi_mem_min)  then if  mem[ r].hh.b0 =ligature_node then
   if  mem[ r].hh.b1 >1 then goto found2;
  j:=1; s:=ha; init_list:=-{0xfffffff=}268435455  ; goto common_ending;
  end;
s:=cur_p; {we have |cur_p<>ha| because |type(cur_p)=glue_node|}
while  mem[ s].hh.rh <>ha do s:= mem[ s].hh.rh ;
j:=0; goto common_ending;
found2: s:=ha; j:=0; hu[0]:=256; init_lig:=false; init_list:=-{0xfffffff=}268435455  ;
common_ending: flush_node_list(r);

{ Reconstitute nodes for the hyphenated word, inserting discretionary hyphens }
repeat l:=j; j:=reconstitute(j,hn,bchar, hyf_char )+1;
if hyphen_passed=0 then
  begin  mem[ s].hh.rh := mem[ mem_top-4 ].hh.rh ;
  while  mem[ s].hh.rh >-{0xfffffff=}268435455   do s:= mem[ s].hh.rh ;
  if odd(hyf[j-1]) then
    begin l:=j; hyphen_passed:=j-1;  mem[ mem_top-4 ].hh.rh :=-{0xfffffff=}268435455  ;
    end;
  end;
if hyphen_passed>0 then
  
{ Create and append a discretionary node as an alternative to the unhyphenated word, and continue to develop both branches until they become equivalent }
repeat r:=get_node(small_node_size);
 mem[ r].hh.rh := mem[ mem_top-4 ].hh.rh ;  mem[ r].hh.b0 :=disc_node;
major_tail:=r; r_count:=0;
while  mem[ major_tail].hh.rh >-{0xfffffff=}268435455   do begin major_tail:= mem[ major_tail].hh.rh ; incr(r_count); end ;
i:=hyphen_passed; hyf[i]:=0;

{ Put the \(c)characters |hu[l..i]| and a hyphen into |pre_break(r)| }
minor_tail:=-{0xfffffff=}268435455  ;   mem[  r+ 1].hh.lh  :=-{0xfffffff=}268435455  ; hyf_node:=new_character(hf,hyf_char);
if hyf_node<>-{0xfffffff=}268435455   then
  begin incr(i); c:=hu[i]; hu[i]:=hyf_char;  begin  mem[  hyf_node].hh.rh :=avail; avail:= hyf_node; ifdef('STAT')  decr(dyn_used); endif('STAT')  end ;
  end;
while l<=i do
  begin l:=reconstitute(l,i,font_bchar[hf], 256  )+1;
  if  mem[ mem_top-4 ].hh.rh >-{0xfffffff=}268435455   then
    begin if minor_tail=-{0xfffffff=}268435455   then   mem[  r+ 1].hh.lh  := mem[ mem_top-4 ].hh.rh 
    else  mem[ minor_tail].hh.rh := mem[ mem_top-4 ].hh.rh ;
    minor_tail:= mem[ mem_top-4 ].hh.rh ;
    while  mem[ minor_tail].hh.rh >-{0xfffffff=}268435455   do minor_tail:= mem[ minor_tail].hh.rh ;
    end;
  end;
if hyf_node<>-{0xfffffff=}268435455   then
  begin hu[i]:=c; {restore the character in the hyphen position}
  l:=i; decr(i);
  end

;

{ Put the \(c)characters |hu[i+1..@,]| into |post_break(r)|, appending to this list and to |major_tail| until synchronization has been achieved }
minor_tail:=-{0xfffffff=}268435455  ;   mem[  r+ 1].hh.rh  :=-{0xfffffff=}268435455  ; c_loc:=0;
if bchar_label[hf]<>non_address then {put left boundary at beginning of new line}
  begin decr(l); c:=hu[l]; c_loc:=l; hu[l]:=256;
  end;
while l<j do
  begin repeat l:=reconstitute(l,hn,bchar, 256  )+1;
  if c_loc>0 then
    begin hu[c_loc]:=c; c_loc:=0;
    end;
  if  mem[ mem_top-4 ].hh.rh >-{0xfffffff=}268435455   then
    begin if minor_tail=-{0xfffffff=}268435455   then   mem[  r+ 1].hh.rh  := mem[ mem_top-4 ].hh.rh 
    else  mem[ minor_tail].hh.rh := mem[ mem_top-4 ].hh.rh ;
    minor_tail:= mem[ mem_top-4 ].hh.rh ;
    while  mem[ minor_tail].hh.rh >-{0xfffffff=}268435455   do minor_tail:= mem[ minor_tail].hh.rh ;
    end;
  until l>=j;
  while l>j do
    
{ Append characters of |hu[j..@,]| to |major_tail|, advancing~|j| }
begin j:=reconstitute(j,hn,bchar, 256  )+1;
 mem[ major_tail].hh.rh := mem[ mem_top-4 ].hh.rh ;
while  mem[ major_tail].hh.rh >-{0xfffffff=}268435455   do begin major_tail:= mem[ major_tail].hh.rh ; incr(r_count); end ;
end

;
  end

;

{ Move pointer |s| to the end of the current list, and set |replace_count(r)| appropriately }
if r_count>127 then {we have to forget the discretionary hyphen}
  begin  mem[ s].hh.rh := mem[ r].hh.rh ;  mem[ r].hh.rh :=-{0xfffffff=}268435455  ; flush_node_list(r);
  end
else begin  mem[ s].hh.rh :=r;  mem[ r].hh.b1 :=r_count;
  end;
s:=major_tail

;
hyphen_passed:=j-1;  mem[ mem_top-4 ].hh.rh :=-{0xfffffff=}268435455  ;
until not odd(hyf[j-1])

;
until j>hn;
 mem[ s].hh.rh :=q

;
flush_list(init_list)

;
exit:end;


 ifdef('INITEX')  
{ Declare procedures for preprocessing hyphenation patterns }
function new_trie_op( d, n:small_number; v:trie_opcode):trie_opcode;
label exit;
var h:neg_trie_op_size..trie_op_size; {trial hash location}
 u:trie_opcode; {trial op code}
 l:0..trie_op_size; {pointer to stored data}
begin h:=abs(n+313*d+361*v+1009*cur_lang) mod (trie_op_size-neg_trie_op_size)
  + neg_trie_op_size;
 while true do    begin l:=trie_op_hash[h];
  if l=0 then {empty position found for a new op}
    begin if trie_op_ptr=trie_op_size then
      overflow({"pattern memory ops"=}963,trie_op_size);
    u:=trie_used[cur_lang];
    if u=max_trie_op then
      overflow({"pattern memory ops per language"=}964,
      max_trie_op-min_trie_op);
    incr(trie_op_ptr); incr(u); trie_used[cur_lang]:=u;
    if u>max_op_used then max_op_used:=u;
    hyf_distance[trie_op_ptr]:=d;
    hyf_num[trie_op_ptr]:=n; hyf_next[trie_op_ptr]:=v;
    trie_op_lang[trie_op_ptr]:=cur_lang; trie_op_hash[h]:=trie_op_ptr;
    trie_op_val[trie_op_ptr]:=u; new_trie_op:=u;  goto exit ;
    end;
  if (hyf_distance[l]=d)and(hyf_num[l]=n)and(hyf_next[l]=v)
   and(trie_op_lang[l]=cur_lang) then
    begin new_trie_op:=trie_op_val[l];  goto exit ;
    end;
  if h>-trie_op_size then decr(h) else h:=trie_op_size;
  end;
exit:end;


function trie_node( p:trie_pointer):trie_pointer; {converts
  to a canonical form}
label exit;
var h:trie_pointer; {trial hash location}
 q:trie_pointer; {trial trie node}
begin h:=abs(trie_c[p]+1009*trie_o[p]+ 
    2718*trie_l[p]+3142*trie_r[p]) mod trie_size;
 while true do    begin q:=trie_hash[h];
  if q=0 then
    begin trie_hash[h]:=p; trie_node:=p;  goto exit ;
    end;
  if (trie_c[q]=trie_c[p])and(trie_o[q]=trie_o[p])and 
    (trie_l[q]=trie_l[p])and(trie_r[q]=trie_r[p]) then
    begin trie_node:=q;  goto exit ;
    end;
  if h>0 then decr(h) else h:=trie_size;
  end;
exit:end;


function compress_trie( p:trie_pointer):trie_pointer;
begin if p=0 then compress_trie:=0
else  begin trie_l[p]:=compress_trie(trie_l[p]);
  trie_r[p]:=compress_trie(trie_r[p]);
  compress_trie:=trie_node(p);
  end;
end;


procedure first_fit( p:trie_pointer); {packs a family into |trie|}
label not_found,found;
var h:trie_pointer; {candidate for |trie_ref[p]|}
 z:trie_pointer; {runs through holes}
 q:trie_pointer; {runs through the family starting at |p|}
 c:ASCII_code; {smallest character in the family}
 l, r:trie_pointer; {left and right neighbors}
 ll:1..256; {upper limit of |trie_min| updating}
begin c:=  trie_c[ p] ;
z:=trie_min[c]; {get the first conceivably good hole}
 while true do    begin h:=z-c;

  
{ Ensure that |trie_max>=h+256| }
if trie_max<h+256 then
  begin if trie_size<=h+256 then overflow({"pattern memory"=}965,trie_size);
{ \xref[TeX capacity exceeded pattern memory][\quad pattern memory] }
  repeat incr(trie_max); trie_taken[trie_max]:=false;
  trie_trl[ trie_max] :=trie_max+1; trie_tro[ trie_max] :=trie_max-1;
  until trie_max=h+256;
  end

;
  if trie_taken[h] then goto not_found;
  
{ If all characters of the family fit relative to |h|, then |goto found|,\30\ otherwise |goto not_found| }
q:=trie_r[p];
while q>0 do
  begin if trie_trl[ h+    trie_c[  q] ] =0 then goto not_found;
  q:=trie_r[q];
  end;
goto found

;
  not_found: z:=trie_trl[ z] ; {move to the next hole}
  end;
found: 
{ Pack the family into |trie| relative to |h| }
trie_taken[h]:=true; trie_hash [p]:=h; q:=p;
repeat z:=h+  trie_c[ q] ; l:=trie_tro[ z] ; r:=trie_trl[ z] ;
trie_tro[ r] :=l; trie_trl[ l] :=r; trie_trl[ z] :=0;
if l<256 then
  begin if z<256 then ll:=z  else ll:=256;
  repeat trie_min[l]:=r; incr(l);
  until l=ll;
  end;
q:=trie_r[q];
until q=0

;
end;


procedure trie_pack( p:trie_pointer); {pack subtries of a family}
var q:trie_pointer; {a local variable that need not be saved on recursive calls}
begin repeat q:=trie_l[p];
if (q>0)and(trie_hash [q]=0) then
  begin first_fit(q); trie_pack(q);
  end;
p:=trie_r[p];
until p=0;
end;


procedure trie_fix( p:trie_pointer); {moves |p| and its siblings into |trie|}
var q:trie_pointer; {a local variable that need not be saved on recursive calls}
 c:ASCII_code; {another one that need not be saved}
 z:trie_pointer; {|trie| reference; this local variable must be saved}
begin z:=trie_hash [p];
repeat q:=trie_l[p]; c:=  trie_c[ p] ;
trie_trl[ z+ c] :=trie_hash [q]; trie_trc[ z+ c] := c ; trie_tro[ z+ c] :=trie_o[p];
if q>0 then trie_fix(q);
p:=trie_r[p];
until p=0;
end;


procedure new_patterns; {initializes the hyphenation pattern data}
label done, done1;
var k, l:0..64; {indices into |hc| and |hyf|;
                  not always in |small_number| range}
 digit_sensed:boolean; {should the next digit be treated as a letter?}
 v:trie_opcode; {trie op code}
 p, q:trie_pointer; {nodes of trie traversed during insertion}
 first_child:boolean; {is |p=trie_l[q]|?}
 c:ASCII_code; {character being inserted}
begin if trie_not_ready then
  begin if eqtb[int_base+ language_code].int  <=0 then cur_lang:=0 else if eqtb[int_base+ language_code].int  >255 then cur_lang:=0 else cur_lang:=eqtb[int_base+ language_code].int   ; scan_left_brace; {a left brace must follow \.[\\patterns]}
  
{ Enter all of the patterns into a linked trie, until coming to a right brace }
k:=0; hyf[0]:=0; digit_sensed:=false;
 while true do    begin get_x_token;
  case cur_cmd of
  letter,other_char:
{ Append a new letter or a hyphen level }
if digit_sensed or(cur_chr<{"0"=}48)or(cur_chr>{"9"=}57) then
  begin if cur_chr={"."=}46 then cur_chr:=0 {edge-of-word delimiter}
  else  begin cur_chr:= eqtb[  lc_code_base+   cur_chr].hh.rh   ;
    if cur_chr=0 then
      begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Nonletter"=} 971); end ;
{ \xref[Nonletter] }
       begin help_ptr:=1; help_line[0]:={"(See Appendix H.)"=} 970; end ; error;
      end;
    end;
  if k<63 then
    begin incr(k); hc[k]:=cur_chr; hyf[k]:=0; digit_sensed:=false;
    end;
  end
else if k<63 then
  begin hyf[k]:=cur_chr-{"0"=}48; digit_sensed:=true;
  end

;
  spacer,right_brace: begin if k>0 then
      
{ Insert a new pattern into the linked trie }
begin 
{ Compute the trie op code, |v|, and set |l:=0| }
if hc[1]=0 then hyf[0]:=0;
if hc[k]=0 then hyf[k]:=0;
l:=k; v:=min_trie_op;
 while true do    begin if hyf[l]<>0 then v:=new_trie_op(k-l,hyf[l],v);
  if l>0 then decr(l) else goto done1;
  end;
done1:

;
q:=0; hc[0]:=cur_lang;
while l<=k do
  begin c:=hc[l]; incr(l); p:=trie_l[q]; first_child:=true;
  while (p>0)and(c>  trie_c[ p] ) do
    begin q:=p; p:=trie_r[q]; first_child:=false;
    end;
  if (p=0)or(c<  trie_c[ p] ) then
    
{ Insert a new trie node between |q| and |p|, and make |p| point to it }
begin if trie_ptr=trie_size then overflow({"pattern memory"=}965,trie_size);
{ \xref[TeX capacity exceeded pattern memory][\quad pattern memory] }
incr(trie_ptr); trie_r[trie_ptr]:=p; p:=trie_ptr; trie_l[p]:=0;
if first_child then trie_l[q]:=p else trie_r[q]:=p;
trie_c[p]:=  c ; trie_o[p]:=min_trie_op;
end

;
  q:=p; {now node |q| represents $p_1\ldots p_[l-1]$}
  end;
if trie_o[q]<>min_trie_op then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Duplicate pattern"=} 972); end ;
{ \xref[Duplicate pattern] }
   begin help_ptr:=1; help_line[0]:={"(See Appendix H.)"=} 970; end ; error;
  end;
trie_o[q]:=v;
end

;
    if cur_cmd=right_brace then goto done;
    k:=0; hyf[0]:=0; digit_sensed:=false;
    end;
   else  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Bad "=} 969); end ; print_esc({"patterns"=}967);
{ \xref[Bad \\patterns] }
     begin help_ptr:=1; help_line[0]:={"(See Appendix H.)"=} 970; end ; error;
    end
   end ;
  end;
done:

;
  end
else begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Too late for "=} 966); end ; print_esc({"patterns"=}967);
   begin help_ptr:=1; help_line[0]:={"All patterns must be given before typesetting begins."=} 968; end ;
  error;  mem[ mem_top-12 ].hh.rh :=scan_toks(false,false); flush_list(def_ref);
  end;
end;


procedure init_trie;
var  p:trie_pointer; {pointer for initialization}
 j, k, t:integer; {all-purpose registers for initialization}
 r, s:trie_pointer; {used to clean up the packed |trie|}

begin 
{ Get ready to compress the trie }

{ Sort \(t)the hyphenation... }
op_start[0]:=-min_trie_op;
for j:=1 to 255 do op_start[j]:=op_start[j-1]+ trie_used[ j- 1] ;
for j:=1 to trie_op_ptr do
  trie_op_hash[j]:=op_start[trie_op_lang[j]]+trie_op_val[j]; {destination}
for j:=1 to trie_op_ptr do while trie_op_hash[j]>j do
  begin k:=trie_op_hash[j];

  t:=hyf_distance[k]; hyf_distance[k]:=hyf_distance[j]; hyf_distance[j]:=t;

  t:=hyf_num[k]; hyf_num[k]:=hyf_num[j]; hyf_num[j]:=t;

  t:=hyf_next[k]; hyf_next[k]:=hyf_next[j]; hyf_next[j]:=t;

  trie_op_hash[j]:=trie_op_hash[k]; trie_op_hash[k]:=k;
  end

;
for p:=0 to trie_size do trie_hash[p]:=0;
trie_l[0] :=compress_trie(trie_l[0] ); {identify equivalent subtries}
for p:=0 to trie_ptr do trie_hash [p]:=0;
for p:=0 to 255 do trie_min[p]:=p+1;
trie_trl[ 0] :=1; trie_max:=0

;
if trie_l[0] <>0 then
  begin first_fit(trie_l[0] ); trie_pack(trie_l[0] );
  end;

{ Move the data into |trie| }
if trie_l[0] =0 then {no patterns were given}
  begin for r:=0 to 256 do  begin trie_trl[ r] :=0; trie_tro[ r] :=min_trie_op; trie_trc[ r] :=min_quarterword; end ;
  trie_max:=256;
  end
else begin trie_fix(trie_l[0] ); {this fixes the non-holes in |trie|}
  r:=0; {now we will zero out all the holes}
  repeat s:=trie_trl[ r] ;  begin trie_trl[ r] :=0; trie_tro[ r] :=min_trie_op; trie_trc[ r] :=min_quarterword; end ; r:=s;
  until r>trie_max;
  end;
trie_trc[ 0] :={"?"=} 63 ; {make |trie_char(c)<>c| for all |c|}

;
trie_not_ready:=false;
end;

 
endif('INITEX') 


procedure line_break( final_widow_penalty:integer);
label done,done1,done2,done3,done4,done5,continue;
var 
{ Local variables for line breaking }
 auto_breaking:boolean; {is node |cur_p| outside a formula?}
 prev_p:halfword ; {helps to determine when glue nodes are breakpoints}
 q, r, s, prev_s:halfword ; {miscellaneous nodes of temporary interest}
 f:internal_font_number; {used when calculating character widths}


 j:small_number; {an index into |hc| or |hu|}
 c:0..255; {character being considered for hyphenation}

 
begin pack_begin_line:=cur_list.ml_field ; {this is for over/underfull box messages}

{ Get ready to start line breaking }
 mem[ mem_top-3 ].hh.rh := mem[ cur_list.head_field ].hh.rh ;
if  ( cur_list.tail_field >=hi_mem_min)  then begin  mem[ cur_list.tail_field ].hh.rh := new_penalty( inf_penalty); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end 
else if  mem[ cur_list.tail_field ].hh.b0 <>glue_node then begin  mem[ cur_list.tail_field ].hh.rh := new_penalty( inf_penalty); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end 
else  begin  mem[ cur_list.tail_field ].hh.b0 :=penalty_node; delete_glue_ref(  mem[  cur_list.tail_field + 1].hh.lh  );
  flush_node_list(  mem[  cur_list.tail_field + 1].hh.rh  );  mem[ cur_list.tail_field +1].int :=inf_penalty;
  end;
 mem[ cur_list.tail_field ].hh.rh :=new_param_glue(par_fill_skip_code);
init_cur_lang:=cur_list.pg_field  mod {0200000=}65536;
init_l_hyf:=cur_list.pg_field  div {020000000=}4194304;
init_r_hyf:=(cur_list.pg_field  div {0200000=}65536) mod {0100=}64;
pop_nest;


no_shrink_error_yet:=true;

if (  mem[   eqtb[  glue_base+   left_skip_code].hh.rh    ].hh.b1 <>normal)and( mem[   eqtb[  glue_base+   left_skip_code].hh.rh    +3].int  <>0) then begin   eqtb[  glue_base+   left_skip_code].hh.rh    :=finite_shrink(  eqtb[  glue_base+   left_skip_code].hh.rh    ); end ; if (  mem[   eqtb[  glue_base+   right_skip_code].hh.rh    ].hh.b1 <>normal)and( mem[   eqtb[  glue_base+   right_skip_code].hh.rh    +3].int  <>0) then begin   eqtb[  glue_base+   right_skip_code].hh.rh    :=finite_shrink(  eqtb[  glue_base+   right_skip_code].hh.rh    ); end ;

q:= eqtb[  glue_base+   left_skip_code].hh.rh    ; r:= eqtb[  glue_base+   right_skip_code].hh.rh    ; background[1]:= mem[ q+width_offset].int  + mem[ r+width_offset].int  ;

background[2]:=0; background[3]:=0; background[4]:=0; background[5]:=0;

background[2+  mem[ q].hh.b0 ]:= mem[ q+2].int  ;

background[2+  mem[ r].hh.b0 ]:= background[2+  mem[ r].hh.b0 ]+ mem[ r+2].int  ;

background[6]:= mem[ q+3].int  + mem[ r+3].int  ;


minimum_demerits:={07777777777=}1073741823 ;
minimal_demerits[tight_fit]:={07777777777=}1073741823 ;
minimal_demerits[decent_fit]:={07777777777=}1073741823 ;
minimal_demerits[loose_fit]:={07777777777=}1073741823 ;
minimal_demerits[very_loose_fit]:={07777777777=}1073741823 ;


if  eqtb[  par_shape_loc].hh.rh   =-{0xfffffff=}268435455   then
  if eqtb[dimen_base+ hang_indent_code].int   =0 then
    begin last_special_line:=0; second_width:=eqtb[dimen_base+ hsize_code].int   ;
    second_indent:=0;
    end
  else 
{ Set line length parameters in preparation for hanging indentation }
begin last_special_line:=abs(eqtb[int_base+ hang_after_code].int  );
if eqtb[int_base+ hang_after_code].int  <0 then
  begin first_width:=eqtb[dimen_base+ hsize_code].int   -abs(eqtb[dimen_base+ hang_indent_code].int   );
  if eqtb[dimen_base+ hang_indent_code].int   >=0 then first_indent:=eqtb[dimen_base+ hang_indent_code].int   
  else first_indent:=0;
  second_width:=eqtb[dimen_base+ hsize_code].int   ; second_indent:=0;
  end
else  begin first_width:=eqtb[dimen_base+ hsize_code].int   ; first_indent:=0;
  second_width:=eqtb[dimen_base+ hsize_code].int   -abs(eqtb[dimen_base+ hang_indent_code].int   );
  if eqtb[dimen_base+ hang_indent_code].int   >=0 then second_indent:=eqtb[dimen_base+ hang_indent_code].int   
  else second_indent:=0;
  end;
end


else  begin last_special_line:= mem[  eqtb[  par_shape_loc].hh.rh   ].hh.lh -1;
  second_width:=mem[ eqtb[  par_shape_loc].hh.rh   +2*(last_special_line+1)].int ;
  second_indent:=mem[ eqtb[  par_shape_loc].hh.rh   +2*last_special_line+1].int ;
  end;
if eqtb[int_base+ looseness_code].int  =0 then easy_line:=last_special_line
else easy_line:={0xfffffff=}268435455 

;

{ Find optimal breakpoints }
threshold:=eqtb[int_base+ pretolerance_code].int  ;
if threshold>=0 then
  begin  ifdef('STAT')  if eqtb[int_base+ tracing_paragraphs_code].int  >0 then
    begin begin_diagnostic; print_nl({"@firstpass"=}947); end;  endif('STAT')  

  second_pass:=false; final_pass:=false;
  end
else  begin threshold:=eqtb[int_base+ tolerance_code].int  ; second_pass:=true;
  final_pass:=(eqtb[dimen_base+ emergency_stretch_code].int   <=0);
   ifdef('STAT')  if eqtb[int_base+ tracing_paragraphs_code].int  >0 then begin_diagnostic; endif('STAT')  
  end;
 while true do    begin if threshold>inf_bad then threshold:=inf_bad;
  if second_pass then 
{ Initialize for hyphenating a paragraph }
begin  ifdef('INITEX')  if trie_not_ready then init_trie; endif('INITEX')  

cur_lang:=init_cur_lang; l_hyf:=init_l_hyf; r_hyf:=init_r_hyf;
end

;
  
{ Create an active breakpoint representing the beginning of the paragraph }
q:=get_node(active_node_size);
 mem[ q].hh.b0 :=unhyphenated;  mem[ q].hh.b1 :=decent_fit;
 mem[ q].hh.rh :=mem_top-7  ;   mem[  q+ 1].hh.rh  :=-{0xfffffff=}268435455  ;
  mem[  q+ 1].hh.lh  :=cur_list.pg_field +1; mem[ q+2].int :=0;  mem[ mem_top-7 ].hh.rh :=q;
 active_width[ 1]:=background[ 1] ; active_width[ 2]:=background[ 2] ; active_width[ 3]:=background[ 3] ; active_width[ 4]:=background[ 4] ; active_width[ 5]:=background[ 5] ; active_width[ 6]:=background[ 6]  ;

passive:=-{0xfffffff=}268435455  ; printed_node:=mem_top-3 ; pass_number:=0;
font_in_short_display:=font_base 

;
  cur_p:= mem[ mem_top-3 ].hh.rh ; auto_breaking:=true;

  prev_p:=cur_p; {glue at beginning is not a legal breakpoint}
  while (cur_p<>-{0xfffffff=}268435455  )and( mem[ mem_top-7 ].hh.rh <>mem_top-7  ) do
    
{ Call |try_break| if |cur_p| is a legal breakpoint; on the second pass, also try to hyphenate the next word, if |cur_p| is a glue node; then advance |cur_p| to the next node of the paragraph that could possibly be a legal breakpoint }
begin if  ( cur_p>=hi_mem_min)  then
  
{ Advance \(c)|cur_p| to the node following the present string of characters }
begin prev_p:=cur_p;
repeat f:=  mem[ cur_p].hh.b0 ;
active_width[1] :=active_width[1] +font_info[width_base[ f]+  font_info[char_base[  f]+effective_char(true,  f,     mem[   cur_p].hh.b1 )].qqqq .b0].int  ;
cur_p:= mem[ cur_p].hh.rh ;
until not  ( cur_p>=hi_mem_min) ;
end

;
case  mem[ cur_p].hh.b0  of
hlist_node,vlist_node,rule_node: active_width[1] :=active_width[1] + mem[ cur_p+width_offset].int  ;
whatsit_node: 
{ Advance \(p)past a whatsit node in the \(l)|line_break| loop } 
 if  mem[  cur_p].hh.b1 =language_node then begin cur_lang:= mem[   cur_p+ 1].hh.rh  ; l_hyf:= mem[   cur_p+ 1].hh.b0  ; r_hyf:= mem[   cur_p+ 1].hh.b1  ; end 

;
glue_node: begin 
{ If node |cur_p| is a legal breakpoint, call |try_break|; then update the active widths by including the glue in |glue_ptr(cur_p)| }
if auto_breaking then
  begin if  ( prev_p>=hi_mem_min)  then try_break(0,unhyphenated)
  else if ( mem[  prev_p].hh.b0 <math_node)  then try_break(0,unhyphenated)
  else if ( mem[ prev_p].hh.b0 =kern_node)and( mem[ prev_p].hh.b1 <>explicit) then
    try_break(0,unhyphenated);
  end;
if (  mem[    mem[    cur_p+ 1].hh.lh  ].hh.b1 <>normal)and( mem[    mem[    cur_p+ 1].hh.lh  +3].int  <>0) then begin    mem[   cur_p+ 1].hh.lh  :=finite_shrink(   mem[   cur_p+ 1].hh.lh  ); end ; q:=  mem[  cur_p+ 1].hh.lh  ;
active_width[1] :=active_width[1] + mem[ q+width_offset].int  ; 
active_width[2+  mem[ q].hh.b0 ]:= 
  active_width[2+  mem[ q].hh.b0 ]+ mem[ q+2].int  ;

active_width[6]:=active_width[6]+ mem[ q+3].int  

;
  if second_pass and auto_breaking then
    
{ Try to hyphenate the following word }
begin prev_s:=cur_p; s:= mem[ prev_s].hh.rh ;
if s<>-{0xfffffff=}268435455   then
  begin 
{ Skip to node |ha|, or |goto done1| if no hyphenation should be attempted }
 while true do    begin if  ( s>=hi_mem_min)  then
    begin c:=   mem[  s].hh.b1  ; hf:=  mem[ s].hh.b0 ;
    end
  else if  mem[ s].hh.b0 =ligature_node then
    if  mem[    s+1 ].hh.rh  =-{0xfffffff=}268435455   then goto continue
    else begin q:= mem[    s+1 ].hh.rh  ; c:=   mem[  q].hh.b1  ; hf:=  mem[ q].hh.b0 ;
      end
  else if ( mem[ s].hh.b0 =kern_node)and( mem[ s].hh.b1 =normal) then goto continue
  else if  mem[ s].hh.b0 =whatsit_node then
    begin 
{ Advance \(p)past a whatsit node in the \(p)pre-hyphenation loop } 
 if  mem[  s].hh.b1 =language_node then begin cur_lang:= mem[   s+ 1].hh.rh  ; l_hyf:= mem[   s+ 1].hh.b0  ; r_hyf:= mem[   s+ 1].hh.b1  ; end 

;
    goto continue;
    end
  else goto done1;
  if  eqtb[  lc_code_base+   c].hh.rh   <>0 then
    if ( eqtb[  lc_code_base+   c].hh.rh   =c)or(eqtb[int_base+ uc_hyph_code].int  >0) then goto done2
    else goto done1;
continue: prev_s:=s; s:= mem[ prev_s].hh.rh ;
  end;
done2: hyf_char:=hyphen_char[hf];
if hyf_char<0 then goto done1;
if hyf_char>255 then goto done1;
ha:=prev_s

;
  if l_hyf+r_hyf>63 then goto done1;
  
{ Skip to node |hb|, putting letters into |hu| and |hc| }
hn:=0;
 while true do    begin if  ( s>=hi_mem_min)  then
    begin if   mem[ s].hh.b0 <>hf then goto done3;
    hyf_bchar:=  mem[ s].hh.b1 ; c:= hyf_bchar ;
    if  eqtb[  lc_code_base+   c].hh.rh   =0 then goto done3;
    if hn=63 then goto done3;
    hb:=s; incr(hn); hu[hn]:=c; hc[hn]:= eqtb[  lc_code_base+   c].hh.rh   ; hyf_bchar:= 256  ;
    end
  else if  mem[ s].hh.b0 =ligature_node then
    
{ Move the characters of a ligature node to |hu| and |hc|; but |goto done3| if they are not all letters }
begin if   mem[   s+1 ].hh.b0 <>hf then goto done3;
j:=hn; q:= mem[    s+1 ].hh.rh  ; if q>-{0xfffffff=}268435455   then hyf_bchar:=  mem[ q].hh.b1 ;
while q>-{0xfffffff=}268435455   do
  begin c:=   mem[  q].hh.b1  ;
  if  eqtb[  lc_code_base+   c].hh.rh   =0 then goto done3;
  if j=63 then goto done3;
  incr(j); hu[j]:=c; hc[j]:= eqtb[  lc_code_base+   c].hh.rh   ;

  q:= mem[ q].hh.rh ;
  end;
hb:=s; hn:=j;
if odd( mem[ s].hh.b1 ) then hyf_bchar:=font_bchar[hf] else hyf_bchar:= 256  ;
end


  else if ( mem[ s].hh.b0 =kern_node)and( mem[ s].hh.b1 =normal) then
    begin hb:=s;
    hyf_bchar:=font_bchar[hf];
    end
  else goto done3;
  s:= mem[ s].hh.rh ;
  end;
done3:

;
  
{ Check that the nodes following |hb| permit hyphenation and that at least |l_hyf+r_hyf| letters have been found, otherwise |goto done1| }
if hn<l_hyf+r_hyf then goto done1; {|l_hyf| and |r_hyf| are |>=1|}
 while true do    begin if not( ( s>=hi_mem_min) ) then
    case  mem[ s].hh.b0  of
    ligature_node:  ;
    kern_node: if  mem[ s].hh.b1 <>normal then goto done4;
    whatsit_node,glue_node,penalty_node,ins_node,adjust_node,mark_node:
      goto done4;
     else  goto done1
     end ;
  s:= mem[ s].hh.rh ;
  end;
done4:

;
  hyphenate;
  end;
done1: end

;
  end;
kern_node: if  mem[ cur_p].hh.b1 =explicit then begin if not  (  mem[  cur_p].hh.rh >=hi_mem_min)  and auto_breaking then if  mem[  mem[  cur_p].hh.rh ].hh.b0 =glue_node then try_break(0,unhyphenated); active_width[1] :=active_width[1] + mem[ cur_p+width_offset].int  ; end 
  else active_width[1] :=active_width[1] + mem[ cur_p+width_offset].int  ;
ligature_node: begin f:=  mem[   cur_p+1 ].hh.b0 ;
  active_width[1] :=active_width[1] +font_info[width_base[ f]+  font_info[char_base[  f]+effective_char(true,  f,     mem[       cur_p+1 ].hh.b1 )].qqqq .b0].int  ;
  end;
disc_node: 
{ Try to break after a discretionary fragment, then |goto done5| }
begin s:=  mem[  cur_p+ 1].hh.lh  ; disc_width:=0;
if s=-{0xfffffff=}268435455   then try_break(eqtb[int_base+ ex_hyphen_penalty_code].int  ,hyphenated)
else  begin repeat 
{ Add the width of node |s| to |disc_width| }
if  ( s>=hi_mem_min)  then
  begin f:=  mem[ s].hh.b0 ;
  disc_width:=disc_width+font_info[width_base[ f]+  font_info[char_base[  f]+effective_char(true,  f,     mem[   s].hh.b1 )].qqqq .b0].int  ;
  end
else  case  mem[ s].hh.b0  of
  ligature_node: begin f:=  mem[   s+1 ].hh.b0 ;
    disc_width:=disc_width+
      font_info[width_base[ f]+  font_info[char_base[  f]+effective_char(true,  f,     mem[       s+1 ].hh.b1 )].qqqq .b0].int  ;
    end;
  hlist_node,vlist_node,rule_node,kern_node:
    disc_width:=disc_width+ mem[ s+width_offset].int  ;
   else  confusion({"disc3"=}951)
{ \xref[this can't happen disc3][\quad disc3] }
   end 

;
    s:= mem[ s].hh.rh ;
  until s=-{0xfffffff=}268435455  ;
  active_width[1] :=active_width[1] +disc_width;
  try_break(eqtb[int_base+ hyphen_penalty_code].int  ,hyphenated);
  active_width[1] :=active_width[1] -disc_width;
  end;
r:= mem[ cur_p].hh.b1 ; s:= mem[ cur_p].hh.rh ;
while r>0 do
  begin 
{ Add the width of node |s| to |act_width| }
if  ( s>=hi_mem_min)  then
  begin f:=  mem[ s].hh.b0 ;
  active_width[1] :=active_width[1] +font_info[width_base[ f]+  font_info[char_base[  f]+effective_char(true,  f,     mem[   s].hh.b1 )].qqqq .b0].int  ;
  end
else  case  mem[ s].hh.b0  of
  ligature_node: begin f:=  mem[   s+1 ].hh.b0 ;
    active_width[1] :=active_width[1] +
      font_info[width_base[ f]+  font_info[char_base[  f]+effective_char(true,  f,     mem[       s+1 ].hh.b1 )].qqqq .b0].int  ;
    end;
  hlist_node,vlist_node,rule_node,kern_node:
    active_width[1] :=active_width[1] + mem[ s+width_offset].int  ;
   else  confusion({"disc4"=}952)
{ \xref[this can't happen disc4][\quad disc4] }
   end 

;
  decr(r); s:= mem[ s].hh.rh ;
  end;
prev_p:=cur_p; cur_p:=s; goto done5;
end

;
math_node: begin auto_breaking:=( mem[ cur_p].hh.b1 =after); begin if not  (  mem[  cur_p].hh.rh >=hi_mem_min)  and auto_breaking then if  mem[  mem[  cur_p].hh.rh ].hh.b0 =glue_node then try_break(0,unhyphenated); active_width[1] :=active_width[1] + mem[ cur_p+width_offset].int  ; end ;
  end;
penalty_node: try_break( mem[ cur_p+1].int ,unhyphenated);
mark_node,ins_node,adjust_node:  ;
 else  confusion({"paragraph"=}950)
{ \xref[this can't happen paragraph][\quad paragraph] }
 end ;

prev_p:=cur_p; cur_p:= mem[ cur_p].hh.rh ;
done5:end

;
  if cur_p=-{0xfffffff=}268435455   then
    
{ Try the final line break at the end of the paragraph, and |goto done| if the desired breakpoints have been found }
begin try_break(eject_penalty,hyphenated);
if  mem[ mem_top-7 ].hh.rh <>mem_top-7   then
  begin 
{ Find an active node with fewest demerits }
r:= mem[ mem_top-7 ].hh.rh ; fewest_demerits:={07777777777=}1073741823 ;
repeat if  mem[ r].hh.b0 <>delta_node then if mem[ r+2].int <fewest_demerits then
  begin fewest_demerits:=mem[ r+2].int ; best_bet:=r;
  end;
r:= mem[ r].hh.rh ;
until r=mem_top-7  ;
best_line:=  mem[  best_bet+ 1].hh.lh  

;
  if eqtb[int_base+ looseness_code].int  =0 then goto done;
  
{ Find the best active node for the desired looseness }
begin r:= mem[ mem_top-7 ].hh.rh ; actual_looseness:=0;
repeat if  mem[ r].hh.b0 <>delta_node then
  begin line_diff:=  mem[  r+ 1].hh.lh  -best_line;
  if ((line_diff<actual_looseness)and(eqtb[int_base+ looseness_code].int  <=line_diff))or 
  ((line_diff>actual_looseness)and(eqtb[int_base+ looseness_code].int  >=line_diff)) then
    begin best_bet:=r; actual_looseness:=line_diff;
    fewest_demerits:=mem[ r+2].int ;
    end
  else if (line_diff=actual_looseness)and 
    (mem[ r+2].int <fewest_demerits) then
    begin best_bet:=r; fewest_demerits:=mem[ r+2].int ;
    end;
  end;
r:= mem[ r].hh.rh ;
until r=mem_top-7  ;
best_line:=  mem[  best_bet+ 1].hh.lh  ;
end

;
  if (actual_looseness=eqtb[int_base+ looseness_code].int  )or final_pass then goto done;
  end;
end

;
  
{ Clean up the memory by removing the break nodes }
q:= mem[ mem_top-7 ].hh.rh ;
while q<>mem_top-7   do
  begin cur_p:= mem[ q].hh.rh ;
  if  mem[ q].hh.b0 =delta_node then free_node(q,delta_node_size)
  else free_node(q,active_node_size);
  q:=cur_p;
  end;
q:=passive;
while q<>-{0xfffffff=}268435455   do
  begin cur_p:= mem[ q].hh.rh ;
  free_node(q,passive_node_size);
  q:=cur_p;
  end

;
  if not second_pass then
    begin ifdef('STAT')  if eqtb[int_base+ tracing_paragraphs_code].int  >0 then print_nl({"@secondpass"=}948);  endif('STAT') 

    threshold:=eqtb[int_base+ tolerance_code].int  ; second_pass:=true; final_pass:=(eqtb[dimen_base+ emergency_stretch_code].int   <=0);
    end {if at first you don't succeed, \dots}
  else begin  ifdef('STAT')  if eqtb[int_base+ tracing_paragraphs_code].int  >0 then
      print_nl({"@emergencypass"=}949);  endif('STAT') 

    background[2]:=background[2]+eqtb[dimen_base+ emergency_stretch_code].int   ; final_pass:=true;
    end;
  end;
done:  ifdef('STAT')  if eqtb[int_base+ tracing_paragraphs_code].int  >0 then
  begin end_diagnostic(true); normalize_selector;
  end; endif('STAT') 


;

{ Break the paragraph at the chosen breakpoints, justify the resulting lines to the correct widths, and append them to the current vertical list }
post_line_break(final_widow_penalty)

;

{ Clean up the memory by removing the break nodes }
q:= mem[ mem_top-7 ].hh.rh ;
while q<>mem_top-7   do
  begin cur_p:= mem[ q].hh.rh ;
  if  mem[ q].hh.b0 =delta_node then free_node(q,delta_node_size)
  else free_node(q,active_node_size);
  q:=cur_p;
  end;
q:=passive;
while q<>-{0xfffffff=}268435455   do
  begin cur_p:= mem[ q].hh.rh ;
  free_node(q,passive_node_size);
  q:=cur_p;
  end

;
pack_begin_line:=0;
end;



{ 817. }

{tangle:pos tex.web:16090:1: }

{ When looking for optimal line breaks, \TeX\ creates a ``break node'' for
each break that is [\sl feasible], in the sense that there is a way to end
a line at the given place without requiring any line to stretch more than
a given tolerance. A break node is characterized by three things: the position
of the break (which is a pointer to a |glue_node|, |math_node|, |penalty_node|,
or |disc_node|); the ordinal number of the line that will follow this
breakpoint; and the fitness classification of the line that has just
ended, i.e., |tight_fit|, |decent_fit|, |loose_fit|, or |very_loose_fit|. }

{ 818. }

{tangle:pos tex.web:16107:1: }

{ The algorithm essentially determines the best possible way to achieve
each feasible combination of position, line, and fitness. Thus, it answers
questions like, ``What is the best way to break the opening part of the
paragraph so that the fourth line is a tight line ending at such-and-such
a place?'' However, the fact that all lines are to be the same length
after a certain point makes it possible to regard all sufficiently large
line numbers as equivalent, when the looseness parameter is zero, and this
makes it possible for the algorithm to save space and time.

An ``active node'' and a ``passive node'' are created in |mem| for each
feasible breakpoint that needs to be considered. Active nodes are three
words long and passive nodes are two words long. We need active nodes only
for breakpoints near the place in the paragraph that is currently being
examined, so they are recycled within a comparatively short time after
they are created. }

{ 819. }

{tangle:pos tex.web:16123:1: }

{ An active node for a given breakpoint contains six fields:

\yskip\hang|link| points to the next node in the list of active nodes; the
last active node has |link=last_active|.

\yskip\hang|break_node| points to the passive node associated with this
breakpoint.

\yskip\hang|line_number| is the number of the line that follows this
breakpoint.

\yskip\hang|fitness| is the fitness classification of the line ending at this
breakpoint.

\yskip\hang|type| is either |hyphenated| or |unhyphenated|, depending on
whether this breakpoint is a |disc_node|.

\yskip\hang|total_demerits| is the minimum possible sum of demerits over all
lines leading from the beginning of the paragraph to this breakpoint.

\yskip\noindent
The value of |link(active)| points to the first active node on a linked list
of all currently active nodes. This list is in order by |line_number|,
except that nodes with |line_number>easy_line| may be in any order relative
to each other. }

{ 822. }

{tangle:pos tex.web:16193:1: }

{ The active list also contains ``delta'' nodes that help the algorithm
compute the badness of individual lines. Such nodes appear only between two
active nodes, and they have |type=delta_node|. If |p| and |r| are active nodes
and if |q| is a delta node between them, so that |link(p)=q| and |link(q)=r|,
then |q| tells the space difference between lines in the horizontal list that
start after breakpoint |p| and lines that start after breakpoint |r|. In
other words, if we know the length of the line that starts after |p| and
ends at our current position, then the corresponding length of the line that
starts after |r| is obtained by adding the amounts in node~|q|. A delta node
contains six scaled numbers, since it must record the net change in glue
stretchability with respect to all orders of infinity. The natural width
difference appears in |mem[q+1].sc|; the stretch differences in units of
pt, fil, fill, and filll appear in |mem[q+2..q+5].sc|; and the shrink difference
appears in |mem[q+6].sc|. The |subtype| field of a delta node is not used. }

{ 824. }

{tangle:pos tex.web:16234:1: }

{ Let's state the principles of the delta nodes more precisely and concisely,
so that the following programs will be less obscure. For each legal
breakpoint~|p| in the paragraph, we define two quantities $\alpha(p)$ and
$\beta(p)$ such that the length of material in a line from breakpoint~|p|
to breakpoint~|q| is $\gamma+\beta(q)-\alpha(p)$, for some fixed $\gamma$.
Intuitively, $\alpha(p)$ and $\beta(q)$ are the total length of material from
the beginning of the paragraph to a point ``after'' a break at |p| and to a
point ``before'' a break at |q|; and $\gamma$ is the width of an empty line,
namely the length contributed by \.[\\leftskip] and \.[\\rightskip].

Suppose, for example, that the paragraph consists entirely of alternating
boxes and glue skips; let the boxes have widths $x_1\ldots x_n$ and
let the skips have widths $y_1\ldots y_n$, so that the paragraph can be
represented by $x_1y_1\ldots x_ny_n$. Let $p_i$ be the legal breakpoint
at $y_i$; then $\alpha(p_i)=x_1+y_1+\cdots+x_i+y_i$, and $\beta(p_i)=
x_1+y_1+\cdots+x_i$. To check this, note that the length of material from
$p_2$ to $p_5$, say, is $\gamma+x_3+y_3+x_4+y_4+x_5=\gamma+\beta(p_5)
-\alpha(p_2)$.

The quantities $\alpha$, $\beta$, $\gamma$ involve glue stretchability and
shrinkability as well as a natural width. If we were to compute $\alpha(p)$
and $\beta(p)$ for each |p|, we would need multiple precision arithmetic, and
the multiprecise numbers would have to be kept in the active nodes.
\TeX\ avoids this problem by working entirely with relative differences
or ``deltas.'' Suppose, for example, that the active list contains
$a_1\,\delta_1\,a_2\,\delta_2\,a_3$, where the |a|'s are active breakpoints
and the $\delta$'s are delta nodes. Then $\delta_1=\alpha(a_1)-\alpha(a_2)$
and $\delta_2=\alpha(a_2)-\alpha(a_3)$. If the line breaking algorithm is
currently positioned at some other breakpoint |p|, the |active_width| array
contains the value $\gamma+\beta(p)-\alpha(a_1)$. If we are scanning through
the list of active nodes and considering a tentative line that runs from
$a_2$ to~|p|, say, the |cur_active_width| array will contain the value
$\gamma+\beta(p)-\alpha(a_2)$. Thus, when we move from $a_2$ to $a_3$,
we want to add $\alpha(a_2)-\alpha(a_3)$ to |cur_active_width|; and this
is just $\delta_2$, which appears in the active list between $a_2$ and
$a_3$. The |background| array contains $\gamma$. The |break_width| array
will be used to calculate values of new delta nodes when the active
list is being updated. }

{ 904. }

{tangle:pos tex.web:17720:1: }

{ We must now face the fact that the battle is not over, even though the
[\def\![\kern-1pt]%
hyphens have been found: The process of reconstituting a word can be nontrivial
because ligatures might change when a hyphen is present. [\sl The \TeX book\/]
discusses the difficulties of the word ``difficult'', and
the discretionary material surrounding a
hyphen can be considerably more complex than that. Suppose
\.[abcdef] is a word in a font for which the only ligatures are \.[b\!c],
\.[c\!d], \.[d\!e], and \.[e\!f]. If this word permits hyphenation
between \.b and \.c, the two patterns with and without hyphenation are
$\.a\,\.b\,\.-\,\.[c\!d]\,\.[e\!f]$ and $\.a\,\.[b\!c]\,\.[d\!e]\,\.f$.
Thus the insertion of a hyphen might cause effects to ripple arbitrarily
far into the rest of the word. A further complication arises if additional
hyphens appear together with such rippling, e.g., if the word in the
example just given could also be hyphenated between \.c and \.d; \TeX\
avoids this by simply ignoring the additional hyphens in such weird cases.]

Still further complications arise in the presence of ligatures that do not
delete the original characters. When punctuation precedes the word being
hyphenated, \TeX's method is not perfect under all possible scenarios,
because punctuation marks and letters can propagate information back and forth.
For example, suppose the original pre-hyphenation pair
\.[*a] changes to \.[*y] via a \.[\?=:] ligature, which changes to \.[xy]
via a \.[=:\?] ligature; if $p_[a-1]=\.x$ and $p_a=\.y$, the reconstitution
procedure isn't smart enough to obtain \.[xy] again. In such cases the
font designer should include a ligature that goes from \.[xa] to \.[xy]. }

{ 919. \[42] Hyphenation }

{tangle:pos tex.web:18072:19: }

{ When a word |hc[1..hn]| has been set up to contain a candidate for hyphenation,
\TeX\ first looks to see if it is in the user's exception dictionary. If not,
hyphens are inserted based on patterns that appear within the given word,
using an algorithm due to Frank~M. Liang.
\xref[Liang, Franklin Mark]

Let's consider Liang's method first, since it is much more interesting than the
exception-lookup routine.  The algorithm begins by setting |hyf[j]| to zero
for all |j|, and invalid characters are inserted into |hc[0]|
and |hc[hn+1]| to serve as delimiters. Then a reasonably fast method is
used to see which of a given set of patterns occurs in the word
|hc[0..(hn+1)]|. Each pattern $p_1\ldots p_k$ of length |k| has an associated
sequence of |k+1| numbers $n_0\ldots n_k$; and if the pattern occurs in
|hc[(j+1)..(j+k)]|, \TeX\ will set |hyf[j+i]:=max(hyf[j+i],$n_i$)| for
|0<=i<=k|. After this has been done for each pattern that occurs, a
discretionary hyphen will be inserted between |hc[j]| and |hc[j+1]| when
|hyf[j]| is odd, as we have already seen.

The set of patterns $p_1\ldots p_k$ and associated numbers $n_0\ldots n_k$
depends, of course, on the language whose words are being hyphenated, and
on the degree of hyphenation that is desired. A method for finding
appropriate |p|'s and |n|'s, from a given dictionary of words and acceptable
hyphenations, is discussed in Liang's Ph.D. thesis (Stanford University,
1983); \TeX\ simply starts with the patterns and works from there. }

{ 934. }

{tangle:pos tex.web:18253:1: }

{ We have now completed the hyphenation routine, so the |line_break| procedure
is finished at last. Since the hyphenation exception table is fresh in our
minds, it's a good time to deal with the routine that adds new entries to it.

When \TeX\ has scanned `\.[\\hyphenation]', it calls on a procedure named
|new_hyph_exceptions| to do the right thing. } procedure new_hyph_exceptions; {enters new exceptions}
label reswitch, exit, found, not_found;
var n:0..64; {length of current word; not always a |small_number|}
 j:0..64; {an index into |hc|}
 h:hyph_pointer; {an index into |hyph_word| and |hyph_list|}
 k:str_number; {an index into |str_start|}
 p:halfword ; {head of a list of hyphen positions}
 q:halfword ; {used when creating a new node for list |p|}
 s:str_number; {strings being compared or stored}
 u, v:pool_pointer; {indices into |str_pool|}
begin scan_left_brace; {a left brace must follow \.[\\hyphenation]}
if eqtb[int_base+ language_code].int  <=0 then cur_lang:=0 else if eqtb[int_base+ language_code].int  >255 then cur_lang:=0 else cur_lang:=eqtb[int_base+ language_code].int   ;

{ Enter as many hyphenation exceptions as are listed, until coming to a right brace; then |return| }
n:=0; p:=-{0xfffffff=}268435455  ;
 while true do    begin get_x_token;
  reswitch: case cur_cmd of
  letter,other_char,char_given:
{ Append a new letter or hyphen }
if cur_chr={"-"=}45 then 
{ Append the value |n| to list |p| }
begin if n<63 then
  begin q:=get_avail;  mem[ q].hh.rh :=p;  mem[ q].hh.lh :=n; p:=q;
  end;
end


else  begin if  eqtb[  lc_code_base+   cur_chr].hh.rh   =0 then
    begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Not a letter"=} 959); end ;
{ \xref[Not a letter] }
     begin help_ptr:=2; help_line[1]:={"Letters in \hyphenation words must have \lccode>0."=} 960; help_line[0]:={"Proceed; I'll ignore the character I just read."=} 961; end ;
    error;
    end
  else if n<63 then
    begin incr(n); hc[n]:= eqtb[  lc_code_base+   cur_chr].hh.rh   ;
    end;
  end

;
  char_num: begin scan_char_num; cur_chr:=cur_val; cur_cmd:=char_given;
    goto reswitch;
    end;
  spacer,right_brace: begin if n>1 then 
{ Enter a hyphenation exception }
begin incr(n); hc[n]:=cur_lang;  begin if pool_ptr+ n > pool_size then overflow({"pool size"=}257,pool_size-init_pool_ptr); { \xref[TeX capacity exceeded pool size][\quad pool size] } end ; h:=0;
for j:=1 to n do
  begin h:=(h+h+hc[j]) mod hyph_prime;
   begin str_pool[pool_ptr]:=   hc[  j] ; incr(pool_ptr); end ;
  end;
s:=make_string;

{ Insert the \(p)pair |(s,p)| into the exception table }
  if hyph_next <= hyph_prime then
     while (hyph_next>0) and (hyph_word[hyph_next-1]>0) do decr(hyph_next);
if (hyph_count=hyph_size)or(hyph_next=0) then
   overflow({"exception dictionary"=}962,hyph_size);
{ \xref[TeX capacity exceeded exception dictionary][\quad exception dictionary] }
incr(hyph_count);
while hyph_word[h]<>0 do
  begin 
{ If the string |hyph_word[h]| is less than \(or)or equal to |s|, interchange |(hyph_word[h],hyph_list[h])| with |(s,p)| }
{This is now a simple hash list, not an ordered one, so
the module title is no longer descriptive.}
k:=hyph_word[h];
if (str_start[ k+1]-str_start[ k]) <>(str_start[ s+1]-str_start[ s])  then goto not_found;
u:=str_start[k]; v:=str_start[s];
repeat if str_pool[u]<>str_pool[v] then goto not_found;
incr(u); incr(v);
until u=str_start[k+1];
{repeat hyphenation exception; flushing old data}
begin decr(str_ptr); pool_ptr:=str_start[str_ptr]; end ; s:=hyph_word[h]; {avoid |slow_make_string|!}
decr(hyph_count);
{ We could also |flush_list(hyph_list[h]);|, but it interferes
  with \.[trip.log]. }
goto found;
not_found:

;
  if hyph_link[h]=0 then
  begin
    hyph_link[h]:=hyph_next;
    if hyph_next >= hyph_size then hyph_next:=hyph_prime;
    if hyph_next > hyph_prime then incr(hyph_next);
  end;
  h:=hyph_link[h]-1;
  end;

found: hyph_word[h]:=s; hyph_list[h]:=p

;
end

;
    if cur_cmd=right_brace then  goto exit ;
    n:=0; p:=-{0xfffffff=}268435455  ;
    end;
   else  
{ Give improper \.[\\hyphenation] error }
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Improper "=} 690); end ; print_esc({"hyphenation"=}955);
{ \xref[Improper \\hyphenation...] }
  print({" will be flushed"=}956);
 begin help_ptr:=2; help_line[1]:={"Hyphenation exceptions must contain only letters"=} 957; help_line[0]:={"and hyphens. But continue; I'll forgive and forget."=} 958; end ;
error;
end


   end ;
  end

;
exit:end;



{ 967. \[44] Breaking vertical lists into pages }

{tangle:pos tex.web:18850:42: }

{ The |vsplit| procedure, which implements \TeX's \.[\\vsplit] operation,
is considerably simpler than |line_break| because it doesn't have to
worry about hyphenation, and because its mission is to discover a single
break instead of an optimum sequence of breakpoints.  But before we get
into the details of |vsplit|, we need to consider a few more basic things. }

{ 968. }

{tangle:pos tex.web:18856:1: }

{ A subroutine called |prune_page_top| takes a pointer to a vlist and
returns a pointer to a modified vlist in which all glue, kern, and penalty nodes
have been deleted before the first box or rule node. However, the first
box or rule is actually preceded by a newly created glue node designed so that
the topmost baseline will be at distance |split_top_skip| from the top,
whenever this is possible without backspacing.

In this routine and those that follow, we make use of the fact that a
vertical list contains no character nodes, hence the |type| field exists
for each node in the list.
\xref[data structure assumptions] } function prune_page_top( p:halfword ):halfword ; {adjust top after page break}
var prev_p:halfword ; {lags one step behind |p|}
 q:halfword ; {temporary variable for list manipulation}
begin prev_p:=mem_top-3 ;  mem[ mem_top-3 ].hh.rh :=p;
while p<>-{0xfffffff=}268435455   do
  case  mem[ p].hh.b0  of
  hlist_node,vlist_node,rule_node:
{ Insert glue for |split_top_skip| and set~|p:=null| }
begin q:=new_skip_param(split_top_skip_code);  mem[ prev_p].hh.rh :=q;  mem[ q].hh.rh :=p;
  {now |temp_ptr=glue_ptr(q)|}
if  mem[ temp_ptr+width_offset].int  > mem[ p+height_offset].int   then  mem[ temp_ptr+width_offset].int  := mem[ temp_ptr+width_offset].int  - mem[ p+height_offset].int  
else  mem[ temp_ptr+width_offset].int  :=0;
p:=-{0xfffffff=}268435455  ;
end

;
  whatsit_node,mark_node,ins_node: begin prev_p:=p; p:= mem[ prev_p].hh.rh ;
    end;
  glue_node,kern_node,penalty_node: begin q:=p; p:= mem[ q].hh.rh ;  mem[ q].hh.rh :=-{0xfffffff=}268435455  ;
     mem[ prev_p].hh.rh :=p; flush_node_list(q);
    end;
   else  confusion({"pruning"=}973)
{ \xref[this can't happen pruning][\quad pruning] }
   end ;
prune_page_top:= mem[ mem_top-3 ].hh.rh ;
end;



{ 970. }

{tangle:pos tex.web:18895:1: }

{ The next subroutine finds the best place to break a given vertical list
so as to obtain a box of height~|h|, with maximum depth~|d|.
A pointer to the beginning of the vertical list is given,
and a pointer to the optimum breakpoint is returned. The list is effectively
followed by a forced break, i.e., a penalty node with the |eject_penalty|;
if the best break occurs at this artificial node, the value |null| is returned.

An array of six |scaled| distances is used to keep track of the height
from the beginning of the list to the current place, just as in |line_break|.
In fact, we use one of the same arrays, only changing its name to reflect
its new significance. } function vert_break( p:halfword ;  h, d:scaled):halfword ;
  {finds optimum page break}
label done,not_found,update_heights;
var prev_p:halfword ; {if |p| is a glue node, |type(prev_p)| determines
  whether |p| is a legal breakpoint}
 q, r:halfword ; {glue specifications}
 pi:integer; {penalty value}
 b:integer; {badness at a trial breakpoint}
 least_cost:integer; {the smallest badness plus penalties found so far}
 best_place:halfword ; {the most recent break that leads to |least_cost|}
 prev_dp:scaled; {depth of previous box in the list}
 t:small_number; {|type| of the node following a kern}
begin prev_p:=p; {an initial glue node is not a legal breakpoint}
least_cost:={07777777777=}1073741823 ;  active_width [ 1]:=0 ; active_width [ 2]:=0 ; active_width [ 3]:=0 ; active_width [ 4]:=0 ; active_width [ 5]:=0 ; active_width [ 6]:=0  ; prev_dp:=0;
 while true do    begin 
{ If node |p| is a legal breakpoint, check if this break is the best known, and |goto done| if |p| is null or if the page-so-far is already too full to accept more stuff }
if p=-{0xfffffff=}268435455   then pi:=eject_penalty
else  
{ Use node |p| to update the current height and depth measurements; if this node is not a legal breakpoint, |goto not_found| or |update_heights|, otherwise set |pi| to the associated penalty at the break }
case  mem[ p].hh.b0  of
hlist_node,vlist_node,rule_node: begin{  } 

  active_width [1] :=active_width [1] +prev_dp+ mem[ p+height_offset].int  ; prev_dp:= mem[ p+depth_offset].int  ;
  goto not_found;
  end;
whatsit_node:
{ Process whatsit |p| in |vert_break| loop, |goto not_found| }
goto not_found

;
glue_node: if ( mem[  prev_p].hh.b0 <math_node)  then pi:=0
  else goto update_heights;
kern_node: begin if  mem[ p].hh.rh =-{0xfffffff=}268435455   then t:=penalty_node
  else t:= mem[  mem[  p].hh.rh ].hh.b0 ;
  if t=glue_node then pi:=0 else goto update_heights;
  end;
penalty_node: pi:= mem[ p+1].int ;
mark_node,ins_node: goto not_found;
 else  confusion({"vertbreak"=}974)
{ \xref[this can't happen vertbreak][\quad vertbreak] }
 end 

;

{ Check if node |p| is a new champion breakpoint; then \(go)|goto done| if |p| is a forced break or if the page-so-far is already too full }
if pi<inf_penalty then
  begin 
{ Compute the badness, |b|, using |awful_bad| if the box is too full }
if active_width [1] <h then
  if (active_width [3]<>0) or (active_width [4]<>0) or
    (active_width [5]<>0) then b:=0
  else b:=badness(h-active_width [1] ,active_width [2])
else if active_width [1] -h>active_width [6] then b:={07777777777=}1073741823 
else b:=badness(active_width [1] -h,active_width [6])

;
  if b<{07777777777=}1073741823  then
    if pi<=eject_penalty then b:=pi
    else if b<inf_bad then b:=b+pi
      else b:=100000 ;
  if b<=least_cost then
    begin best_place:=p; least_cost:=b;
    best_height_plus_depth:=active_width [1] +prev_dp;
    end;
  if (b={07777777777=}1073741823 )or(pi<=eject_penalty) then goto done;
  end

;
if ( mem[ p].hh.b0 <glue_node)or( mem[ p].hh.b0 >kern_node) then goto not_found;
update_heights: 
{ Update the current height and depth measurements with respect to a glue or kern node~|p| }
if  mem[ p].hh.b0 =kern_node then q:=p
else  begin q:=  mem[  p+ 1].hh.lh  ;
  active_width [2+  mem[ q].hh.b0 ]:= 
    active_width [2+  mem[ q].hh.b0 ]+ mem[ q+2].int  ;

  active_width [6]:=active_width [6]+ mem[ q+3].int  ;
  if (  mem[ q].hh.b1 <>normal)and( mem[ q+3].int  <>0) then
    begin{  } 

    begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Infinite glue shrinkage found in box being split"=} 975); end ;

{ \xref[Infinite glue shrinkage...] }
     begin help_ptr:=4; help_line[3]:={"The box you are \vsplitting contains some infinitely"=} 976; help_line[2]:={"shrinkable glue, e.g., `\vss' or `\vskip 0pt minus 1fil'."=} 977; help_line[1]:={"Such glue doesn't belong there; but you can safely proceed,"=} 978; help_line[0]:={"since the offensive shrinkability has been made finite."=} 936; end ;
    error; r:=new_spec(q);   mem[ r].hh.b1 :=normal; delete_glue_ref(q);
      mem[  p+ 1].hh.lh  :=r; q:=r;
    end;
  end;
active_width [1] :=active_width [1] +prev_dp+ mem[ q+width_offset].int  ; prev_dp:=0

;
not_found: if prev_dp>d then
    begin active_width [1] :=active_width [1] +prev_dp-d;
    prev_dp:=d;
    end;

;
  prev_p:=p; p:= mem[ prev_p].hh.rh ;
  end;
done: vert_break:=best_place;
end;



{ 977. }

{tangle:pos tex.web:19031:1: }

{ Now we are ready to consider |vsplit| itself. Most of
its work is accomplished by the two subroutines that we have just considered.

Given the number of a vlist box |n|, and given a desired page height |h|,
the |vsplit| function finds the best initial segment of the vlist and
returns a box for a page of height~|h|. The remainder of the vlist, if
any, replaces the original box, after removing glue and penalties and
adjusting for |split_top_skip|. Mark nodes in the split-off box are used to
set the values of |split_first_mark| and |split_bot_mark|; we use the
fact that |split_first_mark=null| if and only if |split_bot_mark=null|.

The original box becomes ``void'' if and only if it has been entirely
extracted.  The extracted box is ``void'' if and only if the original
box was void (or if it was, erroneously, an hlist box). } function vsplit( n:eight_bits;  h:scaled):halfword ;
  {extracts a page of height |h| from box |n|}
label exit,done;
var v:halfword ; {the box to be split}
p:halfword ; {runs through the vlist}
q:halfword ; {points to where the break occurs}
begin v:= eqtb[  box_base+   n].hh.rh   ;
if cur_mark[split_first_mark_code] <>-{0xfffffff=}268435455   then
  begin delete_token_ref(cur_mark[split_first_mark_code] ); cur_mark[split_first_mark_code] :=-{0xfffffff=}268435455  ;
  delete_token_ref(cur_mark[split_bot_mark_code] ); cur_mark[split_bot_mark_code] :=-{0xfffffff=}268435455  ;
  end;

{ Dispense with trivial cases of void or bad boxes }
if v=-{0xfffffff=}268435455   then
  begin vsplit:=-{0xfffffff=}268435455  ;  goto exit ;
  end;
if  mem[ v].hh.b0 <>vlist_node then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({""=} 335); end ; print_esc({"vsplit"=}979); print({" needs a "=}980);
  print_esc({"vbox"=}981);
{ \xref[vsplit_][\.[\\vsplit needs a \\vbox]] }
   begin help_ptr:=2; help_line[1]:={"The box you are trying to split is an \hbox."=} 982; help_line[0]:={"I can't split such a box, so I'll leave it alone."=} 983; end ;
  error; vsplit:=-{0xfffffff=}268435455  ;  goto exit ;
  end

;
q:=vert_break(  mem[  v+ list_offset].hh.rh  ,h,eqtb[dimen_base+ split_max_depth_code].int   );

{ Look at all the marks in nodes before the break, and set the final link to |null| at the break }
p:=  mem[  v+ list_offset].hh.rh  ;
if p=q then   mem[  v+ list_offset].hh.rh  :=-{0xfffffff=}268435455  
else  while true do  begin if  mem[ p].hh.b0 =mark_node then
    if cur_mark[split_first_mark_code] =-{0xfffffff=}268435455   then
      begin cur_mark[split_first_mark_code] :=mem[ p+1].int ;
      cur_mark[split_bot_mark_code] :=cur_mark[split_first_mark_code] ;
        mem[  cur_mark[split_first_mark_code] ].hh.lh  := 
          mem[  cur_mark[split_first_mark_code] ].hh.lh  +2;
      end
    else  begin delete_token_ref(cur_mark[split_bot_mark_code] );
      cur_mark[split_bot_mark_code] :=mem[ p+1].int ;
      incr(  mem[   cur_mark[split_bot_mark_code] ].hh.lh  ) ;
      end;
  if  mem[ p].hh.rh =q then
    begin  mem[ p].hh.rh :=-{0xfffffff=}268435455  ; goto done;
    end;
  p:= mem[ p].hh.rh ;
  end;
done:

;
q:=prune_page_top(q); p:=  mem[  v+ list_offset].hh.rh  ; free_node(v,box_node_size);
if q=-{0xfffffff=}268435455   then  eqtb[  box_base+   n].hh.rh   :=-{0xfffffff=}268435455   {the |eq_level| of the box stays the same}
else  eqtb[  box_base+   n].hh.rh   :=vpackage( q, 0,additional ,{07777777777=}1073741823 ) ;
vsplit:=vpackage(p,h,exactly,eqtb[dimen_base+ split_max_depth_code].int   );
exit: end;



{ 985. } procedure print_totals;
begin print_scaled(page_so_far[1] );
if page_so_far[ 2]<>0 then begin print({" plus "=}310); print_scaled(page_so_far[ 2]); print({""=} 335); end ;
if page_so_far[ 3]<>0 then begin print({" plus "=}310); print_scaled(page_so_far[ 3]); print({"fil"=} 309); end ;
if page_so_far[ 4]<>0 then begin print({" plus "=}310); print_scaled(page_so_far[ 4]); print({"fill"=} 992); end ;
if page_so_far[ 5]<>0 then begin print({" plus "=}310); print_scaled(page_so_far[ 5]); print({"filll"=} 993); end ;
if page_so_far[6] <>0 then
  begin print({" minus "=}311); print_scaled(page_so_far[6] );
  end;
end;



{ 987. }

{tangle:pos tex.web:19313:1: }

{ Here is a procedure that is called when the |page_contents| is changing
from |empty| to |inserts_only| or |box_there|. } procedure freeze_page_specs( s:small_number);
begin page_contents:=s;
page_so_far[0] :=eqtb[dimen_base+ vsize_code].int   ; page_max_depth:=eqtb[dimen_base+ max_depth_code].int   ;
page_so_far[7] :=0;  page_so_far[ 1]:=0 ; page_so_far[ 2]:=0 ; page_so_far[ 3]:=0 ; page_so_far[ 4]:=0 ; page_so_far[ 5]:=0 ; page_so_far[ 6]:=0  ;
least_page_cost:={07777777777=}1073741823 ;
 ifdef('STAT')  if eqtb[int_base+ tracing_pages_code].int  >0 then
  begin begin_diagnostic;
  print_nl({"%% goal height="=}1001); print_scaled(page_so_far[0] );
{ \xref[goal height] }
  print({", max depth="=}1002); print_scaled(page_max_depth);
  end_diagnostic(false);
  end;  endif('STAT')  

end;



{ 992. }

{tangle:pos tex.web:19377:1: }

{ At certain times box 255 is supposed to be void (i.e., |null|),
or an insertion box is supposed to be ready to accept a vertical list.
If not, an error message is printed, and the following subroutine
flushes the unwanted contents, reporting them to the user. } procedure box_error( n:eight_bits);
begin error; begin_diagnostic;
print_nl({"The following box has been deleted:"=}849);
{ \xref[The following...deleted] }
show_box( eqtb[  box_base+   n].hh.rh   ); end_diagnostic(true);
flush_node_list( eqtb[  box_base+   n].hh.rh   );  eqtb[  box_base+   n].hh.rh   :=-{0xfffffff=}268435455  ;
end;



{ 993. }

{tangle:pos tex.web:19390:1: }

{ The following procedure guarantees that a given box register
does not contain an \.[\\hbox]. } procedure ensure_vbox( n:eight_bits);
var p:halfword ; {the box register contents}
begin p:= eqtb[  box_base+   n].hh.rh   ;
if p<>-{0xfffffff=}268435455   then if  mem[ p].hh.b0 =hlist_node then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Insertions can only be added to a vbox"=} 1003); end ;
{ \xref[Insertions can only...] }
   begin help_ptr:=3; help_line[2]:={"Tut tut: You're trying to \insert into a"=} 1004; help_line[1]:={"\box register that now contains an \hbox."=} 1005; help_line[0]:={"Proceed, and I'll discard its present contents."=} 1006; end ;
  box_error(n);
  end;
end;



{ 994. }

{tangle:pos tex.web:19406:1: }

{ \TeX\ is not always in vertical mode at the time |build_page|
is called; the current mode reflects what \TeX\ should return to, after
the contribution list has been emptied. A call on |build_page| should
be immediately followed by `|goto big_switch|', which is \TeX's central
control point. } { \4 }
{ Declare the procedure called |fire_up| }
procedure fire_up( c:halfword );
label exit;
var p, q, r, s:halfword ; {nodes being examined and/or changed}
 prev_p:halfword ; {predecessor of |p|}
 n:min_quarterword..255; {insertion box number}
 wait:boolean; {should the present insertion be held over?}
 save_vbadness:integer; {saved value of |vbadness|}
 save_vfuzz: scaled; {saved value of |vfuzz|}
 save_split_top_skip: halfword ; {saved value of |split_top_skip|}
begin 
{ Set the value of |output_penalty| }
if  mem[ best_page_break].hh.b0 =penalty_node then
  begin geq_word_define(int_base+output_penalty_code, mem[ best_page_break+1].int );
   mem[ best_page_break+1].int :=inf_penalty;
  end
else geq_word_define(int_base+output_penalty_code,inf_penalty)

;
if cur_mark[bot_mark_code] <>-{0xfffffff=}268435455   then
  begin if cur_mark[top_mark_code] <>-{0xfffffff=}268435455   then delete_token_ref(cur_mark[top_mark_code] );
  cur_mark[top_mark_code] :=cur_mark[bot_mark_code] ; incr(  mem[   cur_mark[top_mark_code] ].hh.lh  ) ;
  delete_token_ref(cur_mark[first_mark_code] ); cur_mark[first_mark_code] :=-{0xfffffff=}268435455  ;
  end;

{ Put the \(o)optimal current page into box 255, update |first_mark| and |bot_mark|, append insertions to their boxes, and put the remaining nodes back on the contribution list }
if c=best_page_break then best_page_break:=-{0xfffffff=}268435455  ; {|c| not yet linked in}

{ Ensure that box 255 is empty before output }
if  eqtb[  box_base+   255].hh.rh   <>-{0xfffffff=}268435455   then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({""=} 335); end ; print_esc({"box"=}414); print({"255 is not void"=}1017);
{ \xref[box255][\.[\\box255 is not void]] }
   begin help_ptr:=2; help_line[1]:={"You shouldn't use \box255 except in \output routines."=} 1018; help_line[0]:={"Proceed, and I'll discard its present contents."=} 1006; end ;
  box_error(255);
  end

;
insert_penalties:=0; {this will count the number of insertions held over}
save_split_top_skip:= eqtb[  glue_base+   split_top_skip_code].hh.rh    ;
if eqtb[int_base+ holding_inserts_code].int  <=0 then
  
{ Prepare all the boxes involved in insertions to act as queues }
begin r:= mem[ mem_top ].hh.rh ;
while r<>mem_top  do
  begin if  mem[  r+ 2].hh.lh  <>-{0xfffffff=}268435455   then
    begin n:=  mem[  r].hh.b1  ; ensure_vbox(n);
    if  eqtb[  box_base+   n].hh.rh   =-{0xfffffff=}268435455   then  eqtb[  box_base+   n].hh.rh   :=new_null_box;
    p:= eqtb[  box_base+   n].hh.rh   +list_offset;
    while  mem[ p].hh.rh <>-{0xfffffff=}268435455   do p:= mem[ p].hh.rh ;
     mem[  r+ 2].hh.rh  :=p;
    end;
  r:= mem[ r].hh.rh ;
  end;
end

;
q:=mem_top-4 ;  mem[ q].hh.rh :=-{0xfffffff=}268435455  ; prev_p:=mem_top-2 ; p:= mem[ prev_p].hh.rh ;
while p<>best_page_break do
  begin if  mem[ p].hh.b0 =ins_node then
    begin if eqtb[int_base+ holding_inserts_code].int  <=0 then
       
{ Either insert the material specified by node |p| into the appropriate box, or hold it for the next page; also delete node |p| from the current page }
begin r:= mem[ mem_top ].hh.rh ;
while  mem[ r].hh.b1 <> mem[ p].hh.b1  do r:= mem[ r].hh.rh ;
if  mem[  r+ 2].hh.lh  =-{0xfffffff=}268435455   then wait:=true
else  begin wait:=false; s:= mem[  r+ 2].hh.rh  ;  mem[ s].hh.rh := mem[  p+ 4].hh.lh  ;
  if  mem[  r+ 2].hh.lh  =p then
    
{ Wrap up the box specified by node |r|, splitting node |p| if called for; set |wait:=true| if node |p| holds a remainder after splitting }
begin if  mem[ r].hh.b0 =split_up then
  if ( mem[  r+ 1].hh.lh  =p)and( mem[  r+ 1].hh.rh  <>-{0xfffffff=}268435455  ) then
    begin while  mem[ s].hh.rh <> mem[  r+ 1].hh.rh   do s:= mem[ s].hh.rh ;
     mem[ s].hh.rh :=-{0xfffffff=}268435455  ;
     eqtb[  glue_base+   split_top_skip_code].hh.rh    := mem[  p+ 4].hh.rh  ;
     mem[  p+ 4].hh.lh  :=prune_page_top( mem[  r+ 1].hh.rh  );
    if  mem[  p+ 4].hh.lh  <>-{0xfffffff=}268435455   then
      begin temp_ptr:=vpackage(  mem[   p+ 4].hh.lh  , 0,additional ,{07777777777=}1073741823 ) ;
       mem[ p+height_offset].int  := mem[ temp_ptr+height_offset].int  + mem[ temp_ptr+depth_offset].int  ;
      free_node(temp_ptr,box_node_size); wait:=true;
      end;
    end;
 mem[  r+ 2].hh.lh  :=-{0xfffffff=}268435455  ;
n:=  mem[  r].hh.b1  ;
temp_ptr:=  mem[   eqtb[  box_base+     n].hh.rh   + list_offset].hh.rh  ;
free_node( eqtb[  box_base+   n].hh.rh   ,box_node_size);
 eqtb[  box_base+   n].hh.rh   :=vpackage( temp_ptr, 0,additional ,{07777777777=}1073741823 ) ;
end


  else  begin while  mem[ s].hh.rh <>-{0xfffffff=}268435455   do s:= mem[ s].hh.rh ;
     mem[  r+ 2].hh.rh  :=s;
    end;
  end;

{ Either append the insertion node |p| after node |q|, and remove it from the current page, or delete |node(p)| }
 mem[ prev_p].hh.rh := mem[ p].hh.rh ;  mem[ p].hh.rh :=-{0xfffffff=}268435455  ;
if wait then
  begin  mem[ q].hh.rh :=p; q:=p; incr(insert_penalties);
  end
else  begin delete_glue_ref( mem[  p+ 4].hh.rh  );
  free_node(p,ins_node_size);
  end;
p:=prev_p

;
end

;
    end
  else if  mem[ p].hh.b0 =mark_node then 
{ Update the values of |first_mark| and |bot_mark| }
begin if cur_mark[first_mark_code] =-{0xfffffff=}268435455   then
  begin cur_mark[first_mark_code] :=mem[ p+1].int ;
  incr(  mem[   cur_mark[first_mark_code] ].hh.lh  ) ;
  end;
if cur_mark[bot_mark_code] <>-{0xfffffff=}268435455   then delete_token_ref(cur_mark[bot_mark_code] );
cur_mark[bot_mark_code] :=mem[ p+1].int ; incr(  mem[   cur_mark[bot_mark_code] ].hh.lh  ) ;
end

;
  prev_p:=p; p:= mem[ prev_p].hh.rh ;
  end;
 eqtb[  glue_base+   split_top_skip_code].hh.rh    :=save_split_top_skip;

{ Break the current page at node |p|, put it in box~255, and put the remaining nodes on the contribution list }
if p<>-{0xfffffff=}268435455   then
  begin if  mem[ mem_top-1 ].hh.rh =-{0xfffffff=}268435455   then
    if nest_ptr=0 then cur_list.tail_field :=page_tail
    else nest[0].tail_field :=page_tail;
   mem[ page_tail].hh.rh := mem[ mem_top-1 ].hh.rh ;
   mem[ mem_top-1 ].hh.rh :=p;
   mem[ prev_p].hh.rh :=-{0xfffffff=}268435455  ;
  end;
save_vbadness:=eqtb[int_base+ vbadness_code].int  ; eqtb[int_base+ vbadness_code].int  :=inf_bad;
save_vfuzz:=eqtb[dimen_base+ vfuzz_code].int   ; eqtb[dimen_base+ vfuzz_code].int   :={07777777777=}1073741823 ; {inhibit error messages}
 eqtb[  box_base+   255].hh.rh   :=vpackage( mem[ mem_top-2 ].hh.rh ,best_size,exactly,page_max_depth);
eqtb[int_base+ vbadness_code].int  :=save_vbadness; eqtb[dimen_base+ vfuzz_code].int   :=save_vfuzz;
if last_glue<>{0xfffffff=}268435455  then delete_glue_ref(last_glue);

{ Start a new current page }
page_contents:=empty; page_tail:=mem_top-2 ;  mem[ mem_top-2 ].hh.rh :=-{0xfffffff=}268435455  ;

last_glue:={0xfffffff=}268435455 ; last_penalty:=0; last_kern:=0;
page_so_far[7] :=0; page_max_depth:=0

; {this sets |last_glue:=max_halfword|}
if q<>mem_top-4  then
  begin  mem[ mem_top-2 ].hh.rh := mem[ mem_top-4 ].hh.rh ; page_tail:=q;
  end

;

{ Delete \(t)the page-insertion nodes }
r:= mem[ mem_top ].hh.rh ;
while r<>mem_top  do
  begin q:= mem[ r].hh.rh ; free_node(r,page_ins_node_size); r:=q;
  end;
 mem[ mem_top ].hh.rh :=mem_top 



;
if (cur_mark[top_mark_code] <>-{0xfffffff=}268435455  )and(cur_mark[first_mark_code] =-{0xfffffff=}268435455  ) then
  begin cur_mark[first_mark_code] :=cur_mark[top_mark_code] ; incr(  mem[   cur_mark[top_mark_code] ].hh.lh  ) ;
  end;
if  eqtb[  output_routine_loc].hh.rh   <>-{0xfffffff=}268435455   then
  if dead_cycles>=eqtb[int_base+ max_dead_cycles_code].int   then
    
{ Explain that too many dead cycles have occurred in a row }
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Output loop---"=} 1019); end ; print_int(dead_cycles);
{ \xref[Output loop...] }
print({" consecutive dead cycles"=}1020);
 begin help_ptr:=3; help_line[2]:={"I've concluded that your \output is awry; it never does a"=} 1021; help_line[1]:={"\shipout, so I'm shipping \box255 out myself. Next time"=} 1022; help_line[0]:={"increase \maxdeadcycles if you want me to be more patient!"=} 1023; end ; error;
end


  else 
{ Fire up the user's output routine and |return| }
begin output_active:=true;
incr(dead_cycles);
push_nest; cur_list.mode_field :=-vmode; cur_list.aux_field .int  :=-65536000 ; cur_list.ml_field :=-line;
begin_token_list( eqtb[  output_routine_loc].hh.rh   ,output_text);
new_save_level(output_group); normal_paragraph;
scan_left_brace;
 goto exit ;
end

;

{ Perform the default output routine }
begin if  mem[ mem_top-2 ].hh.rh <>-{0xfffffff=}268435455   then
  begin if  mem[ mem_top-1 ].hh.rh =-{0xfffffff=}268435455   then
    if nest_ptr=0 then cur_list.tail_field :=page_tail else nest[0].tail_field :=page_tail
  else  mem[ page_tail].hh.rh := mem[ mem_top-1 ].hh.rh ;
   mem[ mem_top-1 ].hh.rh := mem[ mem_top-2 ].hh.rh ;
   mem[ mem_top-2 ].hh.rh :=-{0xfffffff=}268435455  ; page_tail:=mem_top-2 ;
  end;
ship_out( eqtb[  box_base+   255].hh.rh   );  eqtb[  box_base+   255].hh.rh   :=-{0xfffffff=}268435455  ;
end

;
exit:end;

 

procedure build_page; {append contributions to the current page}
label exit,done,done1,continue,contribute,update_heights;
var p:halfword ; {the node being appended}
 q, r:halfword ; {nodes being examined}
 b, c:integer; {badness and cost of current page}
 pi:integer; {penalty to be added to the badness}
 n:min_quarterword..255; {insertion box number}
 delta, h, w:scaled; {sizes used for insertion calculations}
begin if ( mem[ mem_top-1 ].hh.rh =-{0xfffffff=}268435455  )or output_active then  goto exit ;
repeat continue: p:= mem[ mem_top-1 ].hh.rh ;


{ Update the values of |last_glue|, |last_penalty|, and |last_kern| }
if last_glue<>{0xfffffff=}268435455  then delete_glue_ref(last_glue);
last_penalty:=0; last_kern:=0;
if  mem[ p].hh.b0 =glue_node then
  begin last_glue:=  mem[  p+ 1].hh.lh  ; incr(  mem[   last_glue].hh.rh  ) ;
  end
else  begin last_glue:={0xfffffff=}268435455 ;
  if  mem[ p].hh.b0 =penalty_node then last_penalty:= mem[ p+1].int 
  else if  mem[ p].hh.b0 =kern_node then last_kern:= mem[ p+width_offset].int  ;
  end

;

{ Move node |p| to the current page; if it is time for a page break, put the nodes following the break back onto the contribution list, and |return| to the user's output routine if there is one }

{ If the current page is empty and node |p| is to be deleted, |goto done1|; otherwise use node |p| to update the state of the current page; if this node is an insertion, |goto contribute|; otherwise if this node is not a legal breakpoint, |goto contribute| or |update_heights|; otherwise set |pi| to the penalty associated with this breakpoint }
case  mem[ p].hh.b0  of
hlist_node,vlist_node,rule_node: if page_contents<box_there then
    
{ Initialize the current page, insert the \.[\\topskip] glue ahead of |p|, and |goto continue| }
begin if page_contents=empty then freeze_page_specs(box_there)
else page_contents:=box_there;
q:=new_skip_param(top_skip_code); {now |temp_ptr=glue_ptr(q)|}
if  mem[ temp_ptr+width_offset].int  > mem[ p+height_offset].int   then  mem[ temp_ptr+width_offset].int  := mem[ temp_ptr+width_offset].int  - mem[ p+height_offset].int  
else  mem[ temp_ptr+width_offset].int  :=0;
 mem[ q].hh.rh :=p;  mem[ mem_top-1 ].hh.rh :=q; goto continue;
end


  else 
{ Prepare to move a box or rule node to the current page, then |goto contribute| }
begin page_so_far[1] :=page_so_far[1] +page_so_far[7] + mem[ p+height_offset].int  ;
page_so_far[7] := mem[ p+depth_offset].int  ;
goto contribute;
end

;
whatsit_node: 
{ Prepare to move whatsit |p| to the current page, then |goto contribute| }
goto contribute

;
glue_node: if page_contents<box_there then goto done1
  else if ( mem[  page_tail].hh.b0 <math_node)  then pi:=0
  else goto update_heights;
kern_node: if page_contents<box_there then goto done1
  else if  mem[ p].hh.rh =-{0xfffffff=}268435455   then  goto exit 
  else if  mem[  mem[  p].hh.rh ].hh.b0 =glue_node then pi:=0
  else goto update_heights;
penalty_node: if page_contents<box_there then goto done1 else pi:= mem[ p+1].int ;
mark_node: goto contribute;
ins_node: 
{ Append an insertion to the current page and |goto contribute| }
begin if page_contents=empty then freeze_page_specs(inserts_only);
n:= mem[ p].hh.b1 ; r:=mem_top ;
while n>= mem[  mem[  r].hh.rh ].hh.b1  do r:= mem[ r].hh.rh ;
n:= n ;
if  mem[ r].hh.b1 <> n  then
  
{ Create a page insertion node with |subtype(r)=qi(n)|, and include the glue correction for box |n| in the current page state }
begin q:=get_node(page_ins_node_size);  mem[ q].hh.rh := mem[ r].hh.rh ;  mem[ r].hh.rh :=q; r:=q;
 mem[ r].hh.b1 := n ;  mem[ r].hh.b0 :=inserting; ensure_vbox(n);
if  eqtb[  box_base+   n].hh.rh   =-{0xfffffff=}268435455   then  mem[ r+height_offset].int  :=0
else  mem[ r+height_offset].int  := mem[  eqtb[  box_base+    n].hh.rh   +height_offset].int  + mem[  eqtb[  box_base+    n].hh.rh   +depth_offset].int  ;
 mem[  r+ 2].hh.lh  :=-{0xfffffff=}268435455  ;

q:= eqtb[  skip_base+   n].hh.rh   ;
if eqtb[count_base+ n].int =1000 then h:= mem[ r+height_offset].int  
else h:=x_over_n( mem[ r+height_offset].int  ,1000)*eqtb[count_base+ n].int ;
page_so_far[0] :=page_so_far[0] -h- mem[ q+width_offset].int  ;

page_so_far[2+  mem[ q].hh.b0 ]:= page_so_far[2+  mem[ q].hh.b0 ]+ mem[ q+2].int  ;

page_so_far[6] :=page_so_far[6] + mem[ q+3].int  ;
if (  mem[ q].hh.b1 <>normal)and( mem[ q+3].int  <>0) then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Infinite glue shrinkage inserted from "=} 1012); end ; print_esc({"skip"=}400);
{ \xref[Infinite glue shrinkage...] }
  print_int(n);
   begin help_ptr:=3; help_line[2]:={"The correction glue for page breaking with insertions"=} 1013; help_line[1]:={"must have finite shrinkability. But you may proceed,"=} 1014; help_line[0]:={"since the offensive shrinkability has been made finite."=} 936; end ;
  error;
  end;
end

;
if  mem[ r].hh.b0 =split_up then insert_penalties:=insert_penalties+mem[ p+1].int 
else  begin  mem[  r+ 2].hh.rh  :=p;
  delta:=page_so_far[0] -page_so_far[1] -page_so_far[7] +page_so_far[6] ;
    {this much room is left if we shrink the maximum}
  if eqtb[count_base+ n].int =1000 then h:= mem[ p+height_offset].int  
  else h:=x_over_n( mem[ p+height_offset].int  ,1000)*eqtb[count_base+ n].int ; {this much room is needed}
  if ((h<=0)or(h<=delta))and( mem[ p+height_offset].int  + mem[ r+height_offset].int  <=eqtb[scaled_base+ n].int  ) then
    begin page_so_far[0] :=page_so_far[0] -h;  mem[ r+height_offset].int  := mem[ r+height_offset].int  + mem[ p+height_offset].int  ;
    end
  else 
{ Find the best way to split the insertion, and change |type(r)| to |split_up| }
begin if eqtb[count_base+ n].int <=0 then w:={07777777777=}1073741823 
else  begin w:=page_so_far[0] -page_so_far[1] -page_so_far[7] ;
  if eqtb[count_base+ n].int <>1000 then w:=x_over_n(w,eqtb[count_base+ n].int )*1000;
  end;
if w>eqtb[scaled_base+ n].int  - mem[ r+height_offset].int   then w:=eqtb[scaled_base+ n].int  - mem[ r+height_offset].int  ;
q:=vert_break( mem[  p+ 4].hh.lh  ,w, mem[ p+depth_offset].int  );
 mem[ r+height_offset].int  := mem[ r+height_offset].int  +best_height_plus_depth;
 ifdef('STAT')  if eqtb[int_base+ tracing_pages_code].int  >0 then 
{ Display the insertion split cost }
begin begin_diagnostic; print_nl({"% split"=}1015); print_int(n);
{ \xref[split] }
print({" to "=}1016); print_scaled(w);
print_char({","=}44); print_scaled(best_height_plus_depth);

print({" p="=}945);
if q=-{0xfffffff=}268435455   then print_int(eject_penalty)
else if  mem[ q].hh.b0 =penalty_node then print_int( mem[ q+1].int )
else print_char({"0"=}48);
end_diagnostic(false);
end

; endif('STAT')  

if eqtb[count_base+ n].int <>1000 then
  best_height_plus_depth:=x_over_n(best_height_plus_depth,1000)*eqtb[count_base+ n].int ;
page_so_far[0] :=page_so_far[0] -best_height_plus_depth;
 mem[ r].hh.b0 :=split_up;  mem[  r+ 1].hh.rh  :=q;  mem[  r+ 1].hh.lh  :=p;
if q=-{0xfffffff=}268435455   then insert_penalties:=insert_penalties+eject_penalty
else if  mem[ q].hh.b0 =penalty_node then insert_penalties:=insert_penalties+ mem[ q+1].int ;
end

;
  end;
goto contribute;
end

;
 else  confusion({"page"=}1007)
{ \xref[this can't happen page][\quad page] }
 end 

;

{ Check if node |p| is a new champion breakpoint; then \(if)if it is time for a page break, prepare for output, and either fire up the user's output routine and |return| or ship out the page and |goto done| }
if pi<inf_penalty then
  begin 
{ Compute the badness, |b|, of the current page, using |awful_bad| if the box is too full }
if page_so_far[1] <page_so_far[0]  then
  if (page_so_far[3]<>0) or (page_so_far[4]<>0) or 
    (page_so_far[5]<>0) then b:=0
  else b:=badness(page_so_far[0] -page_so_far[1] ,page_so_far[2])
else if page_so_far[1] -page_so_far[0] >page_so_far[6]  then b:={07777777777=}1073741823 
else b:=badness(page_so_far[1] -page_so_far[0] ,page_so_far[6] )

;
  if b<{07777777777=}1073741823  then
    if pi<=eject_penalty then c:=pi
    else  if b<inf_bad then c:=b+pi+insert_penalties
      else c:=100000 
  else c:=b;
  if insert_penalties>=10000 then c:={07777777777=}1073741823 ;
   ifdef('STAT')  if eqtb[int_base+ tracing_pages_code].int  >0 then 
{ Display the page break cost }
begin begin_diagnostic; print_nl({"%"=}37);
print({" t="=}941); print_totals;

print({" g="=}1010); print_scaled(page_so_far[0] );

print({" b="=}944);
if b={07777777777=}1073741823  then print_char({"*"=}42) else print_int(b);
{ \xref[*\relax] }
print({" p="=}945); print_int(pi);
print({" c="=}1011);
if c={07777777777=}1073741823  then print_char({"*"=}42) else print_int(c);
if c<=least_page_cost then print_char({"#"=}35);
end_diagnostic(false);
end

; endif('STAT')  

  if c<=least_page_cost then
    begin best_page_break:=p; best_size:=page_so_far[0] ;
    least_page_cost:=c;
    r:= mem[ mem_top ].hh.rh ;
    while r<>mem_top  do
      begin  mem[  r+ 2].hh.lh  := mem[  r+ 2].hh.rh  ;
      r:= mem[ r].hh.rh ;
      end;
    end;
  if (c={07777777777=}1073741823 )or(pi<=eject_penalty) then
    begin fire_up(p); {output the current page at the best place}
    if output_active then  goto exit ; {user's output routine will act}
    goto done; {the page has been shipped out by default output routine}
    end;
  end

;
if ( mem[ p].hh.b0 <glue_node)or( mem[ p].hh.b0 >kern_node) then goto contribute;
update_heights:
{ Update the current page measurements with respect to the glue or kern specified by node~|p| }
if  mem[ p].hh.b0 =kern_node then q:=p
else begin q:=  mem[  p+ 1].hh.lh  ;
  page_so_far[2+  mem[ q].hh.b0 ]:= 
    page_so_far[2+  mem[ q].hh.b0 ]+ mem[ q+2].int  ;

  page_so_far[6] :=page_so_far[6] + mem[ q+3].int  ;
  if (  mem[ q].hh.b1 <>normal)and( mem[ q+3].int  <>0) then
    begin{  } 

    begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Infinite glue shrinkage found on current page"=} 1008); end ;

{ \xref[Infinite glue shrinkage...] }
     begin help_ptr:=4; help_line[3]:={"The page about to be output contains some infinitely"=} 1009; help_line[2]:={"shrinkable glue, e.g., `\vss' or `\vskip 0pt minus 1fil'."=} 977; help_line[1]:={"Such glue doesn't belong there; but you can safely proceed,"=} 978; help_line[0]:={"since the offensive shrinkability has been made finite."=} 936; end ;
    error;
    r:=new_spec(q);   mem[ r].hh.b1 :=normal; delete_glue_ref(q);
      mem[  p+ 1].hh.lh  :=r; q:=r;
    end;
  end;
page_so_far[1] :=page_so_far[1] +page_so_far[7] + mem[ q+width_offset].int  ; page_so_far[7] :=0

;
contribute: 
{ Make sure that |page_max_depth| is not exceeded }
if page_so_far[7] >page_max_depth then
  begin page_so_far[1] := 
    page_so_far[1] +page_so_far[7] -page_max_depth;

  page_so_far[7] :=page_max_depth;
  end;

;

{ Link node |p| into the current page and |goto done| }
 mem[ page_tail].hh.rh :=p; page_tail:=p;
 mem[ mem_top-1 ].hh.rh := mem[ p].hh.rh ;  mem[ p].hh.rh :=-{0xfffffff=}268435455  ; goto done

;
done1:
{ Recycle node |p| }
 mem[ mem_top-1 ].hh.rh := mem[ p].hh.rh ;  mem[ p].hh.rh :=-{0xfffffff=}268435455  ; flush_node_list(p)

;
done:

;
until  mem[ mem_top-1 ].hh.rh =-{0xfffffff=}268435455  ;

{ Make the contribution list empty by setting its tail to |contrib_head| }
if nest_ptr=0 then cur_list.tail_field :=mem_top-1  {vertical mode}
else nest[0].tail_field :=mem_top-1  {other modes}

;
exit:end;



{ 1029. \[46] The chief executive }

{tangle:pos tex.web:19977:26: }

{ We come now to the |main_control| routine, which contains the master
switch that causes all the various pieces of \TeX\ to do their things,
in the right order.

In a sense, this is the grand climax of the program: It applies all the
tools that we have worked so hard to construct. In another sense, this is
the messiest part of the program: It necessarily refers to other pieces
of code all over the place, so that a person can't fully understand what is
going on without paging back and forth to be reminded of conventions that
are defined elsewhere. We are now at the hub of the web, the central nervous
system that touches most of the other parts and ties them together.
\xref[brain]

The structure of |main_control| itself is quite simple. There's a label
called |big_switch|, at which point the next token of input is fetched
using |get_x_token|. Then the program branches at high speed into one of
about 100 possible directions, based on the value of the current
mode and the newly fetched command code; the sum |abs(mode)+cur_cmd|
indicates what to do next. For example, the case `|vmode+letter|' arises
when a letter occurs in vertical mode (or internal vertical mode); this
case leads to instructions that initialize a new paragraph and enter
horizontal mode.

The big |case| statement that contains this multiway switch has been labeled
|reswitch|, so that the program can |goto reswitch| when the next token
has already been fetched. Most of the cases are quite short; they call
an ``action procedure'' that does the work for that case, and then they
either |goto reswitch| or they ``fall through'' to the end of the |case|
statement, which returns control back to |big_switch|. Thus, |main_control|
is not an extremely large procedure, in spite of the multiplicity of things
it must do; it is small enough to be handled by \PASCAL\ compilers that put
severe restrictions on procedure size.
 \xref[action procedure]

One case is singled out for special treatment, because it accounts for most
of \TeX's activities in typical applications. The process of reading simple
text and converting it into |char_node| records, while looking for ligatures
and kerns, is part of \TeX's ``inner loop''; the whole program runs
efficiently when its inner loop is fast, so this part has been written
with particular care. }

{ 1030. }

{tangle:pos tex.web:20017:22: }

{ We shall concentrate first on the inner loop of |main_control|, deferring
consideration of the other cases until later.
\xref[inner loop] } { \4 }
{ Declare action procedures for use by |main_control| }
procedure app_space; {handle spaces when |space_factor<>1000|}
var q:halfword ; {glue node}
begin if (cur_list.aux_field .hh.lh >=2000)and( eqtb[  glue_base+   xspace_skip_code].hh.rh    <>mem_bot ) then
  q:=new_param_glue(xspace_skip_code)
else  begin if  eqtb[  glue_base+   space_skip_code].hh.rh    <>mem_bot  then main_p:= eqtb[  glue_base+   space_skip_code].hh.rh    
  else 
{ Find the glue specification... }
begin main_p:=font_glue[ eqtb[  cur_font_loc].hh.rh   ];
if main_p=-{0xfffffff=}268435455   then
  begin main_p:=new_spec(mem_bot ); main_k:=param_base[ eqtb[  cur_font_loc].hh.rh   ]+space_code;
   mem[ main_p+width_offset].int  :=font_info[main_k].int ; {that's |space(cur_font)|}
   mem[ main_p+2].int  :=font_info[main_k+1].int ; {and |space_stretch(cur_font)|}
   mem[ main_p+3].int  :=font_info[main_k+2].int ; {and |space_shrink(cur_font)|}
  font_glue[ eqtb[  cur_font_loc].hh.rh   ]:=main_p;
  end;
end

;
  main_p:=new_spec(main_p);
  
{ Modify the glue specification in |main_p| according to the space factor }
if cur_list.aux_field .hh.lh >=2000 then  mem[ main_p+width_offset].int  := mem[ main_p+width_offset].int  +font_info[ extra_space_code+param_base[  eqtb[  cur_font_loc].hh.rh   ]].int  ;
 mem[ main_p+2].int  :=xn_over_d( mem[ main_p+2].int  ,cur_list.aux_field .hh.lh ,1000);
 mem[ main_p+3].int  :=xn_over_d( mem[ main_p+3].int  ,1000,cur_list.aux_field .hh.lh )

;
  q:=new_glue(main_p);   mem[  main_p].hh.rh  :=-{0xfffffff=}268435455  ;
  end;
 mem[ cur_list.tail_field ].hh.rh :=q; cur_list.tail_field :=q;
end;


procedure insert_dollar_sign;
begin back_input; cur_tok:=math_shift_token+{"$"=}36;
begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Missing $ inserted"=} 1031); end ;
{ \xref[Missing \$ inserted] }
 begin help_ptr:=2; help_line[1]:={"I've inserted a begin-math/end-math symbol since I think"=} 1032; help_line[0]:={"you left one out. Proceed, with fingers crossed."=} 1033; end ; ins_error;
end;


procedure you_cant;
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"You can't use `"=} 695); end ;
{ \xref[You can't use x in y mode] }
print_cmd_chr(cur_cmd,cur_chr);
print_in_mode(cur_list.mode_field );
end;


procedure report_illegal_case;
begin you_cant;
 begin help_ptr:=4; help_line[3]:={"Sorry, but I'm not programmed to handle this case;"=} 1034; help_line[2]:={"I'll just pretend that you didn't ask for it."=} 1035; help_line[1]:={"If you're in the wrong mode, you might be able to"=} 1036; help_line[0]:={"return to the right one by typing `I]' or `I$' or `I\par'."=} 1037; end ;

error;
end;


function privileged:boolean;
begin if cur_list.mode_field >0 then privileged:=true
else  begin report_illegal_case; privileged:=false;
  end;
end;


function its_all_over:boolean; {do this when \.[\\end] or \.[\\dump] occurs}
label exit;
begin if privileged then
  begin if (mem_top-2 =page_tail)and(cur_list.head_field =cur_list.tail_field )and(dead_cycles=0) then
    begin its_all_over:=true;  goto exit ;
    end;
  back_input; {we will try to end again after ejecting residual material}
  begin  mem[ cur_list.tail_field ].hh.rh := new_null_box; cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
   mem[ cur_list.tail_field +width_offset].int  :=eqtb[dimen_base+ hsize_code].int   ;
  begin  mem[ cur_list.tail_field ].hh.rh := new_glue( mem_bot +glue_spec_size +glue_spec_size ); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
  begin  mem[ cur_list.tail_field ].hh.rh := new_penalty(-{010000000000=} 1073741824); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;

  build_page; {append \.[\\hbox to \\hsize\[\]\\vfill\\penalty-'10000000000]}
  end;
its_all_over:=false;
exit:end;


procedure append_glue;
var s:small_number; {modifier of skip command}
begin s:=cur_chr;
case s of
fil_code: cur_val:=mem_bot +glue_spec_size ;
fill_code: cur_val:=mem_bot +glue_spec_size +glue_spec_size ;
ss_code: cur_val:=mem_bot +glue_spec_size +glue_spec_size +glue_spec_size ;
fil_neg_code: cur_val:=mem_bot +glue_spec_size +glue_spec_size +glue_spec_size +glue_spec_size ;
skip_code: scan_glue(glue_val);
mskip_code: scan_glue(mu_val);
end; {now |cur_val| points to the glue specification}
begin  mem[ cur_list.tail_field ].hh.rh := new_glue( cur_val); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
if s>=skip_code then
  begin decr(  mem[  cur_val].hh.rh  );
  if s>skip_code then  mem[ cur_list.tail_field ].hh.b1 :=mu_glue;
  end;
end;


procedure append_kern;
var s:quarterword; {|subtype| of the kern node}
begin s:=cur_chr; scan_dimen(s=mu_glue,false,false);
begin  mem[ cur_list.tail_field ].hh.rh := new_kern( cur_val); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;  mem[ cur_list.tail_field ].hh.b1 :=s;
end;


procedure off_save;
var p:halfword ; {inserted token}
begin if cur_group=bottom_level then
  
{ Drop current token and complain that it was unmatched }
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Extra "=} 787); end ; print_cmd_chr(cur_cmd,cur_chr);
{ \xref[Extra x] }
 begin help_ptr:=1; help_line[0]:={"Things are pretty mixed up, but I think the worst is over."=} 1056; end ;

error;
end


else  begin back_input; p:=get_avail;  mem[ mem_top-3 ].hh.rh :=p;
  begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Missing "=} 635); end ;
  
{ Prepare to insert a token that matches |cur_group|, and print what it is }
case cur_group of
semi_simple_group: begin  mem[ p].hh.lh :={07777=}4095 +frozen_end_group;
  print_esc({"endgroup"=}524);
{ \xref[Missing \\endgroup inserted] }
  end;
math_shift_group: begin  mem[ p].hh.lh :=math_shift_token+{"$"=}36; print_char({"$"=}36);
{ \xref[Missing \$ inserted] }
  end;
math_left_group: begin  mem[ p].hh.lh :={07777=}4095 +frozen_right;  mem[ p].hh.rh :=get_avail;
  p:= mem[ p].hh.rh ;  mem[ p].hh.lh :=other_token+{"."=}46; print_esc({"right."=}1055);
{ \xref[Missing \\right\hbox[.] inserted] }
{ \xref[null delimiter] }
  end;
 else  begin  mem[ p].hh.lh :=right_brace_token+{"]"=}125; print_char({"]"=}125);
{ \xref[Missing \] inserted] }
  end
 end 

;
  print({" inserted"=}636); begin_token_list(  mem[  mem_top-3 ].hh.rh ,inserted) ;
   begin help_ptr:=5; help_line[4]:={"I've inserted something that you may have forgotten."=} 1050; help_line[3]:={"(See the <inserted text> above.)"=} 1051; help_line[2]:={"With luck, this will get me unwedged. But if you"=} 1052; help_line[1]:={"really didn't forget anything, try typing `2' now; then"=} 1053; help_line[0]:={"my insertion and my current dilemma will both disappear."=} 1054; end ;
  error;
  end;
end;


procedure extra_right_brace;
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Extra ], or forgotten "=} 1061); end ;
{ \xref[Extra \], or forgotten x] }
case cur_group of
semi_simple_group: print_esc({"endgroup"=}524);
math_shift_group: print_char({"$"=}36);
math_left_group: print_esc({"right"=}891);
end;

 begin help_ptr:=5; help_line[4]:={"I've deleted a group-closing symbol because it seems to be"=} 1062; help_line[3]:={"spurious, as in `$x]$'. But perhaps the ] is legitimate and"=} 1063; help_line[2]:={"you forgot something else, as in `\hbox[$x]'. In such cases"=} 1064; help_line[1]:={"the way to recover is to insert both the forgotten and the"=} 1065; help_line[0]:={"deleted material, e.g., by typing `I$]'."=} 1066; end ; error;
incr(align_state);
end;


procedure normal_paragraph;
begin if eqtb[int_base+ looseness_code].int  <>0 then eq_word_define(int_base+looseness_code,0);
if eqtb[dimen_base+ hang_indent_code].int   <>0 then eq_word_define(dimen_base+hang_indent_code,0);
if eqtb[int_base+ hang_after_code].int  <>1 then eq_word_define(int_base+hang_after_code,1);
if  eqtb[  par_shape_loc].hh.rh   <>-{0xfffffff=}268435455   then eq_define(par_shape_loc,shape_ref,-{0xfffffff=}268435455  );
end;


procedure box_end( box_context:integer);
var p:halfword ; {|ord_noad| for new box in math mode}
begin if box_context<{010000000000=}1073741824  then 
{ Append box |cur_box| to the current list, shifted by |box_context| }
begin if cur_box<>-{0xfffffff=}268435455   then
  begin  mem[ cur_box+4].int  :=box_context;
  if abs(cur_list.mode_field )=vmode then
    begin append_to_vlist(cur_box);
    if adjust_tail<>-{0xfffffff=}268435455   then
      begin if mem_top-5 <>adjust_tail then
        begin  mem[ cur_list.tail_field ].hh.rh := mem[ mem_top-5 ].hh.rh ; cur_list.tail_field :=adjust_tail;
        end;
      adjust_tail:=-{0xfffffff=}268435455  ;
      end;
    if cur_list.mode_field >0 then build_page;
    end
  else  begin if abs(cur_list.mode_field )=hmode then cur_list.aux_field .hh.lh :=1000
    else  begin p:=new_noad;
       mem[   p+1 ].hh.rh :=sub_box;
       mem[   p+1 ].hh.lh :=cur_box; cur_box:=p;
      end;
     mem[ cur_list.tail_field ].hh.rh :=cur_box; cur_list.tail_field :=cur_box;
    end;
  end;
end


else if box_context<{010000000000=}1073741824 +512  then 
{ Store \(c)|cur_box| in a box register }
if box_context<{010000000000=}1073741824 +256 then
  eq_define(box_base-{010000000000=}1073741824 +box_context,box_ref,cur_box)
else geq_define(box_base-{010000000000=}1073741824 -256+box_context,box_ref,cur_box)


else if cur_box<>-{0xfffffff=}268435455   then
  if box_context>{010000000000=}1073741824 +512  then 
{ Append a new leader node that uses |cur_box| }
begin 
{ Get the next non-blank non-relax... }
repeat get_x_token;
until (cur_cmd<>spacer)and(cur_cmd<>relax)

;
if ((cur_cmd=hskip)and(abs(cur_list.mode_field )<>vmode))or 
   ((cur_cmd=vskip)and(abs(cur_list.mode_field )=vmode)) then
  begin append_glue;  mem[ cur_list.tail_field ].hh.b1 :=box_context-({010000000000=}1073741824 +513 -a_leaders);
    mem[  cur_list.tail_field + 1].hh.rh  :=cur_box;
  end
else  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Leaders not followed by proper glue"=} 1079); end ;
{ \xref[Leaders not followed by...] }
   begin help_ptr:=3; help_line[2]:={"You should say `\leaders <box or rule><hskip or vskip>'."=} 1080; help_line[1]:={"I found the <box or rule>, but there's no suitable"=} 1081; help_line[0]:={"<hskip or vskip>, so I'm ignoring these leaders."=} 1082; end ; back_error;
  flush_node_list(cur_box);
  end;
end


  else ship_out(cur_box);
end;


procedure begin_box( box_context:integer);
label exit, done;
var  p, q:halfword ; {run through the current list}
 m:quarterword; {the length of a replacement list}
 k:halfword; {0 or |vmode| or |hmode|}
 n:eight_bits; {a box number}
begin case cur_chr of
box_code: begin scan_eight_bit_int; cur_box:= eqtb[  box_base+   cur_val].hh.rh   ;
   eqtb[  box_base+   cur_val].hh.rh   :=-{0xfffffff=}268435455  ; {the box becomes void, at the same level}
  end;
copy_code: begin scan_eight_bit_int; cur_box:=copy_node_list( eqtb[  box_base+   cur_val].hh.rh   );
  end;
last_box_code: 
{ If the current list ends with a box node, delete it from the list and make |cur_box| point to it; otherwise set |cur_box:=null| }
begin cur_box:=-{0xfffffff=}268435455  ;
if abs(cur_list.mode_field )=mmode then
  begin you_cant;  begin help_ptr:=1; help_line[0]:={"Sorry; this \lastbox will be void."=} 1083; end ; error;
  end
else if (cur_list.mode_field =vmode)and(cur_list.head_field =cur_list.tail_field ) then
  begin you_cant;
   begin help_ptr:=2; help_line[1]:={"Sorry...I usually can't take things from the current page."=} 1084; help_line[0]:={"This \lastbox will therefore be void."=} 1085; end ; error;
  end
else  begin if not  ( cur_list.tail_field >=hi_mem_min)  then
    if ( mem[ cur_list.tail_field ].hh.b0 =hlist_node)or( mem[ cur_list.tail_field ].hh.b0 =vlist_node) then
      
{ Remove the last box, unless it's part of a discretionary }
begin q:=cur_list.head_field ;
repeat p:=q;
if not  ( q>=hi_mem_min)  then if  mem[ q].hh.b0 =disc_node then
  begin for m:=1 to  mem[ q].hh.b1  do p:= mem[ p].hh.rh ;
  if p=cur_list.tail_field  then goto done;
  end;
q:= mem[ p].hh.rh ;
until q=cur_list.tail_field ;
cur_box:=cur_list.tail_field ;  mem[ cur_box+4].int  :=0;
cur_list.tail_field :=p;  mem[ p].hh.rh :=-{0xfffffff=}268435455  ;
done:end

;
  end;
end

;
vsplit_code: 
{ Split off part of a vertical box, make |cur_box| point to it }
begin scan_eight_bit_int; n:=cur_val;
if not scan_keyword({"to"=}856) then
{ \xref[to] }
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Missing `to' inserted"=} 1086); end ;
{ \xref[Missing `to' inserted] }
   begin help_ptr:=2; help_line[1]:={"I'm working on `\vsplit<box number> to <dimen>';"=} 1087; help_line[0]:={"will look for the <dimen> next."=} 1088; end ; error;
  end;
scan_dimen(false,false,false) ;
cur_box:=vsplit(n,cur_val);
end

;
 else  
{ Initiate the construction of an hbox or vbox, then |return| }
begin k:=cur_chr-vtop_code; save_stack[save_ptr+ 0].int :=box_context;
if k=hmode then
  if (box_context<{010000000000=}1073741824 )and(abs(cur_list.mode_field )=vmode) then
    scan_spec(adjusted_hbox_group,true)
  else scan_spec(hbox_group,true)
else  begin if k=vmode then scan_spec(vbox_group,true)
  else  begin scan_spec(vtop_group,true); k:=vmode;
    end;
  normal_paragraph;
  end;
push_nest; cur_list.mode_field :=-k;
if k=vmode then
  begin cur_list.aux_field .int  :=-65536000 ;
  if  eqtb[  every_vbox_loc].hh.rh   <>-{0xfffffff=}268435455   then begin_token_list( eqtb[  every_vbox_loc].hh.rh   ,every_vbox_text);
  end
else  begin cur_list.aux_field .hh.lh :=1000;
  if  eqtb[  every_hbox_loc].hh.rh   <>-{0xfffffff=}268435455   then begin_token_list( eqtb[  every_hbox_loc].hh.rh   ,every_hbox_text);
  end;
 goto exit ;
end


 end ;

box_end(box_context); {in simple cases, we use the box immediately}
exit:end;


procedure scan_box( box_context:integer);
  {the next input should specify a box or perhaps a rule}
begin 
{ Get the next non-blank non-relax... }
repeat get_x_token;
until (cur_cmd<>spacer)and(cur_cmd<>relax)

;
if cur_cmd=make_box then begin_box(box_context)
else if (box_context>={010000000000=}1073741824 +513 )and((cur_cmd=hrule)or(cur_cmd=vrule)) then
  begin cur_box:=scan_rule_spec; box_end(box_context);
  end
else  begin{  } 

  begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"A <box> was supposed to be here"=} 1089); end ;

{ \xref[A <box> was supposed to...] }
   begin help_ptr:=3; help_line[2]:={"I was expecting to see \hbox or \vbox or \copy or \box or"=} 1090; help_line[1]:={"something like that. So you might find something missing in"=} 1091; help_line[0]:={"your output. But keep trying; you can fix this later."=} 1092; end ; back_error;
  end;
end;


procedure package( c:small_number);
var h:scaled; {height of box}
 p:halfword ; {first node in a box}
 d:scaled; {max depth}
begin d:=eqtb[dimen_base+ box_max_depth_code].int   ; unsave; save_ptr:=save_ptr-3;
if cur_list.mode_field =-hmode then cur_box:=hpack( mem[ cur_list.head_field ].hh.rh ,save_stack[save_ptr+ 2].int ,save_stack[save_ptr+ 1].int )
else  begin cur_box:=vpackage( mem[ cur_list.head_field ].hh.rh ,save_stack[save_ptr+ 2].int ,save_stack[save_ptr+ 1].int ,d);
  if c=vtop_code then 
{ Readjust the height and depth of |cur_box|, for \.[\\vtop] }
begin h:=0; p:=  mem[  cur_box+ list_offset].hh.rh  ;
if p<>-{0xfffffff=}268435455   then if  mem[ p].hh.b0 <=rule_node then h:= mem[ p+height_offset].int  ;
 mem[ cur_box+depth_offset].int  := mem[ cur_box+depth_offset].int  -h+ mem[ cur_box+height_offset].int  ;  mem[ cur_box+height_offset].int  :=h;
end

;
  end;
pop_nest; box_end(save_stack[save_ptr+ 0].int );
end;


function norm_min( h:integer):small_number;
begin if h<=0 then norm_min:=1 else if h>=63 then norm_min:=63 
else norm_min:=h;
end;


procedure new_graf( indented:boolean);
begin cur_list.pg_field :=0;
if (cur_list.mode_field =vmode)or(cur_list.head_field <>cur_list.tail_field ) then
  begin  mem[ cur_list.tail_field ].hh.rh := new_param_glue( par_skip_code); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
push_nest; cur_list.mode_field :=hmode; cur_list.aux_field .hh.lh :=1000; if eqtb[int_base+ language_code].int  <=0 then cur_lang:=0 else if eqtb[int_base+ language_code].int  >255 then cur_lang:=0 else cur_lang:=eqtb[int_base+ language_code].int   ; cur_list.aux_field .hh.rh :=cur_lang;
cur_list.pg_field :=(norm_min(eqtb[int_base+ left_hyphen_min_code].int  )*{0100=}64+norm_min(eqtb[int_base+ right_hyphen_min_code].int  ))
             *{0200000=}65536+cur_lang;
if indented then
  begin cur_list.tail_field :=new_null_box;  mem[ cur_list.head_field ].hh.rh :=cur_list.tail_field ;  mem[ cur_list.tail_field +width_offset].int  :=eqtb[dimen_base+ par_indent_code].int   ;
  if (insert_src_special_every_par) then insert_src_special; 
  end;
if  eqtb[  every_par_loc].hh.rh   <>-{0xfffffff=}268435455   then begin_token_list( eqtb[  every_par_loc].hh.rh   ,every_par_text);
if nest_ptr=1 then build_page; {put |par_skip| glue on current page}
end;


procedure indent_in_hmode;
var p, q:halfword ;
begin if cur_chr>0 then {\.[\\indent]}
  begin p:=new_null_box;  mem[ p+width_offset].int  :=eqtb[dimen_base+ par_indent_code].int   ;
  if abs(cur_list.mode_field )=hmode then cur_list.aux_field .hh.lh :=1000
  else  begin q:=new_noad;  mem[   q+1 ].hh.rh :=sub_box;
     mem[   q+1 ].hh.lh :=p; p:=q;
    end;
  begin  mem[ cur_list.tail_field ].hh.rh := p; cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
  end;
end;


procedure head_for_vmode;
begin if cur_list.mode_field <0 then
  if cur_cmd<>hrule then off_save
  else  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"You can't use `"=} 695); end ;
    print_esc({"hrule"=}529); print({"' here except with leaders"=}1095);
{ \xref[You can't use \\hrule...] }
     begin help_ptr:=2; help_line[1]:={"To put a horizontal rule in an hbox or an alignment,"=} 1096; help_line[0]:={"you should use \leaders or \hrulefill (see The TeXbook)."=} 1097; end ;
    error;
    end
else  begin back_input; cur_tok:=par_token; back_input; cur_input.index_field  :=inserted;
  end;
end;


procedure end_graf;
begin if cur_list.mode_field =hmode then
  begin if cur_list.head_field =cur_list.tail_field  then pop_nest {null paragraphs are ignored}
  else line_break(eqtb[int_base+ widow_penalty_code].int  );
  normal_paragraph;
  error_count:=0;
  end;
end;


procedure begin_insert_or_adjust;
begin if cur_cmd=vadjust then cur_val:=255
else  begin scan_eight_bit_int;
  if cur_val=255 then
    begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"You can't "=} 1098); end ; print_esc({"insert"=}327); print_int(255);
{ \xref[You can't \\insert255] }
     begin help_ptr:=1; help_line[0]:={"I'm changing to \insert0; box 255 is special."=} 1099; end ;
    error; cur_val:=0;
    end;
  end;
save_stack[save_ptr+ 0].int :=cur_val; incr(save_ptr);
new_save_level(insert_group); scan_left_brace; normal_paragraph;
push_nest; cur_list.mode_field :=-vmode; cur_list.aux_field .int  :=-65536000 ;
end;


procedure make_mark;
var p:halfword ; {new node}
begin p:=scan_toks(false,true); p:=get_node(small_node_size);
 mem[ p].hh.b0 :=mark_node;  mem[ p].hh.b1 :=0; {the |subtype| is not used}
mem[ p+1].int :=def_ref;  mem[ cur_list.tail_field ].hh.rh :=p; cur_list.tail_field :=p;
end;


procedure append_penalty;
begin scan_int; begin  mem[ cur_list.tail_field ].hh.rh := new_penalty( cur_val); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
if cur_list.mode_field =vmode then build_page;
end;


procedure delete_last;
label exit;
var  p, q:halfword ; {run through the current list}
 m:quarterword; {the length of a replacement list}
begin if (cur_list.mode_field =vmode)and(cur_list.tail_field =cur_list.head_field ) then
  
{ Apologize for inability to do the operation now, unless \.[\\unskip] follows non-glue }
begin if (cur_chr<>glue_node)or(last_glue<>{0xfffffff=}268435455 ) then
  begin you_cant;
   begin help_ptr:=2; help_line[1]:={"Sorry...I usually can't take things from the current page."=} 1084; help_line[0]:={"Try `I\vskip-\lastskip' instead."=} 1100; end ;
  if cur_chr=kern_node then help_line[0]:=
    ({"Try `I\kern-\lastkern' instead."=}1101)
  else if cur_chr<>glue_node then help_line[0]:= 
    ({"Perhaps you can make the output routine do it."=}1102);
  error;
  end;
end


else  begin if not  ( cur_list.tail_field >=hi_mem_min)  then if  mem[ cur_list.tail_field ].hh.b0 =cur_chr then
    begin q:=cur_list.head_field ;
    repeat p:=q;
    if not  ( q>=hi_mem_min)  then if  mem[ q].hh.b0 =disc_node then
      begin for m:=1 to  mem[ q].hh.b1  do p:= mem[ p].hh.rh ;
      if p=cur_list.tail_field  then  goto exit ;
      end;
    q:= mem[ p].hh.rh ;
    until q=cur_list.tail_field ;
     mem[ p].hh.rh :=-{0xfffffff=}268435455  ; flush_node_list(cur_list.tail_field ); cur_list.tail_field :=p;
    end;
  end;
exit:end;


procedure unpackage;
label exit;
var p:halfword ; {the box}
 c:box_code..copy_code; {should we copy?}
begin c:=cur_chr; scan_eight_bit_int; p:= eqtb[  box_base+   cur_val].hh.rh   ;
if p=-{0xfffffff=}268435455   then  goto exit ;
if (abs(cur_list.mode_field )=mmode)or((abs(cur_list.mode_field )=vmode)and( mem[ p].hh.b0 <>vlist_node))or 
   ((abs(cur_list.mode_field )=hmode)and( mem[ p].hh.b0 <>hlist_node)) then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Incompatible list can't be unboxed"=} 1110); end ;
{ \xref[Incompatible list...] }
   begin help_ptr:=3; help_line[2]:={"Sorry, Pandora. (You sneaky devil.)"=} 1111; help_line[1]:={"I refuse to unbox an \hbox in vertical mode or vice versa."=} 1112; help_line[0]:={"And I can't open any boxes in math mode."=} 1113; end ;

  error;  goto exit ;
  end;
if c=copy_code then  mem[ cur_list.tail_field ].hh.rh :=copy_node_list(  mem[  p+ list_offset].hh.rh  )
else  begin  mem[ cur_list.tail_field ].hh.rh :=  mem[  p+ list_offset].hh.rh  ;  eqtb[  box_base+   cur_val].hh.rh   :=-{0xfffffff=}268435455  ;
  free_node(p,box_node_size);
  end;
while  mem[ cur_list.tail_field ].hh.rh <>-{0xfffffff=}268435455   do cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ;
exit:end;


procedure append_italic_correction;
label exit;
var p:halfword ; {|char_node| at the tail of the current list}
 f:internal_font_number; {the font in the |char_node|}
begin if cur_list.tail_field <>cur_list.head_field  then
  begin if  ( cur_list.tail_field >=hi_mem_min)  then p:=cur_list.tail_field 
  else if  mem[ cur_list.tail_field ].hh.b0 =ligature_node then p:= cur_list.tail_field +1 
  else  goto exit ;
  f:=  mem[ p].hh.b0 ;
  begin  mem[ cur_list.tail_field ].hh.rh := new_kern( font_info[italic_base[  f]+(    font_info[char_base[    f]+effective_char(true,    f,       mem[     p].hh.b1 )].qqqq . b2 ) div 4].int  ); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
   mem[ cur_list.tail_field ].hh.b1 :=explicit;
  end;
exit:end;


procedure append_discretionary;
var c:integer; {hyphen character}
begin begin  mem[ cur_list.tail_field ].hh.rh := new_disc; cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
if cur_chr=1 then
  begin c:=hyphen_char[ eqtb[  cur_font_loc].hh.rh   ];
  if c>=0 then if c<256 then   mem[  cur_list.tail_field + 1].hh.lh  :=new_character( eqtb[  cur_font_loc].hh.rh   ,c);
  end
else  begin incr(save_ptr); save_stack[save_ptr+- 1].int :=0; new_save_level(disc_group);
  scan_left_brace; push_nest; cur_list.mode_field :=-hmode; cur_list.aux_field .hh.lh :=1000;
  end;
end;


procedure build_discretionary;
label done,exit;
var p, q:halfword ; {for link manipulation}
 n:integer; {length of discretionary list}
begin unsave;

{ Prune the current list, if necessary, until it contains only |char_node|, |kern_node|, |hlist_node|, |vlist_node|, |rule_node|, and |ligature_node| items; set |n| to the length of the list, and set |q| to the list's tail }
q:=cur_list.head_field ; p:= mem[ q].hh.rh ; n:=0;
while p<>-{0xfffffff=}268435455   do
  begin if not  ( p>=hi_mem_min)  then if  mem[ p].hh.b0 >rule_node then
    if  mem[ p].hh.b0 <>kern_node then if  mem[ p].hh.b0 <>ligature_node then
      begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Improper discretionary list"=} 1120); end ;
{ \xref[Improper discretionary list] }
       begin help_ptr:=1; help_line[0]:={"Discretionary lists must contain only boxes and kerns."=} 1121; end ;

      error;
      begin_diagnostic;
      print_nl({"The following discretionary sublist has been deleted:"=}1122);
{ \xref[The following...deleted] }
      show_box(p);
      end_diagnostic(true);
      flush_node_list(p);  mem[ q].hh.rh :=-{0xfffffff=}268435455  ; goto done;
      end;
  q:=p; p:= mem[ q].hh.rh ; incr(n);
  end;
done:

;
p:= mem[ cur_list.head_field ].hh.rh ; pop_nest;
case save_stack[save_ptr+- 1].int  of
0:  mem[  cur_list.tail_field + 1].hh.lh  :=p;
1:  mem[  cur_list.tail_field + 1].hh.rh  :=p;
2:
{ Attach list |p| to the current list, and record its length; then finish up and |return| }
begin if (n>0)and(abs(cur_list.mode_field )=mmode) then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Illegal math "=} 1114); end ; print_esc({"discretionary"=}346);
{ \xref[Illegal math \\disc...] }
   begin help_ptr:=2; help_line[1]:={"Sorry: The third part of a discretionary break must be"=} 1115; help_line[0]:={"empty, in math formulas. I had to delete your third part."=} 1116; end ;
  flush_node_list(p); n:=0; error;
  end
else  mem[ cur_list.tail_field ].hh.rh :=p;
if n<=max_quarterword then  mem[ cur_list.tail_field ].hh.b1 :=n
else  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Discretionary list is too long"=} 1117); end ;
{ \xref[Discretionary list is too long] }
   begin help_ptr:=2; help_line[1]:={"Wow---I never thought anybody would tweak me here."=} 1118; help_line[0]:={"You can't seriously need such a huge discretionary list?"=} 1119; end ;
  error;
  end;
if n>0 then cur_list.tail_field :=q;
decr(save_ptr);  goto exit ;
end

;
end; {there are no other cases}
incr(save_stack[save_ptr+- 1].int ); new_save_level(disc_group); scan_left_brace;
push_nest; cur_list.mode_field :=-hmode; cur_list.aux_field .hh.lh :=1000;
exit:end;


procedure make_accent;
var s, t: real; {amount of slant}
 p, q, r:halfword ; {character, box, and kern nodes}
 f:internal_font_number; {relevant font}
 a, h, x, w, delta:scaled; {heights and widths, as explained above}
 i:four_quarters; {character information}
begin scan_char_num; f:= eqtb[  cur_font_loc].hh.rh   ; p:=new_character(f,cur_val);
if p<>-{0xfffffff=}268435455   then
  begin x:=font_info[ x_height_code+param_base[ f]].int  ; s:=font_info[ slant_code+param_base[ f]].int  /  65536.0 ;
{ \xref[real division] }
  a:=font_info[width_base[ f]+  font_info[char_base[  f]+effective_char(true,  f,     mem[   p].hh.b1 )].qqqq .b0].int  ;

  do_assignments;

  
{ Create a character node |q| for the next character, but set |q:=null| if problems arise }
q:=-{0xfffffff=}268435455  ; f:= eqtb[  cur_font_loc].hh.rh   ;
if (cur_cmd=letter)or(cur_cmd=other_char)or(cur_cmd=char_given) then
  q:=new_character(f,cur_chr)
else if cur_cmd=char_num then
  begin scan_char_num; q:=new_character(f,cur_val);
  end
else back_input

;
  if q<>-{0xfffffff=}268435455   then 
{ Append the accent with appropriate kerns, then set |p:=q| }
begin t:=font_info[ slant_code+param_base[ f]].int  /  65536.0 ;
{ \xref[real division] }
i:= font_info[char_base[ f]+effective_char(true, f,    mem[  q].hh.b1 )].qqqq ;
w:=font_info[width_base[ f]+ i.b0].int  ; h:=font_info[height_base[ f]+(    i. b1  ) div 16].int  ;
if h<>x then {the accent must be shifted up or down}
  begin p:=hpack(p,0,additional );  mem[ p+4].int  :=x-h;
  end;
delta:=round((w-a)/  2.0 +h*t-x*s);
{ \xref[real multiplication] }
{ \xref[real addition] }
r:=new_kern(delta);  mem[ r].hh.b1 :=acc_kern;  mem[ cur_list.tail_field ].hh.rh :=r;  mem[ r].hh.rh :=p;
cur_list.tail_field :=new_kern(-a-delta);  mem[ cur_list.tail_field ].hh.b1 :=acc_kern;  mem[ p].hh.rh :=cur_list.tail_field ; p:=q;
end

;
   mem[ cur_list.tail_field ].hh.rh :=p; cur_list.tail_field :=p; cur_list.aux_field .hh.lh :=1000;
  end;
end;


procedure align_error;
begin if abs(align_state)>2 then
  
{ Express consternation over the fact that no alignment is in progress }
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Misplaced "=} 1127); end ; print_cmd_chr(cur_cmd,cur_chr);
{ \xref[Misplaced \&] }
{ \xref[Misplaced \\span] }
{ \xref[Misplaced \\cr] }
if cur_tok=tab_token+{"&"=}38 then
  begin  begin help_ptr:=6; help_line[5]:={"I can't figure out why you would want to use a tab mark"=} 1128; help_line[4]:={"here. If you just want an ampersand, the remedy is"=} 1129; help_line[3]:={"simple: Just type `I\&' now. But if some right brace"=} 1130; help_line[2]:={"up above has ended a previous alignment prematurely,"=} 1131; help_line[1]:={"you're probably due for more error messages, and you"=} 1132; help_line[0]:={"might try typing `S' now just to see what is salvageable."=} 1133; end ;
  end
else  begin  begin help_ptr:=5; help_line[4]:={"I can't figure out why you would want to use a tab mark"=} 1128; help_line[3]:={"or \cr or \span just now. If something like a right brace"=} 1134; help_line[2]:={"up above has ended a previous alignment prematurely,"=} 1131; help_line[1]:={"you're probably due for more error messages, and you"=} 1132; help_line[0]:={"might try typing `S' now just to see what is salvageable."=} 1133; end ;
  end;
error;
end


else  begin back_input;
  if align_state<0 then
    begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Missing [ inserted"=} 667); end ;
{ \xref[Missing \[ inserted] }
    incr(align_state); cur_tok:=left_brace_token+{"["=}123;
    end
  else  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Missing ] inserted"=} 1123); end ;
{ \xref[Missing \] inserted] }
    decr(align_state); cur_tok:=right_brace_token+{"]"=}125;
    end;
   begin help_ptr:=3; help_line[2]:={"I've put in what seems to be necessary to fix"=} 1124; help_line[1]:={"the current column of the current alignment."=} 1125; help_line[0]:={"Try to go on, since this might almost work."=} 1126; end ; ins_error;
  end;
end;


procedure no_align_error;
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Misplaced "=} 1127); end ; print_esc({"noalign"=}535);
{ \xref[Misplaced \\noalign] }
 begin help_ptr:=2; help_line[1]:={"I expect to see \noalign only after the \cr of"=} 1135; help_line[0]:={"an alignment. Proceed, and I'll ignore this case."=} 1136; end ; error;
end;
procedure omit_error;
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Misplaced "=} 1127); end ; print_esc({"omit"=}538);
{ \xref[Misplaced \\omit] }
 begin help_ptr:=2; help_line[1]:={"I expect to see \omit only after tab marks or the \cr of"=} 1137; help_line[0]:={"an alignment. Proceed, and I'll ignore this case."=} 1136; end ; error;
end;


procedure do_endv;
begin base_ptr:=input_ptr; input_stack[base_ptr]:=cur_input;
while (input_stack[base_ptr].index_field<>v_template) and
      (input_stack[base_ptr].loc_field=-{0xfffffff=}268435455  ) and
      (input_stack[base_ptr].state_field=token_list) do decr(base_ptr);
if (input_stack[base_ptr].index_field<>v_template) or
      (input_stack[base_ptr].loc_field<>-{0xfffffff=}268435455  ) or
      (input_stack[base_ptr].state_field<>token_list) then
  fatal_error({"(interwoven alignment preambles are not allowed)"=}602);
{ \xref[interwoven alignment preambles...] }
 if cur_group=align_group then
  begin end_graf;
  if fin_col then fin_row;
  end
else off_save;
end;


procedure cs_error;
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Extra "=} 787); end ; print_esc({"endcsname"=}513);
{ \xref[Extra \\endcsname] }
 begin help_ptr:=1; help_line[0]:={"I'm ignoring this, since I wasn't doing a \csname."=} 1139; end ;
error;
end;


procedure push_math( c:group_code);
begin push_nest; cur_list.mode_field :=-mmode; cur_list.aux_field .int :=-{0xfffffff=}268435455  ; new_save_level(c);
end;


procedure init_math;
label reswitch,found,not_found,done;
var w:scaled; {new or partial |pre_display_size|}
 l:scaled; {new |display_width|}
 s:scaled; {new |display_indent|}
 p:halfword ; {current node when calculating |pre_display_size|}
 q:halfword ; {glue specification when calculating |pre_display_size|}
 f:internal_font_number; {font in current |char_node|}
 n:integer; {scope of paragraph shape specification}
 v:scaled; {|w| plus possible glue amount}
 d:scaled; {increment to |v|}
begin get_token; {|get_x_token| would fail on \.[\\ifmmode]\thinspace!}
if (cur_cmd=math_shift)and(cur_list.mode_field >0) then 
{ Go into display math mode }
begin if cur_list.head_field =cur_list.tail_field  then {`\.[\\noindent\$\$]' or `\.[\$\$[ ]\$\$]'}
  begin pop_nest; w:=-{07777777777=}1073741823 ;
  end
else  begin line_break(eqtb[int_base+ display_widow_penalty_code].int  );

  
{ Calculate the natural width, |w|, by which the characters of the final line extend to the right of the reference point, plus two ems; or set |w:=max_dimen| if the non-blank information on that line is affected by stretching or shrinking }
v:= mem[ just_box+4].int  +2*font_info[ quad_code+param_base[  eqtb[  cur_font_loc].hh.rh   ]].int  ; w:=-{07777777777=}1073741823 ;
p:=  mem[  just_box+ list_offset].hh.rh  ;
while p<>-{0xfffffff=}268435455   do
  begin 
{ Let |d| be the natural width of node |p|; if the node is ``visible,'' |goto found|; if the node is glue that stretches or shrinks, set |v:=max_dimen| }
reswitch: if  ( p>=hi_mem_min)  then
  begin f:=  mem[ p].hh.b0 ; d:=font_info[width_base[ f]+  font_info[char_base[  f]+effective_char(true,  f,     mem[   p].hh.b1 )].qqqq .b0].int  ;
  goto found;
  end;
case  mem[ p].hh.b0  of
hlist_node,vlist_node,rule_node: begin d:= mem[ p+width_offset].int  ; goto found;
  end;
ligature_node:
{ Make node |p| look like a |char_node|... }
begin mem[mem_top-12 ]:=mem[ p+1 ];  mem[ mem_top-12 ].hh.rh := mem[ p].hh.rh ;
p:=mem_top-12 ; goto reswitch;
end

;
kern_node,math_node: d:= mem[ p+width_offset].int  ;
glue_node:
{ Let |d| be the natural width of this glue; if stretching or shrinking, set |v:=max_dimen|; |goto found| in the case of leaders }
begin q:=  mem[  p+ 1].hh.lh  ; d:= mem[ q+width_offset].int  ;
if   mem[  just_box+ list_offset].hh.b0  =stretching then
  begin if (  mem[  just_box+ list_offset].hh.b1  =  mem[ q].hh.b0 )and 
     ( mem[ q+2].int  <>0) then
    v:={07777777777=}1073741823 ;
  end
else if   mem[  just_box+ list_offset].hh.b0  =shrinking then
  begin if (  mem[  just_box+ list_offset].hh.b1  =  mem[ q].hh.b1 )and 
     ( mem[ q+3].int  <>0) then
    v:={07777777777=}1073741823 ;
  end;
if  mem[ p].hh.b1 >=a_leaders then goto found;
end

;
whatsit_node: 
{ Let |d| be the width of the whatsit |p| }d:=0

;
 else  d:=0
 end 

;
  if v<{07777777777=}1073741823  then v:=v+d;
  goto not_found;
  found: if v<{07777777777=}1073741823  then
    begin v:=v+d; w:=v;
    end
  else  begin w:={07777777777=}1073741823 ; goto done;
    end;
  not_found: p:= mem[ p].hh.rh ;
  end;
done:

;
  end;
{now we are in vertical mode, working on the list that will contain the display}

{ Calculate the length, |l|, and the shift amount, |s|, of the display lines }
if  eqtb[  par_shape_loc].hh.rh   =-{0xfffffff=}268435455   then
  if (eqtb[dimen_base+ hang_indent_code].int   <>0)and 
   (((eqtb[int_base+ hang_after_code].int  >=0)and(cur_list.pg_field +2>eqtb[int_base+ hang_after_code].int  ))or 
    (cur_list.pg_field +1<-eqtb[int_base+ hang_after_code].int  )) then
    begin l:=eqtb[dimen_base+ hsize_code].int   -abs(eqtb[dimen_base+ hang_indent_code].int   );
    if eqtb[dimen_base+ hang_indent_code].int   >0 then s:=eqtb[dimen_base+ hang_indent_code].int    else s:=0;
    end
  else  begin l:=eqtb[dimen_base+ hsize_code].int   ; s:=0;
    end
else  begin n:= mem[  eqtb[  par_shape_loc].hh.rh   ].hh.lh ;
  if cur_list.pg_field +2>=n then p:= eqtb[  par_shape_loc].hh.rh   +2*n
  else p:= eqtb[  par_shape_loc].hh.rh   +2*(cur_list.pg_field +2);
  s:=mem[p-1].int ; l:=mem[p].int ;
  end

;
push_math(math_shift_group); cur_list.mode_field :=mmode;
eq_word_define(int_base+cur_fam_code,-1);

eq_word_define(dimen_base+pre_display_size_code,w);
eq_word_define(dimen_base+display_width_code,l);
eq_word_define(dimen_base+display_indent_code,s);
if  eqtb[  every_display_loc].hh.rh   <>-{0xfffffff=}268435455   then begin_token_list( eqtb[  every_display_loc].hh.rh   ,every_display_text);
if nest_ptr=1 then build_page;
end


else  begin back_input; 
{ Go into ordinary math mode }
begin push_math(math_shift_group); eq_word_define(int_base+cur_fam_code,-1);
if (insert_src_special_every_math) then insert_src_special;
if  eqtb[  every_math_loc].hh.rh   <>-{0xfffffff=}268435455   then begin_token_list( eqtb[  every_math_loc].hh.rh   ,every_math_text);
end

;
  end;
end;


procedure start_eq_no;
begin save_stack[save_ptr+ 0].int :=cur_chr; incr(save_ptr);

{ Go into ordinary math mode }
begin push_math(math_shift_group); eq_word_define(int_base+cur_fam_code,-1);
if (insert_src_special_every_math) then insert_src_special;
if  eqtb[  every_math_loc].hh.rh   <>-{0xfffffff=}268435455   then begin_token_list( eqtb[  every_math_loc].hh.rh   ,every_math_text);
end

;
end;


procedure scan_math( p:halfword );
label restart,reswitch,exit;
var c:integer; {math character code}
begin restart:
{ Get the next non-blank non-relax... }
repeat get_x_token;
until (cur_cmd<>spacer)and(cur_cmd<>relax)

;
reswitch:case cur_cmd of
letter,other_char,char_given: begin c:=  eqtb[  math_code_base+    cur_chr].hh.rh    ;
    if c={0100000=}32768 then
      begin 
{ Treat |cur_chr| as an active character }
begin cur_cs:=cur_chr+active_base;
cur_cmd:= eqtb[  cur_cs].hh.b0  ; cur_chr:= eqtb[  cur_cs].hh.rh  ;
x_token; back_input;
end

;
      goto restart;
      end;
    end;
char_num: begin scan_char_num; cur_chr:=cur_val; cur_cmd:=char_given;
  goto reswitch;
  end;
math_char_num: begin scan_fifteen_bit_int; c:=cur_val;
  end;
math_given: c:=cur_chr;
delim_num: begin scan_twenty_seven_bit_int; c:=cur_val div {010000=}4096;
  end;
 else  
{ Scan a subformula enclosed in braces and |return| }
begin back_input; scan_left_brace;

save_stack[save_ptr+ 0].int :=p; incr(save_ptr); push_math(math_group);  goto exit ;
end


 end ;

 mem[ p].hh.rh :=math_char;   mem[ p].hh.b1 := c  mod  256 ;
if (c>={070000=}28672 )and ((eqtb[int_base+ cur_fam_code].int  >=0)and(eqtb[int_base+ cur_fam_code].int  <16))  then   mem[ p].hh.b0 :=eqtb[int_base+ cur_fam_code].int  
else   mem[ p].hh.b0 :=(c div 256) mod 16;
exit:end;


procedure set_math_char( c:integer);
var p:halfword ; {the new noad}
begin if c>={0100000=}32768 then
  
{ Treat |cur_chr|... }
begin cur_cs:=cur_chr+active_base;
cur_cmd:= eqtb[  cur_cs].hh.b0  ; cur_chr:= eqtb[  cur_cs].hh.rh  ;
x_token; back_input;
end


else  begin p:=new_noad;  mem[   p+1 ].hh.rh :=math_char;
    mem[   p+1 ].hh.b1 := c  mod  256 ;
    mem[   p+1 ].hh.b0 :=(c div 256) mod 16;
  if c>={070000=}28672  then
    begin if ((eqtb[int_base+ cur_fam_code].int  >=0)and(eqtb[int_base+ cur_fam_code].int  <16))  then   mem[   p+1 ].hh.b0 :=eqtb[int_base+ cur_fam_code].int  ;
     mem[ p].hh.b0 :=ord_noad;
    end
  else   mem[ p].hh.b0 :=ord_noad+(c div {010000=}4096);
   mem[ cur_list.tail_field ].hh.rh :=p; cur_list.tail_field :=p;
  end;
end;


procedure math_limit_switch;
label exit;
begin if cur_list.head_field <>cur_list.tail_field  then if  mem[ cur_list.tail_field ].hh.b0 =op_noad then
  begin  mem[ cur_list.tail_field ].hh.b1 :=cur_chr;  goto exit ;
  end;
begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Limit controls must follow a math operator"=} 1143); end ;
{ \xref[Limit controls must follow...] }
 begin help_ptr:=1; help_line[0]:={"I'm ignoring this misplaced \limits or \nolimits command."=} 1144; end ; error;
exit:end;


procedure scan_delimiter( p:halfword ; r:boolean);
begin if r then scan_twenty_seven_bit_int
else  begin 
{ Get the next non-blank non-relax... }
repeat get_x_token;
until (cur_cmd<>spacer)and(cur_cmd<>relax)

;
  case cur_cmd of
  letter,other_char: cur_val:=eqtb[del_code_base+ cur_chr].int ;
  delim_num: scan_twenty_seven_bit_int;
   else  cur_val:=-1
   end ;
  end;
if cur_val<0 then 
{ Report that an invalid delimiter code is being changed to null; set~|cur_val:=0| }
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Missing delimiter (. inserted)"=} 1145); end ;
{ \xref[Missing delimiter...] }
 begin help_ptr:=6; help_line[5]:={"I was expecting to see something like `(' or `\[' or"=} 1146; help_line[4]:={"`\]' here. If you typed, e.g., `[' instead of `\[', you"=} 1147; help_line[3]:={"should probably delete the `[' by typing `1' now, so that"=} 1148; help_line[2]:={"braces don't get unbalanced. Otherwise just proceed."=} 1149; help_line[1]:={"Acceptable delimiters are characters whose \delcode is"=} 1150; help_line[0]:={"nonnegative, or you can use `\delimiter <delimiter code>'."=} 1151; end ;
back_error; cur_val:=0;
end

;
mem[ p].qqqq.b0 :=(cur_val div {04000000=}1048576) mod 16;
mem[ p].qqqq.b1 :=( cur_val  div {010000=} 4096)  mod  256 ;
mem[ p].qqqq.b2 :=(cur_val div 256) mod 16;
mem[ p].qqqq.b3 := cur_val  mod  256 ;
end;


procedure math_radical;
begin begin  mem[ cur_list.tail_field ].hh.rh := get_node( radical_noad_size); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
 mem[ cur_list.tail_field ].hh.b0 :=radical_noad;  mem[ cur_list.tail_field ].hh.b1 :=normal;
mem[ cur_list.tail_field +1 ].hh:=empty_field;
mem[ cur_list.tail_field +3 ].hh:=empty_field;
mem[ cur_list.tail_field +2 ].hh:=empty_field;
scan_delimiter( cur_list.tail_field +4 ,true); scan_math( cur_list.tail_field +1 );
end;


procedure math_ac;
begin if cur_cmd=accent then
  
{ Complain that the user should have said \.[\\mathaccent] }
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Please use "=} 1152); end ; print_esc({"mathaccent"=}531);
print({" for accents in math mode"=}1153);
{ \xref[Please use \\mathaccent...] }
 begin help_ptr:=2; help_line[1]:={"I'm changing \accent to \mathaccent here; wish me luck."=} 1154; help_line[0]:={"(Accents are not the same in formulas as they are in text.)"=} 1155; end ;
error;
end

;
begin  mem[ cur_list.tail_field ].hh.rh := get_node( accent_noad_size); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
 mem[ cur_list.tail_field ].hh.b0 :=accent_noad;  mem[ cur_list.tail_field ].hh.b1 :=normal;
mem[ cur_list.tail_field +1 ].hh:=empty_field;
mem[ cur_list.tail_field +3 ].hh:=empty_field;
mem[ cur_list.tail_field +2 ].hh:=empty_field;
 mem[   cur_list.tail_field +4 ].hh.rh :=math_char;
scan_fifteen_bit_int;
  mem[   cur_list.tail_field +4 ].hh.b1 := cur_val  mod  256 ;
if (cur_val>={070000=}28672 )and ((eqtb[int_base+ cur_fam_code].int  >=0)and(eqtb[int_base+ cur_fam_code].int  <16))  then   mem[   cur_list.tail_field +4 ].hh.b0 :=eqtb[int_base+ cur_fam_code].int  
else   mem[   cur_list.tail_field +4 ].hh.b0 :=(cur_val div 256) mod 16;
scan_math( cur_list.tail_field +1 );
end;


procedure append_choices;
begin begin  mem[ cur_list.tail_field ].hh.rh := new_choice; cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ; incr(save_ptr); save_stack[save_ptr+- 1].int :=0;
push_math(math_choice_group); scan_left_brace;
end;


{ \4 }
{ Declare the function called |fin_mlist| }
function fin_mlist( p:halfword ):halfword ;
var q:halfword ; {the mlist to return}
begin if cur_list.aux_field .int <>-{0xfffffff=}268435455   then 
{ Compleat the incompleat noad }
begin  mem[   cur_list.aux_field .int +3 ].hh.rh :=sub_mlist;
 mem[   cur_list.aux_field .int +3 ].hh.lh := mem[ cur_list.head_field ].hh.rh ;
if p=-{0xfffffff=}268435455   then q:=cur_list.aux_field .int 
else  begin q:= mem[   cur_list.aux_field .int +2 ].hh.lh ;
  if  mem[ q].hh.b0 <>left_noad then confusion({"right"=}891);
{ \xref[this can't happen right][\quad right] }
   mem[   cur_list.aux_field .int +2 ].hh.lh := mem[ q].hh.rh ;
   mem[ q].hh.rh :=cur_list.aux_field .int ;  mem[ cur_list.aux_field .int ].hh.rh :=p;
  end;
end


else  begin  mem[ cur_list.tail_field ].hh.rh :=p; q:= mem[ cur_list.head_field ].hh.rh ;
  end;
pop_nest; fin_mlist:=q;
end;

{  } 

procedure build_choices;
label exit;
var p:halfword ; {the current mlist}
begin unsave; p:=fin_mlist(-{0xfffffff=}268435455  );
case save_stack[save_ptr+- 1].int  of
0: mem[  cur_list.tail_field + 1].hh.lh  :=p;
1: mem[  cur_list.tail_field + 1].hh.rh  :=p;
2: mem[  cur_list.tail_field + 2].hh.lh  :=p;
3:begin  mem[  cur_list.tail_field + 2].hh.rh  :=p; decr(save_ptr);  goto exit ;
  end;
end; {there are no other cases}
incr(save_stack[save_ptr+- 1].int ); push_math(math_choice_group); scan_left_brace;
exit:end;


procedure sub_sup;
var t:small_number; {type of previous sub/superscript}
 p:halfword ; {field to be filled by |scan_math|}
begin t:=empty; p:=-{0xfffffff=}268435455  ;
if cur_list.tail_field <>cur_list.head_field  then if ( mem[  cur_list.tail_field ].hh.b0 >=ord_noad)and( mem[  cur_list.tail_field ].hh.b0 <left_noad)  then
  begin p:= cur_list.tail_field +2 +cur_cmd-sup_mark; {|supscr| or |subscr|}
  t:= mem[ p].hh.rh ;
  end;
if (p=-{0xfffffff=}268435455  )or(t<>empty) then 
{ Insert a dummy noad to be sub/superscripted }
begin begin  mem[ cur_list.tail_field ].hh.rh := new_noad; cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
p:= cur_list.tail_field +2 +cur_cmd-sup_mark; {|supscr| or |subscr|}
if t<>empty then
  begin if cur_cmd=sup_mark then
    begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Double superscript"=} 1156); end ;
{ \xref[Double superscript] }
     begin help_ptr:=1; help_line[0]:={"I treat `x^1^2' essentially like `x^1[]^2'."=} 1157; end ;
    end
  else  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Double subscript"=} 1158); end ;
{ \xref[Double subscript] }
     begin help_ptr:=1; help_line[0]:={"I treat `x_1_2' essentially like `x_1[]_2'."=} 1159; end ;
    end;
  error;
  end;
end

;
scan_math(p);
end;


procedure math_fraction;
var c:small_number; {the type of generalized fraction we are scanning}
begin c:=cur_chr;
if cur_list.aux_field .int <>-{0xfffffff=}268435455   then
  
{ Ignore the fraction operation and complain about this ambiguous case }
begin if c>=delimited_code then
  begin scan_delimiter(mem_top-12 ,false); scan_delimiter(mem_top-12 ,false);
  end;
if c mod delimited_code=above_code then scan_dimen(false,false,false) ;
begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Ambiguous; you need another [ and ]"=} 1166); end ;
{ \xref[Ambiguous...] }
 begin help_ptr:=3; help_line[2]:={"I'm ignoring this fraction specification, since I don't"=} 1167; help_line[1]:={"know whether a construction like `x \over y \over z'"=} 1168; help_line[0]:={"means `[x \over y] \over z' or `x \over [y \over z]'."=} 1169; end ;
error;
end


else  begin cur_list.aux_field .int :=get_node(fraction_noad_size);
   mem[ cur_list.aux_field .int ].hh.b0 :=fraction_noad;
   mem[ cur_list.aux_field .int ].hh.b1 :=normal;
   mem[   cur_list.aux_field .int +2 ].hh.rh :=sub_mlist;
   mem[   cur_list.aux_field .int +2 ].hh.lh := mem[ cur_list.head_field ].hh.rh ;
  mem[ cur_list.aux_field .int +3 ].hh:=empty_field;
  mem[ cur_list.aux_field .int +4 ].qqqq:=null_delimiter;
  mem[ cur_list.aux_field .int +5 ].qqqq:=null_delimiter;

   mem[ cur_list.head_field ].hh.rh :=-{0xfffffff=}268435455  ; cur_list.tail_field :=cur_list.head_field ;
  
{ Use code |c| to distinguish between generalized fractions }
if c>=delimited_code then
  begin scan_delimiter( cur_list.aux_field .int +4 ,false);
  scan_delimiter( cur_list.aux_field .int +5 ,false);
  end;
case c mod delimited_code of
above_code: begin scan_dimen(false,false,false) ;
   mem[ cur_list.aux_field .int +width_offset].int  :=cur_val;
  end;
over_code:  mem[ cur_list.aux_field .int +width_offset].int  :={010000000000=}1073741824 ;
atop_code:  mem[ cur_list.aux_field .int +width_offset].int  :=0;
end {there are no other cases}

;
  end;
end;


procedure math_left_right;
var t:small_number; {|left_noad| or |right_noad|}
 p:halfword ; {new noad}
begin t:=cur_chr;
if (t=right_noad)and(cur_group<>math_left_group) then
  
{ Try to recover from mismatched \.[\\right] }
begin if cur_group=math_shift_group then
  begin scan_delimiter(mem_top-12 ,false);
  begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Extra "=} 787); end ; print_esc({"right"=}891);
{ \xref[Extra \\right.] }
   begin help_ptr:=1; help_line[0]:={"I'm ignoring a \right that had no matching \left."=} 1170; end ;
  error;
  end
else off_save;
end


else  begin p:=new_noad;  mem[ p].hh.b0 :=t;
  scan_delimiter( p+1 ,false);
  if t=left_noad then
    begin push_math(math_left_group);  mem[ cur_list.head_field ].hh.rh :=p; cur_list.tail_field :=p;
    end
  else  begin p:=fin_mlist(p); unsave; {end of |math_left_group|}
    begin  mem[ cur_list.tail_field ].hh.rh := new_noad; cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;  mem[ cur_list.tail_field ].hh.b0 :=inner_noad;
     mem[   cur_list.tail_field +1 ].hh.rh :=sub_mlist;
     mem[   cur_list.tail_field +1 ].hh.lh :=p;
    end;
  end;
end;


procedure after_math;
var l:boolean; {`\.[\\leqno]' instead of `\.[\\eqno]'}
 danger:boolean; {not enough symbol fonts are present}
 m:integer; {|mmode| or |-mmode|}
 p:halfword ; {the formula}
 a:halfword ; {box containing equation number}

{ Local variables for finishing a displayed formula }
 b:halfword ; {box containing the equation}
 w:scaled; {width of the equation}
 z:scaled; {width of the line}
 e:scaled; {width of equation number}
 q:scaled; {width of equation number plus space to separate from equation}
 d:scaled; {displacement of equation in the line}
 s:scaled; {move the line right this much}
 g1, g2:small_number; {glue parameter codes for before and after}
 r:halfword ; {kern node used to position the display}
 t:halfword ; {tail of adjustment list}

 
begin danger:=false;

{ Check that the necessary fonts for math symbols are present; if not, flush the current math lists and set |danger:=true| }
if (font_params[ eqtb[  math_font_base+   2+   text_size].hh.rh   ]<total_mathsy_params)or 
   (font_params[ eqtb[  math_font_base+   2+   script_size].hh.rh   ]<total_mathsy_params)or 
   (font_params[ eqtb[  math_font_base+   2+   script_script_size].hh.rh   ]<total_mathsy_params) then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Math formula deleted: Insufficient symbol fonts"=} 1171); end ;

{ \xref[Math formula deleted...] }
   begin help_ptr:=3; help_line[2]:={"Sorry, but I can't typeset math unless \textfont 2"=} 1172; help_line[1]:={"and \scriptfont 2 and \scriptscriptfont 2 have all"=} 1173; help_line[0]:={"the \fontdimen values needed in math symbol fonts."=} 1174; end ;
  error; flush_math; danger:=true;
  end
else if (font_params[ eqtb[  math_font_base+   3+   text_size].hh.rh   ]<total_mathex_params)or 
   (font_params[ eqtb[  math_font_base+   3+   script_size].hh.rh   ]<total_mathex_params)or 
   (font_params[ eqtb[  math_font_base+   3+   script_script_size].hh.rh   ]<total_mathex_params) then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Math formula deleted: Insufficient extension fonts"=} 1175); end ;

   begin help_ptr:=3; help_line[2]:={"Sorry, but I can't typeset math unless \textfont 3"=} 1176; help_line[1]:={"and \scriptfont 3 and \scriptscriptfont 3 have all"=} 1177; help_line[0]:={"the \fontdimen values needed in math extension fonts."=} 1178; end ;
  error; flush_math; danger:=true;
  end

;
m:=cur_list.mode_field ; l:=false; p:=fin_mlist(-{0xfffffff=}268435455  ); {this pops the nest}
if cur_list.mode_field =-m then {end of equation number}
  begin 
{ Check that another \.\$ follows }
begin get_x_token;
if cur_cmd<>math_shift then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Display math should end with $$"=} 1179); end ;
{ \xref[Display math...with \$\$] }
   begin help_ptr:=2; help_line[1]:={"The `$' that I just saw supposedly matches a previous `$$'."=} 1180; help_line[0]:={"So I shall assume that you typed `$$' both times."=} 1181; end ;
  back_error;
  end;
end

;
  cur_mlist:=p; cur_style:=text_style; mlist_penalties:=false;
  mlist_to_hlist; a:=hpack( mem[ mem_top-3 ].hh.rh ,0,additional );
  unsave; decr(save_ptr); {now |cur_group=math_shift_group|}
  if save_stack[save_ptr+ 0].int =1 then l:=true;
  danger:=false;
  
{ Check that the necessary fonts for math symbols are present; if not, flush the current math lists and set |danger:=true| }
if (font_params[ eqtb[  math_font_base+   2+   text_size].hh.rh   ]<total_mathsy_params)or 
   (font_params[ eqtb[  math_font_base+   2+   script_size].hh.rh   ]<total_mathsy_params)or 
   (font_params[ eqtb[  math_font_base+   2+   script_script_size].hh.rh   ]<total_mathsy_params) then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Math formula deleted: Insufficient symbol fonts"=} 1171); end ;

{ \xref[Math formula deleted...] }
   begin help_ptr:=3; help_line[2]:={"Sorry, but I can't typeset math unless \textfont 2"=} 1172; help_line[1]:={"and \scriptfont 2 and \scriptscriptfont 2 have all"=} 1173; help_line[0]:={"the \fontdimen values needed in math symbol fonts."=} 1174; end ;
  error; flush_math; danger:=true;
  end
else if (font_params[ eqtb[  math_font_base+   3+   text_size].hh.rh   ]<total_mathex_params)or 
   (font_params[ eqtb[  math_font_base+   3+   script_size].hh.rh   ]<total_mathex_params)or 
   (font_params[ eqtb[  math_font_base+   3+   script_script_size].hh.rh   ]<total_mathex_params) then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Math formula deleted: Insufficient extension fonts"=} 1175); end ;

   begin help_ptr:=3; help_line[2]:={"Sorry, but I can't typeset math unless \textfont 3"=} 1176; help_line[1]:={"and \scriptfont 3 and \scriptscriptfont 3 have all"=} 1177; help_line[0]:={"the \fontdimen values needed in math extension fonts."=} 1178; end ;
  error; flush_math; danger:=true;
  end

;
  m:=cur_list.mode_field ; p:=fin_mlist(-{0xfffffff=}268435455  );
  end
else a:=-{0xfffffff=}268435455  ;
if m<0 then 
{ Finish math in text }
begin begin  mem[ cur_list.tail_field ].hh.rh := new_math( eqtb[dimen_base+ math_surround_code].int   , before); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
cur_mlist:=p; cur_style:=text_style; mlist_penalties:=(cur_list.mode_field >0); mlist_to_hlist;
 mem[ cur_list.tail_field ].hh.rh := mem[ mem_top-3 ].hh.rh ;
while  mem[ cur_list.tail_field ].hh.rh <>-{0xfffffff=}268435455   do cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ;
begin  mem[ cur_list.tail_field ].hh.rh := new_math( eqtb[dimen_base+ math_surround_code].int   , after); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
cur_list.aux_field .hh.lh :=1000; unsave;
end


else  begin if a=-{0xfffffff=}268435455   then 
{ Check that another \.\$ follows }
begin get_x_token;
if cur_cmd<>math_shift then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Display math should end with $$"=} 1179); end ;
{ \xref[Display math...with \$\$] }
   begin help_ptr:=2; help_line[1]:={"The `$' that I just saw supposedly matches a previous `$$'."=} 1180; help_line[0]:={"So I shall assume that you typed `$$' both times."=} 1181; end ;
  back_error;
  end;
end

;
  
{ Finish displayed math }
cur_mlist:=p; cur_style:=display_style; mlist_penalties:=false;
mlist_to_hlist; p:= mem[ mem_top-3 ].hh.rh ;

adjust_tail:=mem_top-5 ; b:=hpack(p,0,additional ); p:=  mem[  b+ list_offset].hh.rh  ;
t:=adjust_tail; adjust_tail:=-{0xfffffff=}268435455  ;

w:= mem[ b+width_offset].int  ; z:=eqtb[dimen_base+ display_width_code].int   ; s:=eqtb[dimen_base+ display_indent_code].int   ;
if (a=-{0xfffffff=}268435455  )or danger then
  begin e:=0; q:=0;
  end
else  begin e:= mem[ a+width_offset].int  ; q:=e+font_info[ 6+param_base[ eqtb[  math_font_base+   2+    text_size].hh.rh   ]].int  ;
  end;
if w+q>z then
  
{ Squeeze the equation as much as possible; if there is an equation number that should go on a separate line by itself, set~|e:=0| }
begin if (e<>0)and((w-total_shrink[normal]+q<=z)or 
   (total_shrink[fil]<>0)or(total_shrink[fill]<>0)or
   (total_shrink[filll]<>0)) then
  begin free_node(b,box_node_size);
  b:=hpack(p,z-q,exactly);
  end
else  begin e:=0;
  if w>z then
    begin free_node(b,box_node_size);
    b:=hpack(p,z,exactly);
    end;
  end;
w:= mem[ b+width_offset].int  ;
end

;

{ Determine the displacement, |d|, of the left edge of the equation, with respect to the line size |z|, assuming that |l=false| }
d:=half(z-w);
if (e>0)and(d<2*e) then {too close}
  begin d:=half(z-w-e);
  if p<>-{0xfffffff=}268435455   then if not  ( p>=hi_mem_min)  then if  mem[ p].hh.b0 =glue_node then d:=0;
  end

;

{ Append the glue or equation number preceding the display }
begin  mem[ cur_list.tail_field ].hh.rh := new_penalty( eqtb[int_base+ pre_display_penalty_code].int  ); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;

if (d+s<=eqtb[dimen_base+ pre_display_size_code].int   )or l then {not enough clearance}
  begin g1:=above_display_skip_code; g2:=below_display_skip_code;
  end
else  begin g1:=above_display_short_skip_code;
  g2:=below_display_short_skip_code;
  end;
if l and(e=0) then {it follows that |type(a)=hlist_node|}
  begin  mem[ a+4].int  :=s; append_to_vlist(a);
  begin  mem[ cur_list.tail_field ].hh.rh := new_penalty( inf_penalty); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
  end
else begin  mem[ cur_list.tail_field ].hh.rh := new_param_glue( g1); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end 

;

{ Append the display and perhaps also the equation number }
if e<>0 then
  begin r:=new_kern(z-w-e-d);
  if l then
    begin  mem[ a].hh.rh :=r;  mem[ r].hh.rh :=b; b:=a; d:=0;
    end
  else  begin  mem[ b].hh.rh :=r;  mem[ r].hh.rh :=a;
    end;
  b:=hpack(b,0,additional );
  end;
 mem[ b+4].int  :=s+d; append_to_vlist(b)

;

{ Append the glue or equation number following the display }
if (a<>-{0xfffffff=}268435455  )and(e=0)and not l then
  begin begin  mem[ cur_list.tail_field ].hh.rh := new_penalty( inf_penalty); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
   mem[ a+4].int  :=s+z- mem[ a+width_offset].int  ;
  append_to_vlist(a);
  g2:=0;
  end;
if t<>mem_top-5  then {migrating material comes after equation number}
  begin  mem[ cur_list.tail_field ].hh.rh := mem[ mem_top-5 ].hh.rh ; cur_list.tail_field :=t;
  end;
begin  mem[ cur_list.tail_field ].hh.rh := new_penalty( eqtb[int_base+ post_display_penalty_code].int  ); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
if g2>0 then begin  mem[ cur_list.tail_field ].hh.rh := new_param_glue( g2); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end 

;
resume_after_display

;
  end;
end;


procedure resume_after_display;
begin if cur_group<>math_shift_group then confusion({"display"=}1182);
{ \xref[this can't happen display][\quad display] }
unsave; cur_list.pg_field :=cur_list.pg_field +3;
push_nest; cur_list.mode_field :=hmode; cur_list.aux_field .hh.lh :=1000; if eqtb[int_base+ language_code].int  <=0 then cur_lang:=0 else if eqtb[int_base+ language_code].int  >255 then cur_lang:=0 else cur_lang:=eqtb[int_base+ language_code].int   ; cur_list.aux_field .hh.rh :=cur_lang;
cur_list.pg_field :=(norm_min(eqtb[int_base+ left_hyphen_min_code].int  )*{0100=}64+norm_min(eqtb[int_base+ right_hyphen_min_code].int  ))
             *{0200000=}65536+cur_lang;

{ Scan an optional space }
begin get_x_token; if cur_cmd<>spacer then back_input;
end

;
if nest_ptr=1 then build_page;
end;


{ \4 }
{ Declare subprocedures for |prefixed_command| }
procedure get_r_token;
label restart;
begin restart: repeat get_token;
until cur_tok<>space_token;
if (cur_cs=0)or(cur_cs>eqtb_top)or
  ((cur_cs>frozen_control_sequence)and(cur_cs<=eqtb_size)) then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Missing control sequence inserted"=} 1197); end ;
{ \xref[Missing control...] }
   begin help_ptr:=5; help_line[4]:={"Please don't say `\def cs[...]', say `\def\cs[...]'."=} 1198; help_line[3]:={"I've inserted an inaccessible control sequence so that your"=} 1199; help_line[2]:={"definition will be completed without mixing me up too badly."=} 1200; help_line[1]:={"You can recover graciously from this error, if you're"=} 1201; help_line[0]:={"careful; see exercise 27.2 in The TeXbook."=} 1202; end ;
{ \xref[TeXbook][\sl The \TeX book] }
  if cur_cs=0 then back_input;
  cur_tok:={07777=}4095 +frozen_protection; ins_error; goto restart;
  end;
end;


procedure trap_zero_glue;
begin if ( mem[ cur_val+width_offset].int  =0)and( mem[ cur_val+2].int  =0)and( mem[ cur_val+3].int  =0) then
  begin incr(  mem[   mem_bot ].hh.rh  ) ;
  delete_glue_ref(cur_val); cur_val:=mem_bot ;
  end;
end;


procedure do_register_command( a:small_number);
label found,exit;
var l, q, r, s:halfword ; {for list manipulation}
 p:int_val..mu_val; {type of register involved}
begin q:=cur_cmd;

{ Compute the register location |l| and its type |p|; but |return| if invalid }
begin if q<>register then
  begin get_x_token;
  if (cur_cmd>=assign_int)and(cur_cmd<=assign_mu_glue) then
    begin l:=cur_chr; p:=cur_cmd-assign_int; goto found;
    end;
  if cur_cmd<>register then
    begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"You can't use `"=} 695); end ; print_cmd_chr(cur_cmd,cur_chr);
{ \xref[You can't use x after ...] }
    print({"' after "=}696); print_cmd_chr(q,0);
     begin help_ptr:=1; help_line[0]:={"I'm forgetting what you said and not changing anything."=} 1226; end ;
    error;  goto exit ;
    end;
  end;
p:=cur_chr; scan_eight_bit_int;
case p of
int_val: l:=cur_val+count_base;
dimen_val: l:=cur_val+scaled_base;
glue_val: l:=cur_val+skip_base;
mu_val: l:=cur_val+mu_skip_base;
end; {there are no other cases}
end;
found:

;
if q=register then scan_optional_equals
else if scan_keyword({"by"=}1222) then  ; {optional `\.[by]'}
{ \xref[by] }
arith_error:=false;
if q<multiply then 
{ Compute result of |register| or |advance|, put it in |cur_val| }
if p<glue_val then
  begin if p=int_val then scan_int else scan_dimen(false,false,false) ;
  if q=advance then cur_val:=cur_val+eqtb[l].int;
  end
else  begin scan_glue(p);
  if q=advance then 
{ Compute the sum of two glue specs }
begin q:=new_spec(cur_val); r:= eqtb[  l].hh.rh  ;
delete_glue_ref(cur_val);
 mem[ q+width_offset].int  := mem[ q+width_offset].int  + mem[ r+width_offset].int  ;
if  mem[ q+2].int  =0 then   mem[ q].hh.b0 :=normal;
if   mem[ q].hh.b0 =  mem[ r].hh.b0  then  mem[ q+2].int  := mem[ q+2].int  + mem[ r+2].int  
else if (  mem[ q].hh.b0 <  mem[ r].hh.b0 )and( mem[ r+2].int  <>0) then
  begin  mem[ q+2].int  := mem[ r+2].int  ;   mem[ q].hh.b0 :=  mem[ r].hh.b0 ;
  end;
if  mem[ q+3].int  =0 then   mem[ q].hh.b1 :=normal;
if   mem[ q].hh.b1 =  mem[ r].hh.b1  then  mem[ q+3].int  := mem[ q+3].int  + mem[ r+3].int  
else if (  mem[ q].hh.b1 <  mem[ r].hh.b1 )and( mem[ r+3].int  <>0) then
  begin  mem[ q+3].int  := mem[ r+3].int  ;   mem[ q].hh.b1 :=  mem[ r].hh.b1 ;
  end;
cur_val:=q;
end

;
  end


else 
{ Compute result of |multiply| or |divide|, put it in |cur_val| }
begin scan_int;
if p<glue_val then
  if q=multiply then
    if p=int_val then cur_val:=mult_and_add( eqtb[ l]. int, cur_val,0,{017777777777=}2147483647) 
    else cur_val:=mult_and_add( eqtb[ l]. int, cur_val, 0,{07777777777=}1073741823) 
  else cur_val:=x_over_n(eqtb[l].int,cur_val)
else  begin s:= eqtb[  l].hh.rh  ; r:=new_spec(s);
  if q=multiply then
    begin  mem[ r+width_offset].int  :=mult_and_add(  mem[  s+width_offset].int  , cur_val, 0,{07777777777=}1073741823) ;
     mem[ r+2].int  :=mult_and_add(  mem[  s+2].int  , cur_val, 0,{07777777777=}1073741823) ;
     mem[ r+3].int  :=mult_and_add(  mem[  s+3].int  , cur_val, 0,{07777777777=}1073741823) ;
    end
  else  begin  mem[ r+width_offset].int  :=x_over_n( mem[ s+width_offset].int  ,cur_val);
     mem[ r+2].int  :=x_over_n( mem[ s+2].int  ,cur_val);
     mem[ r+3].int  :=x_over_n( mem[ s+3].int  ,cur_val);
    end;
  cur_val:=r;
  end;
end

;
if arith_error then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Arithmetic overflow"=} 1223); end ;
{ \xref[Arithmetic overflow] }
   begin help_ptr:=2; help_line[1]:={"I can't carry out that multiplication or division,"=} 1224; help_line[0]:={"since the result is out of range."=} 1225; end ;
  if p>=glue_val then delete_glue_ref(cur_val);
  error;  goto exit ;
  end;
if p<glue_val then if (a>=4)  then geq_word_define( l, cur_val) else eq_word_define( l, cur_val) 
else  begin trap_zero_glue; if (a>=4)  then geq_define( l, glue_ref, cur_val) else eq_define( l, glue_ref, cur_val) ;
  end;
exit: end;


procedure alter_aux;
var c:halfword; {|hmode| or |vmode|}
begin if cur_chr<>abs(cur_list.mode_field ) then report_illegal_case
else  begin c:=cur_chr; scan_optional_equals;
  if c=vmode then
    begin scan_dimen(false,false,false) ; cur_list.aux_field .int  :=cur_val;
    end
  else  begin scan_int;
    if (cur_val<=0)or(cur_val>32767) then
      begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Bad space factor"=} 1229); end ;
{ \xref[Bad space factor] }
       begin help_ptr:=1; help_line[0]:={"I allow only values in the range 1..32767 here."=} 1230; end ;
      int_error(cur_val);
      end
    else cur_list.aux_field .hh.lh :=cur_val;
    end;
  end;
end;


procedure alter_prev_graf;
var p:0..nest_size; {index into |nest|}
begin nest[nest_ptr]:=cur_list; p:=nest_ptr;
while abs(nest[p].mode_field)<>vmode do decr(p);
scan_optional_equals; scan_int;
if cur_val<0 then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Bad "=} 969); end ; print_esc({"prevgraf"=}540);
{ \xref[Bad \\prevgraf] }
   begin help_ptr:=1; help_line[0]:={"I allow only nonnegative values here."=} 1231; end ;
  int_error(cur_val);
  end
else  begin nest[p].pg_field:=cur_val; cur_list:=nest[nest_ptr];
  end;
end;


procedure alter_page_so_far;
var c:0..7; {index into |page_so_far|}
begin c:=cur_chr; scan_optional_equals; scan_dimen(false,false,false) ;
page_so_far[c]:=cur_val;
end;


procedure alter_integer;
var c:0..1; {0 for \.[\\deadcycles], 1 for \.[\\insertpenalties]}
begin c:=cur_chr; scan_optional_equals; scan_int;
if c=0 then dead_cycles:=cur_val
else insert_penalties:=cur_val;
end;


procedure alter_box_dimen;
var c:small_number; {|width_offset| or |height_offset| or |depth_offset|}
 b:eight_bits; {box number}
begin c:=cur_chr; scan_eight_bit_int; b:=cur_val; scan_optional_equals;
scan_dimen(false,false,false) ;
if  eqtb[  box_base+   b].hh.rh   <>-{0xfffffff=}268435455   then mem[ eqtb[  box_base+   b].hh.rh   +c].int :=cur_val;
end;


procedure new_font( a:small_number);
label common_ending;
var u:halfword ; {user's font identifier}
 s:scaled; {stated ``at'' size, or negative of scaled magnification}
 f:internal_font_number; {runs through existing fonts}
 t:str_number; {name for the frozen font identifier}
 old_setting:0..max_selector; {holds |selector| setting}

begin if job_name=0 then open_log_file;
  {avoid confusing \.[texput] with the font name}
{ \xref[texput] }
get_r_token; u:=cur_cs;
if u>=hash_base then t:= hash[ u].rh 
else if u>=single_base then
  if u=null_cs then t:={"FONT"=}1235 else t:=u-single_base
else  begin old_setting:=selector; selector:=new_string;
  print({"FONT"=}1235); print(u-active_base); selector:=old_setting;
{ \xref[FONTx] }
   begin if pool_ptr+ 1 > pool_size then overflow({"pool size"=}257,pool_size-init_pool_ptr); { \xref[TeX capacity exceeded pool size][\quad pool size] } end ; t:=make_string;
  end;
if (a>=4)  then geq_define( u, set_font, font_base ) else eq_define( u, set_font, font_base ) ; scan_optional_equals; scan_file_name;

{ Scan the font size specification }
name_in_progress:=true; {this keeps |cur_name| from being changed}
if scan_keyword({"at"=}1236) then 
{ Put the \(p)(positive) `at' size into |s| }
begin scan_dimen(false,false,false) ; s:=cur_val;
if (s<=0)or(s>={01000000000=}134217728) then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Improper `at' size ("=} 1238); end ;
  print_scaled(s); print({"pt), replaced by 10pt"=}1239);
{ \xref[Improper `at' size...] }
   begin help_ptr:=2; help_line[1]:={"I can only handle fonts at positive sizes that are"=} 1240; help_line[0]:={"less than 2048pt, so I've changed what you said to 10pt."=} 1241; end ;
  error; s:=10* {0200000=}65536 ;
  end;
end


{ \xref[at] }
else if scan_keyword({"scaled"=}1237) then
{ \xref[scaled] }
  begin scan_int; s:=-cur_val;
  if (cur_val<=0)or(cur_val>32768) then
    begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Illegal magnification has been changed to 1000"=} 560); end ;

{ \xref[Illegal magnification...] }
     begin help_ptr:=1; help_line[0]:={"The magnification ratio must be between 1 and 32768."=} 561; end ;
    int_error(cur_val); s:=-1000;
    end;
  end
else s:=-1000;
name_in_progress:=false

;

{ If this font has already been loaded, set |f| to the internal font number and |goto common_ending| }

for f:=font_base+1 to font_ptr do
  if str_eq_str(font_name[f],cur_name)and str_eq_str(font_area[f],cur_area) then
    begin if s>0 then
      begin if s=font_size[f] then goto common_ending;
      end
    else begin arith_error:=false;
      if font_size[f]=xn_over_d(font_dsize[f],-s,1000)
      then if not arith_error
        then goto common_ending;
      end;
    end

;
f:=read_font_info(u,cur_name,cur_area,s);
common_ending:  eqtb[  u].hh.rh  :=f; eqtb[font_id_base+f]:=eqtb[u];   hash[ font_id_base+  f].rh  :=t;
end;


procedure new_interaction;
begin print_ln;
interaction:=cur_chr;
if interaction = batch_mode
then kpse_make_tex_discard_errors := 1
else kpse_make_tex_discard_errors := 0;

{ Initialize the print |selector| based on |interaction| }
if interaction=batch_mode then selector:=no_print else selector:=term_only

;
if log_opened then selector:=selector+2;
end;

{  } 

procedure prefixed_command;
label done,exit;
var a:small_number; {accumulated prefix codes so far}
 f:internal_font_number; {identifies a font}
 j:halfword; {index into a \.[\\parshape] specification}
 k:font_index; {index into |font_info|}
 p, q:halfword ; {for temporary short-term use}
 n:integer; {ditto}
 e:boolean; {should a definition be expanded? or was \.[\\let] not done?}
begin a:=0;
while cur_cmd=prefix do
  begin if not odd(a div cur_chr) then a:=a+cur_chr;
  
{ Get the next non-blank non-relax... }
repeat get_x_token;
until (cur_cmd<>spacer)and(cur_cmd<>relax)

;
  if cur_cmd<=max_non_prefixed_command then
    
{ Discard erroneous prefixes and |return| }
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"You can't use a prefix with `"=} 1192); end ;
{ \xref[You can't use a prefix with x] }
print_cmd_chr(cur_cmd,cur_chr); print_char({"'"=}39);
 begin help_ptr:=1; help_line[0]:={"I'll pretend you didn't say \long or \outer or \global."=} 1193; end ;
back_error;  goto exit ;
end

;
  end;

{ Discard the prefixes \.[\\long] and \.[\\outer] if they are irrelevant }
if (cur_cmd<>def)and(a mod 4<>0) then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"You can't use `"=} 695); end ; print_esc({"long"=}1184); print({"' or `"=}1194);
  print_esc({"outer"=}1185); print({"' with `"=}1195);
{ \xref[You can't use \\long...] }
  print_cmd_chr(cur_cmd,cur_chr); print_char({"'"=}39);
   begin help_ptr:=1; help_line[0]:={"I'll pretend you didn't say \long or \outer here."=} 1196; end ;
  error;
  end

;

{ Adjust \(f)for the setting of \.[\\globaldefs] }
if eqtb[int_base+ global_defs_code].int  <>0 then
  if eqtb[int_base+ global_defs_code].int  <0 then
    begin if (a>=4)  then a:=a-4;
    end
  else  begin if not (a>=4)  then a:=a+4;
    end

;
case cur_cmd of
{ \4 }
{ Assignments }
set_font: if (a>=4)  then geq_define( cur_font_loc, data, cur_chr) else eq_define( cur_font_loc, data, cur_chr) ;


def: begin if odd(cur_chr)and not (a>=4)  and(eqtb[int_base+ global_defs_code].int  >=0) then a:=a+4;
  e:=(cur_chr>=2); get_r_token; p:=cur_cs;
  q:=scan_toks(true,e); if (a>=4)  then geq_define( p, call+( a  mod  4), def_ref) else eq_define( p, call+( a  mod  4), def_ref) ;
  end;


let:  begin n:=cur_chr;
  get_r_token; p:=cur_cs;
  if n=normal then
    begin repeat get_token;
    until cur_cmd<>spacer;
    if cur_tok=other_token+{"="=}61 then
      begin get_token;
      if cur_cmd=spacer then get_token;
      end;
    end
  else  begin get_token; q:=cur_tok; get_token; back_input;
    cur_tok:=q; back_input; {look ahead, then back up}
    end; {note that |back_input| doesn't affect |cur_cmd|, |cur_chr|}
  if cur_cmd>=call then incr(  mem[   cur_chr].hh.lh  ) ;
  if (a>=4)  then geq_define( p, cur_cmd, cur_chr) else eq_define( p, cur_cmd, cur_chr) ;
  end;


shorthand_def: if cur_chr=char_sub_def_code then
 begin scan_char_num; p:=char_sub_code_base+cur_val; scan_optional_equals;
  scan_char_num; n:=cur_val; {accent character in substitution}
  scan_char_num;
  if (eqtb[int_base+ tracing_char_sub_def_code].int  >0) then
    begin begin_diagnostic; print_nl({"New character substitution: "=}1214);
     print (p-char_sub_code_base); print({" = "=}1215);
     print (n); print_char({" "=}32);
     print (cur_val); end_diagnostic(false);
    end;
  n:=n*256+cur_val;
  if (a>=4)  then geq_define( p, data,   n ) else eq_define( p, data,   n ) ;
  if (p-char_sub_code_base)<eqtb[int_base+ char_sub_def_min_code].int   then
    if (a>=4)  then geq_word_define( int_base+ char_sub_def_min_code, p- char_sub_code_base) else eq_word_define( int_base+ char_sub_def_min_code, p- char_sub_code_base) ;
  if (p-char_sub_code_base)>eqtb[int_base+ char_sub_def_max_code].int   then
    if (a>=4)  then geq_word_define( int_base+ char_sub_def_max_code, p- char_sub_code_base) else eq_word_define( int_base+ char_sub_def_max_code, p- char_sub_code_base) ;
 end
else begin n:=cur_chr; get_r_token; p:=cur_cs; if (a>=4)  then geq_define( p, relax, 256) else eq_define( p, relax, 256) ;
  scan_optional_equals;
  case n of
  char_def_code: begin scan_char_num; if (a>=4)  then geq_define( p, char_given, cur_val) else eq_define( p, char_given, cur_val) ;
    end;
  math_char_def_code: begin scan_fifteen_bit_int; if (a>=4)  then geq_define( p, math_given, cur_val) else eq_define( p, math_given, cur_val) ;
    end;
   else  begin scan_eight_bit_int;
    case n of
    count_def_code: if (a>=4)  then geq_define( p, assign_int, count_base+ cur_val) else eq_define( p, assign_int, count_base+ cur_val) ;
    dimen_def_code: if (a>=4)  then geq_define( p, assign_dimen, scaled_base+ cur_val) else eq_define( p, assign_dimen, scaled_base+ cur_val) ;
    skip_def_code: if (a>=4)  then geq_define( p, assign_glue, skip_base+ cur_val) else eq_define( p, assign_glue, skip_base+ cur_val) ;
    mu_skip_def_code: if (a>=4)  then geq_define( p, assign_mu_glue, mu_skip_base+ cur_val) else eq_define( p, assign_mu_glue, mu_skip_base+ cur_val) ;
    toks_def_code: if (a>=4)  then geq_define( p, assign_toks, toks_base+ cur_val) else eq_define( p, assign_toks, toks_base+ cur_val) ;
    end; {there are no other cases}
    end
   end ;
  end;


read_to_cs: begin scan_int; n:=cur_val;
  if not scan_keyword({"to"=}856) then
{ \xref[to] }
    begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Missing `to' inserted"=} 1086); end ;
{ \xref[Missing `to'...] }
     begin help_ptr:=2; help_line[1]:={"You should have said `\read<number> to \cs'."=} 1216; help_line[0]:={"I'm going to look for the \cs now."=} 1217; end ; error;
    end;
  get_r_token;
  p:=cur_cs; read_toks(n,p); if (a>=4)  then geq_define( p, call, cur_val) else eq_define( p, call, cur_val) ;
  end;


toks_register,assign_toks: begin q:=cur_cs;
  if cur_cmd=toks_register then
    begin scan_eight_bit_int; p:=toks_base+cur_val;
    end
  else p:=cur_chr; {|p=every_par_loc| or |output_routine_loc| or \dots}
  scan_optional_equals;
  
{ Get the next non-blank non-relax non-call token }
repeat get_x_token;
until (cur_cmd<>spacer)and(cur_cmd<>relax)

;
  if cur_cmd<>left_brace then 
{ If the right-hand side is a token parameter or token register, finish the assignment and |goto done| }
begin if cur_cmd=toks_register then
  begin scan_eight_bit_int; cur_cmd:=assign_toks; cur_chr:=toks_base+cur_val;
  end;
if cur_cmd=assign_toks then
  begin q:= eqtb[  cur_chr].hh.rh  ;
  if q=-{0xfffffff=}268435455   then if (a>=4)  then geq_define( p, undefined_cs, -{0xfffffff=}268435455  ) else eq_define( p, undefined_cs, -{0xfffffff=}268435455  ) 
  else  begin incr(  mem[   q].hh.lh  ) ; if (a>=4)  then geq_define( p, call, q) else eq_define( p, call, q) ;
    end;
  goto done;
  end;
end

;
  back_input; cur_cs:=q; q:=scan_toks(false,false);
  if  mem[ def_ref].hh.rh =-{0xfffffff=}268435455   then {empty list: revert to the default}
    begin if (a>=4)  then geq_define( p, undefined_cs, -{0xfffffff=}268435455  ) else eq_define( p, undefined_cs, -{0xfffffff=}268435455  ) ;  begin  mem[  def_ref].hh.rh :=avail; avail:= def_ref; ifdef('STAT')  decr(dyn_used); endif('STAT')  end ;
    end
  else  begin if p=output_routine_loc then {enclose in curlies}
      begin  mem[ q].hh.rh :=get_avail; q:= mem[ q].hh.rh ;
       mem[ q].hh.lh :=right_brace_token+{"]"=}125;
      q:=get_avail;  mem[ q].hh.lh :=left_brace_token+{"["=}123;
       mem[ q].hh.rh := mem[ def_ref].hh.rh ;  mem[ def_ref].hh.rh :=q;
      end;
    if (a>=4)  then geq_define( p, call, def_ref) else eq_define( p, call, def_ref) ;
    end;
  end;


assign_int: begin p:=cur_chr; scan_optional_equals; scan_int;
  if (a>=4)  then geq_word_define( p, cur_val) else eq_word_define( p, cur_val) ;
  end;
assign_dimen: begin p:=cur_chr; scan_optional_equals;
  scan_dimen(false,false,false) ; if (a>=4)  then geq_word_define( p, cur_val) else eq_word_define( p, cur_val) ;
  end;
assign_glue,assign_mu_glue: begin p:=cur_chr; n:=cur_cmd; scan_optional_equals;
  if n=assign_mu_glue then scan_glue(mu_val) else scan_glue(glue_val);
  trap_zero_glue;
  if (a>=4)  then geq_define( p, glue_ref, cur_val) else eq_define( p, glue_ref, cur_val) ;
  end;


def_code: begin 
{ Let |n| be the largest legal code value, based on |cur_chr| }
if cur_chr=cat_code_base then n:=max_char_code
else if cur_chr=math_code_base then n:={0100000=}32768
else if cur_chr=sf_code_base then n:={077777=}32767
else if cur_chr=del_code_base then n:={077777777=}16777215
else n:=255

;
  p:=cur_chr; scan_char_num; p:=p+cur_val; scan_optional_equals;
  scan_int;
  if ((cur_val<0)and(p<del_code_base))or(cur_val>n) then
    begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Invalid code ("=} 1218); end ; print_int(cur_val);
{ \xref[Invalid code] }
    if p<del_code_base then print({"), should be in the range 0.."=}1219)
    else print({"), should be at most "=}1220);
    print_int(n);
     begin help_ptr:=1; help_line[0]:={"I'm going to use 0 instead of that illegal code value."=} 1221; end ;

    error; cur_val:=0;
    end;
  if p<math_code_base then if (a>=4)  then geq_define( p, data, cur_val) else eq_define( p, data, cur_val) 
  else if p<del_code_base then if (a>=4)  then geq_define( p, data,   cur_val ) else eq_define( p, data,   cur_val ) 
  else if (a>=4)  then geq_word_define( p, cur_val) else eq_word_define( p, cur_val) ;
  end;


def_family: begin p:=cur_chr; scan_four_bit_int; p:=p+cur_val;
  scan_optional_equals; scan_font_ident; if (a>=4)  then geq_define( p, data, cur_val) else eq_define( p, data, cur_val) ;
  end;


register,advance,multiply,divide: do_register_command(a);


set_box: begin scan_eight_bit_int;
  if (a>=4)  then n:=256+cur_val else n:=cur_val;
  scan_optional_equals;
  if set_box_allowed then scan_box({010000000000=}1073741824 +n)
  else begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Improper "=} 690); end ; print_esc({"setbox"=}544);
{ \xref[Improper \\setbox] }
     begin help_ptr:=2; help_line[1]:={"Sorry, \setbox is not allowed after \halign in a display,"=} 1227; help_line[0]:={"or between \accent and an accented character."=} 1228; end ; error;
    end;
  end;


set_aux:alter_aux;
set_prev_graf:alter_prev_graf;
set_page_dimen:alter_page_so_far;
set_page_int:alter_integer;
set_box_dimen:alter_box_dimen;


set_shape: begin scan_optional_equals; scan_int; n:=cur_val;
  if n<=0 then p:=-{0xfffffff=}268435455  
  else  begin p:=get_node(2*n+1);  mem[ p].hh.lh :=n;
    for j:=1 to n do
      begin scan_dimen(false,false,false) ;
      mem[p+2*j-1].int :=cur_val; {indentation}
      scan_dimen(false,false,false) ;
      mem[p+2*j].int :=cur_val; {width}
      end;
    end;
  if (a>=4)  then geq_define( par_shape_loc, shape_ref, p) else eq_define( par_shape_loc, shape_ref, p) ;
  end;


hyph_data: if cur_chr=1 then
    begin  ifdef('INITEX')  if ini_version then begin  new_patterns; goto done;  end; endif('INITEX')  

    begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Patterns can be loaded only by INITEX"=} 1232); end ;
{ \xref[Patterns can be...] }
    help_ptr:=0 ; error;
    repeat get_token; until cur_cmd=right_brace; {flush the patterns}
     goto exit ;
    end
  else  begin new_hyph_exceptions; goto done;
    end;


assign_font_dimen: begin find_font_dimen(true); k:=cur_val;
  scan_optional_equals; scan_dimen(false,false,false) ; font_info[k].int :=cur_val;
  end;
assign_font_int: begin n:=cur_chr; scan_font_ident; f:=cur_val;
  scan_optional_equals; scan_int;
  if n=0 then hyphen_char[f]:=cur_val else skew_char[f]:=cur_val;
  end;


def_font: new_font(a);


set_interaction: new_interaction;

 
 else  confusion({"prefix"=}1191)
{ \xref[this can't happen prefix][\quad prefix] }
 end ;
done: 
{ Insert a token saved by \.[\\afterassignment], if any }
if after_token<>0 then
  begin cur_tok:=after_token; back_input; after_token:=0;
  end

;
exit:end;


procedure do_assignments;
label exit;
begin  while true do  begin 
{ Get the next non-blank non-relax... }
repeat get_x_token;
until (cur_cmd<>spacer)and(cur_cmd<>relax)

;
  if cur_cmd<=max_non_prefixed_command then  goto exit ;
  set_box_allowed:=false; prefixed_command; set_box_allowed:=true;
  end;
exit:end;


procedure open_or_close_in;
var c:0..1; {1 for \.[\\openin], 0 for \.[\\closein]}
 n:0..15; {stream number}
begin c:=cur_chr; scan_four_bit_int; n:=cur_val;
if read_open[n]<>closed then
  begin a_close(read_file[n]); read_open[n]:=closed;
  end;
if c<>0 then
  begin scan_optional_equals; scan_file_name;
  pack_file_name(cur_name,cur_area,cur_ext) ;
  tex_input_type:=0; {Tell |open_input| we are \.[\\openin].}
  if kpse_in_name_ok(stringcast(name_of_file+1))
     and a_open_in(read_file[n], kpse_tex_format) then
    read_open[n]:=just_open;
  end;
end;


procedure issue_message;
var old_setting:0..max_selector; {holds |selector| setting}
 c:0..1; {identifies \.[\\message] and \.[\\errmessage]}
 s:str_number; {the message}
begin c:=cur_chr;  mem[ mem_top-12 ].hh.rh :=scan_toks(false,true);
old_setting:=selector; selector:=new_string;
token_show(def_ref); selector:=old_setting;
flush_list(def_ref);
 begin if pool_ptr+ 1 > pool_size then overflow({"pool size"=}257,pool_size-init_pool_ptr); { \xref[TeX capacity exceeded pool size][\quad pool size] } end ; s:=make_string;
if c=0 then 
{ Print string |s| on the terminal }
begin if term_offset+(str_start[ s+1]-str_start[ s]) >max_print_line-2 then print_ln
else if (term_offset>0)or(file_offset>0) then print_char({" "=}32);
slow_print(s);  fflush (stdout ) ;
end


else 
{ Print string |s| as an error message }
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({""=} 335); end ; slow_print(s);
if  eqtb[  err_help_loc].hh.rh   <>-{0xfffffff=}268435455   then use_err_help:=true
else if long_help_seen then  begin help_ptr:=1; help_line[0]:={"(That was another \errmessage.)"=} 1248; end 
else  begin if interaction<error_stop_mode then long_help_seen:=true;
   begin help_ptr:=4; help_line[3]:={"This error message was generated by an \errmessage"=} 1249; help_line[2]:={"command, so I can't give any explicit help."=} 1250; help_line[1]:={"Pretend that you're Hercule Poirot: Examine all clues,"=} 1251; help_line[0]:={"and deduce the truth by order and method."=} 1252; end ;
  end;
error; use_err_help:=false;
end

;
begin decr(str_ptr); pool_ptr:=str_start[str_ptr]; end ;
end;


procedure shift_case;
var b:halfword ; {|lc_code_base| or |uc_code_base|}
 p:halfword ; {runs through the token list}
 t:halfword; {token}
 c:eight_bits; {character code}
begin b:=cur_chr; p:=scan_toks(false,false); p:= mem[ def_ref].hh.rh ;
while p<>-{0xfffffff=}268435455   do
  begin 
{ Change the case of the token in |p|, if a change is appropriate }
t:= mem[ p].hh.lh ;
if t<{07777=}4095 +single_base then
  begin c:=t mod 256;
  if  eqtb[  b+  c].hh.rh  <>0 then  mem[ p].hh.lh :=t-c+ eqtb[  b+  c].hh.rh  ;
  end

;
  p:= mem[ p].hh.rh ;
  end;
begin_token_list(  mem[  def_ref].hh.rh ,backed_up) ;  begin  mem[  def_ref].hh.rh :=avail; avail:= def_ref; ifdef('STAT')  decr(dyn_used); endif('STAT')  end ; {omit reference count}
end;


procedure show_whatever;
label common_ending;
var p:halfword ; {tail of a token list to show}
begin
if p=0 then;
case cur_chr of
show_lists_code: begin begin_diagnostic; show_activities;
  end;
show_box_code: 
{ Show the current contents of a box }
begin scan_eight_bit_int; begin_diagnostic;
print_nl({"> \box"=}1270); print_int(cur_val); print_char({"="=}61);
if  eqtb[  box_base+   cur_val].hh.rh   =-{0xfffffff=}268435455   then print({"void"=}415)
else show_box( eqtb[  box_base+   cur_val].hh.rh   );
end

;
show_code: 
{ Show the current meaning of a token, then |goto common_ending| }
begin get_token;
if interaction=error_stop_mode then    ;
print_nl({"> "=}1264);
if cur_cs<>0 then
  begin sprint_cs(cur_cs); print_char({"="=}61);
  end;
print_meaning; goto common_ending;
end

;
 else  
{ Show the current value of some parameter or register, then |goto common_ending| }
begin p:=the_toks;
if interaction=error_stop_mode then    ;
print_nl({"> "=}1264); token_show(mem_top-3 );
flush_list( mem[ mem_top-3 ].hh.rh ); goto common_ending;
end


 end ;


{ Complete a potentially long \.[\\show] command }
end_diagnostic(true); begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"OK"=} 1271); end ;
{ \xref[OK] }
if selector=term_and_log then if eqtb[int_base+ tracing_online_code].int  <=0 then
  begin selector:=term_only; print({" (see the transcript file)"=}1272);
  selector:=term_and_log;
  end

;
common_ending: if interaction<error_stop_mode then
  begin help_ptr:=0 ; decr(error_count);
  end
else if eqtb[int_base+ tracing_online_code].int  >0 then
  begin{  } 

   begin help_ptr:=3; help_line[2]:={"This isn't an error message; I'm just \showing something."=} 1259; help_line[1]:={"Type `I\show...' to show more (e.g., \show\cs,"=} 1260; help_line[0]:={"\showthe\count10, \showbox255, \showlists)."=} 1261; end ;
  end
else  begin{  } 

   begin help_ptr:=5; help_line[4]:={"This isn't an error message; I'm just \showing something."=} 1259; help_line[3]:={"Type `I\show...' to show more (e.g., \show\cs,"=} 1260; help_line[2]:={"\showthe\count10, \showbox255, \showlists)."=} 1261; help_line[1]:={"And type `I\tracingonline=1\show...' to show boxes and"=} 1262; help_line[0]:={"lists on your terminal as well as in the transcript file."=} 1263; end ;
  end;
error;
end;


 ifdef('INITEX')  procedure store_fmt_file;
label found1,found2,done1,done2;
var j, k, l:integer; {all-purpose indices}
 p, q: halfword ; {all-purpose pointers}
 x: integer; {something to dump}
 format_engine: ^ ASCII_code ;
begin 
{ If dumping is not allowed, abort }
if save_ptr<>0 then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"You can't dump inside a group"=} 1274); end ;
{ \xref[You can't dump...] }
   begin help_ptr:=1; help_line[0]:={"`[...\dump]' is a no-no."=} 1275; end ; begin if interaction=error_stop_mode then interaction:=scroll_mode; if log_opened then error; ifdef('TEXMF_DEBUG')  if interaction>batch_mode then debug_help; endif('TEXMF_DEBUG')  history:=fatal_error_stop; jump_out; end ;
  end

;

{ Create the |format_ident|, open the format file, and inform the user that dumping has begun }
selector:=new_string;
print({" (preloaded format="=}1291); print(job_name); print_char({" "=}32);
print_int(eqtb[int_base+ year_code].int  ); print_char({"."=}46);
print_int(eqtb[int_base+ month_code].int  ); print_char({"."=}46); print_int(eqtb[int_base+ day_code].int  ); print_char({")"=}41);
if interaction=batch_mode then selector:=log_only
else selector:=term_and_log;
 begin if pool_ptr+ 1 > pool_size then overflow({"pool size"=}257,pool_size-init_pool_ptr); { \xref[TeX capacity exceeded pool size][\quad pool size] } end ;
format_ident:=make_string;
pack_job_name(format_extension);
while not w_open_out(fmt_file) do
  prompt_file_name({"format file name"=}1292,format_extension);
print_nl({"Beginning to dump on file "=}1293);
{ \xref[Beginning to dump...] }
slow_print(w_make_name_string(fmt_file)); begin decr(str_ptr); pool_ptr:=str_start[str_ptr]; end ;
print_nl({""=}335); slow_print(format_ident)

;

{ Dump constants for consistency check }
dump_int({0x57325458=}1462916184);  {Web2C \TeX's magic constant: "W2TX"}
{Align engine to 4 bytes with one or more trailing NUL}
x:=strlen(engine_name);
format_engine:=xmalloc_array( ASCII_code ,x+4);
strcpy(stringcast(format_engine), engine_name);
for k:=x to x+3 do format_engine[k]:=0;
x:=x+4-(x mod 4);
dump_int(x);dump_things(format_engine[0], x);
libc_free(format_engine);

dump_int(@$);


{ Dump |xord|, |xchr|, and |xprn| }
dump_things(xord[0], 256);
dump_things(xchr[0], 256);
dump_things(xprn[0], 256);

;
dump_int({0xfffffff=}268435455 );

dump_int(hash_high);
dump_int(mem_bot);

dump_int(mem_top);

dump_int(eqtb_size);

dump_int(hash_prime);

dump_int(hyph_prime)

;

{ Dump ML\TeX-specific data }
dump_int({0x4d4c5458=}1296847960);  {ML\TeX's magic constant: "MLTX"}
if mltex_p then dump_int(1)
else dump_int(0);

;

{ Dump the string pool }
dump_int(pool_ptr);
dump_int(str_ptr);
dump_things(str_start[0], str_ptr+1);
dump_things(str_pool[0], pool_ptr);
print_ln; print_int(str_ptr); print({" strings of total length "=}1276);
print_int(pool_ptr)

;

{ Dump the dynamic memory }
sort_avail; var_used:=0;
dump_int(lo_mem_max); dump_int(rover);
p:=mem_bot; q:=rover; x:=0;
repeat dump_things(mem[p], q+2-p);
x:=x+q+2-p; var_used:=var_used+q-p;
p:=q+  mem[ q].hh.lh ; q:=  mem[  q+ 1].hh.rh  ;
until q=rover;
var_used:=var_used+lo_mem_max-p; dyn_used:=mem_end+1-hi_mem_min;

dump_things(mem[p], lo_mem_max+1-p);
x:=x+lo_mem_max+1-p;
dump_int(hi_mem_min); dump_int(avail);
dump_things(mem[hi_mem_min], mem_end+1-hi_mem_min);
x:=x+mem_end+1-hi_mem_min;
p:=avail;
while p<>-{0xfffffff=}268435455   do
  begin decr(dyn_used); p:= mem[ p].hh.rh ;
  end;
dump_int(var_used); dump_int(dyn_used);
print_ln; print_int(x);
print({" memory locations dumped; current usage is "=}1277);
print_int(var_used); print_char({"&"=}38); print_int(dyn_used)

;

{ Dump the table of equivalents }

{ Dump regions 1 to 4 of |eqtb| }
k:=active_base;
repeat j:=k;
while j<int_base-1 do
  begin if ( eqtb[  j].hh.rh  = eqtb[  j+  1].hh.rh  )and( eqtb[  j].hh.b0  = eqtb[  j+  1].hh.b0  )and 
    ( eqtb[  j].hh.b1  = eqtb[  j+  1].hh.b1  ) then goto found1;
  incr(j);
  end;
l:=int_base; goto done1; {|j=int_base-1|}
found1: incr(j); l:=j;
while j<int_base-1 do
  begin if ( eqtb[  j].hh.rh  <> eqtb[  j+  1].hh.rh  )or( eqtb[  j].hh.b0  <> eqtb[  j+  1].hh.b0  )or 
    ( eqtb[  j].hh.b1  <> eqtb[  j+  1].hh.b1  ) then goto done1;
  incr(j);
  end;
done1:dump_int(l-k);
dump_things(eqtb[k], l-k);
k:=j+1; dump_int(k-l);
until k=int_base

;

{ Dump regions 5 and 6 of |eqtb| }
repeat j:=k;
while j<eqtb_size do
  begin if eqtb[j].int=eqtb[j+1].int then goto found2;
  incr(j);
  end;
l:=eqtb_size+1; goto done2; {|j=eqtb_size|}
found2: incr(j); l:=j;
while j<eqtb_size do
  begin if eqtb[j].int<>eqtb[j+1].int then goto done2;
  incr(j);
  end;
done2:dump_int(l-k);
dump_things(eqtb[k], l-k);
k:=j+1; dump_int(k-l);
until k>eqtb_size;
if hash_high>0 then dump_things(eqtb[eqtb_size+1],hash_high);
  {dump |hash_extra| part}

;
dump_int(par_loc); dump_int(write_loc);


{ Dump the hash table }
dump_int(hash_used); cs_count:=frozen_control_sequence-1-hash_used+hash_high;
for p:=hash_base to hash_used do if  hash[ p].rh <>0 then
  begin dump_int(p); dump_hh(hash[p]); incr(cs_count);
  end;
dump_things(hash[hash_used+1], undefined_control_sequence-1-hash_used);
if hash_high>0 then dump_things(hash[eqtb_size+1], hash_high);
dump_int(cs_count);

print_ln; print_int(cs_count); print({" multiletter control sequences"=}1278)



;

{ Dump the font information }
dump_int(fmem_ptr);
dump_things(font_info[0], fmem_ptr);
dump_int(font_ptr);

{ Dump the array info for internal font number |k| }
begin
dump_things(font_check[font_base ], font_ptr+1-font_base );
dump_things(font_size[font_base ], font_ptr+1-font_base );
dump_things(font_dsize[font_base ], font_ptr+1-font_base );
dump_things(font_params[font_base ], font_ptr+1-font_base );
dump_things(hyphen_char[font_base ], font_ptr+1-font_base );
dump_things(skew_char[font_base ], font_ptr+1-font_base );
dump_things(font_name[font_base ], font_ptr+1-font_base );
dump_things(font_area[font_base ], font_ptr+1-font_base );
dump_things(font_bc[font_base ], font_ptr+1-font_base );
dump_things(font_ec[font_base ], font_ptr+1-font_base );
dump_things(char_base[font_base ], font_ptr+1-font_base );
dump_things(width_base[font_base ], font_ptr+1-font_base );
dump_things(height_base[font_base ], font_ptr+1-font_base );
dump_things(depth_base[font_base ], font_ptr+1-font_base );
dump_things(italic_base[font_base ], font_ptr+1-font_base );
dump_things(lig_kern_base[font_base ], font_ptr+1-font_base );
dump_things(kern_base[font_base ], font_ptr+1-font_base );
dump_things(exten_base[font_base ], font_ptr+1-font_base );
dump_things(param_base[font_base ], font_ptr+1-font_base );
dump_things(font_glue[font_base ], font_ptr+1-font_base );
dump_things(bchar_label[font_base ], font_ptr+1-font_base );
dump_things(font_bchar[font_base ], font_ptr+1-font_base );
dump_things(font_false_bchar[font_base ], font_ptr+1-font_base );
for k:=font_base  to font_ptr do
  begin print_nl({"\font"=}1282); print_esc(  hash[ font_id_base+  k].rh  ); print_char({"="=}61);
  print_file_name(font_name[k],font_area[k],{""=}335);
  if font_size[k]<>font_dsize[k] then
    begin print({" at "=}751); print_scaled(font_size[k]); print({"pt"=}402);
    end;
  end;
end

;
print_ln; print_int(fmem_ptr-7); print({" words of font info for "=}1279);
print_int(font_ptr-font_base);
if font_ptr<>font_base+1 then print({" preloaded fonts"=}1280)
else print({" preloaded font"=}1281)

;

{ Dump the hyphenation tables }
dump_int(hyph_count);
if hyph_next <= hyph_prime then hyph_next:=hyph_size;
dump_int(hyph_next);{minimum value of |hyphen_size| needed}
for k:=0 to hyph_size do if hyph_word[k]<>0 then
  begin dump_int(k+65536*hyph_link[k]);
        {assumes number of hyphen exceptions does not exceed 65535}
   dump_int(hyph_word[k]); dump_int(hyph_list[k]);
  end;
print_ln; print_int(hyph_count);
if hyph_count<>1 then print({" hyphenation exceptions"=}1283)
else print({" hyphenation exception"=}1284);
if trie_not_ready then init_trie;
dump_int(trie_max);
dump_things(trie_trl[0], trie_max+1);
dump_things(trie_tro[0], trie_max+1);
dump_things(trie_trc[0], trie_max+1);
dump_int(trie_op_ptr);
dump_things(hyf_distance[1], trie_op_ptr);
dump_things(hyf_num[1], trie_op_ptr);
dump_things(hyf_next[1], trie_op_ptr);
print_nl({"Hyphenation trie of length "=}1285); print_int(trie_max);
{ \xref[Hyphenation trie...] }
print({" has "=}1286); print_int(trie_op_ptr);
if trie_op_ptr<>1 then print({" ops"=}1287)
else print({" op"=}1288);
print({" out of "=}1289); print_int(trie_op_size);
for k:=255 downto 0 do if trie_used[k]>min_quarterword then
  begin print_nl({"  "=}810); print_int( trie_used[ k] );
  print({" for language "=}1290); print_int(k);
  dump_int(k); dump_int( trie_used[ k] );
  end

;

{ Dump a couple more things and the closing check word }
dump_int(interaction); dump_int(format_ident); dump_int(69069);
eqtb[int_base+ tracing_stats_code].int  :=0

;

{ Close the format file }
w_close(fmt_file)

;
end;
endif('INITEX') 


{ \4 }
{ Declare procedures needed in |do_extension| }
procedure new_whatsit( s:small_number; w:small_number);
var p:halfword ; {the new node}
begin p:=get_node(w);  mem[ p].hh.b0 :=whatsit_node;  mem[ p].hh.b1 :=s;
 mem[ cur_list.tail_field ].hh.rh :=p; cur_list.tail_field :=p;
end;


procedure new_write_whatsit( w:small_number);
begin new_whatsit(cur_chr,w);
if w<>write_node_size then scan_four_bit_int
else  begin scan_int;
  if cur_val<0 then cur_val:=17
  else if (cur_val>15) and (cur_val <> 18) then cur_val:=16;
  end;
  mem[  cur_list.tail_field + 1].hh.lh  :=cur_val;
end;

 
procedure do_extension;
var k:integer; {all-purpose integers}
 p:halfword ; {all-purpose pointers}
begin case cur_chr of
open_node:
{ Implement \.[\\openout] }
begin new_write_whatsit(open_node_size);
scan_optional_equals; scan_file_name;

  mem[  cur_list.tail_field + 1].hh.rh  :=cur_name;   mem[  cur_list.tail_field + 2].hh.lh  :=cur_area;   mem[  cur_list.tail_field + 2].hh.rh  :=cur_ext;
end

;
write_node:
{ Implement \.[\\write] }
begin k:=cur_cs; new_write_whatsit(write_node_size);

cur_cs:=k; p:=scan_toks(false,false);   mem[  cur_list.tail_field + 1].hh.rh  :=def_ref;
end

;
close_node:
{ Implement \.[\\closeout] }
begin new_write_whatsit(write_node_size);   mem[  cur_list.tail_field + 1].hh.rh  :=-{0xfffffff=}268435455  ;
end

;
special_node:
{ Implement \.[\\special] }
begin new_whatsit(special_node,write_node_size);   mem[  cur_list.tail_field + 1].hh.lh  :=-{0xfffffff=}268435455  ;
p:=scan_toks(false,true);   mem[  cur_list.tail_field + 1].hh.rh  :=def_ref;
end

;
immediate_code:
{ Implement \.[\\immediate] }
begin get_x_token;
if (cur_cmd=extension)and(cur_chr<=close_node) then
  begin p:=cur_list.tail_field ; do_extension; {append a whatsit node}
  out_what(cur_list.tail_field ); {do the action immediately}
  flush_node_list(cur_list.tail_field ); cur_list.tail_field :=p;  mem[ p].hh.rh :=-{0xfffffff=}268435455  ;
  end
else back_input;
end

;
set_language_code:
{ Implement \.[\\setlanguage] }
if abs(cur_list.mode_field )<>hmode then report_illegal_case
else begin new_whatsit(language_node,small_node_size);
  scan_int;
  if cur_val<=0 then cur_list.aux_field .hh.rh :=0
  else if cur_val>255 then cur_list.aux_field .hh.rh :=0
  else cur_list.aux_field .hh.rh :=cur_val;
   mem[  cur_list.tail_field + 1].hh.rh  :=cur_list.aux_field .hh.rh ;
   mem[  cur_list.tail_field + 1].hh.b0  :=norm_min(eqtb[int_base+ left_hyphen_min_code].int  );
   mem[  cur_list.tail_field + 1].hh.b1  :=norm_min(eqtb[int_base+ right_hyphen_min_code].int  );
  end

;
 else  confusion({"ext1"=}1310)
{ \xref[this can't happen ext1][\quad ext1] }
 end ;
end;


procedure fix_language;
var  l:ASCII_code; {the new current language}
begin if eqtb[int_base+ language_code].int  <=0 then l:=0
else if eqtb[int_base+ language_code].int  >255 then l:=0
else l:=eqtb[int_base+ language_code].int  ;
if l<>cur_list.aux_field .hh.rh  then
  begin new_whatsit(language_node,small_node_size);
   mem[  cur_list.tail_field + 1].hh.rh  :=l; cur_list.aux_field .hh.rh :=l;

   mem[  cur_list.tail_field + 1].hh.b0  :=norm_min(eqtb[int_base+ left_hyphen_min_code].int  );
   mem[  cur_list.tail_field + 1].hh.b1  :=norm_min(eqtb[int_base+ right_hyphen_min_code].int  );
  end;
end;



procedure insert_src_special;
var toklist, p, q : halfword ;
begin
  if (source_filename_stack[in_open] > 0 and is_new_source (source_filename_stack[in_open]
, line)) then begin
    toklist := get_avail;
    p := toklist;
     mem[ p].hh.lh  := {07777=}4095 +frozen_special;
     mem[ p].hh.rh  := get_avail; p :=  mem[ p].hh.rh ;
     mem[ p].hh.lh  := left_brace_token+{"["=}123;
    q := str_toks (make_src_special (source_filename_stack[in_open], line));
     mem[ p].hh.rh  :=  mem[ mem_top-3 ].hh.rh ;
    p := q;
     mem[ p].hh.rh  := get_avail; p :=  mem[ p].hh.rh ;
     mem[ p].hh.lh  := right_brace_token+{"]"=}125;
    begin_token_list( toklist,inserted) ;
    remember_source_info (source_filename_stack[in_open], line);
  end;
end;

procedure append_src_special;
var q : halfword ;
begin
  if (source_filename_stack[in_open] > 0 and is_new_source (source_filename_stack[in_open]
, line)) then begin
    new_whatsit (special_node, write_node_size);
      mem[  cur_list.tail_field + 1].hh.lh   := 0;
    def_ref := get_avail;
      mem[  def_ref].hh.lh   := -{0xfffffff=}268435455  ;
    q := str_toks (make_src_special (source_filename_stack[in_open], line));
     mem[ def_ref].hh.rh  :=  mem[ mem_top-3 ].hh.rh ;
      mem[  cur_list.tail_field + 1].hh.rh   := def_ref;
    remember_source_info (source_filename_stack[in_open], line);
  end;
end;

 
{ \4 }
{ Declare the procedure called |handle_right_brace| }
procedure handle_right_brace;
var p, q:halfword ; {for short-term use}
 d:scaled; {holds |split_max_depth| in |insert_group|}
 f:integer; {holds |floating_penalty| in |insert_group|}
begin case cur_group of
simple_group: unsave;
bottom_level: begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Too many ]'s"=} 1057); end ;
{ \xref[Too many \]'s] }
   begin help_ptr:=2; help_line[1]:={"You've closed more groups than you opened."=} 1058; help_line[0]:={"Such booboos are generally harmless, so keep going."=} 1059; end ; error;
  end;
semi_simple_group,math_shift_group,math_left_group: extra_right_brace;
{ \4 }
{ Cases of |handle_right_brace| where a |right_brace| triggers a delayed action }
hbox_group: package(0);
adjusted_hbox_group: begin adjust_tail:=mem_top-5 ; package(0);
  end;
vbox_group: begin end_graf; package(0);
  end;
vtop_group: begin end_graf; package(vtop_code);
  end;


insert_group: begin end_graf; q:= eqtb[  glue_base+   split_top_skip_code].hh.rh    ; incr(  mem[   q].hh.rh  ) ;
  d:=eqtb[dimen_base+ split_max_depth_code].int   ; f:=eqtb[int_base+ floating_penalty_code].int  ; unsave; decr(save_ptr);
  {now |saved(0)| is the insertion number, or 255 for |vadjust|}
  p:=vpackage(  mem[  cur_list.head_field ].hh.rh , 0,additional ,{07777777777=}1073741823 ) ; pop_nest;
  if save_stack[save_ptr+ 0].int <255 then
    begin begin  mem[ cur_list.tail_field ].hh.rh := get_node( ins_node_size); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
     mem[ cur_list.tail_field ].hh.b0 :=ins_node;  mem[ cur_list.tail_field ].hh.b1 := save_stack[save_ptr+  0].int  ;
     mem[ cur_list.tail_field +height_offset].int  := mem[ p+height_offset].int  + mem[ p+depth_offset].int  ;  mem[  cur_list.tail_field + 4].hh.lh  :=  mem[  p+ list_offset].hh.rh  ;
     mem[  cur_list.tail_field + 4].hh.rh  :=q;  mem[ cur_list.tail_field +depth_offset].int  :=d; mem[ cur_list.tail_field +1].int :=f;
    end
  else  begin begin  mem[ cur_list.tail_field ].hh.rh := get_node( small_node_size); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
     mem[ cur_list.tail_field ].hh.b0 :=adjust_node;

     mem[ cur_list.tail_field ].hh.b1 :=0; {the |subtype| is not used}
    mem[ cur_list.tail_field +1].int :=  mem[  p+ list_offset].hh.rh  ; delete_glue_ref(q);
    end;
  free_node(p,box_node_size);
  if nest_ptr=0 then build_page;
  end;
output_group: 
{ Resume the page builder... }
begin if (cur_input.loc_field <>-{0xfffffff=}268435455  ) or
 ((cur_input.index_field  <>output_text)and(cur_input.index_field  <>backed_up)) then
  
{ Recover from an unbalanced output routine }
begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Unbalanced output routine"=} 1024); end ;
{ \xref[Unbalanced output routine] }
 begin help_ptr:=2; help_line[1]:={"Your sneaky output routine has problematic ['s and/or ]'s."=} 1025; help_line[0]:={"I can't handle that very well; good luck."=} 1026; end ; error;
repeat get_token;
until cur_input.loc_field =-{0xfffffff=}268435455  ;
end {loops forever if reading from a file, since |null=min_halfword<=0|}

;
end_token_list; {conserve stack space in case more outputs are triggered}
end_graf; unsave; output_active:=false; insert_penalties:=0;


{ Ensure that box 255 is empty after output }
if  eqtb[  box_base+   255].hh.rh   <>-{0xfffffff=}268435455   then
  begin begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Output routine didn't use all of "=} 1027); end ;
  print_esc({"box"=}414); print_int(255);
{ \xref[Output routine didn't use...] }
   begin help_ptr:=3; help_line[2]:={"Your \output commands should empty \box255,"=} 1028; help_line[1]:={"e.g., by saying `\shipout\box255'."=} 1029; help_line[0]:={"Proceed; I'll discard its present contents."=} 1030; end ;
  box_error(255);
  end

;
if cur_list.tail_field <>cur_list.head_field  then {current list goes after heldover insertions}
  begin  mem[ page_tail].hh.rh := mem[ cur_list.head_field ].hh.rh ;
  page_tail:=cur_list.tail_field ;
  end;
if  mem[ mem_top-2 ].hh.rh <>-{0xfffffff=}268435455   then {and both go before heldover contributions}
  begin if  mem[ mem_top-1 ].hh.rh =-{0xfffffff=}268435455   then nest[0].tail_field :=page_tail;
   mem[ page_tail].hh.rh := mem[ mem_top-1 ].hh.rh ;
   mem[ mem_top-1 ].hh.rh := mem[ mem_top-2 ].hh.rh ;
   mem[ mem_top-2 ].hh.rh :=-{0xfffffff=}268435455  ; page_tail:=mem_top-2 ;
  end;
pop_nest; build_page;
end

;


disc_group: build_discretionary;


align_group: begin back_input; cur_tok:={07777=}4095 +frozen_cr;
  begin if interaction=error_stop_mode then    ; if file_line_error_style_p then print_file_line else print_nl({"! "=}262); print({"Missing "=} 635); end ; print_esc({"cr"=}913); print({" inserted"=}636);
{ \xref[Missing \\cr inserted] }
   begin help_ptr:=1; help_line[0]:={"I'm guessing that you meant to end an alignment here."=} 1138; end ;
  ins_error;
  end;


no_align_group: begin end_graf; unsave; align_peek;
  end;


vcenter_group: begin end_graf; unsave; save_ptr:=save_ptr-2;
  p:=vpackage(  mem[  cur_list.head_field ].hh.rh , save_stack[save_ptr+  1].int , save_stack[save_ptr+  0].int ,{07777777777=}1073741823 ) ; pop_nest;
  begin  mem[ cur_list.tail_field ].hh.rh := new_noad; cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;  mem[ cur_list.tail_field ].hh.b0 :=vcenter_noad;
   mem[   cur_list.tail_field +1 ].hh.rh :=sub_box;  mem[   cur_list.tail_field +1 ].hh.lh :=p;
  end;


math_choice_group: build_choices;


math_group: begin unsave; decr(save_ptr);

   mem[ save_stack[save_ptr+  0].int ].hh.rh :=sub_mlist; p:=fin_mlist(-{0xfffffff=}268435455  );  mem[ save_stack[save_ptr+  0].int ].hh.lh :=p;
  if p<>-{0xfffffff=}268435455   then if  mem[ p].hh.rh =-{0xfffffff=}268435455   then
   if  mem[ p].hh.b0 =ord_noad then
    begin if  mem[   p+3 ].hh.rh =empty then
     if  mem[   p+2 ].hh.rh =empty then
      begin mem[save_stack[save_ptr+ 0].int ].hh:=mem[ p+1 ].hh;
      free_node(p,noad_size);
      end;
    end
  else if  mem[ p].hh.b0 =accent_noad then if save_stack[save_ptr+ 0].int = cur_list.tail_field +1  then
   if  mem[ cur_list.tail_field ].hh.b0 =ord_noad then 
{ Replace the tail of the list by |p| }
begin q:=cur_list.head_field ; while  mem[ q].hh.rh <>cur_list.tail_field  do q:= mem[ q].hh.rh ;
 mem[ q].hh.rh :=p; free_node(cur_list.tail_field ,noad_size); cur_list.tail_field :=p;
end

;
  end;

 
 else  confusion({"rightbrace"=}1060)
{ \xref[this can't happen rightbrace][\quad rightbrace] }
 end ;
end;

 
procedure main_control; {governs \TeX's activities}
label big_switch,reswitch,main_loop,main_loop_wrapup,
  main_loop_move,main_loop_move+1,main_loop_move+2,main_loop_move_lig,
  main_loop_lookahead,main_loop_lookahead+1,
  main_lig_loop,main_lig_loop+1,main_lig_loop+2,
  append_normal_space,exit;
var t:integer; {general-purpose temporary variable}
begin if  eqtb[  every_job_loc].hh.rh   <>-{0xfffffff=}268435455   then begin_token_list( eqtb[  every_job_loc].hh.rh   ,every_job_text);
big_switch: get_x_token;

reswitch: 
{ Give diagnostic information, if requested }
if interrupt<>0 then if OK_to_interrupt then
  begin back_input; begin if interrupt<>0 then pause_for_instructions; end ; goto big_switch;
  end;
 ifdef('TEXMF_DEBUG')  if panicking then check_mem(false);   endif('TEXMF_DEBUG') 
if eqtb[int_base+ tracing_commands_code].int  >0 then show_cur_cmd_chr

;
case abs(cur_list.mode_field )+cur_cmd of
hmode+letter,hmode+other_char,hmode+char_given: goto main_loop;
hmode+char_num: begin scan_char_num; cur_chr:=cur_val; goto main_loop; end;
hmode+no_boundary: begin get_x_token;
  if (cur_cmd=letter)or(cur_cmd=other_char)or(cur_cmd=char_given)or
   (cur_cmd=char_num) then cancel_boundary:=true;
  goto reswitch;
  end;
hmode+spacer: if cur_list.aux_field .hh.lh =1000 then goto append_normal_space
  else app_space;
hmode+ex_space,mmode+ex_space: goto append_normal_space;
{ \4 }
{ Cases of |main_control| that are not part of the inner loop }
vmode+ relax,hmode+ relax,mmode+ relax ,vmode+spacer,mmode+spacer,mmode+no_boundary: ;
vmode+ ignore_spaces,hmode+ ignore_spaces,mmode+ ignore_spaces : begin 
{ Get the next non-blank non-call... }
repeat get_x_token;
until cur_cmd<>spacer

;
  goto reswitch;
  end;
vmode+stop: if its_all_over then  goto exit ; {this is the only way out}
{ \4 }
{ Forbidden cases detected in |main_control| }
vmode+vmove,hmode+hmove,mmode+hmove,vmode+ last_item,hmode+ last_item,mmode+ last_item ,


vmode+vadjust,

vmode+ital_corr,

vmode+ eq_no,hmode+ eq_no ,

  vmode+ mac_param,hmode+ mac_param,mmode+ mac_param :
  report_illegal_case;

{ Math-only cases in non-math modes, or vice versa }
vmode+ sup_mark,hmode+ sup_mark , vmode+ sub_mark,hmode+ sub_mark , vmode+ math_char_num,hmode+ math_char_num ,
vmode+ math_given,hmode+ math_given , vmode+ math_comp,hmode+ math_comp , vmode+ delim_num,hmode+ delim_num ,
vmode+ left_right,hmode+ left_right , vmode+ above,hmode+ above , vmode+ radical,hmode+ radical ,
vmode+ math_style,hmode+ math_style , vmode+ math_choice,hmode+ math_choice , vmode+ vcenter,hmode+ vcenter ,
vmode+ non_script,hmode+ non_script , vmode+ mkern,hmode+ mkern , vmode+ limit_switch,hmode+ limit_switch ,
vmode+ mskip,hmode+ mskip , vmode+ math_accent,hmode+ math_accent ,
mmode+endv, mmode+par_end, mmode+stop, mmode+vskip, mmode+un_vbox,
mmode+valign, mmode+hrule

: insert_dollar_sign;
{ \4 }
{ Cases of |main_control| that build boxes and lists }
vmode+hrule,hmode+vrule,mmode+vrule: begin begin  mem[ cur_list.tail_field ].hh.rh := scan_rule_spec; cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
  if abs(cur_list.mode_field )=vmode then cur_list.aux_field .int  :=-65536000 
  else if abs(cur_list.mode_field )=hmode then cur_list.aux_field .hh.lh :=1000;
  end;


vmode+vskip,hmode+hskip,mmode+hskip,mmode+mskip: append_glue;
vmode+ kern,hmode+ kern,mmode+ kern ,mmode+mkern: append_kern;


vmode+ left_brace,hmode+ left_brace : new_save_level(simple_group);
vmode+ begin_group,hmode+ begin_group,mmode+ begin_group : new_save_level(semi_simple_group);
vmode+ end_group,hmode+ end_group,mmode+ end_group : if cur_group=semi_simple_group then unsave
  else off_save;


vmode+ right_brace,hmode+ right_brace,mmode+ right_brace : handle_right_brace;


vmode+hmove,hmode+vmove,mmode+vmove: begin t:=cur_chr;
  scan_dimen(false,false,false) ;
  if t=0 then scan_box(cur_val) else scan_box(-cur_val);
  end;
vmode+ leader_ship,hmode+ leader_ship,mmode+ leader_ship : scan_box({010000000000=}1073741824 +513 -a_leaders+cur_chr);
vmode+ make_box,hmode+ make_box,mmode+ make_box : begin_box(0);


vmode+start_par: new_graf(cur_chr>0);
vmode+letter,vmode+other_char,vmode+char_num,vmode+char_given,
   vmode+math_shift,vmode+un_hbox,vmode+vrule,
   vmode+accent,vmode+discretionary,vmode+hskip,vmode+valign,
   vmode+ex_space,vmode+no_boundary:{  } 

  begin back_input; new_graf(true);
  end;


hmode+start_par,mmode+start_par: indent_in_hmode;


vmode+par_end: begin normal_paragraph;
  if cur_list.mode_field >0 then build_page;
  end;
hmode+par_end: begin if align_state<0 then off_save; {this tries to
    recover from an alignment that didn't end properly}
  end_graf; {this takes us to the enclosing mode, if |mode>0|}
  if cur_list.mode_field =vmode then build_page;
  end;
hmode+stop,hmode+vskip,hmode+hrule,hmode+un_vbox,hmode+halign: head_for_vmode;


vmode+ insert,hmode+ insert,mmode+ insert ,hmode+vadjust,mmode+vadjust: begin_insert_or_adjust;
vmode+ mark,hmode+ mark,mmode+ mark : make_mark;


vmode+ break_penalty,hmode+ break_penalty,mmode+ break_penalty : append_penalty;


vmode+ remove_item,hmode+ remove_item,mmode+ remove_item : delete_last;


vmode+un_vbox,hmode+un_hbox,mmode+un_hbox: unpackage;


hmode+ital_corr: append_italic_correction;
mmode+ital_corr: begin  mem[ cur_list.tail_field ].hh.rh := new_kern( 0); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;


hmode+discretionary,mmode+discretionary: append_discretionary;


hmode+accent: make_accent;


vmode+ car_ret,hmode+ car_ret,mmode+ car_ret , vmode+ tab_mark,hmode+ tab_mark,mmode+ tab_mark : align_error;
vmode+ no_align,hmode+ no_align,mmode+ no_align : no_align_error;
vmode+ omit,hmode+ omit,mmode+ omit : omit_error;


vmode+halign,hmode+valign:init_align;
mmode+halign: if privileged then
  if cur_group=math_shift_group then init_align
  else off_save;
vmode+endv,hmode+endv: do_endv;


vmode+ end_cs_name,hmode+ end_cs_name,mmode+ end_cs_name : cs_error;


hmode+math_shift:init_math;


mmode+eq_no: if privileged then
  if cur_group=math_shift_group then start_eq_no
  else off_save;


mmode+left_brace: begin begin  mem[ cur_list.tail_field ].hh.rh := new_noad; cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
  back_input; scan_math( cur_list.tail_field +1 );
  end;


mmode+letter,mmode+other_char,mmode+char_given:
  set_math_char(  eqtb[  math_code_base+    cur_chr].hh.rh    );
mmode+char_num: begin scan_char_num; cur_chr:=cur_val;
  set_math_char(  eqtb[  math_code_base+    cur_chr].hh.rh    );
  end;
mmode+math_char_num: begin scan_fifteen_bit_int; set_math_char(cur_val);
  end;
mmode+math_given: set_math_char(cur_chr);
mmode+delim_num: begin scan_twenty_seven_bit_int;
  set_math_char(cur_val div {010000=}4096);
  end;


mmode+math_comp: begin begin  mem[ cur_list.tail_field ].hh.rh := new_noad; cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
   mem[ cur_list.tail_field ].hh.b0 :=cur_chr; scan_math( cur_list.tail_field +1 );
  end;
mmode+limit_switch: math_limit_switch;


mmode+radical:math_radical;


mmode+accent,mmode+math_accent:math_ac;


mmode+vcenter: begin scan_spec(vcenter_group,false); normal_paragraph;
  push_nest; cur_list.mode_field :=-vmode; cur_list.aux_field .int  :=-65536000 ;
  if (insert_src_special_every_vbox) then insert_src_special;
  if  eqtb[  every_vbox_loc].hh.rh   <>-{0xfffffff=}268435455   then begin_token_list( eqtb[  every_vbox_loc].hh.rh   ,every_vbox_text);
  end;


mmode+math_style: begin  mem[ cur_list.tail_field ].hh.rh := new_style( cur_chr); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
mmode+non_script: begin begin  mem[ cur_list.tail_field ].hh.rh := new_glue( mem_bot ); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ;
   mem[ cur_list.tail_field ].hh.b1 :=cond_math_glue;
  end;
mmode+math_choice: append_choices;


mmode+sub_mark,mmode+sup_mark: sub_sup;


mmode+above: math_fraction;


mmode+left_right: math_left_right;


mmode+math_shift: if cur_group=math_shift_group then after_math
  else off_save;

 
{ \4 }
{ Cases of |main_control| that don't depend on |mode| }
vmode+ toks_register,hmode+ toks_register,mmode+ toks_register ,
vmode+ assign_toks,hmode+ assign_toks,mmode+ assign_toks ,
vmode+ assign_int,hmode+ assign_int,mmode+ assign_int ,
vmode+ assign_dimen,hmode+ assign_dimen,mmode+ assign_dimen ,
vmode+ assign_glue,hmode+ assign_glue,mmode+ assign_glue ,
vmode+ assign_mu_glue,hmode+ assign_mu_glue,mmode+ assign_mu_glue ,
vmode+ assign_font_dimen,hmode+ assign_font_dimen,mmode+ assign_font_dimen ,
vmode+ assign_font_int,hmode+ assign_font_int,mmode+ assign_font_int ,
vmode+ set_aux,hmode+ set_aux,mmode+ set_aux ,
vmode+ set_prev_graf,hmode+ set_prev_graf,mmode+ set_prev_graf ,
vmode+ set_page_dimen,hmode+ set_page_dimen,mmode+ set_page_dimen ,
vmode+ set_page_int,hmode+ set_page_int,mmode+ set_page_int ,
vmode+ set_box_dimen,hmode+ set_box_dimen,mmode+ set_box_dimen ,
vmode+ set_shape,hmode+ set_shape,mmode+ set_shape ,
vmode+ def_code,hmode+ def_code,mmode+ def_code ,
vmode+ def_family,hmode+ def_family,mmode+ def_family ,
vmode+ set_font,hmode+ set_font,mmode+ set_font ,
vmode+ def_font,hmode+ def_font,mmode+ def_font ,
vmode+ register,hmode+ register,mmode+ register ,
vmode+ advance,hmode+ advance,mmode+ advance ,
vmode+ multiply,hmode+ multiply,mmode+ multiply ,
vmode+ divide,hmode+ divide,mmode+ divide ,
vmode+ prefix,hmode+ prefix,mmode+ prefix ,
vmode+ let,hmode+ let,mmode+ let ,
vmode+ shorthand_def,hmode+ shorthand_def,mmode+ shorthand_def ,
vmode+ read_to_cs,hmode+ read_to_cs,mmode+ read_to_cs ,
vmode+ def,hmode+ def,mmode+ def ,
vmode+ set_box,hmode+ set_box,mmode+ set_box ,
vmode+ hyph_data,hmode+ hyph_data,mmode+ hyph_data ,
vmode+ set_interaction,hmode+ set_interaction,mmode+ set_interaction :prefixed_command;


vmode+ after_assignment,hmode+ after_assignment,mmode+ after_assignment :begin get_token; after_token:=cur_tok;
  end;


vmode+ after_group,hmode+ after_group,mmode+ after_group :begin get_token; save_for_after(cur_tok);
  end;


vmode+ in_stream,hmode+ in_stream,mmode+ in_stream : open_or_close_in;


vmode+ message,hmode+ message,mmode+ message :issue_message;


vmode+ case_shift,hmode+ case_shift,mmode+ case_shift :shift_case;


vmode+ xray,hmode+ xray,mmode+ xray : show_whatever;

 
{ \4 }
{ Cases of |main_control| that are for extensions to \TeX }
vmode+ extension,hmode+ extension,mmode+ extension :do_extension;

 

 
end; {of the big |case| statement}
goto big_switch;
main_loop:
{ Append character |cur_chr| and the following characters (if~any) to the current hlist in the current font; |goto reswitch| when a non-character has been fetched }
if ((cur_list.head_field =cur_list.tail_field ) and (cur_list.mode_field >0)) then begin
  if (insert_src_special_auto) then append_src_special;
end;
{  } main_s:= eqtb[  sf_code_base+   cur_chr].hh.rh   ; if main_s=1000 then cur_list.aux_field .hh.lh :=1000 else if main_s<1000 then begin if main_s>0 then cur_list.aux_field .hh.lh :=main_s; end else if cur_list.aux_field .hh.lh <1000 then cur_list.aux_field .hh.lh :=1000 else cur_list.aux_field .hh.lh :=main_s ;

main_f:= eqtb[  cur_font_loc].hh.rh   ;
bchar:=font_bchar[main_f]; false_bchar:=font_false_bchar[main_f];
if cur_list.mode_field >0 then if eqtb[int_base+ language_code].int  <>cur_list.aux_field .hh.rh  then fix_language;
{  } begin  lig_stack:=avail; if  lig_stack=-{0xfffffff=}268435455   then  lig_stack:=get_avail else begin avail:= mem[  lig_stack].hh.rh ;  mem[  lig_stack].hh.rh :=-{0xfffffff=}268435455  ; ifdef('STAT')  incr(dyn_used); endif('STAT')  end; end ;   mem[ lig_stack].hh.b0 :=main_f; cur_l:= cur_chr ;
  mem[ lig_stack].hh.b1 :=cur_l;

cur_q:=cur_list.tail_field ;
if cancel_boundary then
  begin cancel_boundary:=false; main_k:=non_address;
  end
else main_k:=bchar_label[main_f];
if main_k=non_address then goto main_loop_move+2; {no left boundary processing}
cur_r:=cur_l; cur_l:= 256  ;
goto main_lig_loop+1; {begin with cursor after left boundary}


main_loop_wrapup:
{ Make a ligature node, if |ligature_present|; insert a null discretionary, if appropriate }
if cur_l< 256   then begin if  mem[ cur_q].hh.rh >-{0xfffffff=}268435455   then if   mem[ cur_list.tail_field ].hh.b1 = hyphen_char[ main_f]  then ins_disc:=true; if ligature_present then  begin main_p:=new_ligature(main_f,cur_l, mem[ cur_q].hh.rh ); if lft_hit then begin  mem[ main_p].hh.b1 :=2; lft_hit:=false; end; if   rt_hit then if lig_stack=-{0xfffffff=}268435455   then begin incr( mem[ main_p].hh.b1 ); rt_hit:=false; end;  mem[ cur_q].hh.rh :=main_p; cur_list.tail_field :=main_p; ligature_present:=false; end ; if ins_disc then begin ins_disc:=false; if cur_list.mode_field >0 then begin  mem[ cur_list.tail_field ].hh.rh := new_disc; cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ; end; end 

;
main_loop_move:
{ If the cursor is immediately followed by the right boundary, |goto reswitch|; if it's followed by an invalid character, |goto big_switch|; otherwise move the cursor one step to the right and |goto main_lig_loop| }
{ \xref[inner loop] }
if lig_stack=-{0xfffffff=}268435455   then goto reswitch;
cur_q:=cur_list.tail_field ; cur_l:=  mem[ lig_stack].hh.b1 ;
main_loop_move+1:if not  ( lig_stack>=hi_mem_min)  then goto main_loop_move_lig;
main_loop_move+2:
if( effective_char( false, main_f,   cur_chr ) >font_ec[main_f])or
  ( effective_char( false, main_f,   cur_chr ) <font_bc[main_f]) then
  begin char_warning(main_f,cur_chr);  begin  mem[  lig_stack].hh.rh :=avail; avail:= lig_stack; ifdef('STAT')  decr(dyn_used); endif('STAT')  end ; goto big_switch;
  end;
main_i:=effective_char_info(main_f,cur_l);
if not ( main_i.b0>min_quarterword)  then
  begin char_warning(main_f,cur_chr);  begin  mem[  lig_stack].hh.rh :=avail; avail:= lig_stack; ifdef('STAT')  decr(dyn_used); endif('STAT')  end ; goto big_switch;
  end;
 mem[ cur_list.tail_field ].hh.rh :=lig_stack; cur_list.tail_field :=lig_stack {|main_loop_lookahead| is next}

;
main_loop_lookahead:
{ Look ahead for another character, or leave |lig_stack| empty if there's none there }
get_next; {set only |cur_cmd| and |cur_chr|, for speed}
if cur_cmd=letter then goto main_loop_lookahead+1;
if cur_cmd=other_char then goto main_loop_lookahead+1;
if cur_cmd=char_given then goto main_loop_lookahead+1;
x_token; {now expand and set |cur_cmd|, |cur_chr|, |cur_tok|}
if cur_cmd=letter then goto main_loop_lookahead+1;
if cur_cmd=other_char then goto main_loop_lookahead+1;
if cur_cmd=char_given then goto main_loop_lookahead+1;
if cur_cmd=char_num then
  begin scan_char_num; cur_chr:=cur_val; goto main_loop_lookahead+1;
  end;
if cur_cmd=no_boundary then bchar:= 256  ;
cur_r:=bchar; lig_stack:=-{0xfffffff=}268435455  ; goto main_lig_loop;
main_loop_lookahead+1: {  } main_s:= eqtb[  sf_code_base+   cur_chr].hh.rh   ; if main_s=1000 then cur_list.aux_field .hh.lh :=1000 else if main_s<1000 then begin if main_s>0 then cur_list.aux_field .hh.lh :=main_s; end else if cur_list.aux_field .hh.lh <1000 then cur_list.aux_field .hh.lh :=1000 else cur_list.aux_field .hh.lh :=main_s ;
{  } begin  lig_stack:=avail; if  lig_stack=-{0xfffffff=}268435455   then  lig_stack:=get_avail else begin avail:= mem[  lig_stack].hh.rh ;  mem[  lig_stack].hh.rh :=-{0xfffffff=}268435455  ; ifdef('STAT')  incr(dyn_used); endif('STAT')  end; end ;   mem[ lig_stack].hh.b0 :=main_f;
cur_r:= cur_chr ;   mem[ lig_stack].hh.b1 :=cur_r;
if cur_r=false_bchar then cur_r:= 256   {this prevents spurious ligatures}

;
main_lig_loop:
{ If there's a ligature/kern command relevant to |cur_l| and |cur_r|, adjust the text appropriately; exit to |main_loop_wrapup| }
if ((  main_i. b2 ) mod 4) <>lig_tag then goto main_loop_wrapup;
if cur_r= 256   then goto main_loop_wrapup;
main_k:=lig_kern_base[ main_f]+ main_i.b3 ; main_j:=font_info[main_k].qqqq;
if  main_j.b0 <= 128   then goto main_lig_loop+2;
main_k:=lig_kern_base[ main_f]+256*  main_j.b2 +  main_j.b3 +32768-256*(128+min_quarterword)  ;
main_lig_loop+1:main_j:=font_info[main_k].qqqq;
main_lig_loop+2:if  main_j.b1 =cur_r then
 if  main_j.b0 <= 128   then
  
{ Do ligature or kern command, returning to |main_lig_loop| or |main_loop_wrapup| or |main_loop_move| }
begin if  main_j.b2 >= 128   then
  begin if cur_l< 256   then begin if  mem[ cur_q].hh.rh >-{0xfffffff=}268435455   then if   mem[ cur_list.tail_field ].hh.b1 = hyphen_char[ main_f]  then ins_disc:=true; if ligature_present then  begin main_p:=new_ligature(main_f,cur_l, mem[ cur_q].hh.rh ); if lft_hit then begin  mem[ main_p].hh.b1 :=2; lft_hit:=false; end; if   rt_hit then if lig_stack=-{0xfffffff=}268435455   then begin incr( mem[ main_p].hh.b1 ); rt_hit:=false; end;  mem[ cur_q].hh.rh :=main_p; cur_list.tail_field :=main_p; ligature_present:=false; end ; if ins_disc then begin ins_disc:=false; if cur_list.mode_field >0 then begin  mem[ cur_list.tail_field ].hh.rh := new_disc; cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ; end; end ;
  begin  mem[ cur_list.tail_field ].hh.rh := new_kern( font_info[kern_base[  main_f]+256*   main_j.b2 +   main_j.b3 ].int  ); cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ; goto main_loop_move;
  end;
if cur_l= 256   then lft_hit:=true
else if lig_stack=-{0xfffffff=}268435455   then rt_hit:=true;
begin if interrupt<>0 then pause_for_instructions; end ; {allow a way out in case there's an infinite ligature loop}
case  main_j.b2  of
 1 , 5 :begin cur_l:= main_j.b3 ; {\.[=:\?], \.[=:\?>]}
  main_i:= font_info[char_base[ main_f]+effective_char(true, main_f,  cur_l)].qqqq ; ligature_present:=true;
  end;
 2 , 6 :begin cur_r:= main_j.b3 ; {\.[\?=:], \.[\?=:>]}
  if lig_stack=-{0xfffffff=}268435455   then {right boundary character is being consumed}
    begin lig_stack:=new_lig_item(cur_r); bchar:= 256  ;
    end
  else if  ( lig_stack>=hi_mem_min)  then {|link(lig_stack)=null|}
    begin main_p:=lig_stack; lig_stack:=new_lig_item(cur_r);
     mem[    lig_stack+1 ].hh.rh  :=main_p;
    end
  else   mem[ lig_stack].hh.b1 :=cur_r;
  end;
 3 :begin cur_r:= main_j.b3 ; {\.[\?=:\?]}
  main_p:=lig_stack; lig_stack:=new_lig_item(cur_r);
   mem[ lig_stack].hh.rh :=main_p;
  end;
 7 , 11 :begin if cur_l< 256   then begin if  mem[ cur_q].hh.rh >-{0xfffffff=}268435455   then if   mem[ cur_list.tail_field ].hh.b1 = hyphen_char[ main_f]  then ins_disc:=true; if ligature_present then  begin main_p:=new_ligature(main_f,cur_l, mem[ cur_q].hh.rh ); if lft_hit then begin  mem[ main_p].hh.b1 :=2; lft_hit:=false; end; if   false then if lig_stack=-{0xfffffff=}268435455   then begin incr( mem[ main_p].hh.b1 ); rt_hit:=false; end;  mem[ cur_q].hh.rh :=main_p; cur_list.tail_field :=main_p; ligature_present:=false; end ; if ins_disc then begin ins_disc:=false; if cur_list.mode_field >0 then begin  mem[ cur_list.tail_field ].hh.rh := new_disc; cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ; end; end ; {\.[\?=:\?>], \.[\?=:\?>>]}
  cur_q:=cur_list.tail_field ; cur_l:= main_j.b3 ;
  main_i:= font_info[char_base[ main_f]+effective_char(true, main_f,  cur_l)].qqqq ; ligature_present:=true;
  end;
 else  begin cur_l:= main_j.b3 ; ligature_present:=true; {\.[=:]}
  if lig_stack=-{0xfffffff=}268435455   then goto main_loop_wrapup
  else goto main_loop_move+1;
  end
 end ;
if  main_j.b2 > 4  then
  if  main_j.b2 <> 7  then goto main_loop_wrapup;
if cur_l< 256   then goto main_lig_loop;
main_k:=bchar_label[main_f]; goto main_lig_loop+1;
end

;
if  main_j.b0 = 0  then incr(main_k)
else begin if  main_j.b0 >= 128   then goto main_loop_wrapup;
  main_k:=main_k+   main_j.b0  +1;
  end;
goto main_lig_loop+1

;
main_loop_move_lig:
{ Move the cursor past a pseudo-ligature, then |goto main_loop_lookahead| or |main_lig_loop| }
main_p:= mem[    lig_stack+1 ].hh.rh  ;
if main_p>-{0xfffffff=}268435455   then begin  mem[ cur_list.tail_field ].hh.rh := main_p; cur_list.tail_field := mem[ cur_list.tail_field ].hh.rh ; end ; {append a single character}
temp_ptr:=lig_stack; lig_stack:= mem[ temp_ptr].hh.rh ;
free_node(temp_ptr,small_node_size);
main_i:= font_info[char_base[ main_f]+effective_char(true, main_f,  cur_l)].qqqq ; ligature_present:=true;
if lig_stack=-{0xfffffff=}268435455   then
  if main_p>-{0xfffffff=}268435455   then goto main_loop_lookahead
  else cur_r:=bchar
else cur_r:=  mem[ lig_stack].hh.b1 ;
goto main_lig_loop



;
append_normal_space:
{ Append a normal inter-word space to the current list, then |goto big_switch| }
if  eqtb[  glue_base+   space_skip_code].hh.rh    =mem_bot  then
  begin 
{ Find the glue specification, |main_p|, for text spaces in the current font }
begin main_p:=font_glue[ eqtb[  cur_font_loc].hh.rh   ];
if main_p=-{0xfffffff=}268435455   then
  begin main_p:=new_spec(mem_bot ); main_k:=param_base[ eqtb[  cur_font_loc].hh.rh   ]+space_code;
   mem[ main_p+width_offset].int  :=font_info[main_k].int ; {that's |space(cur_font)|}
   mem[ main_p+2].int  :=font_info[main_k+1].int ; {and |space_stretch(cur_font)|}
   mem[ main_p+3].int  :=font_info[main_k+2].int ; {and |space_shrink(cur_font)|}
  font_glue[ eqtb[  cur_font_loc].hh.rh   ]:=main_p;
  end;
end

;
  temp_ptr:=new_glue(main_p);
  end
else temp_ptr:=new_param_glue(space_skip_code);
 mem[ cur_list.tail_field ].hh.rh :=temp_ptr; cur_list.tail_field :=temp_ptr;
goto big_switch

;
exit:end;



{ 1055. \[47] Building boxes and lists }

{tangle:pos tex.web:20502:31: }

{ The most important parts of |main_control| are concerned with \TeX's
chief mission of box-making. We need to control the activities that put
entries on vlists and hlists, as well as the activities that convert
those lists into boxes. All of the necessary machinery has already been
developed; it remains for us to ``push the buttons'' at the right times. }

{ 1062. }

{tangle:pos tex.web:20619:5: }

{ Many of the actions related to box-making are triggered by the appearance
of braces in the input. For example, when the user says `\.[\\hbox]
\.[to] \.[100pt\[$\langle\,\hbox[\rm hlist]\,\rangle$\]]' in vertical mode,
the information about the box size (100pt, |exactly|) is put onto |save_stack|
with a level boundary word just above it, and |cur_group:=adjusted_hbox_group|;
\TeX\ enters restricted horizontal mode to process the hlist. The right
brace eventually causes |save_stack| to be restored to its former state,
at which time the information about the box size (100pt, |exactly|) is
available once again; a box is packaged and we leave restricted horizontal
mode, appending the new box to the current list of the enclosing mode
(in this case to the current list of vertical mode), followed by any
vertical adjustments that were removed from the box by |hpack|.

The next few sections of the program are therefore concerned with the
treatment of left and right curly braces. }

{ 1284. }

{tangle:pos tex.web:23583:4: }

{ The |error| routine calls on |give_err_help| if help is requested from
the |err_help| parameter. } procedure give_err_help;
begin token_show( eqtb[  err_help_loc].hh.rh   );
end;



{ 1303. }

{tangle:pos tex.web:23781:5: }

{ Corresponding to the procedure that dumps a format file, we have a function
that reads one in. The function returns |false| if the dumped format is
incompatible with the present \TeX\ table sizes, etc. } { \4 }
{ Declare the function called |open_fmt_file| }
function open_fmt_file:boolean;
label found,exit;
var j:0..buf_size; {the first space after the format file name}
begin j:=cur_input.loc_field ;
if buffer[cur_input.loc_field ]={"&"=}38 then
  begin incr(cur_input.loc_field ); j:=cur_input.loc_field ; buffer[last]:={" "=}32;
  while buffer[j]<>{" "=}32 do incr(j);
  pack_buffered_name(0,cur_input.loc_field ,j-1); {Kpathsea does everything}
  if w_open_in(fmt_file) then goto found;
     ;
  write(stdout ,'Sorry, I can''t find the format `') ;
  fputs (stringcast(name_of_file + 1), stdout);
  write(stdout ,'''; will try `') ;
  fputs (TEX_format_default + 1, stdout);
  writeln( stdout ,'''.')  ;
{ \xref[Sorry, I can't find...] }
   fflush (stdout ) ;
  end;
  {now pull out all the stops: try for the system \.[plain] file}
pack_buffered_name(format_default_length-format_ext_length,1,0);
if not w_open_in(fmt_file) then
  begin    ;
  write(stdout ,'I can''t find the format file `') ;
  fputs (TEX_format_default + 1, stdout);
  writeln( stdout ,'''!')  ;
{ \xref[I can't find the format...] }
{ \xref[plain] }
  open_fmt_file:=false;  goto exit ;
  end;
found:cur_input.loc_field :=j; open_fmt_file:=true;
exit:end;

 
function load_fmt_file:boolean;
label bad_fmt,exit;
var j, k:integer; {all-purpose indices}
 p, q: halfword ; {all-purpose pointers}
 x: integer; {something undumped}
 format_engine: ^ ASCII_code ;
 dummy_xord: ASCII_code;
 dummy_xchr:  ASCII_code ;
 dummy_xprn: ASCII_code;
begin 
{ Undump constants for consistency check }
 ifdef('INITEX')  if ini_version then begin 
libc_free(font_info); libc_free(str_pool); libc_free(str_start);
libc_free(yhash); libc_free(zeqtb); libc_free(yzmem);
 end; endif('INITEX')  
undump_int(x);
 if debug_format_file then begin write (stderr, 'fmtdebug:', 'format magic number');  writeln( stderr, ' = ',   x) ; end; ;
if x<>{0x57325458=}1462916184 then goto bad_fmt; {not a format file}
undump_int(x);
 if debug_format_file then begin write (stderr, 'fmtdebug:', 'engine name size');  writeln( stderr, ' = ',   x) ; end; ;
if (x<0) or (x>256) then goto bad_fmt; {corrupted format file}
format_engine:=xmalloc_array( ASCII_code , x);
undump_things(format_engine[0], x);
format_engine[x-1]:=0; {force string termination, just in case}
if strcmp(engine_name, stringcast(format_engine)) then
  begin    ;
  writeln( stdout ,'---! ',   stringcast(  name_of_file+  1), ' was written by ',   format_engine)  ;
  libc_free(format_engine);
  goto bad_fmt;
end;
libc_free(format_engine);
undump_int(x);
 if debug_format_file then begin write (stderr, 'fmtdebug:', 'string pool checksum');  writeln( stderr, ' = ',   x) ; end; ;
if x<>@$ then begin {check that strings are the same}
     ;
  writeln( stdout ,'---! ',   stringcast(  name_of_file+  1),
           ' made by different executable version, strings are different')  ;
  goto bad_fmt;
end;

{ Undump |xord|, |xchr|, and |xprn| }
if translate_filename then begin
  for k:=0 to 255 do undump_things(dummy_xord, 1);
  for k:=0 to 255 do undump_things(dummy_xchr, 1);
  for k:=0 to 255 do undump_things(dummy_xprn, 1);
  end
else begin
  undump_things(xord[0], 256);
  undump_things(xchr[0], 256);
  undump_things(xprn[0], 256);
  if eight_bit_p then
    for k:=0 to 255 do
      xprn[k]:=1;
end;


;
undump_int(x);
if x<>{0xfffffff=}268435455  then goto bad_fmt; {check |max_halfword|}
undump_int(hash_high);
  if (hash_high<0)or(hash_high>sup_hash_extra) then goto bad_fmt;
  if hash_extra<hash_high then hash_extra:=hash_high;
  eqtb_top:=eqtb_size+hash_extra;
  if hash_extra=0 then hash_top:=undefined_control_sequence else
        hash_top:=eqtb_top;
  yhash:=xmalloc_array(two_halves,1+hash_top-hash_offset);
  hash:=yhash - hash_offset;
   hash[ hash_base].lh :=0;  hash[ hash_base].rh :=0;
  for x:=hash_base+1 to hash_top do hash[x]:=hash[hash_base];
  zeqtb:=xmalloc_array (memory_word,eqtb_top+1);
  eqtb:=zeqtb;

   eqtb[  undefined_control_sequence].hh.b0  :=undefined_cs;
   eqtb[  undefined_control_sequence].hh.rh  :=-{0xfffffff=}268435455  ;
   eqtb[  undefined_control_sequence].hh.b1  :=level_zero;
  for x:=eqtb_size+1 to eqtb_top do
    eqtb[x]:=eqtb[undefined_control_sequence];
undump_int(x);  if debug_format_file then begin write (stderr, 'fmtdebug:', 'mem_bot');  writeln( stderr, ' = ',   x) ; end; ;
if x<>mem_bot then goto bad_fmt;
undump_int(mem_top);  if debug_format_file then begin write (stderr, 'fmtdebug:', 'mem_top');  writeln( stderr, ' = ',   mem_top) ; end; ;
if mem_bot+1100>mem_top then goto bad_fmt;


cur_list.head_field :=mem_top-1 ; cur_list.tail_field :=mem_top-1 ;
     page_tail:=mem_top-2 ;  {page initialization}

mem_min := mem_bot - extra_mem_bot;
mem_max := mem_top + extra_mem_top;

yzmem:=xmalloc_array (memory_word, mem_max - mem_min + 1);
zmem := yzmem - mem_min;   {this pointer arithmetic fails with some compilers}
mem := zmem;
undump_int(x);
if x<>eqtb_size then goto bad_fmt;
undump_int(x);
if x<>hash_prime then goto bad_fmt;
undump_int(x);
if x<>hyph_prime then goto bad_fmt

;

{ Undump ML\TeX-specific data }
undump_int(x);   {check magic constant of ML\TeX}
if x<>{0x4d4c5458=}1296847960 then goto bad_fmt;
undump_int(x);   {undump |mltex_p| flag into |mltex_enabled_p|}
if x=1 then mltex_enabled_p:=true
else if x<>0 then goto bad_fmt;

;

{ Undump the string pool }
begin undump_int(x); if x< 0 then goto bad_fmt; if x> sup_pool_size- pool_free then  begin    ; writeln( stdout ,'---! Must increase the ','string pool size')  ; { \xref[Must increase the x] } goto bad_fmt; end  else  if debug_format_file then begin write (stderr, 'fmtdebug:', 'string pool size');  writeln( stderr, ' = ',   x) ; end; ;  pool_ptr:=x; end ;
if pool_size<pool_ptr+pool_free then
  pool_size:=pool_ptr+pool_free;
begin undump_int(x); if x< 0 then goto bad_fmt; if x> sup_max_strings- strings_free then  begin    ; writeln( stdout ,'---! Must increase the ','sup strings')  ; { \xref[Must increase the x] } goto bad_fmt; end  else  if debug_format_file then begin write (stderr, 'fmtdebug:', 'sup strings');  writeln( stderr, ' = ',   x) ; end; ;  str_ptr:=x; end ;

if max_strings<str_ptr+strings_free then
  max_strings:=str_ptr+strings_free;
str_start:=xmalloc_array(pool_pointer, max_strings);
undump_checked_things(0, pool_ptr, str_start[0], str_ptr+1);

str_pool:=xmalloc_array(packed_ASCII_code, pool_size);
undump_things(str_pool[0], pool_ptr);
init_str_ptr:=str_ptr; init_pool_ptr:=pool_ptr

;

{ Undump the dynamic memory }
begin undump_int(x); if (x< mem_bot +glue_spec_size +glue_spec_size +glue_spec_size +glue_spec_size +glue_spec_size-1 + 1000) or (x> mem_top-13 - 1) then goto bad_fmt else  lo_mem_max:=x; end ;
begin undump_int(x); if (x< mem_bot +glue_spec_size +glue_spec_size +glue_spec_size +glue_spec_size +glue_spec_size-1 + 1) or (x> lo_mem_max) then goto bad_fmt else  rover:=x; end ;
p:=mem_bot; q:=rover;
repeat undump_things(mem[p], q+2-p);
{If the format file is messed up, that addition to |p| might cause it to
 become garbage. Report from Gregory James DUCK to Karl, 14 Sep 2023.
 Also changed in \MF. Fix from DRF, who explains: we test before doing the
 addition to avoid assuming silent wrap-around overflow, and also to to
 catch cases where |node_size| was, say, bogusly the equivalent of $-1$
 and thus |p+node_size| would still look valid.}
if (  mem[ q].hh.lh >lo_mem_max-q) or (  mem[  q+ 1].hh.rh  >lo_mem_max)
   or ((q>=  mem[  q+ 1].hh.rh  )and(  mem[  q+ 1].hh.rh  <>rover))
then goto bad_fmt;
p:=q+  mem[ q].hh.lh ;
q:=  mem[  q+ 1].hh.rh  ;
until q=rover;
undump_things(mem[p], lo_mem_max+1-p);
if mem_min<mem_bot-2 then {make more low memory available}
  begin p:=  mem[  rover+ 1].hh.lh  ; q:=mem_min+1;
   mem[ mem_min].hh.rh :=-{0xfffffff=}268435455  ;  mem[ mem_min].hh.lh :=-{0xfffffff=}268435455  ; {we don't use the bottom word}
    mem[  p+ 1].hh.rh  :=q;   mem[  rover+ 1].hh.lh  :=q;

    mem[  q+ 1].hh.rh  :=rover;   mem[  q+ 1].hh.lh  :=p;  mem[ q].hh.rh := {0xfffffff=}268435455  ;
    mem[ q].hh.lh :=mem_bot-q;
  end;
begin undump_int(x); if (x< lo_mem_max+ 1) or (x> mem_top-13 ) then goto bad_fmt else  hi_mem_min:=x; end ;
begin undump_int(x); if (x< -{0xfffffff=}268435455  ) or (x> mem_top) then goto bad_fmt else  avail:=x; end ; mem_end:=mem_top;
undump_things (mem[hi_mem_min], mem_end+1-hi_mem_min);
undump_int(var_used); undump_int(dyn_used)

;

{ Undump the table of equivalents }

{ Undump regions 1 to 6 of |eqtb| }
k:=active_base;
repeat undump_int(x);
if (x<1)or(k+x>eqtb_size+1) then goto bad_fmt;
undump_things(eqtb[k], x);
k:=k+x;
undump_int(x);
if (x<0)or(k+x>eqtb_size+1) then goto bad_fmt;
for j:=k to k+x-1 do eqtb[j]:=eqtb[k-1];
k:=k+x;
until k>eqtb_size;
if hash_high>0 then undump_things(eqtb[eqtb_size+1],hash_high);
  {undump |hash_extra| part}

;
begin undump_int(x); if (x< hash_base) or (x> hash_top) then goto bad_fmt else  par_loc:=x; end ;
par_token:={07777=}4095 +par_loc;

begin undump_int(x); if (x< hash_base) or (x> hash_top) then goto bad_fmt else  write_loc:=x; end ;


{ Undump the hash table }
begin undump_int(x); if (x< hash_base) or (x> frozen_control_sequence) then goto bad_fmt else  hash_used:=x; end ; p:=hash_base-1;
repeat begin undump_int(x); if (x< p+ 1) or (x> hash_used) then goto bad_fmt else  p:=x; end ; undump_hh(hash[p]);
until p=hash_used;
undump_things (hash[hash_used+1], undefined_control_sequence-1-hash_used);
if debug_format_file then begin
  print_csnames (hash_base, undefined_control_sequence - 1);
end;
if hash_high > 0 then begin
  undump_things (hash[eqtb_size+1], hash_high);
  if debug_format_file then begin
    print_csnames (eqtb_size + 1, hash_high - (eqtb_size + 1));
  end;
end;
undump_int(cs_count)



;

{ Undump the font information }
begin undump_int(x); if x< 7 then goto bad_fmt; if x> sup_font_mem_size then  begin    ; writeln( stdout ,'---! Must increase the ','font mem size')  ; { \xref[Must increase the x] } goto bad_fmt; end  else  if debug_format_file then begin write (stderr, 'fmtdebug:', 'font mem size');  writeln( stderr, ' = ',   x) ; end; ;  fmem_ptr:=x; end ;
if fmem_ptr>font_mem_size then font_mem_size:=fmem_ptr;
font_info:=xmalloc_array(fmemory_word, font_mem_size);
undump_things(font_info[0], fmem_ptr);

begin undump_int(x); if x< font_base then goto bad_fmt; if x> font_base+ max_font_max then  begin    ; writeln( stdout ,'---! Must increase the ','font max')  ; { \xref[Must increase the x] } goto bad_fmt; end  else  if debug_format_file then begin write (stderr, 'fmtdebug:', 'font max');  writeln( stderr, ' = ',   x) ; end; ;  font_ptr:=x; end ;
{This undumps all of the font info, despite the name.}

{ Undump the array info for internal font number |k| }
begin {Allocate the font arrays}
font_check:=xmalloc_array(four_quarters, font_max);
font_size:=xmalloc_array(scaled, font_max);
font_dsize:=xmalloc_array(scaled, font_max);
font_params:=xmalloc_array(font_index, font_max);
font_name:=xmalloc_array(str_number, font_max);
font_area:=xmalloc_array(str_number, font_max);
font_bc:=xmalloc_array(eight_bits, font_max);
font_ec:=xmalloc_array(eight_bits, font_max);
font_glue:=xmalloc_array(halfword, font_max);
hyphen_char:=xmalloc_array(integer, font_max);
skew_char:=xmalloc_array(integer, font_max);
bchar_label:=xmalloc_array(font_index, font_max);
font_bchar:=xmalloc_array(nine_bits, font_max);
font_false_bchar:=xmalloc_array(nine_bits, font_max);
char_base:=xmalloc_array(integer, font_max);
width_base:=xmalloc_array(integer, font_max);
height_base:=xmalloc_array(integer, font_max);
depth_base:=xmalloc_array(integer, font_max);
italic_base:=xmalloc_array(integer, font_max);
lig_kern_base:=xmalloc_array(integer, font_max);
kern_base:=xmalloc_array(integer, font_max);
exten_base:=xmalloc_array(integer, font_max);
param_base:=xmalloc_array(integer, font_max);

undump_things(font_check[font_base ], font_ptr+1-font_base );
undump_things(font_size[font_base ], font_ptr+1-font_base );
undump_things(font_dsize[font_base ], font_ptr+1-font_base );
undump_checked_things(-{0xfffffff=}268435455 , {0xfffffff=}268435455 ,
                      font_params[font_base ], font_ptr+1-font_base );
undump_things(hyphen_char[font_base ], font_ptr+1-font_base );
undump_things(skew_char[font_base ], font_ptr+1-font_base );
undump_upper_check_things(str_ptr, font_name[font_base ], font_ptr+1-font_base );
undump_upper_check_things(str_ptr, font_area[font_base ], font_ptr+1-font_base );
{There's no point in checking these values against the range $[0,255]$,
 since the data type is |unsigned char|, and all values of that type are
 in that range by definition.}
undump_things(font_bc[font_base ], font_ptr+1-font_base );
undump_things(font_ec[font_base ], font_ptr+1-font_base );
undump_things(char_base[font_base ], font_ptr+1-font_base );
undump_things(width_base[font_base ], font_ptr+1-font_base );
undump_things(height_base[font_base ], font_ptr+1-font_base );
undump_things(depth_base[font_base ], font_ptr+1-font_base );
undump_things(italic_base[font_base ], font_ptr+1-font_base );
undump_things(lig_kern_base[font_base ], font_ptr+1-font_base );
undump_things(kern_base[font_base ], font_ptr+1-font_base );
undump_things(exten_base[font_base ], font_ptr+1-font_base );
undump_things(param_base[font_base ], font_ptr+1-font_base );
undump_checked_things(-{0xfffffff=}268435455 , lo_mem_max,
                     font_glue[font_base ], font_ptr+1-font_base );
undump_checked_things(0, fmem_ptr-1,
                     bchar_label[font_base ], font_ptr+1-font_base );
undump_checked_things(min_quarterword,  256  ,
                     font_bchar[font_base ], font_ptr+1-font_base );
undump_checked_things(min_quarterword,  256  ,
                     font_false_bchar[font_base ], font_ptr+1-font_base );
end

;

;

{ Undump the hyphenation tables }
begin undump_int(x); if x< 0 then goto bad_fmt; if x> hyph_size then  begin    ; writeln( stdout ,'---! Must increase the ','hyph_size')  ; { \xref[Must increase the x] } goto bad_fmt; end  else  if debug_format_file then begin write (stderr, 'fmtdebug:', 'hyph_size');  writeln( stderr, ' = ',   x) ; end; ;  hyph_count:=x; end ;
begin undump_int(x); if x< hyph_prime then goto bad_fmt; if x> hyph_size then  begin    ; writeln( stdout ,'---! Must increase the ','hyph_size')  ; { \xref[Must increase the x] } goto bad_fmt; end  else  if debug_format_file then begin write (stderr, 'fmtdebug:', 'hyph_size');  writeln( stderr, ' = ',   x) ; end; ;  hyph_next:=x; end ;
j:=0;
for k:=1 to hyph_count do
  begin undump_int(j); if j<0 then goto bad_fmt;
   if j>65535 then
   begin hyph_next:= j div 65536; j:=j - hyph_next * 65536; end
       else hyph_next:=0;
   if (j>=hyph_size)or(hyph_next>hyph_size) then goto bad_fmt;
   hyph_link[j]:=hyph_next;
  begin undump_int(x); if (x< 0) or (x> str_ptr) then goto bad_fmt else  hyph_word[ j]:=x; end ;
  begin undump_int(x); if (x< -{0xfffffff=}268435455 ) or (x> {0xfffffff=}268435455 ) then goto bad_fmt else  hyph_list[ j]:=x; end ;
  end;
  {|j| is now the largest occupied location in |hyph_word|}
  incr(j);
  if j<hyph_prime then j:=hyph_prime;
  hyph_next:=j;
  if hyph_next >= hyph_size then hyph_next:=hyph_prime else
  if hyph_next >= hyph_prime then incr(hyph_next);
begin undump_int(x); if x< 0 then goto bad_fmt; if x> trie_size then  begin    ; writeln( stdout ,'---! Must increase the ','trie size')  ; { \xref[Must increase the x] } goto bad_fmt; end  else  if debug_format_file then begin write (stderr, 'fmtdebug:', 'trie size');  writeln( stderr, ' = ',   x) ; end; ;  j:=x; end ;  ifdef('INITEX')  trie_max:=j; endif('INITEX') 
{These first three haven't been allocated yet unless we're \.[INITEX];
 we do that precisely so we don't allocate more space than necessary.}
if not trie_trl then trie_trl:=xmalloc_array(trie_pointer,j+1);
undump_things(trie_trl[0], j+1);
if not trie_tro then trie_tro:=xmalloc_array(trie_pointer,j+1);
undump_things(trie_tro[0], j+1);
if not trie_trc then trie_trc:=xmalloc_array(quarterword, j+1);
undump_things(trie_trc[0], j+1);
begin undump_int(x); if x< 0 then goto bad_fmt; if x> trie_op_size then  begin    ; writeln( stdout ,'---! Must increase the ','trie op size')  ; { \xref[Must increase the x] } goto bad_fmt; end  else  if debug_format_file then begin write (stderr, 'fmtdebug:', 'trie op size');  writeln( stderr, ' = ',   x) ; end; ;  j:=x; end ;  ifdef('INITEX')  trie_op_ptr:=j; endif('INITEX') 
{I'm not sure we have such a strict limitation (64) on these values, so
 let's leave them unchecked.}
undump_things(hyf_distance[1], j);
undump_things(hyf_num[1], j);
undump_upper_check_things(max_trie_op, hyf_next[1], j);
ifdef('INITEX')  for k:=0 to 255 do trie_used[k]:=min_quarterword; endif('INITEX')  

k:=256;
while j>0 do
  begin begin undump_int(x); if (x< 0) or (x> k- 1) then goto bad_fmt else  k:=x; end ; begin undump_int(x); if (x< 1) or (x> j) then goto bad_fmt else  x:=x; end ; ifdef('INITEX')  trie_used[k]:= x ; endif('INITEX')  

  j:=j-x; op_start[k]:= j ;
  end;
 ifdef('INITEX')  trie_not_ready:=false  endif('INITEX') 

;

{ Undump a couple more things and the closing check word }
begin undump_int(x); if (x< batch_mode) or (x> error_stop_mode) then goto bad_fmt else  interaction:=x; end ;
if interaction_option<>unspecified_mode then interaction:=interaction_option;
begin undump_int(x); if (x< 0) or (x> str_ptr) then goto bad_fmt else  format_ident:=x; end ;
undump_int(x);
if x<>69069 then goto bad_fmt

;
load_fmt_file:=true;  goto exit ; {it worked!}
bad_fmt:    ;
  writeln( stdout ,'(Fatal format file error; I''m stymied)')  ;
{ \xref[Fatal format file error] }
load_fmt_file:=false;
exit:end;



{ 1306. }

{tangle:pos tex.web:23835:66: }

{ The inverse macros are slightly more complicated, since we need to check
the range of the values we are reading in. We say `|undump(a)(b)(x)|' to
read an integer value |x| that is supposed to be in the range |a<=x<=b|.
System error messages should be suppressed when undumping.
\xref[system dependencies] }

{ 1330. \[51] The main program }

{tangle:pos tex.web:24222:23: }

{ This is it: the part of \TeX\ that executes all those procedures we have
written.

Well---almost. Let's leave space for a few more routines that we may
have forgotten. } 
{ Last-minute procedures }
procedure close_files_and_terminate;
var k:integer; {all-purpose index}
begin 
{ Finish the extensions }
for k:=0 to 15 do if write_open[k] then a_close(write_file[k])

; eqtb[int_base+ new_line_char_code].int  :=-1;
 ifdef('STAT')  if eqtb[int_base+ tracing_stats_code].int  >0 then 
{ Output statistics about this job }
if log_opened then
  begin writeln( log_file,' ')  ;
  writeln( log_file,'Here is how much of TeX''s memory',' you used:')  ;
{ \xref[Here is how much...] }
  write(log_file,' ', str_ptr- init_str_ptr: 1,' string') ;
  if str_ptr<>init_str_ptr+1 then write(log_file,'s') ;
  writeln( log_file,' out of ',   max_strings-  init_str_ptr:  1)  ;

  writeln( log_file,' ',  pool_ptr-  init_pool_ptr:  1,' string characters out of ',
      pool_size-  init_pool_ptr:  1)  ;

  writeln( log_file,' ',  lo_mem_max-  mem_min+  mem_end-  hi_mem_min+  2:  1, 
    ' words of memory out of ',  mem_end+  1-  mem_min:  1)  ;

  writeln( log_file,' ',  cs_count:  1,' multiletter control sequences out of ',
      hash_size:  1, '+',   hash_extra:  1)  ;

  write(log_file,' ', fmem_ptr: 1,' words of font info for ',
     font_ptr- font_base: 1,' font') ;
  if font_ptr<>font_base+1 then write(log_file,'s') ;
  writeln( log_file,', out of ',  font_mem_size:  1,' for ',  font_max-  font_base:  1)  ;

  write(log_file,' ', hyph_count: 1,' hyphenation exception') ;
  if hyph_count<>1 then write(log_file,'s') ;
  writeln( log_file,' out of ',  hyph_size:  1)  ;

  writeln( log_file,' ',  max_in_stack:  1,'i,',  max_nest_stack:  1,'n,', 
      max_param_stack:  1,'p,', 
      max_buf_stack+  1:  1,'b,', 
      max_save_stack+  6:  1,'s stack positions out of ', 
      stack_size:  1,'i,',
      nest_size:  1,'n,',
      param_size:  1,'p,',
      buf_size:  1,'b,',
      save_size:  1,'s')  ;
  end

;  endif('STAT') 

   ; 
{ Finish the \.[DVI] file }
while cur_s>-1 do
  begin if cur_s>0 then  begin dvi_buf[dvi_ptr]:= pop; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end 
  else  begin  begin dvi_buf[dvi_ptr]:= eop; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ; incr(total_pages);
    end;
  decr(cur_s);
  end;
if total_pages=0 then print_nl({"No pages of output."=}850)
{ \xref[No pages of output] }
else if cur_s<>-2 then
  begin  begin dvi_buf[dvi_ptr]:= post; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ; {beginning of the postamble}
  dvi_four(last_bop); last_bop:=dvi_offset+dvi_ptr-5; {|post| location}
  dvi_four(25400000); dvi_four(473628672); {conversion ratio for sp}
  prepare_mag; dvi_four(eqtb[int_base+ mag_code].int  ); {magnification factor}
  dvi_four(max_v); dvi_four(max_h);

   begin dvi_buf[dvi_ptr]:= max_push  div  256; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;  begin dvi_buf[dvi_ptr]:= max_push  mod  256; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;

   begin dvi_buf[dvi_ptr]:=( total_pages  div  256)  mod  256; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;  begin dvi_buf[dvi_ptr]:= total_pages  mod  256; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;

  
{ Output the font definitions for all fonts that were used }
while font_ptr>font_base do
  begin if font_used[font_ptr] then dvi_font_def(font_ptr);
  decr(font_ptr);
  end

;
   begin dvi_buf[dvi_ptr]:= post_post; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ; dvi_four(last_bop);  begin dvi_buf[dvi_ptr]:= id_byte; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ;

  ifdef ('IPC')
  k:=7-((3+dvi_offset+dvi_ptr) mod 4); {the number of 223's}
endif ('IPC')
ifndef ('IPC')
  k:=4+((dvi_buf_size-dvi_ptr) mod 4); {the number of 223's}
endifn ('IPC')
  while k>0 do
    begin  begin dvi_buf[dvi_ptr]:= 223; incr(dvi_ptr); if dvi_ptr=dvi_limit then dvi_swap; end ; decr(k);
    end;
  
{ Empty the last bytes out of |dvi_buf| }
if dvi_limit=half_buf then write_dvi(half_buf,dvi_buf_size-1);
if dvi_ptr>({0x7fffffff=}2147483647-dvi_offset) then
  begin cur_s:=-2;
  fatal_error({"dvi length exceeds ""7FFFFFFF"=}839);
{ \xref[dvi length exceeds...] }
  end;
if dvi_ptr>0 then write_dvi(0,dvi_ptr-1)

;
  print_nl({"Output written on "=}851); print_file_name(0, output_file_name, 0);
{ \xref[Output written on x] }
  print({" ("=}284); print_int(total_pages);
  if total_pages<>1 then print({" pages"=}852)
  else print({" page"=}853);
  print({", "=}854); print_int(dvi_offset+dvi_ptr); print({" bytes)."=}855);
  b_close(dvi_file);
  end

;
if log_opened then
  begin writeln( log_file)  ; a_close(log_file); selector:=selector-2;
  if selector=term_only then
    begin print_nl({"Transcript written on "=}1294);
{ \xref[Transcript written...] }
    print_file_name(0,  texmf_log_name , 0); print_char({"."=}46);
    end;
  end;
print_ln;
if (edit_name_start<>0) and (interaction>batch_mode) then
  call_edit(str_pool,edit_name_start,edit_name_length,edit_line);
end;


procedure final_cleanup;
label exit;
var c:small_number; {0 for \.[\\end], 1 for \.[\\dump]}
begin c:=cur_chr; if c<>1 then eqtb[int_base+ new_line_char_code].int  :=-1;
if job_name=0 then open_log_file;
while input_ptr>0 do
  if cur_input.state_field =token_list then end_token_list else end_file_reading;
while open_parens>0 do
  begin print({" )"=}1295); decr(open_parens);
  end;
if cur_level>level_one then
  begin print_nl({"("=}40); print_esc({"end occurred "=}1296);
  print({"inside a group at level "=}1297);
{ \xref[end_][\.[(\\end occurred...)]] }
  print_int(cur_level-level_one); print_char({")"=}41);
  end;
while cond_ptr<>-{0xfffffff=}268435455   do
  begin print_nl({"("=}40); print_esc({"end occurred "=}1296);
  print({"when "=}1298); print_cmd_chr(if_test,cur_if);
  if if_line<>0 then
    begin print({" on line "=}1299); print_int(if_line);
    end;
  print({" was incomplete)"=}1300);
  if_line:=mem[ cond_ptr+1].int ;
  cur_if:= mem[ cond_ptr].hh.b1 ; temp_ptr:=cond_ptr;
  cond_ptr:= mem[ cond_ptr].hh.rh ; free_node(temp_ptr,if_node_size);
  end;
if history<>spotless then
 if ((history=warning_issued)or(interaction<error_stop_mode)) then
  if selector=term_and_log then
  begin selector:=term_only;
  print_nl({"(see the transcript file for additional information)"=}1301);
{ \xref[see the transcript file...] }
  selector:=term_and_log;
  end;
if c=1 then
  begin  ifdef('INITEX')  if ini_version then begin  for c:=top_mark_code to split_bot_mark_code do
    if cur_mark[c]<>-{0xfffffff=}268435455   then delete_token_ref(cur_mark[c]);
  if last_glue<>{0xfffffff=}268435455  then delete_glue_ref(last_glue);
  store_fmt_file;  goto exit ; end; endif('INITEX')  

  print_nl({"(\dump is performed only by INITEX)"=}1302);  goto exit ;
{ \xref[dump_][\.[\\dump...only by INITEX]] }
  end;
exit:end;


 ifdef('INITEX')  procedure init_prim; {initialize all the primitives}
begin no_new_control_sequence:=false;

{ Put each... }
primitive({"lineskip"=}381,assign_glue,glue_base+line_skip_code);

 { \xref[line_skip_][\.[\\lineskip] primitive] }
primitive({"baselineskip"=}382,assign_glue,glue_base+baseline_skip_code);

 { \xref[baseline_skip_][\.[\\baselineskip] primitive] }
primitive({"parskip"=}383,assign_glue,glue_base+par_skip_code);

 { \xref[par_skip_][\.[\\parskip] primitive] }
primitive({"abovedisplayskip"=}384,assign_glue,glue_base+above_display_skip_code);

 { \xref[above_display_skip_][\.[\\abovedisplayskip] primitive] }
primitive({"belowdisplayskip"=}385,assign_glue,glue_base+below_display_skip_code);

 { \xref[below_display_skip_][\.[\\belowdisplayskip] primitive] }
primitive({"abovedisplayshortskip"=}386,
  assign_glue,glue_base+above_display_short_skip_code);

 { \xref[above_display_short_skip_][\.[\\abovedisplayshortskip] primitive] }
primitive({"belowdisplayshortskip"=}387,
  assign_glue,glue_base+below_display_short_skip_code);

 { \xref[below_display_short_skip_][\.[\\belowdisplayshortskip] primitive] }
primitive({"leftskip"=}388,assign_glue,glue_base+left_skip_code);

 { \xref[left_skip_][\.[\\leftskip] primitive] }
primitive({"rightskip"=}389,assign_glue,glue_base+right_skip_code);

 { \xref[right_skip_][\.[\\rightskip] primitive] }
primitive({"topskip"=}390,assign_glue,glue_base+top_skip_code);

 { \xref[top_skip_][\.[\\topskip] primitive] }
primitive({"splittopskip"=}391,assign_glue,glue_base+split_top_skip_code);

 { \xref[split_top_skip_][\.[\\splittopskip] primitive] }
primitive({"tabskip"=}392,assign_glue,glue_base+tab_skip_code);

 { \xref[tab_skip_][\.[\\tabskip] primitive] }
primitive({"spaceskip"=}393,assign_glue,glue_base+space_skip_code);

 { \xref[space_skip_][\.[\\spaceskip] primitive] }
primitive({"xspaceskip"=}394,assign_glue,glue_base+xspace_skip_code);

 { \xref[xspace_skip_][\.[\\xspaceskip] primitive] }
primitive({"parfillskip"=}395,assign_glue,glue_base+par_fill_skip_code);

 { \xref[par_fill_skip_][\.[\\parfillskip] primitive] }
primitive({"thinmuskip"=}396,assign_mu_glue,glue_base+thin_mu_skip_code);

 { \xref[thin_mu_skip_][\.[\\thinmuskip] primitive] }
primitive({"medmuskip"=}397,assign_mu_glue,glue_base+med_mu_skip_code);

 { \xref[med_mu_skip_][\.[\\medmuskip] primitive] }
primitive({"thickmuskip"=}398,assign_mu_glue,glue_base+thick_mu_skip_code);

 { \xref[thick_mu_skip_][\.[\\thickmuskip] primitive] }


primitive({"output"=}403,assign_toks,output_routine_loc);
 { \xref[output_][\.[\\output] primitive] }
primitive({"everypar"=}404,assign_toks,every_par_loc);
 { \xref[every_par_][\.[\\everypar] primitive] }
primitive({"everymath"=}405,assign_toks,every_math_loc);
 { \xref[every_math_][\.[\\everymath] primitive] }
primitive({"everydisplay"=}406,assign_toks,every_display_loc);
 { \xref[every_display_][\.[\\everydisplay] primitive] }
primitive({"everyhbox"=}407,assign_toks,every_hbox_loc);
 { \xref[every_hbox_][\.[\\everyhbox] primitive] }
primitive({"everyvbox"=}408,assign_toks,every_vbox_loc);
 { \xref[every_vbox_][\.[\\everyvbox] primitive] }
primitive({"everyjob"=}409,assign_toks,every_job_loc);
 { \xref[every_job_][\.[\\everyjob] primitive] }
primitive({"everycr"=}410,assign_toks,every_cr_loc);
 { \xref[every_cr_][\.[\\everycr] primitive] }
primitive({"errhelp"=}411,assign_toks,err_help_loc);
 { \xref[err_help_][\.[\\errhelp] primitive] }


primitive({"pretolerance"=}425,assign_int,int_base+pretolerance_code);

 { \xref[pretolerance_][\.[\\pretolerance] primitive] }
primitive({"tolerance"=}426,assign_int,int_base+tolerance_code);

 { \xref[tolerance_][\.[\\tolerance] primitive] }
primitive({"linepenalty"=}427,assign_int,int_base+line_penalty_code);

 { \xref[line_penalty_][\.[\\linepenalty] primitive] }
primitive({"hyphenpenalty"=}428,assign_int,int_base+hyphen_penalty_code);

 { \xref[hyphen_penalty_][\.[\\hyphenpenalty] primitive] }
primitive({"exhyphenpenalty"=}429,assign_int,int_base+ex_hyphen_penalty_code);

 { \xref[ex_hyphen_penalty_][\.[\\exhyphenpenalty] primitive] }
primitive({"clubpenalty"=}430,assign_int,int_base+club_penalty_code);

 { \xref[club_penalty_][\.[\\clubpenalty] primitive] }
primitive({"widowpenalty"=}431,assign_int,int_base+widow_penalty_code);

 { \xref[widow_penalty_][\.[\\widowpenalty] primitive] }
primitive({"displaywidowpenalty"=}432,
  assign_int,int_base+display_widow_penalty_code);

 { \xref[display_widow_penalty_][\.[\\displaywidowpenalty] primitive] }
primitive({"brokenpenalty"=}433,assign_int,int_base+broken_penalty_code);

 { \xref[broken_penalty_][\.[\\brokenpenalty] primitive] }
primitive({"binoppenalty"=}434,assign_int,int_base+bin_op_penalty_code);

 { \xref[bin_op_penalty_][\.[\\binoppenalty] primitive] }
primitive({"relpenalty"=}435,assign_int,int_base+rel_penalty_code);

 { \xref[rel_penalty_][\.[\\relpenalty] primitive] }
primitive({"predisplaypenalty"=}436,assign_int,int_base+pre_display_penalty_code);

 { \xref[pre_display_penalty_][\.[\\predisplaypenalty] primitive] }
primitive({"postdisplaypenalty"=}437,assign_int,int_base+post_display_penalty_code);

 { \xref[post_display_penalty_][\.[\\postdisplaypenalty] primitive] }
primitive({"interlinepenalty"=}438,assign_int,int_base+inter_line_penalty_code);

 { \xref[inter_line_penalty_][\.[\\interlinepenalty] primitive] }
primitive({"doublehyphendemerits"=}439,
  assign_int,int_base+double_hyphen_demerits_code);

 { \xref[double_hyphen_demerits_][\.[\\doublehyphendemerits] primitive] }
primitive({"finalhyphendemerits"=}440,
  assign_int,int_base+final_hyphen_demerits_code);

 { \xref[final_hyphen_demerits_][\.[\\finalhyphendemerits] primitive] }
primitive({"adjdemerits"=}441,assign_int,int_base+adj_demerits_code);

 { \xref[adj_demerits_][\.[\\adjdemerits] primitive] }
primitive({"mag"=}442,assign_int,int_base+mag_code);

 { \xref[mag_][\.[\\mag] primitive] }
primitive({"delimiterfactor"=}443,assign_int,int_base+delimiter_factor_code);

 { \xref[delimiter_factor_][\.[\\delimiterfactor] primitive] }
primitive({"looseness"=}444,assign_int,int_base+looseness_code);

 { \xref[looseness_][\.[\\looseness] primitive] }
primitive({"time"=}445,assign_int,int_base+time_code);

 { \xref[time_][\.[\\time] primitive] }
primitive({"day"=}446,assign_int,int_base+day_code);

 { \xref[day_][\.[\\day] primitive] }
primitive({"month"=}447,assign_int,int_base+month_code);

 { \xref[month_][\.[\\month] primitive] }
primitive({"year"=}448,assign_int,int_base+year_code);

 { \xref[year_][\.[\\year] primitive] }
primitive({"showboxbreadth"=}449,assign_int,int_base+show_box_breadth_code);

 { \xref[show_box_breadth_][\.[\\showboxbreadth] primitive] }
primitive({"showboxdepth"=}450,assign_int,int_base+show_box_depth_code);

 { \xref[show_box_depth_][\.[\\showboxdepth] primitive] }
primitive({"hbadness"=}451,assign_int,int_base+hbadness_code);

 { \xref[hbadness_][\.[\\hbadness] primitive] }
primitive({"vbadness"=}452,assign_int,int_base+vbadness_code);

 { \xref[vbadness_][\.[\\vbadness] primitive] }
primitive({"pausing"=}453,assign_int,int_base+pausing_code);

 { \xref[pausing_][\.[\\pausing] primitive] }
primitive({"tracingonline"=}454,assign_int,int_base+tracing_online_code);

 { \xref[tracing_online_][\.[\\tracingonline] primitive] }
primitive({"tracingmacros"=}455,assign_int,int_base+tracing_macros_code);

 { \xref[tracing_macros_][\.[\\tracingmacros] primitive] }
primitive({"tracingstats"=}456,assign_int,int_base+tracing_stats_code);

 { \xref[tracing_stats_][\.[\\tracingstats] primitive] }
primitive({"tracingparagraphs"=}457,assign_int,int_base+tracing_paragraphs_code);

 { \xref[tracing_paragraphs_][\.[\\tracingparagraphs] primitive] }
primitive({"tracingpages"=}458,assign_int,int_base+tracing_pages_code);

 { \xref[tracing_pages_][\.[\\tracingpages] primitive] }
primitive({"tracingoutput"=}459,assign_int,int_base+tracing_output_code);

 { \xref[tracing_output_][\.[\\tracingoutput] primitive] }
primitive({"tracinglostchars"=}460,assign_int,int_base+tracing_lost_chars_code);

 { \xref[tracing_lost_chars_][\.[\\tracinglostchars] primitive] }
primitive({"tracingcommands"=}461,assign_int,int_base+tracing_commands_code);

 { \xref[tracing_commands_][\.[\\tracingcommands] primitive] }
primitive({"tracingrestores"=}462,assign_int,int_base+tracing_restores_code);

 { \xref[tracing_restores_][\.[\\tracingrestores] primitive] }
primitive({"uchyph"=}463,assign_int,int_base+uc_hyph_code);

 { \xref[uc_hyph_][\.[\\uchyph] primitive] }
primitive({"outputpenalty"=}464,assign_int,int_base+output_penalty_code);

 { \xref[output_penalty_][\.[\\outputpenalty] primitive] }
primitive({"maxdeadcycles"=}465,assign_int,int_base+max_dead_cycles_code);

 { \xref[max_dead_cycles_][\.[\\maxdeadcycles] primitive] }
primitive({"hangafter"=}466,assign_int,int_base+hang_after_code);

 { \xref[hang_after_][\.[\\hangafter] primitive] }
primitive({"floatingpenalty"=}467,assign_int,int_base+floating_penalty_code);

 { \xref[floating_penalty_][\.[\\floatingpenalty] primitive] }
primitive({"globaldefs"=}468,assign_int,int_base+global_defs_code);

 { \xref[global_defs_][\.[\\globaldefs] primitive] }
primitive({"fam"=}469,assign_int,int_base+cur_fam_code);

 { \xref[fam_][\.[\\fam] primitive] }
primitive({"escapechar"=}470,assign_int,int_base+escape_char_code);

 { \xref[escape_char_][\.[\\escapechar] primitive] }
primitive({"defaulthyphenchar"=}471,assign_int,int_base+default_hyphen_char_code);

 { \xref[default_hyphen_char_][\.[\\defaulthyphenchar] primitive] }
primitive({"defaultskewchar"=}472,assign_int,int_base+default_skew_char_code);

 { \xref[default_skew_char_][\.[\\defaultskewchar] primitive] }
primitive({"endlinechar"=}473,assign_int,int_base+end_line_char_code);

 { \xref[end_line_char_][\.[\\endlinechar] primitive] }
primitive({"newlinechar"=}474,assign_int,int_base+new_line_char_code);

 { \xref[new_line_char_][\.[\\newlinechar] primitive] }
primitive({"language"=}475,assign_int,int_base+language_code);

 { \xref[language_][\.[\\language] primitive] }
primitive({"lefthyphenmin"=}476,assign_int,int_base+left_hyphen_min_code);

 { \xref[left_hyphen_min_][\.[\\lefthyphenmin] primitive] }
primitive({"righthyphenmin"=}477,assign_int,int_base+right_hyphen_min_code);

 { \xref[right_hyphen_min_][\.[\\righthyphenmin] primitive] }
primitive({"holdinginserts"=}478,assign_int,int_base+holding_inserts_code);

 { \xref[holding_inserts_][\.[\\holdinginserts] primitive] }
primitive({"errorcontextlines"=}479,assign_int,int_base+error_context_lines_code);

 { \xref[error_context_lines_][\.[\\errorcontextlines] primitive] }
if mltex_p then
  begin mltex_enabled_p:=true;  {enable character substitution}
  if false then {remove the if-clause to enable \.[\\charsubdefmin]}
  primitive({"charsubdefmin"=}480,assign_int,int_base+char_sub_def_min_code);

 { \xref[char_sub_def_min_][\.[\\charsubdefmin] primitive] }
  primitive({"charsubdefmax"=}481,assign_int,int_base+char_sub_def_max_code);

 { \xref[char_sub_def_max_][\.[\\charsubdefmax] primitive] }
  primitive({"tracingcharsubdef"=}482,assign_int,int_base+tracing_char_sub_def_code);

 { \xref[tracing_char_sub_def_][\.[\\tracingcharsubdef] primitive] }
  end;


primitive({"parindent"=}486,assign_dimen,dimen_base+par_indent_code);

 { \xref[par_indent_][\.[\\parindent] primitive] }
primitive({"mathsurround"=}487,assign_dimen,dimen_base+math_surround_code);

 { \xref[math_surround_][\.[\\mathsurround] primitive] }
primitive({"lineskiplimit"=}488,assign_dimen,dimen_base+line_skip_limit_code);

 { \xref[line_skip_limit_][\.[\\lineskiplimit] primitive] }
primitive({"hsize"=}489,assign_dimen,dimen_base+hsize_code);

 { \xref[hsize_][\.[\\hsize] primitive] }
primitive({"vsize"=}490,assign_dimen,dimen_base+vsize_code);

 { \xref[vsize_][\.[\\vsize] primitive] }
primitive({"maxdepth"=}491,assign_dimen,dimen_base+max_depth_code);

 { \xref[max_depth_][\.[\\maxdepth] primitive] }
primitive({"splitmaxdepth"=}492,assign_dimen,dimen_base+split_max_depth_code);

 { \xref[split_max_depth_][\.[\\splitmaxdepth] primitive] }
primitive({"boxmaxdepth"=}493,assign_dimen,dimen_base+box_max_depth_code);

 { \xref[box_max_depth_][\.[\\boxmaxdepth] primitive] }
primitive({"hfuzz"=}494,assign_dimen,dimen_base+hfuzz_code);

 { \xref[hfuzz_][\.[\\hfuzz] primitive] }
primitive({"vfuzz"=}495,assign_dimen,dimen_base+vfuzz_code);

 { \xref[vfuzz_][\.[\\vfuzz] primitive] }
primitive({"delimitershortfall"=}496,
  assign_dimen,dimen_base+delimiter_shortfall_code);

 { \xref[delimiter_shortfall_][\.[\\delimitershortfall] primitive] }
primitive({"nulldelimiterspace"=}497,
  assign_dimen,dimen_base+null_delimiter_space_code);

 { \xref[null_delimiter_space_][\.[\\nulldelimiterspace] primitive] }
primitive({"scriptspace"=}498,assign_dimen,dimen_base+script_space_code);

 { \xref[script_space_][\.[\\scriptspace] primitive] }
primitive({"predisplaysize"=}499,assign_dimen,dimen_base+pre_display_size_code);

 { \xref[pre_display_size_][\.[\\predisplaysize] primitive] }
primitive({"displaywidth"=}500,assign_dimen,dimen_base+display_width_code);

 { \xref[display_width_][\.[\\displaywidth] primitive] }
primitive({"displayindent"=}501,assign_dimen,dimen_base+display_indent_code);

 { \xref[display_indent_][\.[\\displayindent] primitive] }
primitive({"overfullrule"=}502,assign_dimen,dimen_base+overfull_rule_code);

 { \xref[overfull_rule_][\.[\\overfullrule] primitive] }
primitive({"hangindent"=}503,assign_dimen,dimen_base+hang_indent_code);

 { \xref[hang_indent_][\.[\\hangindent] primitive] }
primitive({"hoffset"=}504,assign_dimen,dimen_base+h_offset_code);

 { \xref[h_offset_][\.[\\hoffset] primitive] }
primitive({"voffset"=}505,assign_dimen,dimen_base+v_offset_code);

 { \xref[v_offset_][\.[\\voffset] primitive] }
primitive({"emergencystretch"=}506,assign_dimen,dimen_base+emergency_stretch_code);

 { \xref[emergency_stretch_][\.[\\emergencystretch] primitive] }


primitive({" "=}32,ex_space,0);

 { \xref[Single-character primitives /][\quad\.[\\\ ]] }
primitive({"/"=}47,ital_corr,0);

 { \xref[Single-character primitives /][\quad\.[\\/]] }
primitive({"accent"=}516,accent,0);

 { \xref[accent_][\.[\\accent] primitive] }
primitive({"advance"=}517,advance,0);

 { \xref[advance_][\.[\\advance] primitive] }
primitive({"afterassignment"=}518,after_assignment,0);

 { \xref[after_assignment_][\.[\\afterassignment] primitive] }
primitive({"aftergroup"=}519,after_group,0);

 { \xref[after_group_][\.[\\aftergroup] primitive] }
primitive({"begingroup"=}520,begin_group,0);

 { \xref[begin_group_][\.[\\begingroup] primitive] }
primitive({"char"=}521,char_num,0);

 { \xref[char_][\.[\\char] primitive] }
primitive({"csname"=}512,cs_name,0);

 { \xref[cs_name_][\.[\\csname] primitive] }
primitive({"delimiter"=}522,delim_num,0);

 { \xref[delimiter_][\.[\\delimiter] primitive] }
primitive({"divide"=}523,divide,0);

 { \xref[divide_][\.[\\divide] primitive] }
primitive({"endcsname"=}513,end_cs_name,0);

 { \xref[end_cs_name_][\.[\\endcsname] primitive] }
primitive({"endgroup"=}524,end_group,0);
 { \xref[end_group_][\.[\\endgroup] primitive] }
 hash[ frozen_end_group].rh :={"endgroup"=}524; eqtb[frozen_end_group]:=eqtb[cur_val];

primitive({"expandafter"=}525,expand_after,0);

 { \xref[expand_after_][\.[\\expandafter] primitive] }
primitive({"font"=}526,def_font,0);

 { \xref[font_][\.[\\font] primitive] }
primitive({"fontdimen"=}527,assign_font_dimen,0);

 { \xref[font_dimen_][\.[\\fontdimen] primitive] }
primitive({"halign"=}528,halign,0);

 { \xref[halign_][\.[\\halign] primitive] }
primitive({"hrule"=}529,hrule,0);

 { \xref[hrule_][\.[\\hrule] primitive] }
primitive({"ignorespaces"=}530,ignore_spaces,0);

 { \xref[ignore_spaces_][\.[\\ignorespaces] primitive] }
primitive({"insert"=}327,insert,0);

 { \xref[insert_][\.[\\insert] primitive] }
primitive({"mark"=}348,mark,0);

 { \xref[mark_][\.[\\mark] primitive] }
primitive({"mathaccent"=}531,math_accent,0);

 { \xref[math_accent_][\.[\\mathaccent] primitive] }
primitive({"mathchar"=}532,math_char_num,0);

 { \xref[math_char_][\.[\\mathchar] primitive] }
primitive({"mathchoice"=}533,math_choice,0);

 { \xref[math_choice_][\.[\\mathchoice] primitive] }
primitive({"multiply"=}534,multiply,0);

 { \xref[multiply_][\.[\\multiply] primitive] }
primitive({"noalign"=}535,no_align,0);

 { \xref[no_align_][\.[\\noalign] primitive] }
primitive({"noboundary"=}536,no_boundary,0);

 { \xref[no_boundary_][\.[\\noboundary] primitive] }
primitive({"noexpand"=}537,no_expand,0);

 { \xref[no_expand_][\.[\\noexpand] primitive] }
primitive({"nonscript"=}332,non_script,0);

 { \xref[non_script_][\.[\\nonscript] primitive] }
primitive({"omit"=}538,omit,0);

 { \xref[omit_][\.[\\omit] primitive] }
primitive({"parshape"=}413,set_shape,0);

 { \xref[par_shape_][\.[\\parshape] primitive] }
primitive({"penalty"=}539,break_penalty,0);

 { \xref[penalty_][\.[\\penalty] primitive] }
primitive({"prevgraf"=}540,set_prev_graf,0);

 { \xref[prev_graf_][\.[\\prevgraf] primitive] }
primitive({"radical"=}541,radical,0);

 { \xref[radical_][\.[\\radical] primitive] }
primitive({"read"=}542,read_to_cs,0);

 { \xref[read_][\.[\\read] primitive] }
primitive({"relax"=}543,relax,256); {cf.\ |scan_file_name|}
 { \xref[relax_][\.[\\relax] primitive] }
 hash[ frozen_relax].rh :={"relax"=}543; eqtb[frozen_relax]:=eqtb[cur_val];

primitive({"setbox"=}544,set_box,0);

 { \xref[set_box_][\.[\\setbox] primitive] }
primitive({"the"=}545,the,0);

 { \xref[the_][\.[\\the] primitive] }
primitive({"toks"=}412,toks_register,0);

 { \xref[toks_][\.[\\toks] primitive] }
primitive({"vadjust"=}349,vadjust,0);

 { \xref[vadjust_][\.[\\vadjust] primitive] }
primitive({"valign"=}546,valign,0);

 { \xref[valign_][\.[\\valign] primitive] }
primitive({"vcenter"=}547,vcenter,0);

 { \xref[vcenter_][\.[\\vcenter] primitive] }
primitive({"vrule"=}548,vrule,0);

 { \xref[vrule_][\.[\\vrule] primitive] }


primitive({"par"=}604,par_end,256); {cf.\ |scan_file_name|}
 { \xref[par_][\.[\\par] primitive] }
par_loc:=cur_val; par_token:={07777=}4095 +par_loc;


primitive({"input"=}639,input,0);

 { \xref[input_][\.[\\input] primitive] }
primitive({"endinput"=}640,input,1);

 { \xref[end_input_][\.[\\endinput] primitive] }


primitive({"topmark"=}641,top_bot_mark,top_mark_code);
 { \xref[top_mark_][\.[\\topmark] primitive] }
primitive({"firstmark"=}642,top_bot_mark,first_mark_code);
 { \xref[first_mark_][\.[\\firstmark] primitive] }
primitive({"botmark"=}643,top_bot_mark,bot_mark_code);
 { \xref[bot_mark_][\.[\\botmark] primitive] }
primitive({"splitfirstmark"=}644,top_bot_mark,split_first_mark_code);
 { \xref[split_first_mark_][\.[\\splitfirstmark] primitive] }
primitive({"splitbotmark"=}645,top_bot_mark,split_bot_mark_code);
 { \xref[split_bot_mark_][\.[\\splitbotmark] primitive] }


primitive({"count"=}484,register,int_val);
 { \xref[count_][\.[\\count] primitive] }
primitive({"dimen"=}508,register,dimen_val);
 { \xref[dimen_][\.[\\dimen] primitive] }
primitive({"skip"=}400,register,glue_val);
 { \xref[skip_][\.[\\skip] primitive] }
primitive({"muskip"=}401,register,mu_val);
 { \xref[mu_skip_][\.[\\muskip] primitive] }


primitive({"spacefactor"=}678,set_aux,hmode);
 { \xref[space_factor_][\.[\\spacefactor] primitive] }
primitive({"prevdepth"=}679,set_aux,vmode);

 { \xref[prev_depth_][\.[\\prevdepth] primitive] }
primitive({"deadcycles"=}680,set_page_int,0);
 { \xref[dead_cycles_][\.[\\deadcycles] primitive] }
primitive({"insertpenalties"=}681,set_page_int,1);
 { \xref[insert_penalties_][\.[\\insertpenalties] primitive] }
primitive({"wd"=}682,set_box_dimen,width_offset);
 { \xref[wd_][\.[\\wd] primitive] }
primitive({"ht"=}683,set_box_dimen,height_offset);
 { \xref[ht_][\.[\\ht] primitive] }
primitive({"dp"=}684,set_box_dimen,depth_offset);
 { \xref[dp_][\.[\\dp] primitive] }
primitive({"lastpenalty"=}685,last_item,int_val);
 { \xref[last_penalty_][\.[\\lastpenalty] primitive] }
primitive({"lastkern"=}686,last_item,dimen_val);
 { \xref[last_kern_][\.[\\lastkern] primitive] }
primitive({"lastskip"=}687,last_item,glue_val);
 { \xref[last_skip_][\.[\\lastskip] primitive] }
primitive({"inputlineno"=}688,last_item,input_line_no_code);
 { \xref[input_line_no_][\.[\\inputlineno] primitive] }
primitive({"badness"=}689,last_item,badness_code);
 { \xref[badness_][\.[\\badness] primitive] }


primitive({"number"=}745,convert,number_code);

 { \xref[number_][\.[\\number] primitive] }
primitive({"romannumeral"=}746,convert,roman_numeral_code);

 { \xref[roman_numeral_][\.[\\romannumeral] primitive] }
primitive({"string"=}747,convert,string_code);

 { \xref[string_][\.[\\string] primitive] }
primitive({"meaning"=}748,convert,meaning_code);

 { \xref[meaning_][\.[\\meaning] primitive] }
primitive({"fontname"=}749,convert,font_name_code);

 { \xref[font_name_][\.[\\fontname] primitive] }
primitive({"jobname"=}750,convert,job_name_code);

 { \xref[job_name_][\.[\\jobname] primitive] }


primitive({"if"=}767,if_test,if_char_code);
 { \xref[if_char_][\.[\\if] primitive] }
primitive({"ifcat"=}768,if_test,if_cat_code);
 { \xref[if_cat_code_][\.[\\ifcat] primitive] }
primitive({"ifnum"=}769,if_test,if_int_code);
 { \xref[if_int_][\.[\\ifnum] primitive] }
primitive({"ifdim"=}770,if_test,if_dim_code);
 { \xref[if_dim_][\.[\\ifdim] primitive] }
primitive({"ifodd"=}771,if_test,if_odd_code);
 { \xref[if_odd_][\.[\\ifodd] primitive] }
primitive({"ifvmode"=}772,if_test,if_vmode_code);
 { \xref[if_vmode_][\.[\\ifvmode] primitive] }
primitive({"ifhmode"=}773,if_test,if_hmode_code);
 { \xref[if_hmode_][\.[\\ifhmode] primitive] }
primitive({"ifmmode"=}774,if_test,if_mmode_code);
 { \xref[if_mmode_][\.[\\ifmmode] primitive] }
primitive({"ifinner"=}775,if_test,if_inner_code);
 { \xref[if_inner_][\.[\\ifinner] primitive] }
primitive({"ifvoid"=}776,if_test,if_void_code);
 { \xref[if_void_][\.[\\ifvoid] primitive] }
primitive({"ifhbox"=}777,if_test,if_hbox_code);
 { \xref[if_hbox_][\.[\\ifhbox] primitive] }
primitive({"ifvbox"=}778,if_test,if_vbox_code);
 { \xref[if_vbox_][\.[\\ifvbox] primitive] }
primitive({"ifx"=}779,if_test,ifx_code);
 { \xref[ifx_][\.[\\ifx] primitive] }
primitive({"ifeof"=}780,if_test,if_eof_code);
 { \xref[if_eof_][\.[\\ifeof] primitive] }
primitive({"iftrue"=}781,if_test,if_true_code);
 { \xref[if_true_][\.[\\iftrue] primitive] }
primitive({"iffalse"=}782,if_test,if_false_code);
 { \xref[if_false_][\.[\\iffalse] primitive] }
primitive({"ifcase"=}783,if_test,if_case_code);
 { \xref[if_case_][\.[\\ifcase] primitive] }


primitive({"fi"=}784,fi_or_else,fi_code);
 { \xref[fi_][\.[\\fi] primitive] }
 hash[ frozen_fi].rh :={"fi"=}784; eqtb[frozen_fi]:=eqtb[cur_val];
primitive({"or"=}785,fi_or_else,or_code);
 { \xref[or_][\.[\\or] primitive] }
primitive({"else"=}786,fi_or_else,else_code);
 { \xref[else_][\.[\\else] primitive] }


primitive({"nullfont"=}811,set_font,font_base );
 { \xref[null_font_][\.[\\nullfont] primitive] }
 hash[ frozen_null_font].rh :={"nullfont"=}811; eqtb[frozen_null_font]:=eqtb[cur_val];


primitive({"span"=}912,tab_mark,span_code);

 { \xref[span_][\.[\\span] primitive] }
primitive({"cr"=}913,car_ret,cr_code);
 { \xref[cr_][\.[\\cr] primitive] }
 hash[ frozen_cr].rh :={"cr"=}913; eqtb[frozen_cr]:=eqtb[cur_val];

primitive({"crcr"=}914,car_ret,cr_cr_code);
 { \xref[cr_cr_][\.[\\crcr] primitive] }
 hash[ frozen_end_template].rh :={"endtemplate"=}915;  hash[ frozen_endv].rh :={"endtemplate"=}915;
{ \xref[endtemplate] }
 eqtb[  frozen_endv].hh.b0  :=endv;  eqtb[  frozen_endv].hh.rh  :=mem_top-11 ;
 eqtb[  frozen_endv].hh.b1  :=level_one;

eqtb[frozen_end_template]:=eqtb[frozen_endv];
 eqtb[  frozen_end_template].hh.b0  :=end_template;


primitive({"pagegoal"=}984,set_page_dimen,0);
 { \xref[page_goal_][\.[\\pagegoal] primitive] }
primitive({"pagetotal"=}985,set_page_dimen,1);
 { \xref[page_total_][\.[\\pagetotal] primitive] }
primitive({"pagestretch"=}986,set_page_dimen,2);
 { \xref[page_stretch_][\.[\\pagestretch] primitive] }
primitive({"pagefilstretch"=}987,set_page_dimen,3);
 { \xref[page_fil_stretch_][\.[\\pagefilstretch] primitive] }
primitive({"pagefillstretch"=}988,set_page_dimen,4);
 { \xref[page_fill_stretch_][\.[\\pagefillstretch] primitive] }
primitive({"pagefilllstretch"=}989,set_page_dimen,5);
 { \xref[page_filll_stretch_][\.[\\pagefilllstretch] primitive] }
primitive({"pageshrink"=}990,set_page_dimen,6);
 { \xref[page_shrink_][\.[\\pageshrink] primitive] }
primitive({"pagedepth"=}991,set_page_dimen,7);
 { \xref[page_depth_][\.[\\pagedepth] primitive] }


primitive({"end"=}1038,stop,0);

 { \xref[end_][\.[\\end] primitive] }
primitive({"dump"=}1039,stop,1);

 { \xref[dump_][\.[\\dump] primitive] }


primitive({"hskip"=}1040,hskip,skip_code);

 { \xref[hskip_][\.[\\hskip] primitive] }
primitive({"hfil"=}1041,hskip,fil_code);
 { \xref[hfil_][\.[\\hfil] primitive] }
primitive({"hfill"=}1042,hskip,fill_code);

 { \xref[hfill_][\.[\\hfill] primitive] }
primitive({"hss"=}1043,hskip,ss_code);
 { \xref[hss_][\.[\\hss] primitive] }
primitive({"hfilneg"=}1044,hskip,fil_neg_code);

 { \xref[hfil_neg_][\.[\\hfilneg] primitive] }
primitive({"vskip"=}1045,vskip,skip_code);

 { \xref[vskip_][\.[\\vskip] primitive] }
primitive({"vfil"=}1046,vskip,fil_code);
 { \xref[vfil_][\.[\\vfil] primitive] }
primitive({"vfill"=}1047,vskip,fill_code);

 { \xref[vfill_][\.[\\vfill] primitive] }
primitive({"vss"=}1048,vskip,ss_code);
 { \xref[vss_][\.[\\vss] primitive] }
primitive({"vfilneg"=}1049,vskip,fil_neg_code);

 { \xref[vfil_neg_][\.[\\vfilneg] primitive] }
primitive({"mskip"=}333,mskip,mskip_code);

 { \xref[mskip_][\.[\\mskip] primitive] }
primitive({"kern"=}337,kern,explicit);
 { \xref[kern_][\.[\\kern] primitive] }
primitive({"mkern"=}339,mkern,mu_glue);

 { \xref[mkern_][\.[\\mkern] primitive] }


primitive({"moveleft"=}1067,hmove,1);
 { \xref[move_left_][\.[\\moveleft] primitive] }
primitive({"moveright"=}1068,hmove,0);

 { \xref[move_right_][\.[\\moveright] primitive] }
primitive({"raise"=}1069,vmove,1);
 { \xref[raise_][\.[\\raise] primitive] }
primitive({"lower"=}1070,vmove,0);
 { \xref[lower_][\.[\\lower] primitive] }


primitive({"box"=}414,make_box,box_code);
 { \xref[box_][\.[\\box] primitive] }
primitive({"copy"=}1071,make_box,copy_code);
 { \xref[copy_][\.[\\copy] primitive] }
primitive({"lastbox"=}1072,make_box,last_box_code);
 { \xref[last_box_][\.[\\lastbox] primitive] }
primitive({"vsplit"=}979,make_box,vsplit_code);
 { \xref[vsplit_][\.[\\vsplit] primitive] }
primitive({"vtop"=}1073,make_box,vtop_code);

 { \xref[vtop_][\.[\\vtop] primitive] }
primitive({"vbox"=}981,make_box,vtop_code+vmode);
 { \xref[vbox_][\.[\\vbox] primitive] }
primitive({"hbox"=}1074,make_box,vtop_code+hmode);

 { \xref[hbox_][\.[\\hbox] primitive] }
primitive({"shipout"=}1075,leader_ship,a_leaders-1); {|ship_out_flag=leader_flag-1|}
 { \xref[ship_out_][\.[\\shipout] primitive] }
primitive({"leaders"=}1076,leader_ship,a_leaders);
 { \xref[leaders_][\.[\\leaders] primitive] }
primitive({"cleaders"=}1077,leader_ship,c_leaders);
 { \xref[c_leaders_][\.[\\cleaders] primitive] }
primitive({"xleaders"=}1078,leader_ship,x_leaders);
 { \xref[x_leaders_][\.[\\xleaders] primitive] }


primitive({"indent"=}1093,start_par,1);
 { \xref[indent_][\.[\\indent] primitive] }
primitive({"noindent"=}1094,start_par,0);
 { \xref[no_indent_][\.[\\noindent] primitive] }


primitive({"unpenalty"=}1103,remove_item,penalty_node);

 { \xref[un_penalty_][\.[\\unpenalty] primitive] }
primitive({"unkern"=}1104,remove_item,kern_node);

 { \xref[un_kern_][\.[\\unkern] primitive] }
primitive({"unskip"=}1105,remove_item,glue_node);

 { \xref[un_skip_][\.[\\unskip] primitive] }
primitive({"unhbox"=}1106,un_hbox,box_code);

 { \xref[un_hbox_][\.[\\unhbox] primitive] }
primitive({"unhcopy"=}1107,un_hbox,copy_code);

 { \xref[un_hcopy_][\.[\\unhcopy] primitive] }
primitive({"unvbox"=}1108,un_vbox,box_code);

 { \xref[un_vbox_][\.[\\unvbox] primitive] }
primitive({"unvcopy"=}1109,un_vbox,copy_code);

 { \xref[un_vcopy_][\.[\\unvcopy] primitive] }


primitive({"-"=}45,discretionary,1);
 { \xref[Single-character primitives -][\quad\.[\\-]] }
primitive({"discretionary"=}346,discretionary,0);
 { \xref[discretionary_][\.[\\discretionary] primitive] }


primitive({"eqno"=}1140,eq_no,0);
 { \xref[eq_no_][\.[\\eqno] primitive] }
primitive({"leqno"=}1141,eq_no,1);
 { \xref[leq_no_][\.[\\leqno] primitive] }


primitive({"mathord"=}880,math_comp,ord_noad);
 { \xref[math_ord_][\.[\\mathord] primitive] }
primitive({"mathop"=}881,math_comp,op_noad);
 { \xref[math_op_][\.[\\mathop] primitive] }
primitive({"mathbin"=}882,math_comp,bin_noad);
 { \xref[math_bin_][\.[\\mathbin] primitive] }
primitive({"mathrel"=}883,math_comp,rel_noad);
 { \xref[math_rel_][\.[\\mathrel] primitive] }
primitive({"mathopen"=}884,math_comp,open_noad);
 { \xref[math_open_][\.[\\mathopen] primitive] }
primitive({"mathclose"=}885,math_comp,close_noad);
 { \xref[math_close_][\.[\\mathclose] primitive] }
primitive({"mathpunct"=}886,math_comp,punct_noad);
 { \xref[math_punct_][\.[\\mathpunct] primitive] }
primitive({"mathinner"=}887,math_comp,inner_noad);
 { \xref[math_inner_][\.[\\mathinner] primitive] }
primitive({"underline"=}889,math_comp,under_noad);
 { \xref[underline_][\.[\\underline] primitive] }
primitive({"overline"=}888,math_comp,over_noad);

 { \xref[overline_][\.[\\overline] primitive] }
primitive({"displaylimits"=}1142,limit_switch,normal);
 { \xref[display_limits_][\.[\\displaylimits] primitive] }
primitive({"limits"=}892,limit_switch,limits);
 { \xref[limits_][\.[\\limits] primitive] }
primitive({"nolimits"=}893,limit_switch,no_limits);
 { \xref[no_limits_][\.[\\nolimits] primitive] }


primitive({"displaystyle"=}875,math_style,display_style);
 { \xref[display_style_][\.[\\displaystyle] primitive] }
primitive({"textstyle"=}876,math_style,text_style);
 { \xref[text_style_][\.[\\textstyle] primitive] }
primitive({"scriptstyle"=}877,math_style,script_style);
 { \xref[script_style_][\.[\\scriptstyle] primitive] }
primitive({"scriptscriptstyle"=}878,math_style,script_script_style);
 { \xref[script_script_style_][\.[\\scriptscriptstyle] primitive] }


primitive({"above"=}1160,above,above_code);

 { \xref[above_][\.[\\above] primitive] }
primitive({"over"=}1161,above,over_code);

 { \xref[over_][\.[\\over] primitive] }
primitive({"atop"=}1162,above,atop_code);

 { \xref[atop_][\.[\\atop] primitive] }
primitive({"abovewithdelims"=}1163,above,delimited_code+above_code);

 { \xref[above_with_delims_][\.[\\abovewithdelims] primitive] }
primitive({"overwithdelims"=}1164,above,delimited_code+over_code);

 { \xref[over_with_delims_][\.[\\overwithdelims] primitive] }
primitive({"atopwithdelims"=}1165,above,delimited_code+atop_code);
 { \xref[atop_with_delims_][\.[\\atopwithdelims] primitive] }


primitive({"left"=}890,left_right,left_noad);
 { \xref[left_][\.[\\left] primitive] }
primitive({"right"=}891,left_right,right_noad);
 { \xref[right_][\.[\\right] primitive] }
 hash[ frozen_right].rh :={"right"=}891; eqtb[frozen_right]:=eqtb[cur_val];


primitive({"long"=}1184,prefix,1);
 { \xref[long_][\.[\\long] primitive] }
primitive({"outer"=}1185,prefix,2);
 { \xref[outer_][\.[\\outer] primitive] }
primitive({"global"=}1186,prefix,4);
 { \xref[global_][\.[\\global] primitive] }
primitive({"def"=}1187,def,0);
 { \xref[def_][\.[\\def] primitive] }
primitive({"gdef"=}1188,def,1);
 { \xref[gdef_][\.[\\gdef] primitive] }
primitive({"edef"=}1189,def,2);
 { \xref[edef_][\.[\\edef] primitive] }
primitive({"xdef"=}1190,def,3);
 { \xref[xdef_][\.[\\xdef] primitive] }


primitive({"let"=}1204,let,normal);

 { \xref[let_][\.[\\let] primitive] }
primitive({"futurelet"=}1205,let,normal+1);

 { \xref[future_let_][\.[\\futurelet] primitive] }


primitive({"chardef"=}1206,shorthand_def,char_def_code);

 { \xref[char_def_][\.[\\chardef] primitive] }
primitive({"mathchardef"=}1207,shorthand_def,math_char_def_code);

 { \xref[math_char_def_][\.[\\mathchardef] primitive] }
primitive({"countdef"=}1208,shorthand_def,count_def_code);

 { \xref[count_def_][\.[\\countdef] primitive] }
primitive({"dimendef"=}1209,shorthand_def,dimen_def_code);

 { \xref[dimen_def_][\.[\\dimendef] primitive] }
primitive({"skipdef"=}1210,shorthand_def,skip_def_code);

 { \xref[skip_def_][\.[\\skipdef] primitive] }
primitive({"muskipdef"=}1211,shorthand_def,mu_skip_def_code);

 { \xref[mu_skip_def_][\.[\\muskipdef] primitive] }
primitive({"toksdef"=}1212,shorthand_def,toks_def_code);

 { \xref[toks_def_][\.[\\toksdef] primitive] }
if mltex_p then
  begin
  primitive({"charsubdef"=}1213,shorthand_def,char_sub_def_code);

 { \xref[char_sub_def_][\.[\\charsubdef] primitive] }
  end;


primitive({"catcode"=}420,def_code,cat_code_base);
 { \xref[cat_code_][\.[\\catcode] primitive] }
primitive({"mathcode"=}424,def_code,math_code_base);
 { \xref[math_code_][\.[\\mathcode] primitive] }
primitive({"lccode"=}421,def_code,lc_code_base);
 { \xref[lc_code_][\.[\\lccode] primitive] }
primitive({"uccode"=}422,def_code,uc_code_base);
 { \xref[uc_code_][\.[\\uccode] primitive] }
primitive({"sfcode"=}423,def_code,sf_code_base);
 { \xref[sf_code_][\.[\\sfcode] primitive] }
primitive({"delcode"=}485,def_code,del_code_base);
 { \xref[del_code_][\.[\\delcode] primitive] }
primitive({"textfont"=}417,def_family,math_font_base);
 { \xref[text_font_][\.[\\textfont] primitive] }
primitive({"scriptfont"=}418,def_family,math_font_base+script_size);
 { \xref[script_font_][\.[\\scriptfont] primitive] }
primitive({"scriptscriptfont"=}419,def_family,math_font_base+script_script_size);
 { \xref[script_script_font_][\.[\\scriptscriptfont] primitive] }


primitive({"hyphenation"=}955,hyph_data,0);
 { \xref[hyphenation_][\.[\\hyphenation] primitive] }
primitive({"patterns"=}967,hyph_data,1);
 { \xref[patterns_][\.[\\patterns] primitive] }


primitive({"hyphenchar"=}1233,assign_font_int,0);
 { \xref[hyphen_char_][\.[\\hyphenchar] primitive] }
primitive({"skewchar"=}1234,assign_font_int,1);
 { \xref[skew_char_][\.[\\skewchar] primitive] }


primitive({"batchmode"=}272,set_interaction,batch_mode);
 { \xref[batch_mode_][\.[\\batchmode] primitive] }
primitive({"nonstopmode"=}273,set_interaction,nonstop_mode);
 { \xref[nonstop_mode_][\.[\\nonstopmode] primitive] }
primitive({"scrollmode"=}274,set_interaction,scroll_mode);
 { \xref[scroll_mode_][\.[\\scrollmode] primitive] }
primitive({"errorstopmode"=}1243,set_interaction,error_stop_mode);
 { \xref[error_stop_mode_][\.[\\errorstopmode] primitive] }


primitive({"openin"=}1244,in_stream,1);
 { \xref[open_in_][\.[\\openin] primitive] }
primitive({"closein"=}1245,in_stream,0);
 { \xref[close_in_][\.[\\closein] primitive] }


primitive({"message"=}1246,message,0);
 { \xref[message_][\.[\\message] primitive] }
primitive({"errmessage"=}1247,message,1);
 { \xref[err_message_][\.[\\errmessage] primitive] }


primitive({"lowercase"=}1253,case_shift,lc_code_base);
 { \xref[lowercase_][\.[\\lowercase] primitive] }
primitive({"uppercase"=}1254,case_shift,uc_code_base);
 { \xref[uppercase_][\.[\\uppercase] primitive] }


primitive({"show"=}1255,xray,show_code);
 { \xref[show_][\.[\\show] primitive] }
primitive({"showbox"=}1256,xray,show_box_code);
 { \xref[show_box_][\.[\\showbox] primitive] }
primitive({"showthe"=}1257,xray,show_the_code);
 { \xref[show_the_][\.[\\showthe] primitive] }
primitive({"showlists"=}1258,xray,show_lists_code);
 { \xref[show_lists_code_][\.[\\showlists] primitive] }


primitive({"openout"=}1304,extension,open_node);

 { \xref[open_out_][\.[\\openout] primitive] }
primitive({"write"=}601,extension,write_node); write_loc:=cur_val;

 { \xref[write_][\.[\\write] primitive] }
primitive({"closeout"=}1305,extension,close_node);

 { \xref[close_out_][\.[\\closeout] primitive] }
primitive({"special"=}1306,extension,special_node);

 hash[ frozen_special].rh :={"special"=}1306; eqtb[frozen_special]:=eqtb[cur_val];

 { \xref[special_][\.[\\special] primitive] }
primitive({"immediate"=}1307,extension,immediate_code);

 { \xref[immediate_][\.[\\immediate] primitive] }
primitive({"setlanguage"=}1308,extension,set_language_code);

 { \xref[set_language_][\.[\\setlanguage] primitive] }

;
no_new_control_sequence:=true;
end;
endif('INITEX') 


 ifdef('TEXMF_DEBUG')  procedure debug_help; {routine to display various things}
label breakpoint,exit;
var k, l, m, n:integer;
begin    ;
   while true do  begin    ;
  print_nl({"debug # (-1 to exit):"=}1303);  fflush (stdout ) ;
{ \xref[debug \#] }
  read(stdin ,m);
  if m<0 then  goto exit 
  else if m=0 then
    dump_core {do something to cause a core dump}
  else  begin read(stdin ,n);
    case m of
    { \4 }
{ Numbered cases for |debug_help| }
1: print_word(mem[n]); {display |mem[n]| in all forms}
2: print_int( mem[ n].hh.lh );
3: print_int( mem[ n].hh.rh );
4: print_word(eqtb[n]);
5: begin print_scaled(font_info[n].int ); print_char({" "=}32);

  print_int(font_info[n].qqqq.b0); print_char({":"=}58);

  print_int(font_info[n].qqqq.b1); print_char({":"=}58);

  print_int(font_info[n].qqqq.b2); print_char({":"=}58);

  print_int(font_info[n].qqqq.b3);
  end;
6: print_word(save_stack[n]);
7: show_box(n);
  {show a box, abbreviated by |show_box_depth| and |show_box_breadth|}
8: begin breadth_max:=10000; depth_threshold:=pool_size-pool_ptr-10;
  show_node_list(n); {show a box in its entirety}
  end;
9: show_token_list(n,-{0xfffffff=}268435455  ,1000);
10: slow_print(n);
11: check_mem(n>0); {check wellformedness; print new busy locations if |n>0|}
12: search_mem(n); {look for pointers to |n|}
13: begin read(stdin ,l); print_cmd_chr(n,l);
  end;
14: for k:=0 to n do print(buffer[k]);
15: begin font_in_short_display:=font_base ; short_display(n);
  end;
16: panicking:=not panicking;

 
     else  print({"?"=}63)
     end ;
    end;
  end;
exit:end;
endif('TEXMF_DEBUG') 





{ 1332. }

{tangle:pos tex.web:24275:61: }

{ Now this is really it: \TeX\ starts and ends here.

The initial test involving |ready_already| should be deleted if the
\PASCAL\ runtime system is smart enough to detect such a ``mistake.''
\xref[system dependencies] } procedure main_body;
begin  {|start_here|}

{Bounds that may be set from the configuration file. We want the user to
 be able to specify the names with underscores, but \.[TANGLE] removes
 underscores, so we're stuck giving the names twice, once as a string,
 once as the identifier. How ugly.}
  bound_default:= 0; bound_name:='mem_bot';  setup_bound_variable(addressof( mem_bot), bound_name, bound_default) ;
  bound_default:= 250000; bound_name:='main_memory';  setup_bound_variable(addressof( main_memory), bound_name, bound_default) ;
    {|memory_word|s for |mem| in \.[INITEX]}
  bound_default:= 0; bound_name:='extra_mem_top';  setup_bound_variable(addressof( extra_mem_top), bound_name, bound_default) ;
    {increase high mem in \.[VIRTEX]}
  bound_default:= 0; bound_name:='extra_mem_bot';  setup_bound_variable(addressof( extra_mem_bot), bound_name, bound_default) ;
    {increase low mem in \.[VIRTEX]}
  bound_default:= 200000; bound_name:='pool_size';  setup_bound_variable(addressof( pool_size), bound_name, bound_default) ;
  bound_default:= 75000; bound_name:='string_vacancies';  setup_bound_variable(addressof( string_vacancies), bound_name, bound_default) ;
  bound_default:= 5000; bound_name:='pool_free';  setup_bound_variable(addressof( pool_free), bound_name, bound_default) ; {min pool avail after fmt}
  bound_default:= 15000; bound_name:='max_strings';  setup_bound_variable(addressof( max_strings), bound_name, bound_default) ;
  bound_default:= 100; bound_name:='strings_free';  setup_bound_variable(addressof( strings_free), bound_name, bound_default) ;
  bound_default:= 100000; bound_name:='font_mem_size';  setup_bound_variable(addressof( font_mem_size), bound_name, bound_default) ;
  bound_default:= 500; bound_name:='font_max';  setup_bound_variable(addressof( font_max), bound_name, bound_default) ;
  bound_default:= 20000; bound_name:='trie_size';  setup_bound_variable(addressof( trie_size), bound_name, bound_default) ;
    {if |ssup_trie_size| increases, recompile}
  bound_default:= 659; bound_name:='hyph_size';  setup_bound_variable(addressof( hyph_size), bound_name, bound_default) ;
  bound_default:= 3000; bound_name:='buf_size';  setup_bound_variable(addressof( buf_size), bound_name, bound_default) ;
  bound_default:= 50; bound_name:='nest_size';  setup_bound_variable(addressof( nest_size), bound_name, bound_default) ;
  bound_default:= 15; bound_name:='max_in_open';  setup_bound_variable(addressof( max_in_open), bound_name, bound_default) ;
  bound_default:= 60; bound_name:='param_size';  setup_bound_variable(addressof( param_size), bound_name, bound_default) ;
  bound_default:= 4000; bound_name:='save_size';  setup_bound_variable(addressof( save_size), bound_name, bound_default) ;
  bound_default:= 300; bound_name:='stack_size';  setup_bound_variable(addressof( stack_size), bound_name, bound_default) ;
  bound_default:= 16384; bound_name:='dvi_buf_size';  setup_bound_variable(addressof( dvi_buf_size), bound_name, bound_default) ;
  bound_default:= 79; bound_name:='error_line';  setup_bound_variable(addressof( error_line), bound_name, bound_default) ;
  bound_default:= 50; bound_name:='half_error_line';  setup_bound_variable(addressof( half_error_line), bound_name, bound_default) ;
  bound_default:= 79; bound_name:='max_print_line';  setup_bound_variable(addressof( max_print_line), bound_name, bound_default) ;
  bound_default:= 0; bound_name:='hash_extra';  setup_bound_variable(addressof( hash_extra), bound_name, bound_default) ;
  bound_default:= 10000; bound_name:='expand_depth';  setup_bound_variable(addressof( expand_depth), bound_name, bound_default) ;

  begin if  mem_bot < infmem_bot then  mem_bot := infmem_bot else if  mem_bot > supmem_bot then  mem_bot := supmem_bot end ;
  begin if  main_memory < infmain_memory then  main_memory := infmain_memory else if  main_memory > supmain_memory then  main_memory := supmain_memory end ;
 ifdef('INITEX')  if ini_version then begin 
  extra_mem_top := 0;
  extra_mem_bot := 0;
 end; endif('INITEX')  
  if extra_mem_bot>sup_main_memory then extra_mem_bot:=sup_main_memory;
  if extra_mem_top>sup_main_memory then extra_mem_top:=sup_main_memory;
  {|mem_top| is an index, |main_memory| a size}
  mem_top := mem_bot + main_memory -1;
  mem_min := mem_bot;
  mem_max := mem_top;

  {Check other constants against their sup and inf.}
  begin if  trie_size < inftrie_size then  trie_size := inftrie_size else if  trie_size > suptrie_size then  trie_size := suptrie_size end ;
  begin if  hyph_size < infhyph_size then  hyph_size := infhyph_size else if  hyph_size > suphyph_size then  hyph_size := suphyph_size end ;
  begin if  buf_size < infbuf_size then  buf_size := infbuf_size else if  buf_size > supbuf_size then  buf_size := supbuf_size end ;
  begin if  nest_size < infnest_size then  nest_size := infnest_size else if  nest_size > supnest_size then  nest_size := supnest_size end ;
  begin if  max_in_open < infmax_in_open then  max_in_open := infmax_in_open else if  max_in_open > supmax_in_open then  max_in_open := supmax_in_open end ;
  begin if  param_size < infparam_size then  param_size := infparam_size else if  param_size > supparam_size then  param_size := supparam_size end ;
  begin if  save_size < infsave_size then  save_size := infsave_size else if  save_size > supsave_size then  save_size := supsave_size end ;
  begin if  stack_size < infstack_size then  stack_size := infstack_size else if  stack_size > supstack_size then  stack_size := supstack_size end ;
  begin if  dvi_buf_size < infdvi_buf_size then  dvi_buf_size := infdvi_buf_size else if  dvi_buf_size > supdvi_buf_size then  dvi_buf_size := supdvi_buf_size end ;
  begin if  pool_size < infpool_size then  pool_size := infpool_size else if  pool_size > suppool_size then  pool_size := suppool_size end ;
  begin if  string_vacancies < infstring_vacancies then  string_vacancies := infstring_vacancies else if  string_vacancies > supstring_vacancies then  string_vacancies := supstring_vacancies end ;
  begin if  pool_free < infpool_free then  pool_free := infpool_free else if  pool_free > suppool_free then  pool_free := suppool_free end ;
  begin if  max_strings < infmax_strings then  max_strings := infmax_strings else if  max_strings > supmax_strings then  max_strings := supmax_strings end ;
  begin if  strings_free < infstrings_free then  strings_free := infstrings_free else if  strings_free > supstrings_free then  strings_free := supstrings_free end ;
  begin if  font_mem_size < inffont_mem_size then  font_mem_size := inffont_mem_size else if  font_mem_size > supfont_mem_size then  font_mem_size := supfont_mem_size end ;
  begin if  font_max < inffont_max then  font_max := inffont_max else if  font_max > supfont_max then  font_max := supfont_max end ;
  begin if  hash_extra < infhash_extra then  hash_extra := infhash_extra else if  hash_extra > suphash_extra then  hash_extra := suphash_extra end ;
  if error_line > ssup_error_line then error_line := ssup_error_line;

  {array memory allocation}
  buffer:=xmalloc_array (ASCII_code, buf_size);
  nest:=xmalloc_array (list_state_record, nest_size);
  save_stack:=xmalloc_array (memory_word, save_size);
  input_stack:=xmalloc_array (in_state_record, stack_size);
  input_file:=xmalloc_array (alpha_file, max_in_open);
  line_stack:=xmalloc_array (integer, max_in_open);
  source_filename_stack:=xmalloc_array (str_number, max_in_open);
  full_source_filename_stack:=xmalloc_array (str_number, max_in_open);
  param_stack:=xmalloc_array (halfword, param_size);
  dvi_buf:=xmalloc_array (eight_bits, dvi_buf_size);
  hyph_word :=xmalloc_array (str_number, hyph_size);
  hyph_list :=xmalloc_array (halfword, hyph_size);
  hyph_link :=xmalloc_array (hyph_pointer, hyph_size);
 ifdef('INITEX')  if ini_version then begin 
  yzmem:=xmalloc_array (memory_word, mem_top - mem_bot + 1);
  zmem := yzmem - mem_bot;   {Some compilers require |mem_bot=0|}
  eqtb_top := eqtb_size+hash_extra;
  if hash_extra=0 then hash_top:=undefined_control_sequence else
        hash_top:=eqtb_top;
  yhash:=xmalloc_array (two_halves,1+hash_top-hash_offset);
  hash:=yhash - hash_offset;   {Some compilers require |hash_offset=0|}
   hash[ hash_base].lh :=0;  hash[ hash_base].rh :=0;
  for hash_used:=hash_base+1 to hash_top do hash[hash_used]:=hash[hash_base];
  zeqtb:=xmalloc_array (memory_word, eqtb_top);
  eqtb:=zeqtb;

  str_start:=xmalloc_array (pool_pointer, max_strings);
  str_pool:=xmalloc_array (packed_ASCII_code, pool_size);
  font_info:=xmalloc_array (fmemory_word, font_mem_size);
 end; endif('INITEX')  
history:=fatal_error_stop; {in case we quit during initialization}
 ; {open the terminal for output}
if ready_already=314159 then goto start_of_TEX;

{ Check the ``constant'' values... }
bad:=0;
if (half_error_line<30)or(half_error_line>error_line-15) then bad:=1;
if max_print_line<60 then bad:=2;
if dvi_buf_size mod 8<>0 then bad:=3;
if mem_bot+1100>mem_top then bad:=4;
if hash_prime>hash_size then bad:=5;
if max_in_open>=128 then bad:=6;
if mem_top<256+11 then bad:=7; {we will want |null_list>255|}

ifdef('INITEX')  if (mem_min<>mem_bot)or(mem_max<>mem_top) then bad:=10; endif('INITEX')  

if (mem_min>mem_bot)or(mem_max<mem_top) then bad:=10;
if (min_quarterword>0)or(max_quarterword<127) then bad:=11;
if (-{0xfffffff=}268435455 >0)or({0xfffffff=}268435455 <32767) then bad:=12;
if (min_quarterword<-{0xfffffff=}268435455 )or 
  (max_quarterword>{0xfffffff=}268435455 ) then bad:=13;
if (mem_bot-sup_main_memory<-{0xfffffff=}268435455 )or 
  (mem_top+sup_main_memory>={0xfffffff=}268435455 ) then bad:=14;
if (max_font_max<-{0xfffffff=}268435455 )or(max_font_max>{0xfffffff=}268435455 ) then bad:=15;
if font_max>font_base+max_font_max then bad:=16;
if (save_size>{0xfffffff=}268435455 )or(max_strings>{0xfffffff=}268435455 ) then bad:=17;
if buf_size>{0xfffffff=}268435455  then bad:=18;
if max_quarterword-min_quarterword<255 then bad:=19;

if {07777=}4095 +eqtb_size+hash_extra>{0xfffffff=}268435455  then bad:=21;
if (hash_offset<0)or(hash_offset>hash_base) then bad:=42;

if format_default_length> maxint  then bad:=31;

if 2*{0xfffffff=}268435455 <mem_top-mem_min then bad:=41;

 
if bad>0 then
  begin writeln( stdout ,'Ouch---my internal constants have been clobbered!',
    '---case ',  bad:  1)  ;
{ \xref[Ouch...clobbered] }
  goto final_end;
  end;
initialize; {set global variables to their starting values}
 ifdef('INITEX')  if ini_version then begin  if not get_strings_started then goto final_end;
init_prim; {call |primitive| for each primitive}
init_str_ptr:=str_ptr; init_pool_ptr:=pool_ptr; fix_date_and_time;
end; endif('INITEX')  

ready_already:=314159;
start_of_TEX: 
{ Initialize the output routines }
selector:=term_only; tally:=0; term_offset:=0; file_offset:=0;

if src_specials_p or file_line_error_style_p or parse_first_line_p then
  write(stdout , 'This is GoTeXk, Version 3.141592653 (gotex v0.0-prerelease)'  ) 
else
  write(stdout , 'This is GoTeX, Version 3.141592653 (gotex v0.0-prerelease)'  ) ;
write(stdout , version_string) ;
if format_ident=0 then writeln( stdout ,' (preloaded format=',  dump_name,')')  
else  begin slow_print(format_ident); print_ln;
  end;
if shellenabledp then begin
  write(stdout ,' ') ;
  if restrictedshell then begin
    write(stdout ,'restricted ') ;
  end;
  writeln( stdout ,'\write18 enabled.')  ;
end;
if src_specials_p then begin
  writeln( stdout ,' Source specials enabled.')  
end;
if translate_filename then begin
  write(stdout ,' (') ;
  fputs(translate_filename, stdout);
  writeln( stdout ,')')  ;
end;
 fflush (stdout ) ;

job_name:=0; name_in_progress:=false; log_opened:=false;

output_file_name:=0;

;

{ Get the first line of input and prepare to start }
begin 
{ Initialize the input routines }
begin input_ptr:=0; max_in_stack:=0;
source_filename_stack[0]:=0;full_source_filename_stack[0]:=0;
in_open:=0; open_parens:=0; max_buf_stack:=0;
param_ptr:=0; max_param_stack:=0;
first:=buf_size; repeat buffer[first]:=0; decr(first); until first=0;
buffer[0]:=0;
scanner_status:=normal; warning_index:=-{0xfffffff=}268435455  ; first:=1;
cur_input.state_field :=new_line; cur_input.start_field :=1; cur_input.index_field :=0; line:=0; cur_input.name_field :=0;
force_eof:=false;
align_state:=1000000;

if not init_terminal then goto final_end;
cur_input.limit_field :=last; first:=last+1; {|init_terminal| has set |loc| and |last|}
end

;
if (format_ident=0)or(buffer[cur_input.loc_field ]={"&"=}38)or dump_line then
  begin if format_ident<>0 then initialize; {erase preloaded format}
  if not open_fmt_file then goto final_end;
  if not load_fmt_file then
    begin w_close(fmt_file);
  eqtb:=zeqtb; goto final_end;
    end;
  w_close(fmt_file);
  while (cur_input.loc_field <cur_input.limit_field )and(buffer[cur_input.loc_field ]={" "=}32) do incr(cur_input.loc_field );
  end;
if  (eqtb[int_base+ end_line_char_code].int  <0)or(eqtb[int_base+ end_line_char_code].int  >255)  then decr(cur_input.limit_field )
else  buffer[cur_input.limit_field ]:=eqtb[int_base+ end_line_char_code].int  ;
if mltex_enabled_p then
  begin writeln( stdout ,'MLTeX v2.2 enabled')  ;
  end;
fix_date_and_time;

 ifdef('INITEX') 
if trie_not_ready then begin {initex without format loaded}
  trie_trl:=xmalloc_array (trie_pointer, trie_size);
  trie_tro:=xmalloc_array (trie_pointer, trie_size);
  trie_trc:=xmalloc_array (quarterword, trie_size);

  trie_c:=xmalloc_array (packed_ASCII_code, trie_size);
  trie_o:=xmalloc_array (trie_opcode, trie_size);
  trie_l:=xmalloc_array (trie_pointer, trie_size);
  trie_r:=xmalloc_array (trie_pointer, trie_size);
  trie_hash:=xmalloc_array (trie_pointer, trie_size);
  trie_taken:=xmalloc_array (boolean, trie_size);

  trie_l[0] :=0; trie_c[0]:=  0 ; trie_ptr:=0;

  {Allocate and initialize font arrays}
  font_check:=xmalloc_array(four_quarters, font_max);
  font_size:=xmalloc_array(scaled, font_max);
  font_dsize:=xmalloc_array(scaled, font_max);
  font_params:=xmalloc_array(font_index, font_max);
  font_name:=xmalloc_array(str_number, font_max);
  font_area:=xmalloc_array(str_number, font_max);
  font_bc:=xmalloc_array(eight_bits, font_max);
  font_ec:=xmalloc_array(eight_bits, font_max);
  font_glue:=xmalloc_array(halfword, font_max);
  hyphen_char:=xmalloc_array(integer, font_max);
  skew_char:=xmalloc_array(integer, font_max);
  bchar_label:=xmalloc_array(font_index, font_max);
  font_bchar:=xmalloc_array(nine_bits, font_max);
  font_false_bchar:=xmalloc_array(nine_bits, font_max);
  char_base:=xmalloc_array(integer, font_max);
  width_base:=xmalloc_array(integer, font_max);
  height_base:=xmalloc_array(integer, font_max);
  depth_base:=xmalloc_array(integer, font_max);
  italic_base:=xmalloc_array(integer, font_max);
  lig_kern_base:=xmalloc_array(integer, font_max);
  kern_base:=xmalloc_array(integer, font_max);
  exten_base:=xmalloc_array(integer, font_max);
  param_base:=xmalloc_array(integer, font_max);

  font_ptr:=font_base ; fmem_ptr:=7;
  font_name[font_base ]:={"nullfont"=}811; font_area[font_base ]:={""=}335;
  hyphen_char[font_base ]:={"-"=}45; skew_char[font_base ]:=-1;
  bchar_label[font_base ]:=non_address;
  font_bchar[font_base ]:= 256  ; font_false_bchar[font_base ]:= 256  ;
  font_bc[font_base ]:=1; font_ec[font_base ]:=0;
  font_size[font_base ]:=0; font_dsize[font_base ]:=0;
  char_base[font_base ]:=0; width_base[font_base ]:=0;
  height_base[font_base ]:=0; depth_base[font_base ]:=0;
  italic_base[font_base ]:=0; lig_kern_base[font_base ]:=0;
  kern_base[font_base ]:=0; exten_base[font_base ]:=0;
  font_glue[font_base ]:=-{0xfffffff=}268435455  ; font_params[font_base ]:=7;
  param_base[font_base ]:=-1;
  for font_k:=0 to 6 do font_info[font_k].int :=0;
  end;
  endif('INITEX') 

  font_used:=xmalloc_array (boolean, font_max);
  for font_k:=font_base to font_max do font_used[font_k]:=false;

{ Compute the magic offset }
magic_offset:=str_start[math_spacing]-9*ord_noad

;

{ Initialize the print |selector|... }
if interaction=batch_mode then selector:=no_print else selector:=term_only

;
if (cur_input.loc_field <cur_input.limit_field )and( eqtb[  cat_code_base+   buffer[   cur_input.loc_field ]].hh.rh   <>escape) then start_input;
  {\.[\\input] assumed}
end

;
history:=spotless; {ready to go!}
main_control; {come to life}
final_cleanup; {prepare for death}
close_files_and_terminate;
final_end: begin  fflush (stdout ) ; ready_already:=0; if (history <> spotless) and (history <> warning_issued) then uexit(1) else uexit(0); end ;
end {|main_body|};

{ 1340. \[53] Extensions }

{tangle:pos tex.web:24529:17: }

{ The program above includes a bunch of ``hooks'' that allow further
capabilities to be added without upsetting \TeX's basic structure.
Most of these hooks are concerned with ``whatsit'' nodes, which are
intended to be used for special purposes; whenever a new extension to
\TeX\ involves a new kind of whatsit node, a corresponding change needs
to be made to the routines below that deal with such nodes,
but it will usually be unnecessary to make many changes to the
other parts of this program.

In order to demonstrate how extensions can be made, we shall treat
`\.[\\write]', `\.[\\openout]', `\.[\\closeout]', `\.[\\immediate]',
`\.[\\special]', and `\.[\\setlanguage]' as if they were extensions.
These commands are actually primitives of \TeX, and they should
appear in all implementations of the system; but let's try to imagine
that they aren't. Then the program below illustrates how a person
could add them.

Sometimes, of course, an extension will require changes to \TeX\ itself;
no system of hooks could be complete enough for all conceivable extensions.
The features associated with `\.[\\write]' are almost all confined to the
following paragraphs, but there are small parts of the |print_ln| and
|print_char| procedures that were introduced specifically to \.[\\write]
characters. Furthermore one of the token lists recognized by the scanner
is a |write_text|; and there are a few other miscellaneous places where we
have already provided for some aspect of \.[\\write].  The goal of a \TeX\
extender should be to minimize alterations to the standard parts of the
program, and to avoid them completely if possible. He or she should also
be quite sure that there's no easy way to accomplish the desired goals
with the standard features that \TeX\ already has. ``Think thrice before
extending,'' because that may save a lot of work, and it will also keep
incompatible extensions of \TeX\ from proliferating.
\xref[system dependencies]
\xref[extensions to \TeX] }

{ 1341. }

{tangle:pos tex.web:24562:23: }

{ First let's consider the format of whatsit nodes that are used to represent
the data associated with \.[\\write] and its relatives. Recall that a whatsit
has |type=whatsit_node|, and the |subtype| is supposed to distinguish
different kinds of whatsits. Each node occupies two or more words; the
exact number is immaterial, as long as it is readily determined from the
|subtype| or other data.

We shall introduce five |subtype| values here, corresponding to the
control sequences \.[\\openout], \.[\\write], \.[\\closeout], \.[\\special], and
\.[\\setlanguage]. The second word of I/O whatsits has a |write_stream| field
that identifies the write-stream number (0 to 15, or 16 for out-of-range and
positive, or 17 for out-of-range and negative).
In the case of \.[\\write] and \.[\\special], there is also a field that
points to the reference count of a token list that should be sent. In the
case of \.[\\openout], we need three words and three auxiliary subfields
to hold the string numbers for name, area, and extension. }

{ 1392. }

{tangle:pos tex.ch:4607:3: }

{ % Related to [29.526] expansion depth check
When |scan_file_name| starts it looks for a |left_brace|
(skipping \.[\\relax]es, as other \.[\\toks]-like primitives).
If a |left_brace| is found, then the procedure scans a file
name contained in a balanced token list, expanding tokens as
it goes. When the scanner finds the balanced token list, it
is converted into a string and fed character-by-character to
|more_name| to do its job the same as in the ``normal'' file
name scanning. } procedure scan_file_name_braced;
var
   save_scanner_status: small_number; {|scanner_status| upon entry}
   save_def_ref: halfword ; {|def_ref| upon entry, important if inside `\.[\\message]}
   save_cur_cs: halfword ;
   s: str_number; {temp string}
   p: halfword ; {temp pointer}
   i: integer; {loop tally}
   save_stop_at_space: boolean; {this should be in tex.ch}
   dummy: boolean;
    {Initializing}
begin save_scanner_status := scanner_status; {|scan_toks| sets |scanner_status| to |absorbing|}
  save_def_ref := def_ref; {|scan_toks| uses |def_ref| to point to the token list just read}
  save_cur_cs := cur_cs; {we set |cur_cs| back a few tokens to use in runaway errors}
    {Scanning a token list}
  cur_cs := warning_index; {for possible runaway error}
  {mimick |call_func| from pdfTeX}
  if scan_toks(false, true) <> 0 then  ; {actually do the scanning}
  {|s := tokens_to_string(def_ref);|}
  old_setting := selector; selector:=new_string;
  show_token_list( mem[ def_ref].hh.rh ,-{0xfffffff=}268435455  ,pool_size-pool_ptr);
  selector := old_setting;
  s := make_string;
  {turns the token list read in a string to input}
    {Restoring some variables}
  delete_token_ref(def_ref); {remove the token list from memory}
  def_ref := save_def_ref; {and restore |def_ref|}
  cur_cs := save_cur_cs; {restore |cur_cs|}
  scanner_status := save_scanner_status; {restore |scanner_status|}
    {Passing the read string to the input machinery}
  save_stop_at_space := stop_at_space; {save |stop_at_space|}
  stop_at_space := false; {set |stop_at_space| to false to allow spaces in file names}
  begin_name;
  for i:=str_start[s] to str_start[s+1]-1 do
    dummy := more_name(str_pool[i]); {add each read character to the current file name}
  stop_at_space := save_stop_at_space; {restore |stop_at_space|}
end;

{ 1406. }

{tangle:pos tex.ch:4970:2: }

{ This function used to be in pdftex, but is useful in tex too. } function get_nullstr: str_number;
begin
    get_nullstr := {""=}335;
end;

{ 1407. \[55] Index }

{tangle:pos tex.web:24996:12: }

{ Here is where you can find all uses of each identifier in the program,
with underlined entries pointing to where the identifier was defined.
If the identifier is only one letter long, however, you get to see only
the underlined entries. [\sl All references are to section numbers instead of
page numbers.]

This index also lists error messages and other aspects of the program
that you might want to look up some day. For example, the entry
for ``system dependencies'' lists all sections that should receive
special attention from people who are installing \TeX\ in a new
operating environment. A list of various things that can't happen appears
under ``this can't happen''. Approximately 40 sections are listed under
``inner loop''; these account for about 60\pct! of \TeX's running time,
exclusive of input and output. }