@x gftype.web:73:
@d banner=='This is GFtype, Version 3.1' {printed when the program starts}
@y
@d banner=='This is GFtype, Version 3.1 (gogftype v0.0-prerelease)' {printed when the program starts}
@z

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

@x gftype.web:723:
wants_mnemonics:=true; wants_pixels:=true;
@y
@z

@x gftype.web:744:
@d update_terminal == break(term_out) {empty the terminal output buffer}
@y
@d update_terminal == {empty the terminal output buffer}
@z

@x gftype.web:778:
@p procedure dialog;
label 1,2;
begin rewrite(term_out); {prepare the terminal for output}
write_ln(term_out,banner);@/
@<Determine whether the user |wants_mnemonics|@>;
@<Determine whether the user |wants_pixels|@>;
@<Print all the selected options@>;
end;

@ @<Determine whether the user |wants_mnemonics|@>=
1: write(term_out,'Mnemonic output? (default=no, ? for help): ');
@.Mnemonic output?@>
input_ln;
buffer[0]:=lower_casify(buffer[0]);
if buffer[0]<>"?" then
  wants_mnemonics:=(buffer[0]="y")or(buffer[0]="1")or(buffer[0]="t")
else  begin write(term_out,'Type Y for complete listing,');
  write_ln(term_out,' N for errors/images only.');
  goto 1;
  end

@ @<Determine whether the user |wants_pixels|@>=
2: write(term_out,'Pixel output? (default=yes, ? for help): ');
@.Pixel output?@>
input_ln;
buffer[0]:=lower_casify(buffer[0]);
if buffer[0]<>"?" then
  wants_pixels:=(buffer[0]="y")or(buffer[0]="1")or(buffer[0]="t")
    or(buffer[0]=" ")
else  begin write(term_out,'Type Y to list characters pictorially');
  write_ln(term_out,' with *''s, N to omit this option.');
  goto 2;
  end
@y
@p procedure dialog;
label 1,2;
begin
@<Print all the selected options@>;
end;
@z


@x gftype.web:1257:
print(' altogether.');
@y
print(' altogether.');
print_nl;
@z

@x gftype.web:1280:
final_end:end.
@y
end.
@z
