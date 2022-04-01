package log;

var blessed = js.Node.require('reblessed');

/**
 * ...
 * @author Kyrasuum
 */
class Log {
    var screen: Dynamic;

    public function new(scr: Dynamic) {
        screen = scr;

        init_gui();
    }

    public function init_gui() {
        //create a box for our ui element
        var box = blessed.box({
            top: 'center',
            left: '0px',
            width: '20%',
            height: '100%',
            content: 'Log',
            tags: true,
            style: {
                fg: 'white',
                bg: 'gray',
            }
        });

        //append our box to the screen.
        screen.append(box);

        
        //click handler
        box.on('click', function(data) {});

        //enter key handler (during focus)
        box.key('enter', function(ch, key) {});

        //capture focus
        box.focus();

        //redraw the screen
        screen.render();
    }
}
