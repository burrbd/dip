@x tex.web:287:
start_of_TEX@t\hskip-2pt@>, end_of_TEX@t\hskip-2pt@>,@,final_end;
@y
start_of_TEX;
@z

@x tex.web:378:
@d othercases == others: {default for cases not listed explicitly}
@y
@d othercases == else {default for cases not listed explicitly}
@d write_ln(#)==writeln(#)
@d read_ln(#)==readln(#)
@z

@x tex.web:1002:
@d update_terminal == break(term_out) {empty the terminal output buffer}
@d clear_terminal == break_in(term_in,true) {clear the terminal input buffer}
@d wake_up_terminal == do_nothing {cancel the user's cancellation of output}
@y
@d update_terminal == {empty the terminal output buffer}
@d clear_terminal == {clear the terminal input buffer}
@d wake_up_terminal == do_nothing {cancel the user's cancellation of output}
@z

@x tex.web:1035:
if format_ident=0 then
  begin write_ln(term_out,'Buffer size exceeded!'); goto final_end;
@y
if format_ident=0 then
  begin write_ln(term_out,'Buffer size exceeded!'); panic(final_end);
@z

@x tex.web:1255:
@p @!init function get_strings_started:boolean; {initializes the string pool,
  but returns |false| if something goes wrong}
label done,exit;
var k,@!l:0..255; {small indices or counters}
@!m,@!n:text_char; {characters input from |pool_file|}
@!g:str_number; {garbage}
@!a:integer; {accumulator for check sum}
@!c:boolean; {check sum has been checked}
begin pool_ptr:=0; str_ptr:=0; str_start[0]:=0;
@y
@p @!init function get_strings_started:boolean; {initializes the string pool,
  but returns |false| if something goes wrong}
label done,exit;
var k,@!l:0..255; {small indices or counters}
@!m,@!n:text_char; {characters input from |pool_file|}
@!g:str_number; {garbage}
@!a:integer; {accumulator for check sum}
@!c:boolean; {check sum has been checked}
begin
if g=0 then;
pool_ptr:=0; str_ptr:=0; str_start[0]:=0;
@z

@x tex.web:1871:
procedure jump_out;
begin goto end_of_TEX;
end;
@y
procedure jump_out;
begin panic(end_of_TEX);
end;
@z

@x tex.web:10431:
@!tfm_file:byte_file;
@y
@!tfm_file:byte_file;
buf:eight_bits; { tfm_file buffer }
@z

@x tex.web:16389:
@ @<Other local variables for |try_break|@>=
@!prev_prev_r:pointer; {a step behind |prev_r|, if |type(prev_r)=delta_node|}
@!s:pointer; {runs through nodes ahead of |cur_p|}
@!q:pointer; {points to a new node being created}
@!v:pointer; {points to a glue specification or a node ahead of |cur_p|}
@!t:integer; {node count, if |cur_p| is a discretionary node}
@!f:internal_font_number; {used in character width calculation}
@!l:halfword; {line number of current active node}
@!node_r_stays_active:boolean; {should node |r| remain in the active list?}
@!line_width:scaled; {the current line will be justified to this width}
@!fit_class:very_loose_fit..tight_fit; {possible fitness class of test line}
@!b:halfword; {badness of test line}
@!d:integer; {demerits of test line}
@!artificial_demerits:boolean; {has |d| been forced to zero?}
@!save_link:pointer; {temporarily holds value of |link(cur_p)|}
@!shortfall:scaled; {used in badness calculations}
@y
@ @<Other local variables for |try_break|@>=
@!prev_prev_r:pointer; {a step behind |prev_r|, if |type(prev_r)=delta_node|}
@!s:pointer; {runs through nodes ahead of |cur_p|}
@!q:pointer; {points to a new node being created}
@!v:pointer; {points to a glue specification or a node ahead of |cur_p|}
@!t:integer; {node count, if |cur_p| is a discretionary node}
@!f:internal_font_number; {used in character width calculation}
@!l:halfword; {line number of current active node}
@!node_r_stays_active:boolean; {should node |r| remain in the active list?}
@!line_width:scaled; {the current line will be justified to this width}
@!fit_class:very_loose_fit..tight_fit; {possible fitness class of test line}
@!b:halfword; {badness of test line}
@!d:integer; {demerits of test line}
@!artificial_demerits:boolean; {has |d| been forced to zero?}
@!shortfall:scaled; {used in badness calculations}
@z


@x tex.web:23662:
procedure show_whatever;
label common_ending;
var p:pointer; {tail of a token list to show}
begin case cur_chr of
@y
procedure show_whatever;
label common_ending;
var p:pointer; {tail of a token list to show}
begin
if p=0 then;
case cur_chr of
@z

@x tex.web:23872:
@<Undump constants for consistency check@>=
x:=fmt_file^.int;
@y
@<Undump constants for consistency check@>=
read(fmt_file,x);
@z

@x tex.web:24303:
main_control; {come to life}
final_cleanup; {prepare for death}
end_of_TEX: close_files_and_terminate;
final_end: ready_already:=0;
end.
@y
main_control; {come to life}
final_cleanup; {prepare for death}
close_files_and_terminate;
final_end: ready_already:=0;
end.
@z

@x tex.web:24656:
procedure do_extension;
var i,@!j,@!k:integer; {all-purpose integers}
@!p,@!q,@!r:pointer; {all-purpose pointers}
begin case cur_chr of
@y
procedure do_extension;
var k:integer; {all-purpose integers}
@!p:pointer; {all-purpose pointers}
begin case cur_chr of
@z
