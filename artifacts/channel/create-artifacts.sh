
chmod -R 0755 ./crypto-config
# Delete existing artifacts
rm -rf ./crypto-config
rm genesis.block mychannel.tx channel1.tx channel2.tx
rm -rf ../../channel-artifacts/*

#Generate Crypto artifactes for organizations
cryptogen generate --config=./crypto-config.yaml --output=./crypto-config/



# System channel
SYS_CHANNEL="sys-channel"

# channel name defaults to "mychannel"
CHANNEL_NAME="mychannel"
CHANNEL_NAME1="channel1"
CHANNEL_NAME2="channel2"

echo $CHANNEL_NAME

# Generate System Genesis block
configtxgen -profile OrdererGenesis -configPath . -channelID $SYS_CHANNEL  -outputBlock ./genesis.block


# Generate channel configuration block
configtxgen -profile BasicChannel -configPath . -outputCreateChannelTx ./$CHANNEL_NAME.tx -channelID $CHANNEL_NAME
configtxgen -profile BasicChannel -configPath . -outputCreateChannelTx ./$CHANNEL_NAME1.tx -channelID $CHANNEL_NAME1
configtxgen -profile BasicChannel -configPath . -outputCreateChannelTx ./$CHANNEL_NAME2.tx -channelID $CHANNEL_NAME2

echo "#######    Generating anchor peer update mychannel  ##########"
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./Org2MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org2MSP
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./Org3MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org3MSP

configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./Org1MSPanchorsa.tx -channelID $CHANNEL_NAME1 -asOrg Org1MSP
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./Org2MSPanchorsa.tx -channelID $CHANNEL_NAME1 -asOrg Org2MSP


configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./Org1MSPanchorsb.tx -channelID $CHANNEL_NAME2 -asOrg Org1MSP
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./Org3MSPanchorsb.tx -channelID $CHANNEL_NAME2 -asOrg Org3MSP