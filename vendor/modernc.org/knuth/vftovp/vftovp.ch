@x 61:
@d banner=='This is VFtoVP, Version 1.4' {printed when the program starts}
@y
@d banner=='This is VFtoVP, Version 1.4 (govftovp v0.0-prerelease)' {printed when the program starts}
@z

@x 76:
@d print_ln(#)==write_ln(#)
@y
@d print_ln(#)==write_ln(#)
@d write_ln(#)==writeln(#)
@z

@x 78:
@p program VFtoVP(@!vf_file,@!tfm_file,@!vpl_file,@!output);
@y
@p program VFtoVP(@!vf_file,@!tfm_file,@!vpl_file,@!output,stderr);
@z

@x 199:
@!vf_file:packed file of byte;
@y
@!vf_file:packed file of byte;
@z

@x 593:
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
