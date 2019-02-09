<?php
/**
 * @author  eku
 * Date: 30.06.15
 * Time: 21:02
 */


require_once('PrologySpeedcam.php');


$psc = new PrologySpeedcam('speedcam.bin');

$res = $psc->getFirst(0);
print_r($res);
print_r($psc->parseSecond($res['position']));


for ($i = 0; $i < 100; $i++) {
    $res = $psc->getFirst($i);
    $point = $psc->parseSecond($res['position']);
    echo "{$point['mLon']}, {$point['mLat']}\t{$point['type']}\n";
}

