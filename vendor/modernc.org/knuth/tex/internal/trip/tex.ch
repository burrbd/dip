@x tex.web:193:
@d banner=='This is TeX, Version 3.141592653' {printed when \TeX\ starts}
@y
@d banner=='This is TeX, Version 3.141592653 (TRIP)' {printed when \TeX\ starts}
@z

@x tex.web:287:
start_of_TEX@t\hskip-2pt@>, end_of_TEX@t\hskip-2pt@>,@,final_end;
@y
start_of_TEX;
@z

@x tex.web:309: TRIP
@d stat==@{ {change this to `$\\{stat}\equiv\null$' when gathering
  usage statistics}
@d tats==@t@>@} {change this to `$\\{tats}\equiv\null$' when gathering
  usage statistics}
@y
@d stat== {change this to `$\\{stat}\equiv\null$' when gathering
  usage statistics}
@d tats== {change this to `$\\{tats}\equiv\null$' when gathering
  usage statistics}
@z

@x tex.web:378:
@d othercases == others: {default for cases not listed explicitly}
@y
@d othercases == else {default for cases not listed explicitly}
@d write_ln(#)==writeln(#)
@d read_ln(#)==readln(#)
@z

@x tex.web:393: TRIP
@!mem_max=30000; {greatest index in \TeX's internal |mem| array;
  must be strictly less than |max_halfword|;
  must be equal to |mem_top| in \.{INITEX}, otherwise |>=mem_top|}
@!mem_min=0; {smallest index in \TeX's internal |mem| array;
  must be |min_halfword| or more;
  must be equal to |mem_bot| in \.{INITEX}, otherwise |<=mem_bot|}
@!buf_size=500; {maximum number of characters simultaneously present in
  current lines of open files and in control sequences between
  \.{\\csname} and \.{\\endcsname}; must not exceed |max_halfword|}
@!error_line=72; {width of context lines on terminal error messages}
@!half_error_line=42; {width of first lines of contexts in terminal
  error messages; should be between 30 and |error_line-15|}
@!max_print_line=79; {width of longest text lines output; should be at least 60}
@y
@!mem_max=3000; {greatest index in \TeX's internal |mem| array;
  must be strictly less than |max_halfword|;
  must be equal to |mem_top| in \.{INITEX}, otherwise |>=mem_top|}
@!mem_min=1; {smallest index in \TeX's internal |mem| array;
  must be |min_halfword| or more;
  must be equal to |mem_bot| in \.{INITEX}, otherwise |<=mem_bot|}
@!buf_size=500; {maximum number of characters simultaneously present in
  current lines of open files and in control sequences between
  \.{\\csname} and \.{\\endcsname}; must not exceed |max_halfword|}
@!error_line=64; {width of context lines on terminal error messages}
@!half_error_line=32; {width of first lines of contexts in terminal
  error messages; should be between 30 and |error_line-15|}
@!max_print_line=72; {width of longest text lines output; should be at least 60}
@z

@x tex.web:441: TRIP
@d mem_bot=0 {smallest index in the |mem| array dumped by \.{INITEX};
  must not be less than |mem_min|}
@d mem_top==30000 {largest index in the |mem| array dumped by \.{INITEX};
@y
@d mem_bot=1 {smallest index in the |mem| array dumped by \.{INITEX};
  must not be less than |mem_min|}
@d mem_top==3000 {largest index in the |mem| array dumped by \.{INITEX};
@z

@x tex.web:948:
@p function input_ln(var f:alpha_file;@!bypass_eoln:boolean):boolean;
  {inputs the next line or returns |false|}
var last_nonblank:0..buf_size; {|last| with trailing blanks removed}
begin if bypass_eoln then if not eof(f) then get(f);
@y
@p function input_ln(var f:alpha_file;@!bypass_eoln:boolean):boolean;
  {inputs the next line or returns |false|}
var last_nonblank:0..buf_size; {|last| with trailing blanks removed}
begin if bypass_eoln then if not eof(f) and eoln(f) then get(f);
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
buf:eight_bits;
@z

@x tex.web:23664:
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

@x tex.web:24324:
procedure close_files_and_terminate;
var k:integer; {all-purpose index}
begin @<Finish the extensions@>; new_line_char:=-1;
@!stat if tracing_stats>0 then @<Output statistics about this job@>;@;@+tats@/
wake_up_terminal; @<Finish the \.{DVI} file@>;
if log_opened then
  begin wlog_cr; a_close(log_file); selector:=selector-2;
  if selector=term_only then
    begin print_nl("Transcript written on ");
@.Transcript written...@>
    slow_print(log_name); print_char(".");
    end;
  end;
end;
@y
procedure close_files_and_terminate;
var k:integer; {all-purpose index}
begin @<Finish the extensions@>; new_line_char:=-1;
@!stat if tracing_stats>0 then @<Output statistics about this job@>;@;@+tats@/
wake_up_terminal; @<Finish the \.{DVI} file@>;
if log_opened then
  begin wlog_cr; a_close(log_file); selector:=selector-2;
  if selector=term_only then
    begin print_nl("Transcript written on ");
@.Transcript written...@>
    slow_print(log_name); print_char(".");
    end;
  end;
  write_ln(term_out);
end;
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
