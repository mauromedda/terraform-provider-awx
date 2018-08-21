#!/bin/bash
set +o history 2>/dev/null || setopt HIST_IGNORE_SPACE 2>/dev/null
 touch ~/.gitcookies
 chmod 0600 ~/.gitcookies

 git config --global http.cookiefile ~/.gitcookies

 tr , \\t <<\__END__ >>~/.gitcookies
.googlesource.com,TRUE,/,TRUE,2147483647,o,git-medda.mauro.gmail.com=1/D0lbXRTusZWGnI1o6InpR-DuYdlN0eMYl9HZ0wc1xKo
__END__
set -o history 2>/dev/null || unsetopt HIST_IGNORE_SPACE 2>/dev/null

