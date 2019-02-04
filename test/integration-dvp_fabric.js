// Copyright (c) 2016-2019 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+

const Web3Utils = require('web3-utils');
const rlp = require('rlp');

// Connect to the Test RPC running
const Web3 = require('web3');
const web3 = new Web3();
web3.setProvider(new web3.providers.HttpProvider('http://localhost:8545'));

const Ion = artifacts.require("Ion");
const ShareSettle = artifacts.require("ShareSettle");
const Token = artifacts.require("Token");
const FabricStore = artifacts.require("FabricStore");
const BaseValidation = artifacts.require("BaseValidation");

require('chai')
 .use(require('chai-as-promised'))
 .should();

const DEPLOYEDCHAINID = "0xab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac8075"
const TESTCHAINID = "0x22b55e8a4f7c03e1689da845dd463b09299cb3a574e64c68eafc4e99077a7254"

let TRANSFER_DATA = [{
    channelId: "sharechannel",
    blocks: [{
        hash: "Qxdh5WEYE9VfjiCpRVihDPd96ViDS6wNanLJSp3D_uI",
        number: 7,
        prevHash: "-AiL5Ej5nPqbkOa3d0FHwjTMqhL5TRKzhyfjBX7UKFQ",
        dataHash: "VqyCgRZOu2ElVUtYD_c_-wljyMgNWA9Vd4DtwcaCZ58",
        timestampS: 1549011931,
        timestampN: 724233992,
        transactions: [{
            txId: "6d3c8ab96271aebf4412d532f87123ac4187725d4442c108e01d5deb151fcab8",
            nsrw: [{
                namespace: "Shares",
                readsets: [{
                    key: "orders", 
                    version: {
                        blockNumber: 6,
                        txNumber: 0
                    }
                },
                {
                    key: "shares", 
                    version: {
                        blockNumber: 4,
                        txNumber: 0
                    }
                },
                {
                    key: "trades", 
                    version: {
                        blockNumber: 3,
                        txNumber: 0
                    }
                }],
                writesets: [{
                    key: "shares",
                    isDelete: "false",
                    value: "0xf842f840f83e87496f6e436f72708203e8f2d894dbad2ab2ff31d2823089447298e1df08dec457bd8201f4d894280d9c579809eda6c0be390fb3c0943689be2d9b8201f4"
                },
                {
                    key: "trades",
                    isDelete: "false",
                    value: "0xf84ef84cf84a89756e69717565526566f83e87496f6e436f727094dbad2ab2ff31d2823089447298e1df08dec457bd94280d9c579809eda6c0be390fb3c0943689be2d9b056489756e69717565526566"
                }]
            }, {
                namespace: "lscc",
                readsets: [{
                    key: "Shares",
                    version: {
                        blockNumber: 3,
                        txNumber: 0
                    } 
                }],
                writesets: []
            }]
        }]
    }]
}]

