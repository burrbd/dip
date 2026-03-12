@x 67:
@d banner=='This is TFtoPL, Version 3.3' {printed when the program starts}
@y
@d banner=='This is TFtoPL, Version 3.3 (gotftopl v0.0-prerelease)' {printed when the program starts}
@z

@x 81:
@d print_ln(#)==write_ln(#)
@y
@d print_ln(#)==write_ln(#)
@d write_ln(#)==writeln(#)
@z

@x 83:
@p program TFtoPL(@!tfm_file,@!pl_file,@!output);
@y;
@p program TFtoPL(@!tfm_file,@!pl_file,@!output,stderr);
@z

@x 431:
@d abort(#)==begin print_ln(#);
  print_ln('Sorry, but I can''t go on; are you sure this is a TFM?');
  goto final_end;
  end
@y
@d abort(#)==begin print_ln(stderr,#);
  print_ln(stderr,'Sorry, but I can''t go on; are you sure this is a TFM?');
  goto final_end;
  end
@z
