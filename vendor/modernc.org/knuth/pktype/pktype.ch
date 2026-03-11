@x pktype:67:
@d banner=='This is PKtype, Version 2.3' {printed when the program starts}
@y
@d banner=='This is PKtype, Version 2.3 (gopktype v0.0-prerelease)' {printed when the program starts}
@z

@x pktype:67:
@d print_ln(#)==write_ln(output,#)
@y
@d print_ln(#)==write_ln(output,#)
@d write_ln(#)==writeln(#)
@d read_ln(#)==readln(#)
@z

@x pktype:72:
@p program PKtype(@!input,@!output);
label @<Labels in the outer block@>@/
@y
@p program PKtype(@!pk_file,@!typ_file,@!output);
@z

@x pktype:88:
@<Labels...@>=final_end;
@y
@z

@x pktype:109:
@d abort(#)==begin print_ln(' ',#); t_print_ln(' ',#); jump_out; end
@y
@d abort(#)==begin write_ln(stderr,' ',#); t_print_ln(' ',#); jump_out; end
@z

@x pktype:111:
@p procedure jump_out;
begin goto final_end;
end;
@y
@p procedure jump_out;
begin
	panic(final_end);
end;
@z


@x pktype:754:
@p procedure open_pk_file; {prepares the input for reading}
begin reset(pk_file,pk_name);
pk_loc := 0 ;
end;
@#
procedure open_typ_file; {prepares to write text data to the |typ_file|}
begin rewrite(typ_file,typ_name);
end;
@y
@p procedure open_pk_file; {prepares the input for reading}
begin reset(pk_file);
pk_loc := 0 ;
end;
@#
procedure open_typ_file; {prepares to write text data to the |typ_file|}
begin rewrite(typ_file);
end;
@z

@x pktype:784:
@<Open files@>=
open_pk_file ;
open_typ_file ;
t_print_ln(banner) ;
t_print('Input file: ') ;
i := 1 ;
while pk_name[i] <> ' ' do begin
   t_print(pk_name[i]) ; incr(i) ;
end ;
t_print_ln(' ')
@y
@<Open files@>=
open_pk_file ;
open_typ_file ;
@z

@x pktype:1109:
@ @p procedure dialog ;
var i : integer ; {index variable}
buffer : packed array [1..name_length] of char; {input buffer}
begin
   for i := 1 to name_length do begin
      typ_name[i] := ' ' ;
      pk_name[i] := ' ' ;
   end;
   print('Input file name:  ') ;
   flush_buffer ;
   get_line(pk_name) ;
   print('Output file name:  ') ;
   flush_buffer ;
   get_line(typ_name) ;
end ;
@y
@ @p procedure dialog ;
begin
end ;
@z

@x pktype:1146:
t_print_ln(pk_loc:1,' bytes read from packed file.');
final_end :
end .
@y
t_print_ln(pk_loc:1,' bytes read from packed file.');
end .
@z
