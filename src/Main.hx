package;

import Console.*;

/**
 * ...
 * @author Kyrasuum
 */
class Main {
    public function new() {
        proc_args();
    }

    private function proc_args() {
        var options = [
            ["--{command}", "-{letter}", "description", "default"   ],
            ["file",        "f",  "file(s) to open on launch",    ""],
            ["directory",   "d",  "working directory",            ""]
        ];
    
        for (arg in Sys.args()){
            switch (arg) {
                case "-h": {
                    println("Strangelet provides the following command line arguements:");
                    print_columns(options);
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

    public function start() {
    }
    
    static function main() {
        var app = new Main();
        app.start();
    }
}
