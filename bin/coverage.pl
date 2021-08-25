#!/usr/bin/env perl
#
# scan the tests and templates and tell us which territories
# don't yet have rules and tests
# 
use strict;
use warnings;
use utf8;
use feature "unicode_strings";

use Data::Dumper;
use Getopt::Long;
use File::Basename qw(dirname);

my $help    = 0;
my $details = 0;
GetOptions (
    'details' => \$details,
    'help'    => \$help,
) or die "invalid options";

if ($help) {
    usage();
    exit(0);
}

# get the list of countries
my %countries;
my $country_file = dirname(__FILE__) . "/../conf/country_codes.yaml";
open my $FH, "<:encoding(UTF-8)", $country_file or die "unable to open $country_file $!";
while (my $line = <$FH>){
    chomp($line);
    if ($line =~ m/^(\w\w): \# (.*)$/){
        $countries{$1} = $2;
    }
}
close $FH;
my $total_countries = scalar(keys %countries); 
print "We are aware of " . $total_countries . " territories \n";


# which countries have tests?
my $test_dir = dirname(__FILE__) . '/../testcases/countries';
opendir(my $dh, $test_dir) || die "Error: Couldn't opendir($test_dir): $!\n";
my @files = grep { -f "$test_dir/$_" } readdir($dh);
closedir($dh);

my %test_countries;
foreach my $f (sort @files){
    $f =~ s/\.yaml//;
    $f = uc($f);
    $test_countries{$f} = 1;
}
my $test_countries = scalar(keys %test_countries);
my $test_perc = int(100 * $test_countries / $total_countries );
print "We have tests for " .  $test_countries . ' ('
    . $test_perc . '%) territories' .  "\n";
if ($details){
    print "We need tests for:\n";
    foreach my $cc (sort keys %countries){
        next if (defined($test_countries{$cc}));
        print "\t" . $cc . "\t". $countries{$cc}. "\n";
    }
}

# which countries have rules?
my $rules_file = dirname(__FILE__) . '/../conf/countries/worldwide.yaml';
open my $RFH, "<:encoding(UTF-8)", $rules_file or die "unable to open $rules_file $!";
my %rules;
while (my $line = <$RFH>){
    chomp($line);
    if ($line =~ m/^"?(\w\w)"?:\s*$/){
        $rules{$1} = 1;
    }
}
close $RFH;
my $rules_countries = scalar(keys %rules);
my $rules_perc = int(100 * $rules_countries / $total_countries );
print "We have rules for " .  $rules_countries . ' ('
    . $rules_perc . '%) territories' . "\n";

if ($details){
    print "We need rules for:\n";
    foreach my $cc (sort keys %countries){
        next if (defined($rules{$cc}));
        print "\t" . $cc . "\t". $countries{$cc}. "\n";
    }
}

# find territories without rules or tests
my %neither;
foreach my $cc (sort keys %countries){
    next if (defined($rules{$cc}));
    next if (defined($test_countries{$cc}));
    $neither{$cc} = 1;
}
my $neither_countries = scalar(keys %neither);
my $neither_perc = int(100 * $neither_countries / $total_countries );
print $neither_countries . ' (' . $neither_perc . '%) territories have neither rules nor tests' . "\n";
if ($details){
    print "Territories with no test and no rules:\n";
    foreach my $cc (sort keys %neither){
        print "\t" . $cc . "\t". $countries{$cc}. "\n";
    }
}


sub usage {
    print "\tHow many territories have formatting rules and tests?\n";
    print "\tBy default prints just a high level summary\n";
    print "\tusage:\n";
    print "\t\t no required parameters\n";
    print "\n";
    print "\t\t optional parameters:\n";
    print "\t\t --detail\t print full list of countries missing rules and tests\n";
    print "\t\t --help\t print this message \n";
    print "\n";
}


