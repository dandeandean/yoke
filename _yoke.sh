#!/usr/bin/env bash

function _yoke() {
  COMPREPLY=($(yoke complete $COMP_LINE));
};

complete -F _yoke yoke
