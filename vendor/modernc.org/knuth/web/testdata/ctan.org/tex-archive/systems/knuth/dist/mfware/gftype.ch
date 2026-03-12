@x gftype.wen:82:
@d othercases == others: {default for cases not listed explicitly}
@y
@d othercases == else {default for cases not listed explicitly}
@z

@x gftype.wen:94:
@d print_nl==write_ln
@y
@d print_nl==write_ln()
@d write_ln(#)==writeln(#)
@d read_ln(#)==readln(#)
@z

@x gftype.web:96:
@p program GF_type(@!gf_file,@!output);
label @<Labels in the outer block@>@/
@y
@p program GF_type(@!gf_file,@!output,stderr);
@z

@x gftype.web:112:
@<Labels...@>=final_end;
@y
@z

@x gftype.wen:144:
@d abort(#)==begin print(' ',#); jump_out;
    end
@d bad_gf(#)==abort('Bad GF file: ',#,'!')
@.Bad GF file@>

@p procedure jump_out;
begin goto final_end;
end;
@y
@d abort(#)==begin print(stderr, ' ',#); jump_out;
    end
@d bad_gf(#)==abort('Bad GF file: ',#,'!')
@.Bad GF file@>

@p procedure jump_out;
begin
	panic(final_end);
end;
@z

@x gftype.web:651:
@!gf_file:byte_file; {the stuff we are \.{GF}typing}
@y
@!gf_file:byte_file; {the stuff we are \.{GF}typing}
stderr:text;
@z

@x gftype.web:744:
@d update_terminal == break(term_out) {empty the terminal output buffer}
@y
@d update_terminal == {empty the terminal output buffer}
@z

@x gftype.web:1280:
final_end:end.
@y
end.
@z
