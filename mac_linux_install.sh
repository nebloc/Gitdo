#!/bin/sh
DIR=$HOME/.gitdo
mkdir $DIR
echo "Copying plugins..."
echo "Copying hooks..."
cp -r ./plugins $DIR
cp -r ./hooks $DIR

GOOS=`uname`

MACHINE_TYPE=`uname -m`
GOARCH=``
if [ ${MACHINE_TYPE} == 'x86_64' ]; then
	GOARCH="64"
else
	GOARCH="32"
fi


VERSIONTOCP="gitdo_${GOOS}_${GOARCH}"
echo "Copying $VERSIONTOCP to your /usr/local/bin/ ..."
cp $VERSIONTOCP /usr/local/bin/gitdo

echo "Done."
