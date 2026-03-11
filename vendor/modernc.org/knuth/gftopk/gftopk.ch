@x gftopk:100:
@d banner=='This is GFtoPK, Version 2.4' {printed when the program starts}
@y
@d banner=='This is GFtoPK, Version 2.4 (gogftopk v0.0-prerelease)' {printed when the program starts}
@z

@x gftopk:116:
@d othercases == others: {default for cases not listed explicitly}
@y
@d othercases == else {default for cases not listed explicitly}
@z

@x
@d print_ln(#)==write_ln(#)
@y
@d print_ln(#)==write_ln(#)
@d write_ln(#)==writeln(#)
@z

@x gftopk:129:
@p program GFtoPK(@!gf_file,@!pk_file,@!output);
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
`|final_end|'.

@d final_end=9999 {label for the end of it all}

@<Labels...@>=final_end;
@y
@p program GFtoPK(@!gf_file,@!pk_file,@!output);
const @<Constants in the outer block@>@/
type @<Types in the outer block@>@/
var @<Globals in the outer block@>@/
procedure initialize; {this procedure gets things started properly}
  var i:integer; {loop index for initializations}
  begin print_ln(banner);@/
  @<Set initial values@>@/
  end;

@ If the program has to stop prematurely, it goes to the
`|final_end|'.

@d final_end=9999 {label for the end of it all}
@z

@x gftopk:172:
@d abort(#)==begin print(' ',#); jump_out;
    end
@d bad_gf(#)==abort('Bad GF file: ',#,'!')
@.Bad GF file@>

@p procedure jump_out;
begin goto final_end;
end;
@y
@d abort(#)==begin print(stderr,' ',#); jump_out;
    end
@d bad_gf(#)==abort('Bad GF file: ',#,'!')
@.Bad GF file@>

@p procedure jump_out;
begin
	panic(final_end);
end;
@z

@x gftopk:2166:
@p begin
  initialize ;
  convert_gf_file ;
  @<Check for unrasterized locators@> ;
  print_ln(gf_len:1,' bytes packed to ',pk_loc:1,' bytes.') ;
final_end : end .
@y
@p begin
  initialize ;
  convert_gf_file ;
  @<Check for unrasterized locators@> ;
  print_ln(gf_len:1,' bytes packed to ',pk_loc:1,' bytes.') ;
end .
@z
