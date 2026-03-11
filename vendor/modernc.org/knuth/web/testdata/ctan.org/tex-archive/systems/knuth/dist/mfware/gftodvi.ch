@x gftodvi.web:88:
@d othercases == others: {default for cases not listed explicitly}
@y
@d othercases == else {default for cases not listed explicitly}
@z

@x
@d print_nl(#)==@+begin write_ln; write(#);@+end
@y
@d print_nl(#)==@+begin write_ln(); write(#);@+end
@d write_ln(#)==writeln(#)
@d read_ln(#)==readln(#)
@z

@x gftodvi.web:108:
@p program GF_to_DVI(@!output);
label @<Labels in the outer block@>@/
const @<Constants in the outer block@>@/
type @<Types in the outer block@>@/
var @<Globals in the outer block@>@/
procedure initialize; {this procedure gets things started properly}
  var @!i,@!j,@!m,@!n:integer; {loop indices for initializations}
  begin print_ln(banner);@/
  @<Set initial values@>@/
  end;

@ If the program has to stop prematurely, it goes to the
`|final_end|'.

@d final_end=9999 {label for the end of it all}

@<Labels...@>=final_end;
@y
@p program GF_to_DVI(@!output);
@<Labels in the outer block@>@/
const @<Constants in the outer block@>@/
type @<Types in the outer block@>@/
var @<Globals in the outer block@>@/
procedure initialize; {this procedure gets things started properly}
  var @!i,@!j,@!m,@!n:integer; {loop indices for initializations}
  begin print_ln(banner);@/
  @<Set initial values@>@/
  end;

@ If the program has to stop prematurely, it goes to the
`|final_end|'.

@d final_end=9999 {label for the end of it all}

@<Labels...@>=
@z

@x gftodvi.web:184:
@d abort(#)==@+begin print(' ',#); jump_out;@+end
@d bad_gf(#)==abort('Bad GF file: ',#,'! (at byte ',cur_loc-1:1,')')
@.Bad GF file@>

@p procedure jump_out;
begin goto final_end;
end;
@y
@d abort(#)==@+begin write_ln(stderr, ' ',#); jump_out;@+end
@d bad_gf(#)==abort('Bad GF file: ',#,'! (at byte ',cur_loc-1:1,')')
@.Bad GF file@>

@p procedure jump_out;
begin panic(final_end);
end;
@z

@x gftodvi.web:383:
@d update_terminal == break(output) {empty the terminal output buffer}
@y
@d update_terminal == {empty the terminal output buffer}
@z

@x gftodvi.web:2802:
@<Finish the \.{DVI} file and |goto final_end|@>=
begin dvi_out(post); {beginning of the postamble}
dvi_four(last_bop); last_bop:=dvi_offset+dvi_ptr-5; {|post| location}
dvi_four(25400000); dvi_four(473628672); {conversion ratio for sp}
dvi_four(1000); {magnification factor}
dvi_four(max_v); dvi_four(max_h);@/
dvi_out(0); dvi_out(3); {`\\{max\_push}' is said to be 3}@/
dvi_out(total_pages div 256); dvi_out(total_pages mod 256);@/
if not fonts_not_loaded then
  for k:=title_font to logo_font do
    if length(font_name[k])>0 then dvi_font_def(k);
dvi_out(post_post); dvi_four(last_bop); dvi_out(dvi_id_byte);@/
k:=4+((dvi_buf_size-dvi_ptr) mod 4); {the number of 223's}
while k>0 do
  begin dvi_out(223); decr(k);
  end;
@<Empty the last bytes out of |dvi_buf|@>;
goto final_end;
end
@y
@<Finish the \.{DVI} file and |goto final_end|@>=
begin dvi_out(post); {beginning of the postamble}
dvi_four(last_bop); last_bop:=dvi_offset+dvi_ptr-5; {|post| location}
dvi_four(25400000); dvi_four(473628672); {conversion ratio for sp}
dvi_four(1000); {magnification factor}
dvi_four(max_v); dvi_four(max_h);@/
dvi_out(0); dvi_out(3); {`\\{max\_push}' is said to be 3}@/
dvi_out(total_pages div 256); dvi_out(total_pages mod 256);@/
if not fonts_not_loaded then
  for k:=title_font to logo_font do
    if length(font_name[k])>0 then dvi_font_def(k);
dvi_out(post_post); dvi_four(last_bop); dvi_out(dvi_id_byte);@/
k:=4+((dvi_buf_size-dvi_ptr) mod 4); {the number of 223's}
while k>0 do
  begin dvi_out(223); decr(k);
  end;
@<Empty the last bytes out of |dvi_buf|@>;
end
@z

@x gftodvi.web:4328:
@p begin initialize; {get all variables initialized}
@<Initialize the strings@>;
start_gf; {open the input and output files}
@<Process the preamble@>;
cur_gf:=get_byte; init_str_ptr:=str_ptr;
loop@+  begin @<Initialize variables for the next character@>;
  while (cur_gf>=xxx1)and(cur_gf<=no_op) do @<Process a no-op command@>;
  if cur_gf=post then @<Finish the \.{DVI} file and |goto final_end|@>;
  if cur_gf<>boc then if cur_gf<>boc1 then abort('Missing boc!');
@.Missing boc@>
  @<Process a character@>;
  cur_gf:=get_byte; str_ptr:=init_str_ptr; pool_ptr:=str_start[str_ptr];
  end;
final_end:end.
@y
@p begin initialize; {get all variables initialized}
@<Initialize the strings@>;
start_gf; {open the input and output files}
@<Process the preamble@>;
cur_gf:=get_byte; init_str_ptr:=str_ptr;
loop@+  begin @<Initialize variables for the next character@>;
  while (cur_gf>=xxx1)and(cur_gf<=no_op) do @<Process a no-op command@>;
  if cur_gf=post then @<Finish the \.{DVI} file and |goto final_end|@>;
  if cur_gf<>boc then if cur_gf<>boc1 then abort('Missing boc!');
@.Missing boc@>
  @<Process a character@>;
  cur_gf:=get_byte; str_ptr:=init_str_ptr; pool_ptr:=str_start[str_ptr];
  end;
end.
@z
