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
@p program PKtype(@!input,@!output);
@z

@x pktype:88:
@<Labels...@>=final_end;
@y
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

@x pktype:773:
@p function pk_byte : eight_bits ;
var temp : eight_bits ;
begin
   temp := pk_file^ ;
   get(pk_file) ;
   incr(pk_loc) ;
   pk_byte := temp ;
end ;
@y
@p function pk_byte : eight_bits ;
var temp : eight_bits ;
begin
   read(pk_file, temp) ;
   incr(pk_loc) ;
   pk_byte := temp ;
end ;
@z

@x pktype:1099:
@d flush_buffer == begin end
@d get_line(#) == if eoln(input) then read_ln(input) ;
   i := 1 ;
   while not (eoln(input) or eof(input)) do begin
      #[i] := input^ ;
      incr(i) ;
      get(input) ;
   end ;
   #[i] := ' '
@y
@d flush_buffer == begin end
@d get_line(#) == if eoln(input) then read_ln(input) ;
   i := 1 ;
   while not (eoln(input) or eof(input)) do begin
      read(input, #[i]);
      incr(i) ;
   end ;
   #[i] := ' '
@z

@x pktype:1109:
@ @p procedure dialog ;
var i : integer ; {index variable}
buffer : packed array [1..name_length] of char; {input buffer}
begin
@y
@ @p procedure dialog ;
var i : integer ; {index variable}
begin
@z


@x pktype:1146:
t_print_ln(pk_loc:1,' bytes read from packed file.');
final_end :
end .
@y
t_print_ln(pk_loc:1,' bytes read from packed file.');
end .
@z