// Create a formatted block from the example block
createData = (DATA) => {
    const formattedData = [[
        DATA[0].channelId,
        [
            DATA[0].blocks[0].hash,
            DATA[0].blocks[0].number,
            DATA[0].blocks[0].prevHash,
            DATA[0].blocks[0].dataHash,
            DATA[0].blocks[0].timestampS,
            DATA[0].blocks[0].timestampN,
            [
                [
                    DATA[0].blocks[0].transactions[0].txId,
                    [
                        [
                            DATA[0].blocks[0].transactions[0].nsrw[0].namespace,
                            [
                                [
                                    DATA[0].blocks[0].transactions[0].nsrw[0].readsets[0].key,
                                    [
                                        DATA[0].blocks[0].transactions[0].nsrw[0].readsets[0].version.blockNumber,
                                        DATA[0].blocks[0].transactions[0].nsrw[0].readsets[0].version.txNumber
                                    ]
                                ],
                                [
                                    DATA[0].blocks[0].transactions[0].nsrw[0].readsets[1].key,
                                    [
                                        DATA[0].blocks[0].transactions[0].nsrw[0].readsets[1].version.blockNumber,
                                        DATA[0].blocks[0].transactions[0].nsrw[0].readsets[1].version.txNumber
                                    ]
                                ],
                                [
                                    DATA[0].blocks[0].transactions[0].nsrw[0].readsets[2].key,
                                    [
                                        DATA[0].blocks[0].transactions[0].nsrw[0].readsets[2].version.blockNumber,
                                        DATA[0].blocks[0].transactions[0].nsrw[0].readsets[2].version.txNumber
                                    ]
                                ]
                            ],
                            [
                                [
                                    DATA[0].blocks[0].transactions[0].nsrw[0].writesets[0].key,
                                    DATA[0].blocks[0].transactions[0].nsrw[0].writesets[0].isDelete,
                                    DATA[0].blocks[0].transactions[0].nsrw[0].writesets[0].value  
                                ],
                                [
                                    DATA[0].blocks[0].transactions[0].nsrw[0].writesets[1].key,
                                    DATA[0].blocks[0].transactions[0].nsrw[0].writesets[1].isDelete,
                                    DATA[0].blocks[0].transactions[0].nsrw[0].writesets[1].value
                                ]
                            ]
                        ],
                        [
                            DATA[0].blocks[0].transactions[0].nsrw[1].namespace,
                            [
                                [
                                    DATA[0].blocks[0].transactions[0].nsrw[1].readsets[0].key,
                                    [
                                            DATA[0].blocks[0].transactions[0].nsrw[1].readsets[0].version.blockNumber,
                                            DATA[0].blocks[0].transactions[0].nsrw[1].readsets[0].version.txNumber
                                    ]
                                ]
                            ],
                            []
                        ]
                    ]
                ]
            ]
        ]
    ]];
    return formattedData;
  }

const cliRlpBlock = "0xf901e2f901df8c73686172656368616e6e656cf901cfab5178646835574559453956666a6943705256696844506439365669445336774e616e4c4a537033445f754907ab2d41694c35456a356e5071626b4f613364304648776a544d71684c3554524b7a6879666a425837554b4651ab5671794367525a4f7532456c56557459445f635f2d776c6a794d674e5741395664344474776361435a3538845c540bdb842b2aef08f9013df9013ab84036643363386162393632373161656266343431326435333266383731323361633431383737323564343434326331303865303164356465623135316663616238f8f6f8e186536861726573e1ca866f7264657273c20680ca86736861726573c20480ca86747261646573c20380f8b6f853867368617265738566616c7365b844f842f840f83e87496f6e436f72708203e8f2d894dbad2ab2ff31d2823089447298e1df08dec457bd8201f4d894280d9c579809eda6c0be390fb3c0943689be2d9b8201f4f85f867472616465738566616c7365b850f84ef84cf84a89756e69717565526566f83e87496f6e436f727094dbad2ab2ff31d2823089447298e1df08dec457bd94280d9c579809eda6c0be390fb3c0943689be2d9b056489756e69717565526566d2846c736363cbca86536861726573c20380c0"

