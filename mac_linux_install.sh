#!/bin/sh
DIR=$HOME/.gitdo
mkdir $DIR
echo "Copying plugins..."
echo "Copying hooks..."
cp -r ./plugins $DIR/
cp -r ./hooks $DIR/

GOOS=`uname | awk '{print tolower($0)}'`
MACHINE_TYPE=`uname -m`
echo ${GOOS}

GOARCH=``
if [ ${MACHINE_TYPE} = 'x86_64' ]; then
	GOARCH="64"
else
	GOARCH="32"
fi


VERSIONTOCP="gitdo_${GOOS}_${GOARCH}"
echo "Copying $VERSIONTOCP to your /usr/local/bin/ ..."
sudo cp $VERSIONTOCP /usr/local/bin/gitdo
sudo chmod +x /usr/local/bin/gitdo

echo "Done."
