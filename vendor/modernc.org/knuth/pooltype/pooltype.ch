@x 53
@p program POOLtype(@!pool_file,@!output);
@y
@p program POOLtype(@!pool_file,@!output,stderr);
@z

@x 284
for i:=0 to @'37 do xchr[i]:=' ';
for i:=@'177 to @'377 do xchr[i]:=' ';
@y
for i:=0 to @'37 do xchr[i]:=chr(i);
for i:=@'177 to @'377 do xchr[i]:=chr(i);
@z

@x 319:
@d abort(#)==begin write_ln(#); goto 9999;
@y
@d write_ln(#)==writeln(#)
@d read_ln(#)==readln(#)
@d abort(#)==begin write_ln(stderr, #); goto 9999;
@z

@x 335: Make the result compatible with pooltype from tex-live.
  begin write(k:3,': "'); l:=k;
@y
  begin write(k,': "'); l:=k;
@z

@x 405: Make the result compatible with pooltype from tex-live.
  write(s:3,': "'); count:=count+l;
@y
  write(s,': "'); count:=count+l;
@z
