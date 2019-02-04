const Ion = artifacts.require("Ion");
const BaseValidation = artifacts.require("BaseValidation");
const Token = artifacts.require("Token");
const FabricStore = artifacts.require("FabricStore");
const ShareSettle = artifacts.require("ShareSettle");

module.exports = async (deployer) => {
  try {
      deployer.deploy(Ion, "0xab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac8075")
      .then(() => Ion.deployed)
      .then(() => deployer.deploy(BaseValidation, Ion.address))
      .then(() => BaseValidation.deployed)
      .then(() => deployer.deploy(FabricStore, Ion.address))
      .then(() => FabricStore.deployed)
      .then(() => deployer.deploy(Token))
      .then(() => Token.deployed)
      .then(() => deployer.deploy(ShareSettle, Token.address, FabricStore.address))
      .then(() => ShareSettle.deployed)
  } catch(err) {
    console.log('ERROR on deploy:',err);
  }

};