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

@x weave.web:515:
@d print_ln(#)==write_ln(term_out,#) {`|print|' and then start new line}
@y
@d print_ln(#)==write_ln(term_out,#) {`|print|' and then start new line}
@d write_ln(#)==writeln(#)
@d read_ln(#)==readln(#)
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
