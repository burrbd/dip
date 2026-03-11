@x mft.web:143:
@p @t\4@>@<Compiler directives@>@/
program MFT(@!mf_file,@!change_file,@!style_file,@!tex_file);
label end_of_MFT; {go here to finish}
@y
@p @t\4@>@<Compiler directives@>@/
program MFT(@!mf_file,@!change_file,@!style_file,@!tex_file);
@z

@x mft.web:211:
@d othercases == others: {default for cases not listed explicitly}
@y
@d othercases == else {default for cases not listed explicitly}
@z

@x mft.web:460:
@d new_line==write_ln(term_out) {start new line on the terminal}
@y
@d new_line==write_ln(term_out) {start new line on the terminal}
@d write_ln(#)==writeln(#)
@d read_ln(#)==readln(#)
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
@<Error handling...@>=
procedure jump_out;
begin goto end_of_MFT;
end;
@y
@<Error handling...@>=
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
