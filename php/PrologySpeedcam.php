<?php
/**
 * @author  eku
 * Date: 30.06.15
 * Time: 20:49
 */

class PrologySpeedcam {

    public $file;

    public function __construct($fileName)
    {
        $this->file = fopen('speedcam.bin', 'r');
    }

    public function getHeader()
    {
        $str = $this->read(0, 20);
        $headers = unpack('Ia/Ifirst/Isecond/Ib/Ic', $str);
        return $headers;
    }

    protected function read($position, $length)
    {
        rewind($this->file);
        fseek($this->file, $position);
        $result = fread($this->file, $length);
        rewind($this->file);
        return $result;
    }

    protected function getRaw($str)
    {
        return join(' ', str_split(unpack('H*', $str)[1], 2));
    }

    public function getFirst($n)
    {
        $header = $this->getHeader();
        if ($n > $header['first']) {
            return null;
        }

        $position = 20 + 10 * $n;

        return $this->parseFirst($position);
    }

    public function parseFirst($position)
    {
        $str = $this->read($position, 10);
        $unp1 = unpack('sa/Ib', $str);
        $unp = unpack('x/x/x/x/x/x/Iposition', $str);
        $unp['raw'] = $this->getRaw($str);
        return array_merge($unp1, $unp);
    }

    public function getSecond($n)
    {
        $header = $this->getHeader();
        if ($n > $header['second']) {
            return null;
        }
        $position = 20 + ($header['first'] * 10) + ($n * 13);
        return $this->parseSecond($position);
    }

    public function parseSecond($position)
    {
        $str = $this->read($position, 13);
        $unp = unpack('ImLon/ImLat/sangle/cdir/ctype/cspeed', $str);
        $unp['mLon'] /= 10000;
        $unp['mLat'] /= 10000;
        return $unp;
    }

}