// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ITreasury} from "./ITreasury.sol";
import {INexus} from "./INexus.sol";

interface ITreasuryNative is ITreasury, INexus {
    function depositNative() external payable;
    function requestNative(uint256 amount) external;
    function approveNative(address user, uint256 amount) external;
    function approveNativeAll(address user) external;
    function withdrawNative(uint256 amount) external;
    function withdrawNativeAll() external;
    function nativeBalance() external view returns (uint256);
    function nativeAvailable() external view returns (uint256);
    function nativeAllocation(address user) external view returns (uint256);
    function nativePending(address user) external view returns (uint256);
}