contract('Fabric SimpleShares Integration', (accounts) => {
    let ion;
    let validation;
    let storage;

    // Update transfer data with addresses of accounts
    let updateValue = TRANSFER_DATA[0].blocks[0].transactions[0].nsrw[0].writesets[0].value;
    updateValue = updateValue.replace("e7cf944311eabff15b1b091422a2ecada1dd053d", accounts[0].toString().slice(2));
    updateValue = updateValue.replace("9ecd4a8ca1560c4bb92ca9ebfa5eab448048db93", accounts[1].toString().slice(2));
    TRANSFER_DATA[0].blocks[0].transactions[0].nsrw[0].writesets[0].value = updateValue;
    const formattedData = createData(TRANSFER_DATA);

    let rlpEncodedBlock = "0x" + rlp.encode(formattedData).toString('hex');

    beforeEach('setup contract for each test', async function () {
        ion = await Ion.new(DEPLOYEDCHAINID);
        validation = await BaseValidation.new(ion.address);
        storage = await FabricStore.new(ion.address);
    })


    describe('Chaincode usage Contract', () => {
        const Bob = accounts[0]
        const Alice = accounts[1]


        it('Submit Block, retrieve state and execute', async () => {
            await validation.register();
            const token = await Token.new()
            const shareSettle = await ShareSettle.new(token.address, storage.address)

            const value = 5;
            const price = 100;
            const reference = Web3Utils.sha3('uniqueRef');

            // Mint ERC223 tokens, funding Bob
            await token.mint(value*price, {from: Bob});

            // Bob initiates a trade agreement which it will settle later
            let tx = await shareSettle.initiateTrade(
                "IonCorp",
                Alice,
                value,
                price,
                reference,
                {
                    from: Bob
                },
            );
            console.log("\tGas used to initiate trade: " + tx.receipt.gasUsed.toString());

            tx = await validation.RegisterChain(TESTCHAINID, storage.address);

            let receipt = await validation.SubmitBlock(TESTCHAINID, rlpEncodedBlock, storage.address);
            console.log("\tGas used to store fabric block: %d", receipt.receipt.gasUsed);

            // Escrow tokens in IonLock contract under specific reference
            let receiptTransfer = await token.metadataTransfer(
                shareSettle.address,
                value*price,
                reference,
                {
                    from: Bob
                },
            );
            console.log("\tGas used to escrow tokens: " + receiptTransfer.receipt.gasUsed.toString());

            let balance = await token.balanceOf(shareSettle.address);
            assert.equal(value*price, balance)
            tx = await shareSettle.retrieveAndExecute(TESTCHAINID, TRANSFER_DATA[0].channelId, TRANSFER_DATA[0].blocks[0].transactions[0].nsrw[0].writesets[1].key);
            console.log("\tGas used to retrieve and execute: " + tx.receipt.gasUsed.toString());
            
            assert.equal(500, await token.balanceOf(Alice));
            assert.equal(0, await token.balanceOf(Bob));
            assert.equal(0, await token.balanceOf(shareSettle.address));
            
        })

        it('Submit Block, encoded by Fabric CLI, retrieve state and execute', async () => {
            await validation.register();
            const token = await Token.new()
            const shareSettle = await ShareSettle.new(token.address, storage.address)

            const value = 5;
            const price = 100;
            const reference = Web3Utils.sha3('uniqueRef');
            console.log(reference)

            // Mint ERC223 tokens, funding Bob
            await token.mint(value*price, {from: Bob});

            // Bob initiates a trade agreement which it will settle later
            let tx = await shareSettle.initiateTrade(
                "IonCorp",
                Alice,
                value,
                price,
                reference,
                {
                    from: Bob
                },
            );
            console.log("\tGas used to initiate trade: " + tx.receipt.gasUsed.toString());

            tx = await validation.RegisterChain(TESTCHAINID, storage.address);

            let receipt = await validation.SubmitBlock(TESTCHAINID, cliRlpBlock, storage.address);
            console.log("\tGas used to store fabric block: %d", receipt.receipt.gasUsed);

            // Escrow tokens in IonLock contract under specific reference
            let receiptTransfer = await token.metadataTransfer(
                shareSettle.address,
                value*price,
                reference,
                {
                    from: Bob
                },
            );
            console.log("\tGas used to escrow tokens: " + receiptTransfer.receipt.gasUsed.toString());

            let balance = await token.balanceOf(shareSettle.address);
            assert.equal(value*price, balance)
            console.log(TESTCHAINID, TRANSFER_DATA[0].channelId, TRANSFER_DATA[0].blocks[0].transactions[0].nsrw[0].writesets[1].key);
            tx = await shareSettle.retrieveAndExecute(TESTCHAINID, TRANSFER_DATA[0].channelId, TRANSFER_DATA[0].blocks[0].transactions[0].nsrw[0].writesets[1].key);
            console.log("\tGas used to retrieve and execute: " + tx.receipt.gasUsed.toString());
            
            assert.equal(500, await token.balanceOf(Alice));
            assert.equal(0, await token.balanceOf(Bob));
            assert.equal(0, await token.balanceOf(shareSettle.address));
        })

    })

})