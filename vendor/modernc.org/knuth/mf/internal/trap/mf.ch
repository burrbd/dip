@x mf.web:160:
@d banner=='This is METAFONT, Version 2.71828182' {printed when \MF\ starts}
@y
@d banner=='This is METAFONT, Version 2.71828182 (TRAP)' {printed when \MF\ starts}
@z

@x mf.web:250:
start_of_MF@t\hskip-2pt@>, end_of_MF@t\hskip-2pt@>,@,final_end;
@y
start_of_MF;
@z

@x mf.web:272: TRAP
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

@x mf.web:337:
@d othercases == others: {default for cases not listed explicitly}
@y
@d othercases == else {default for cases not listed explicitly}
@z

@x mf.web:349: TRAP
@!mem_max=30000; {greatest index in \MF's internal |mem| array;
  must be strictly less than |max_halfword|;
  must be equal to |mem_top| in \.{INIMF}, otherwise |>=mem_top|}
@!max_internal=100; {maximum number of internal quantities}
@!buf_size=500; {maximum number of characters simultaneously present in
  current lines of open files; must not exceed |max_halfword|}
@!error_line=72; {width of context lines on terminal error messages}
@!half_error_line=42; {width of first lines of contexts in terminal
  error messages; should be between 30 and |error_line-15|}
