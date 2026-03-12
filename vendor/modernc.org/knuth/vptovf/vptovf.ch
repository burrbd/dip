@x 59
@d banner=='This is VPtoVF, Version 1.6' {printed when the program starts}
@y
@d banner=='This is VPtoVF, Version 1.6 (govptovf v0.0-prerelease)' {printed when the program starts}
@z

@x 75
@d print_ln(#)==write_ln(#)
@y
@d print_ln(#)==write_ln(#)
@d write_ln(#)==writeln(#)
@d read_ln(#)==readln(#)
@z

@x 77:
@p program VPtoVF(@!vpl_file,@!vf_file,@!tfm_file,@!output);
@y
@p program VPtoVF(@!vpl_file,@!vf_file,@!tfm_file,@!output,stderr);
@z
