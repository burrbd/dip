@x tangle.web:83:
@d banner=='This is TANGLE, Version 4.6'
@y
@d banner=='This is TANGLE, Version 4.6 (gotangle v0.0-prereleaase)'
@z

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

@x tangle.web:196:
@!max_id_length=12; {long identifiers are chopped to this length, which must
  not exceed |line_length|}
@!unambig_length=7; {identifiers must be unique if chopped to this length}
  {note that 7 is more strict than \PASCAL's 8, but this can be varied}
@y
@!max_id_length=32; {long identifiers are chopped to this length, which must
  not exceed |line_length|}
@!unambig_length=32; {identifiers must be unique if chopped to this length}
@z

@x tangle.web:509:
@d print(#)==write(term_out,#) {`|print|' means write on the terminal}
@d print_ln(#)==write_ln(term_out,#) {`|print|' and then start new line}
@d new_line==write_ln(term_out) {start new line}
@d print_nl(#)==  {print information starting on a new line}
  begin new_line; print(#);
  end
@y
@d print(#)==write(#) {`|print|' means write on the terminal}
@d print_ln(#)==write_ln(#) {`|print|' and then start new line}
@d new_line==write_ln() {start new line}
@d print_nl(#)==  {print information starting on a new line}
  begin new_line; print(#);
  end
@d write_ln(#)==writeln(#)
@d read_ln(#)==readln(#)
@z

@x tangle.web:525:
rewrite(term_out,'TTY:'); {send |term_out| output to the terminal}

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


@x tangle.web:3300:
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
end {there are no other cases}
;write_ln(output)
@z