@!max_print_line=79; {width of longest text lines output; should be at least 60}
@!screen_width=768; {number of pixels in each row of screen display}
@!screen_depth=1024; {number of pixels in each column of screen display}
@!stack_size=30; {maximum number of simultaneous input sources}
@!max_strings=2000; {maximum number of strings; must not exceed |max_halfword|}
@!string_vacancies=8000; {the minimum number of characters that should be
  available for the user's identifier names and strings,
  after \MF's own error messages are stored}
@!pool_size=32000; {maximum number of characters in strings, including all
  error messages and help texts, and the names of all identifiers;
  must exceed |string_vacancies| by the total
  length of \MF's own strings, which is currently about 22000}
@!move_size=5000; {space for storing moves in a single octant}
@!max_wiggle=300; {number of autorounded points per cycle}
@!gf_buf_size=800; {size of the output buffer, must be a multiple of 8}
@y
@!mem_max=3000; {greatest index in \MF's internal |mem| array;
  must be strictly less than |max_halfword|;
  must be equal to |mem_top| in \.{INIMF}, otherwise |>=mem_top|}
@!max_internal=100; {maximum number of internal quantities}
@!buf_size=500; {maximum number of characters simultaneously present in
  current lines of open files; must not exceed |max_halfword|}
@!error_line=64; {width of context lines on terminal error messages}
@!half_error_line=32; {width of first lines of contexts in terminal
  error messages; should be between 30 and |error_line-15|}
@!max_print_line=72; {width of longest text lines output; should be at least 60}
@!screen_width=100; {number of pixels in each row of screen display}
@!screen_depth=200; {number of pixels in each column of screen display}
@!stack_size=30; {maximum number of simultaneous input sources}
@!max_strings=2000; {maximum number of strings; must not exceed |max_halfword|}
@!string_vacancies=8000; {the minimum number of characters that should be
  available for the user's identifier names and strings,
  after \MF's own error messages are stored}
@!pool_size=32000; {maximum number of characters in strings, including all
  error messages and help texts, and the names of all identifiers;
  must exceed |string_vacancies| by the total
  length of \MF's own strings, which is currently about 22000}
@!move_size=5000; {space for storing moves in a single octant}
@!max_wiggle=300; {number of autorounded points per cycle}
@!gf_buf_size=8; {size of the output buffer, must be a multiple of 8}
@z

@x mf.web:399: TRAP
@d mem_top==30000 {largest index in the |mem| array dumped by \.{INIMF};
  must be substantially larger than |mem_min|
  and not greater than |mem_max|}
@y
@d mem_top==3000 {largest index in the |mem| array dumped by \.{INIMF};
  must be substantially larger than |mem_min|
  and not greater than |mem_max|}
@z

@x mf.web:922:
@d update_terminal == break(term_out) {empty the terminal output buffer}
@d clear_terminal == break_in(term_in,true) {clear the terminal input buffer}
@y
@d update_terminal == {empty the terminal output buffer}
@d clear_terminal == {clear the terminal input buffer}
@z

@x mf.web:955:
  begin write_ln(term_out,'Buffer size exceeded!'); goto final_end;
@y
  begin write_ln(term_out,'Buffer size exceeded!'); panic(final_end);
@z

@x mf.web:997:
@p function init_terminal:boolean; {gets the terminal input started}
label exit;
begin t_open_in;
loop@+begin wake_up_terminal; write(term_out,'**'); update_terminal;
@.**@>
  if not input_ln(term_in,true) then {this shouldn't happen}
@y
@p function init_terminal:boolean; {gets the terminal input started}
label exit;
begin t_open_in;
loop@+begin wake_up_terminal; write(term_out,'**'); update_terminal;
@.**@>
  if not input_ln(term_in,false) then {this shouldn't happen}
@z

@x mf.web:1413:
@d wterm_ln(#)==write_ln(term_out,#)
@y
@d wterm_ln(#)==write_ln(term_out,#)
@d write_ln(#)==writeln(#)
@d read_ln(#)==readln(#)
@z

@x mf.web:1738:
procedure jump_out;
begin goto end_of_MF;
end;
@y
procedure jump_out;
begin panic(end_of_MF);
end;
@z

@x mf.web:12090 TRAP:
begin init_screen:=false;
@y
begin init_screen:=true; {screen instructions will be logged}
@z

@x mf.web:12138:
@p procedure blank_rectangle(@!left_col,@!right_col:screen_col;
  @!top_row,@!bot_row:screen_row);
var @!r:screen_row;
@!c:screen_col;
@y
@p procedure blank_rectangle(@!left_col,@!right_col:screen_col;
  @!top_row,@!bot_row:screen_row);
@z

@x mf.web:12165:
@p procedure paint_row(@!r:screen_row;@!b:pixel_color;var @!a:trans_spec;
  @!n:screen_col);
var @!k:screen_col; {an index into |a|}
@!c:screen_col; {an index into |screen_pixel|}
@y
@p procedure paint_row(@!r:screen_row;@!b:pixel_color;var @!a:trans_spec;
  @!n:screen_col);
var @!k:screen_col; {an index into |a|}
@z

@x mf.web:15878:
@p procedure open_log_file;
var @!old_setting:0..max_selector; {previous |selector| setting}
@!k:0..buf_size; {index into |months| and |buffer|}
@!l:0..buf_size; {end of first input line}
@!m:integer; {the current month}
@y
@p procedure open_log_file;
var @!old_setting:0..max_selector; {previous |selector| setting}
@!k:0..buf_size; {index into |months| and |buffer|}
@!l:0..buf_size; {end of first input line}
@z

@x mf.web:22883:
main_control; {come to life}
final_cleanup; {prepare for death}
end_of_MF: close_files_and_terminate;
final_end: ready_already:=0;
end.
@y
main_control; {come to life}
final_cleanup; {prepare for death}
close_files_and_terminate;
final_end: ready_already:=0;
end.
@z

@x mf.web:22901:
procedure close_files_and_terminate;
var @!k:integer; {all-purpose index}
@!lh:integer; {the length of the \.{TFM} header, in words}
@!lk_offset:0..256; {extra words inserted at beginning of |lig_kern| array}
@!p:pointer; {runs through a list of \.{TFM} dimensions}
@!x:scaled; {a |tfm_width| value being output to the \.{GF} file}
begin
@!stat if internal[tracing_stats]>0 then
  @<Output statistics about this job@>;@;@+tats@/
wake_up_terminal; @<Finish the \.{TFM} and \.{GF} files@>;
if log_opened then
  begin wlog_cr;
  a_close(log_file); selector:=selector-2;
  if selector=term_only then
    begin print_nl("Transcript written on ");
@.Transcript written...@>
    slow_print(log_name); print_char(".");
    end;
  end;
end;
@y
procedure close_files_and_terminate;
var @!k:integer; {all-purpose index}
@!lh:integer; {the length of the \.{TFM} header, in words}
@!lk_offset:0..256; {extra words inserted at beginning of |lig_kern| array}
@!p:pointer; {runs through a list of \.{TFM} dimensions}
@!x:scaled; {a |tfm_width| value being output to the \.{GF} file}
begin
@!stat if internal[tracing_stats]>0 then
  @<Output statistics about this job@>;@;@+tats@/
wake_up_terminal; @<Finish the \.{TFM} and \.{GF} files@>;
if log_opened then
  begin wlog_cr;
  a_close(log_file); selector:=selector-2;
  if selector=term_only then
    begin print_nl("Transcript written on ");
@.Transcript written...@>
    slow_print(log_name); print_char(".");
    end;
  end;
  write_ln(term_out);
end;
@z
