@x dvitype.web:90:
@d banner=='This is DVItype, Version 3.6' {printed when the program starts}
@y
@d banner=='This is DVItype, Version 3.6 (godvitype v0.0-prerelease)' {printed when the program starts}
@z

@x dvitype.web:116:
@d othercases == others: {default for cases not listed explicitly}
@y
@d othercases == else {default for cases not listed explicitly}
@z

@x dvitype.web:127:
@d print_ln(#)==write_ln(#)
@y
@d print_ln(#)==write_ln(#)
@d write_ln(#)==writeln(#)
@d read_ln(#)==readln(#)
@z

@x dvitype.web:129:
@p program DVI_type(@!dvi_file,@!output);
label @<Labels in the outer block@>@/
const @<Constants in the outer block@>@/
type @<Types in the outer block@>@/
var @<Globals in the outer block@>@/
procedure initialize; {this procedure gets things started properly}
  var i:integer; {loop index for initializations}
  begin print_ln(banner);@/
  @<Set initial values@>@/
  end;

@ If the program has to stop prematurely, it goes to the
`|final_end|'. Another label, |done|, is used when stopping normally.

@d final_end=9999 {label for the end of it all}
@d done=30 {go here when finished with a subtask}

@<Labels...@>=final_end,done;
@y
@p program DVI_type(@!dvi_file,@!output);
label @<Labels in the outer block@>@/
const @<Constants in the outer block@>@/
type @<Types in the outer block@>@/
var @<Globals in the outer block@>@/
procedure initialize; {this procedure gets things started properly}
  var i:integer; {loop index for initializations}
  begin print_ln(banner);@/
  @<Set initial values@>@/
  end;

@ If the program has to stop prematurely, it goes to the
`|final_end|'. Another label, |done|, is used when stopping normally.

@d final_end=9999 {label for the end of it all}
@d done=30 {go here when finished with a subtask}

@<Labels...@>=done;
@z

@x dvitype.web:177:
@d abort(#)==begin print(' ',#); jump_out;
    end
@d bad_dvi(#)==abort('Bad DVI file: ',#,'!')
@.Bad DVI file@>

@p procedure jump_out;
begin goto final_end;
end;
@y
@d abort(#)==begin print(stderr,' ',#); jump_out;
    end
@d bad_dvi(#)==abort('Bad DVI file: ',#,'!')
@.Bad DVI file@>

@p procedure jump_out;
begin
	panic(final_end);
end;
@z

@x dvitype.web:1303:
out_mode:=the_works; max_pages:=1000000; start_vals:=0; start_there[0]:=false;
@y
@z

@x dvitype.web:1336:
@d update_terminal == break(term_out) {empty the terminal output buffer}
@y
@d update_terminal == {empty the terminal output buffer}
@z

@x dvitype.web:2166:
@p @t\4@>@<Declare the procedure called |scan_bop|@>@;
procedure skip_pages(@!bop_seen:boolean);
label 9999; {end of this subroutine}
var p:integer; {a parameter}
@!k:0..255; {command code}
@!down_the_drain:integer; {garbage}
begin showing:=false;
while true do
@y
@p @t\4@>@<Declare the procedure called |scan_bop|@>@;
procedure skip_pages(@!bop_seen:boolean);
label 9999; {end of this subroutine}
var p:integer; {a parameter}
@!k:0..255; {command code}
@!down_the_drain:integer; {garbage}
begin showing:=false;
p:=down_the_drain; { Go: use the var }
while true do
@z

@x dvitype.web:2404:
@* The main program.
Now we are ready to put it all together. This is where \.{DVItype} starts,
and where it ends.

@p begin initialize; {get all variables initialized}
dialog; {set up all the options}
@y
@* The main program.
Now we are ready to put it all together. This is where \.{DVItype} starts,
and where it ends.

@p begin initialize; {get all variables initialized}
@<Print all the selected options@>;
@z

@x dvitype.web:2421:
final_end:end.
@y
end.
@z
