@x mft.web:128:
The ``banner line'' defined here should be changed whenever \.{MFT}
is modified.

@d banner=='This is MFT, Version 2.1'
@y
The ``banner line'' defined here should be changed whenever \.{MFT}
is modified.

@d banner=='This is MFT, Version 2.1 (gomft v0.0-prerelease)'
@z

@x mft.web:143:
@p @t\4@>@<Compiler directives@>@/
program MFT(@!mf_file,@!change_file,@!style_file,@!tex_file);
label end_of_MFT; {go here to finish}
@y
@p @t\4@>@<Compiler directives@>@/
program MFT(@!mf_file,@!change_file,@!style_file,@!tex_file,@!output);
@z

@x mft.web:211:
@d othercases == others: {default for cases not listed explicitly}
@y
@d othercases == else {default for cases not listed explicitly}
@z

@x mft.web:224:
@!line_length=80; {lines of \TeX\ output have at most this many characters,
  should be less than 256}
@y TeX-live compatibility
@!line_length=79; {lines of \TeX\ output have at most this many characters,
  should be less than 256}
@z

@x mft.web:458:
@d print(#)==write(term_out,#) {`|print|' means write on the terminal}
@d print_ln(#)==write_ln(term_out,#) {`|print|' and then start new line}
@d new_line==write_ln(term_out) {start new line on the terminal}
@y
@d term_out==output
@d print(#)==write(term_out,#) {`|print|' means write on the terminal}
@d print_ln(#)==write_ln(term_out,#) {`|print|' and then start new line}
@d new_line==write_ln(term_out) {start new line on the terminal}
@d write_ln(#)==writeln(#)
@d read_ln(#)==readln(#)
@z

@x mft.web:466:
@!term_out:text_file; {the terminal as an output file}
@y
@z

@x mft.web:474:
rewrite(term_out,'TTY:'); {send |term_out| output to the terminal}
@y
@z

@x mft.web:476:
@ The |update_terminal| procedure is called when we want
to make sure that everything we have output to the terminal so far has
actually left the computer's internal buffers and been sent.
@^system dependencies@>

@d update_terminal == break(term_out) {empty the terminal output buffer}
@y
@ The |update_terminal| procedure is called when we want
to make sure that everything we have output to the terminal so far has
actually left the computer's internal buffers and been sent.
@^system dependencies@>

@d update_terminal == {empty the terminal output buffer}
@z

@x mft.web:605:
procedure jump_out;
begin goto end_of_MFT;
end;
@y
procedure jump_out;
begin
	panic(end_of_MFT);
end;
@z

@x mft.web:1942:
@* The main program.
Let's put it all together now: \.{MFT} starts and ends here.
@^system dependencies@>

@p begin initialize; {beginning of the main program}
print_ln(banner); {print a ``banner line''}
@<Store all the primitives@>;
@<Store all the translations@>;
@<Initialize the input...@>;
do_the_translation;
@<Check that all changes have been read@>;
end_of_MFT:{here files should be closed if the operating system requires it}
@<Print the job |history|@>;
end.
@y
@* The main program.
Let's put it all together now: \.{MFT} starts and ends here.
@^system dependencies@>

@p begin initialize; {beginning of the main program}
print_ln(banner); {print a ``banner line''}
@<Store all the primitives@>;
@<Store all the translations@>;
@<Initialize the input...@>;
do_the_translation;
@<Check that all changes have been read@>;
@<Print the job |history|@>;
end.
@z


@x mft.web:1942:
@<Print the job |history|@>=
case history of
spotless: print_nl('(No errors were found.)');
harmless_message: print_nl('(Did you see the warning message above?)');
error_message: print_nl('(Pardon me, but I think I spotted something wrong.)');
fatal_message: print_nl('(That was a fatal error, my friend.)');
end {there are no other cases}
@y
@<Print the job |history|@>=
case history of
spotless: print_nl('(No errors were found.)');
harmless_message: print_nl('(Did you see the warning message above?)');
error_message: print_nl('(Pardon me, but I think I spotted something wrong.)');
fatal_message: print_nl('(That was a fatal error, my friend.)');
end; {there are no other cases}
write_ln('');
@z
