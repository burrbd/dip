@x weave.web:77:
@d banner=='This is WEAVE, Version 4.5'
@y
@d banner=='This is WEAVE, Version 4.5 (goweave v0.0-prerelease)'
@z

@x weave.web:90:
program WEAVE(@!web_file,@!change_file,@!tex_file);
label end_of_WEAVE; {go here to finish}
@y
program WEAVE(@!web_file,@!change_file,@!tex_file);
@z

@x weave.web:177:
@d othercases == others: {default for cases not listed explicitly}
@y
@d othercases == else {default for cases not listed explicitly}
@z

@x weave.web:514:
@d print(#)==write(term_out,#) {`|print|' means write on the terminal}
@d print_ln(#)==write_ln(term_out,#) {`|print|' and then start new line}
@d new_line==write_ln(term_out) {start new line}
@y
@d print(#)==write(#) {`|print|' means write on the terminal}
@d print_ln(#)==write_ln(#) {`|print|' and then start new line}
@d new_line==write_ln() {start new line}
@d write_ln(#)==writeln(#)
@d read_ln(#)==readln(#)
@z

@x weave.web:522:
@!term_out:text_file; {the terminal as an output file}

@ Different systems have different ways of specifying that the output on a
certain file will appear on the user's terminal. Here is one way to do this
on the \PASCAL\ system that was used in \.{TANGLE}'s initial development:
@^system dependencies@>

@<Set init...@>=
rewrite(term_out,'TTY:'); {send |term_out| output to the terminal}
@y
 {the terminal as an output file}

@ Different systems have different ways of specifying that the output on a
certain file will appear on the user's terminal. Here is one way to do this
on the \PASCAL\ system that was used in \.{TANGLE}'s initial development:
@^system dependencies@>

@<Set init...@>=
 {send |term_out| output to the terminal}
@z

@x weave.web:532:
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

@x weave.web:693:
procedure jump_out;
begin goto end_of_WEAVE;
end;
@y
procedure jump_out;
begin panic(end_of_WEAVE);
end;
@z


@x weave.web:4855:
@<Check that all changes have been read@>;
end_of_WEAVE:
@y
@<Check that all changes have been read@>;
@z

@x weave.web:4878:
case history of
spotless: print_nl('(No errors were found.)');
harmless_message: print_nl('(Did you see the warning message above?)');
error_message: print_nl('(Pardon me, but I think I spotted something wrong.)');
fatal_message: print_nl('(That was a fatal error, my friend.)');
end {there are no other cases}
@y
case history of
spotless: print_nl('(No errors were found.)');
harmless_message: print_nl('(Did you see the warning message above?)');
error_message: print_nl('(Pardon me, but I think I spotted something wrong.)');
fatal_message: print_nl('(That was a fatal error, my friend.)');
end; {there are no other cases}
write_ln();
@z
