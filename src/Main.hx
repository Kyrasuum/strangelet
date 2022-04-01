package;

import Console.*;
var blessed = js.Node.require('reblessed');

/**
 * ...
 * @author Kyrasuum
 */
class Main {
    //Our blessed screen object.
    var screen = blessed.screen({
      smartCSR: true
    });
    
    var log: log.Log;
    var fbr: fbr.Fbr;

    public function new() {
        proc_args();
        init_screen();

        log = new log.Log(screen);
        fbr = new fbr.Fbr(screen);
    }

    public function exit() {
        Sys.exit(0);
    }

    private function proc_args() {
        //help text
        var options = [
            ["--{long-command}", "-{short-command}", "description", "default"],
            ["file",        "f",  "file(s) to open on launch",      ""],
            ["directory",   "d",  "working directory",              ""]
        ];

        //switch on cli arguements
        for (arg in Sys.args()){
            switch (arg) {
                case "-h", "--help": {
                    println("Strangelet provides the following command line arguements:");
                    print_columns(options);
                    exit();
                }
                default: {
                    println(arg);
                }
            }
        }
    }

    private function print_columns(lines: Array<Array<String>>) {
        //construct column character max widths
        var col_widths: Array<Int> = [];

        for (i in 0...lines.length){
            var line = lines[i];
            for (j in 0... line.length){
                var col = line[j];
                if (i == 0){
                    col_widths.push(col.length);
                } else {
                    if (j > col_widths.length) {
                        col_widths.push(col.length);
                    } else if (col.length > col_widths[j]) {
                        col_widths[j] = col.length;
                    }
                }
            }
        }

        //print formatted columns
        for (i in 0...lines.length){
            var line = lines[i];
            for (j in 0... line.length){
                var col = line[j];

                print(col);
                for (k in col.length... col_widths[j]){
                    print(" ");
                }
                print("\t");
            }
            println("");
        }
    }

    public function init_screen() {
        screen.title = 'Strangelet';

        setup_keys();
    }

    public function setup_keys() {
        //Quit on Control-q
        screen.key(['C-q'], function(ch, key) {
            exit();
        });
    }

    public function start() {
    }
    
	public static function main() {
        var app = new Main();
        app.start();
    }
}
