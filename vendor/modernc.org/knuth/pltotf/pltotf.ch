@x 69:
@d banner=='This is PLtoTF, Version 3.6' {printed when the program starts}
@y
@d banner=='This is PLtoTF, Version 3.6 (gopltotf v0.0-prerelease)' {printed when the program starts}
@z

@x 84:
@d print_ln(#)==write_ln(#)
@y
@d print_ln(#)==write_ln(#)
@d write_ln(#)==writeln(#)
@d read_ln(#)==readln(#)
@z

@x 86:
@p program PLtoTF(@!pl_file,@!tfm_file,@!output);
@y
@p program PLtoTF(@!pl_file,@!tfm_file,@!output,stderr);
@z
