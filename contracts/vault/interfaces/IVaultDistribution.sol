// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface IVaultDistribution {
    error NothingToSpread(address account);
    error AmountTooSmall(uint256 amount, uint256 recipientCount);
    error NothingToCollect();

    event Spread(address indexed from, uint256 indexed totalAmount, uint256 indexed share, uint256 remainder);
    event Collected(address indexed by, uint256 totalAmount);

    function spread() external;
    function collect() external;
}