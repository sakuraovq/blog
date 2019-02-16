<?php
$arra = [3,56,71,1,2];


for ($i = 0;$i<count($arra);$i++){
	// 因为下标从0开始-1
	for ($j=0;$j<count($arra)-$i-1;$j++){
		
		if($arra[$j]>$arra[$j+1]){
			$item = $arra[$j+1];
			$arra[$j+1] = $arra[$j];
			$arra[$j] = $item;
			
		}
	}
	
}
var_dump($arra);
