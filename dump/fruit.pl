#!/usr/bin/perl
use strict;
use warnings;

#perl fruit.pl fruits.html fruity.txt
#pulls all the items from a list 
#reference URL https://www.liveeatlearn.com/list-of-fruits/ 

# Check if a file was provided
die "Usage: $0 filename.html\n" unless @ARGV == 1;

my $filename = $ARGV[0];

# Read the entire HTML file
open my $fh, '<', $filename or die "Could not open '$filename': $!";
local $/;  # Enable slurp mode
my $html = <$fh>;
close $fh;

# Find all <li>...</li> elements
while ($html =~ m|<li[^>]*>(.*?)</li>|gis) {
    my $item = $1;
    $item =~ s/<[^>]+>//g;  # Optionally strip nested HTML tags
    $item =~ s/^\s+|\s+$//g;  # Trim whitespace
    print "$item\n" if $item ne '';
}

