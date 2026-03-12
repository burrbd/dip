@x tangle.web:83:
@p @t\4@>@<Compiler directives@>@/
program TANGLE(@!web_file,@!change_file,@!Pascal_file,@!pool);
label end_of_TANGLE; {go here to finish}
const @<Constants in the outer block@>@/
type @<Types in the outer block@>@/
var @<Globals in the outer block@>@/
@<Error handling procedures@>@/
procedure initialize;
  var @<Local variables for initialization@>@/
  begin @<Set initial values@>@/
  end;
@y
@p @t\4@>@<Compiler directives@>@/
program TANGLE(@!web_file,@!change_file,@!Pascal_file,@!pool);
const @<Constants in the outer block@>@/
type @<Types in the outer block@>@/
var @<Globals in the outer block@>@/
@<Error handling procedures@>@/
procedure initialize;
  var @<Local variables for initialization@>@/
  begin @<Set initial values@>@/
  end;
@z

@x tangle.web:174:
@d othercases == others: {default for cases not listed explicitly}
@y
@d othercases == else {default for cases not listed explicitly}
@z

@x tangle.web:510:
@d print_ln(#)==write_ln(term_out,#) {`|print|' and then start new line}
@y
@d print_ln(#)==write_ln(term_out,#) {`|print|' and then start new line}
@d write_ln(#)==writeln(#)
@d read_ln(#)==readln(#)
@z

@x tangle.web:527:
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

@x tangle.web:691:
@d fatal_error(#)==begin new_line; print(#); error; mark_fatal; jump_out;
  end

@<Error handling...@>=
procedure jump_out;
begin goto end_of_TANGLE;
end;
@y
@d fatal_error(#)==begin write_ln(stderr,#); error; mark_fatal; jump_out;
  end

@<Error handling...@>=
procedure jump_out;
begin panic(end_of_TANGLE);
end;
@z

@x tangle.web:3246:
@p begin initialize;
@<Initialize the input system@>;
print_ln(banner); {print a ``banner line''}
@<Phase I: Read all the user's text and compress it into |tok_mem|@>;
stat for ii:=0 to zz-1 do max_tok_ptr[ii]:=tok_ptr[ii];@+tats@;@/
@<Phase II:...@>;
end_of_TANGLE:
if string_ptr>256 then @<Finish off the string pool file@>;
stat @<Print statistics about memory usage@>;@+tats@;@/
@t\4\4@>{here files should be closed if the operating system requires it}
@<Print the job |history|@>;
end.
@y
@p begin initialize;
@<Initialize the input system@>;
print_ln(banner); {print a ``banner line''}
@<Phase I: Read all the user's text and compress it into |tok_mem|@>;
stat for ii:=0 to zz-1 do max_tok_ptr[ii]:=tok_ptr[ii];@+tats@;@/
@<Phase II:...@>;
if string_ptr>256 then @<Finish off the string pool file@>;
stat @<Print statistics about memory usage@>;@+tats@;@/
@t\4\4@>{here files should be closed if the operating system requires it}
@<Print the job |history|@>;
end.
@z
