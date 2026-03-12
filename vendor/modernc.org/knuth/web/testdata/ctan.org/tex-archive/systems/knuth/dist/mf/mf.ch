@x mf.web:250:
start_of_MF@t\hskip-2pt@>, end_of_MF@t\hskip-2pt@>,@,final_end;
@y
start_of_MF;
@z

@x mf.web:337:
@d othercases == others: {default for cases not listed explicitly}
@y
@d othercases == else {default for cases not listed explicitly}
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

@x mf.web:22628:
@<Undump constants for consistency check@>=
x:=base_file^.int;
@y
@<Undump constants for consistency check@>=
read(base_file, x);
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
