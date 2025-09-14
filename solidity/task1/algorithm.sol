// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

contract Algorithm {

    function ReverseString(string calldata str) external pure returns (string memory) {
        bytes calldata bytesStr = bytes(str);
        bytes memory result = new bytes(bytesStr.length);
        for (uint i = 0; i < bytesStr.length; i++) {
            result[i] = bytesStr[bytesStr.length - 1 - i];
        }
        return string(result);
    }

    function BinarySearch(int[] calldata sortArr, int findVal) external pure returns (int) {
        uint left;
        uint right = sortArr.length;
        while (right > left){
            uint mid = (right - left) / 2;
            if (sortArr[mid + left] == findVal ){
                return int(mid + left);
            }
            if (mid == 0){
                return -1;
            }
            if (sortArr[mid + left] < findVal){
                left = mid + left + 1;
            } else {
                right = right - mid;
            }
        }
        return -1;
    }


    function mergeArray(int[] calldata arr1, int[] calldata arr2) external  pure returns (int[] memory) {
        uint length = arr1.length + arr2.length;
        uint arr1Ptr;
        uint arr2Ptr;
        int[] memory result = new int[](length);
        for (uint i = 0; i < length; i++) {
            if ( arr1Ptr < arr1.length && arr2Ptr < arr2.length ) {
                if (arr1[arr1Ptr] > arr2[arr2Ptr]){
                    result[i] = arr2[arr2Ptr];
                    arr2Ptr ++;
                } else {
                    result[i] = arr1[arr1Ptr];
                    arr1Ptr ++;
                }
                continue;
            }
            if (arr1Ptr < arr1.length){
                result[i] = arr1[arr1Ptr];
                arr1Ptr ++;
            } else {
                result[i] = arr2[arr2Ptr];
                arr2Ptr ++;
            }
        }
        return result;
    }

    mapping(bytes1 => int) private romancharMap;
    mapping(uint => string) private romanNumMap;
    uint[] private romanNums = [1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1];
    constructor() {
        romancharMap['I'] = 1;
        romancharMap['V'] = 5;
        romancharMap['X'] = 10;
        romancharMap['L'] = 50;
        romancharMap['C'] = 100;
        romancharMap['D'] = 500;
        romancharMap['M'] = 1000;

        romanNumMap[1] = "I";
        romanNumMap[4] = "IV";
        romanNumMap[5] = "V";
        romanNumMap[9] = "IX";
        romanNumMap[10] = "X";
        romanNumMap[40] = "XL";
        romanNumMap[50] = "L";
        romanNumMap[90] = "XC";
        romanNumMap[100] = "C";
        romanNumMap[400] = "CD";
        romanNumMap[500] = "D";
        romanNumMap[900] = "CM";
        romanNumMap[1000] = "M";

    }
    function romanToInt(string calldata s) external view returns (int) {
        bytes memory b = bytes(s);
        int result;
        for  (uint8 i = 0; i < b.length; i++) {
            bytes1 key = b[i];
            int num = romancharMap[key];
            require(num != 0, 'input error');
            if (i + 1 < b.length && num < romancharMap[b[i+1]]){
                result -= num;
            }
            result += num;
        }
        return result;
    }

    function intToRoman(uint num) external view returns (string memory) {
        string memory result;
        while (num>0){
            for (uint i = 0; i < romanNums.length; i++) {
                if (num >= romanNums[i]) {
                    num -= romanNums[i];
                    result = string.concat(result, romanNumMap[romanNums[i]]);
                    break;
                }
            }
        }
        return result;
    }

}