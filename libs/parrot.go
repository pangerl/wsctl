// Package libs @Author lanpang
// @Date 2024/9/28 下午21:28:00
// @Desc
package libs

import (
	"bytes"
	"slices"
	"time"

	"github.com/nsf/termbox-go"
)

var (
	animation   = "parrot"
	delay       = 75
	loops       = 0
	frame_index = 0
	color_index = 0
)

func RunParrot(orientation string) {
	inventory := NewInventory()

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

	termbox.SetOutputMode(termbox.Output256)

	loop_index := 0
	draw(inventory[animation], orientation)

loop:
	for {
		select {
		case ev := <-event_queue:
			if (ev.Type == termbox.EventKey && (ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC || ev.Ch == 'q')) || ev.Type == termbox.EventInterrupt {
				break loop
			}
		default:
			loop_index++
			if loops > 0 && (loop_index/9) >= loops {
				break loop
			}
			draw(inventory[animation], orientation)
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
	}
}

func draw(animation Animation, orientation string) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	lines := bytes.Split(animation.Frames[frame_index], []byte{'\n'})

	if orientation == "aussie" {
		lines = reverse(lines)
	}

	for x, line := range lines {
		for y, cell := range line {
			termbox.SetCell(y, x, rune(cell), colors[color_index], termbox.ColorDefault)
		}
	}

	termbox.Flush()
	frame_index++
	color_index++
	if frame_index >= len(animation.Frames) {
		frame_index = 0
	}
	if color_index >= len(colors) {
		color_index = 0
	}
}

func reverse(lines [][]byte) [][]byte {
	slices.Reverse(lines)
	return lines
}

func NewInventory() Inventory {
	return Inventory{
		"parrot": Animation{
			Metadata: map[string]string{
				"description": "The classic Party Parrot.",
			},
			Frames: [][]byte{
				[]byte(`
  _                   _____                  
 | |                 |  __ \                 
 | |     __ _ _ __   | |__) |_ _ _ __   __ _ 
 | |    / _  | '_ \  |  ___/ _  | '_  \ / _ |
 | |___| (_| | | | | | |  | (_| | | | | (_| |
 |______\__,_|_| |_| |_|   \__,_|_| |_|\__, |
                                        __/ |
                                       |___/ 
                        .cccc;;cc;';c.
                      .,:dkdc:;;:c:,:d:.
                     .loc'.,cc::::::,..,:.
                   .cl;....;dkdccc::,...c;
                  .c:,';:'..ckc',;::;....;c.
                .c:'.,dkkoc:ok:;llllc,,c,';:.
               .;c,';okkkkkkkk:,lllll,:kd;.;:,.
               co..:kkkkkkkkkk:;llllc':kkc..oNc
             .cl;.,okkkkkkkkkkc,:cll;,okkc'.cO;
             ;k:..ckkkkkkkkkkkl..,;,.;xkko:',l'
            .,...';dkkkkkkkkkkd;.....ckkkl'.cO;
         .,,:,.;oo:ckkkkkkkkkkkdoc;;cdkkkc..cd,
      .cclo;,ccdkkl;llccdkkkkkkkkkkkkkkkd,.c;
     .lol:;;okkkkkxooc::loodkkkkkkkkkkkko'.oc
   .c:'..lkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkd,.oc
  .lo;,ccdkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkd,.c;
,dx:..;lllllllllllllllllllllllllllllllloc'...
cNO;........................................
`),
				[]byte(`                
  _                   _____                  
 | |                 |  __ \                 
 | |     __ _ _ __   | |__) |_ _ _ __   __ _ 
 | |    / _  | '_ \  |  ___/ _  | '_  \ / _ |
 | |___| (_| | | | | | |  | (_| | | | | (_| |
 |______\__,_|_| |_| |_|   \__,_|_| |_|\__, |
                                        __/ |
                                       |___/ 
               .ckx;'........':c.
             .,:c:c:::oxxocoo::::,',.
            .odc'..:lkkoolllllo;..;d,
            ;c..:o:..;:..',;'.......;.
           ,c..:0Xx::o:.,cllc:,'::,.,c.
           ;c;lkXXXXXXl.;lllll;lXXOo;':c.
         ,dc.oXXXXXXXXl.,lllll;lXXXXx,c0:
         ;Oc.oXXXXXXXXo.':ll:;'oXXXXO;,l'
         'l;;OXXXXXXXXd'.'::'..dXXXXO;,l'
         'l;:0XXXXXXXX0x:...,:o0XXXXk,:x,
         'l;;kXXXXXXKXXXkol;oXXXXXXXO;oNc
        ,c'..ckk2XXXXXXXXXX00XXXXXXX0:;o:.
      .':;..:dd::ooooOXXXXXXXXXXXXXXXo..c;
    .',',:co0XX0kkkxx0XXXXXXXXXXXXXXX0c..;l.
  .:;'..oXXXXXXXXXXXXXXXXXXXXXXXXXXXXXko;';:.
.cdc..:oOXXXXXXXXKXXXXXXXXXXXXXXXXXXXXXXo..oc
:0o...:dxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxo,.:,
cNo........................................;'
`),
				[]byte(`
  _                   _____                  
 | |                 |  __ \                 
 | |     __ _ _ __   | |__) |_ _ _ __   __ _ 
 | |    / _  | '_ \  |  ___/ _  | '_  \ / _ |
 | |___| (_| | | | | | |  | (_| | | | | (_| |
 |______\__,_|_| |_| |_|   \__,_|_| |_|\__, |
                                        __/ |
                                       |___/           
            .cc;.  ...  .;c.
         .,,cc:cc:lxxxl:ccc:;,.
        .lo;...lKKklllookl..cO;
      .cl;.,;'.okl;...'.;,..';:.
     .:o;;dkx,.ll..,cc::,..,'.;:,.
     co..lKKKkokl.':lllo;''ol..;dl.
   .,c;.,xKKKKKKo.':llll;.'oOxo,.cl,.
   cNo..lKKKKKKKo'';llll;;okKKKl..oNc
   cNo..lKKKKKKKko;':c:,'lKKKKKo'.oNc
   cNo..lKKKKKKKKKl.....'dKKKKKxc,l0:
   .c:'.lKKKKKKKKKk;....oKKKKKKo'.oNc
     ,:.,oxOKKKKKKKOxxxxOKKKKKKxc,;ol:.
     ;c..'':oookKKKKKKKKKKKKKKKKKk:.'clc.
   ,dl'.,oxo;'';oxOKKKKKKKKKKKKKKKOxxl::;,,.
  .dOc..lKKKkoooookKKKKKKKKKKKKKKKKKKKxl,;ol.
  cx,';okKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKl..;lc.
  co..:dddddddddddddddddddddddddddddddddl:;''::.
  co..........................................."
`),
				[]byte(`
  _                   _____                  
 | |                 |  __ \                 
 | |     __ _ _ __   | |__) |_ _ _ __   __ _ 
 | |    / _  | '_ \  |  ___/ _  | '_  \ / _ |
 | |___| (_| | | | | | |  | (_| | | | | (_| |
 |______\__,_|_| |_| |_|   \__,_|_| |_|\__, |
                                        __/ |
                                       |___/            
        .ccccccc.
      .,,,;cooolccol;;,,.
     .dOx;..;lllll;..;xOd.
   .cdo,',loOXXXXXkll;';odc.
  ,oo:;c,':oko:cccccc,...ckl.
  ;c.;kXo..::..;c::'.......oc
,dc..oXX0kk0o.':lll;..cxxc.,ld,
kNo.'oXXXXXXo'':lll;..oXXOd;cOd.
KOc;oOXXXXXXo.':lol,..dXXXXl';xc
Ol,:k0XXXXXX0c.,clc'.:0XXXXx,.oc
KOc;dOXXXXXXXl..';'..lXXXXXd..oc
dNo..oXXXXXXXOx:..'lxOXXXXXk,.:; ..
cNo..lXXXXXXXXXOolkXXXXXXXXXkl;..;:.;.
.,;'.,dkkkkk0XXXXXXXXXXXXXXXXXOxxl;,;,;l:.
  ;c.;:''''':doOXXXXXXXXXXXXXXXXXXOdo;';clc.
  ;c.lOdood:'''oXXXXXXXXXXXXXXXXXXXXXk,..;ol.
  ';.:xxxxxocccoxxxxxxxxxxxxxxxxxxxxxxl::'.';;.
  ';........................................;l'
`),
				[]byte(`
  _                   _____                  
 | |                 |  __ \                 
 | |     __ _ _ __   | |__) |_ _ _ __   __ _ 
 | |    / _  | '_ \  |  ___/ _  | '_  \ / _ |
 | |___| (_| | | | | | |  | (_| | | | | (_| |
 |______\__,_|_| |_| |_|   \__,_|_| |_|\__, |
                                        __/ |
                                       |___/ 
        .;:;;,.,;;::,.
     .;':;........'co:.
   .clc;'':cllllc::,.':c.
  .lo;;o:coxdlooollc;',::,,.
.c:'.,cl,.'lc',,;;'......cO;
do;';oxoc::l;;llllc'.';;'.';.
c..ckkkkkkkd,;llllc'.:kkd;.':c.
'.,okkkkkkkkc;llllc,.:kkkdl,cO;
..;xkkkkkkkkc,ccll:,;okkkkk:,cl,
..,dkkkkkkkkc..,;,'ckkkkkkkc;ll.
..'okkkkkkkko,....'okkkkkkkc,:c.
c..ckkkkkkkkkdl;,:okkkkkkkkd,.',';.
d..':lxkkkkkkkkxxkkkkkkkkkkkdoc;,;'..'.,.
o...'';llllldkkkkkkkkkkkkkkkkkkdll;..'cdo.
o..,l;'''''';dkkkkkkkkkkkkkkkkkkkkdlc,..;lc.
o..;lc;;;;;;,,;clllllllllllllllllllllc'..,:c.
o..........................................;'
`),
				[]byte(`
  _                   _____                  
 | |                 |  __ \                 
 | |     __ _ _ __   | |__) |_ _ _ __   __ _ 
 | |    / _  | '_ \  |  ___/ _  | '_  \ / _ |
 | |___| (_| | | | | | |  | (_| | | | | (_| |
 |______\__,_|_| |_| |_|   \__,_|_| |_|\__, |
                                        __/ |
                                       |___/ 
           .,,,,,,,,,.
         .ckKxodooxOOdcc.
      .cclooc'....';;cool.
     .loc;;;;clllllc;;;;;:;,.
   .c:'.,okd;;cdo:::::cl,..oc
  .:o;';okkx;';;,';::;'....,;,.
  co..ckkkkkddk:,cclll;.,c:,:o:.
  co..ckkkkkkkk:,cllll;.:kkd,.':c.
.,:;.,okkkkkkkk:,cclll;.:kkkdl;;o:.
cNo..ckkkkkkkkko,.;llc,.ckkkkkc..oc
,dd;.:kkkkkkkkkx;..;:,.'lkkkkko,.:,
  ;c.ckkkkkkkkkkc.....;ldkkkkkk:.,'
,dc..'okkkkkkkkkxoc;;cxkkkkkkkkc..,;,.
kNo..':lllllldkkkkkkkkkkkkkkkkkdcc,.;l.
KOc,l;''''''';lldkkkkkkkkkkkkkkkkkc..;lc.
xx:':;;;;,.,,...,;;cllllllllllllllc;'.;oo,
cNo.....................................oc
`),
				[]byte(`
  _                   _____                  
 | |                 |  __ \                 
 | |     __ _ _ __   | |__) |_ _ _ __   __ _ 
 | |    / _  | '_ \  |  ___/ _  | '_  \ / _ |
 | |___| (_| | | | | | |  | (_| | | | | (_| |
 |______\__,_|_| |_| |_|   \__,_|_| |_|\__, |
                                        __/ |
                                       |___/ 
                   .ccccccc.
               .ccckNKOOOOkdcc.
            .;;cc:ccccccc:,::::,,.
         .c;:;.,cccllxOOOxlllc,;ol.
        .lkc,coxo:;oOOxooooooo;..:,
      .cdc.,dOOOc..cOd,.',,;'....':c.
      cNx'.lOOOOxlldOl..;lll;.....cO;
     ,do;,:dOOOOOOOOOl'':lll;..:d:.'c,
     co..lOOOOOOOOOOOl'':lll;.'lOd,.cd.
     co.,dOOOOOOOOOOOo,.;llc,.,dOOc..dc
     co..lOOOOOOOOOOOOc.';:,..cOOOl..oc
   .,:;.'::lxOOOOOOOOOo:'...,:oOOOc..dc
   ;Oc..cl'':llxOOOOOOOOdcclxOOOOx,.cd.
  .:;';lxl''''':lldOOOOOOOOOOOOOOc..oc
,dl,.'cooc:::,....,::coooooooooooc'.c:
cNo.................................oc
`),
				[]byte(`
  _                   _____                  
 | |                 |  __ \                 
 | |     __ _ _ __   | |__) |_ _ _ __   __ _ 
 | |    / _  | '_ \  |  ___/ _  | '_  \ / _ |
 | |___| (_| | | | | | |  | (_| | | | | (_| |
 |______\__,_|_| |_| |_|   \__,_|_| |_|\__, |
                                        __/ |
                                       |___/ 

                        .cccccccc.
                  .,,,;;cc:cccccc:;;,.
                .cdxo;..,::cccc::,..;l.
               ,oo:,,:c:cdxxdllll:;,';:,.
             .cl;.,oxxc'.,cc,.',;;'...oNc
             ;Oc..cxxxc'.,c;..;lll;...cO;
           .;;',:ldxxxdoldxc..;lll:'...'c,
           ;c..cxxxxkxxkxxxc'.;lll:'','.cdc.
         .c;.;odxxxxxxxxxxxd;.,cll;.,l:.'dNc
        .:,''ccoxkxxkxxxxxxx:..,:;'.:xc..oNc
      .lc,.'lc':dxxxkxxxxxxxdl,...',lx:..dNc
     .:,',coxoc;;ccccoxxxxxxxxo:::oxxo,.cdc.
  .;':;.'oxxxxxc''''';cccoxxxxxxxxxkxc..oc
,do:'..,:llllll:;;;;;;,..,;:lllllllll;..oc
cNo.....................................oc
`),
				[]byte(`
  _                   _____                  
 | |                 |  __ \                 
 | |     __ _ _ __   | |__) |_ _ _ __   __ _ 
 | |    / _  | '_ \  |  ___/ _  | '_  \ / _ |
 | |___| (_| | | | | | |  | (_| | | | | (_| |
 |______\__,_|_| |_| |_|   \__,_|_| |_|\__, |
                                        __/ |
                                       |___/ 
                              .ccccc.
                         .cc;'coooxkl;.
                     .:c:::c:,;,,,;c;;,.'.
                   .clc,',:,..:xxocc;...c;
                  .c:,';:ox:..:c,,,,,,...cd,
                .c:'.,oxxxxl::l:.;loll;..;ol.
                ;Oc..:xxxxxxxxx:.,llll,....oc
             .,;,',:loxxxxxxxxx:.,llll;.,;.'ld,
            .lo;..:xxxxxxxxxxxx:.'cllc,.:l:'cO;
           .:;...'cxxxxxxxxxxxxol;,::,..cdl;;l'
         .cl;':;'';oxxxxxxxxxxxxx:....,cooc,cO;
     .,,,::;,lxoc:,,:lxxxxxxxxxxxo:,,;lxxl;'oNc
   .cdxo;':lxxxxxxc'';cccccoxxxxxxxxxxxxo,.;lc.
  .loc'.'lxxxxxxxxocc;''''';ccoxxxxxxxxx:..oc
occ'..',:cccccccccccc:;;;;;;;;:ccccccccc,.'c,
Ol;......................................;l'
`),
				[]byte(`
  _                   _____                  
 | |                 |  __ \                 
 | |     __ _ _ __   | |__) |_ _ _ __   __ _ 
 | |    / _  | '_ \  |  ___/ _  | '_  \ / _ |
 | |___| (_| | | | | | |  | (_| | | | | (_| |
 |______\__,_|_| |_| |_|   \__,_|_| |_|\__, |
                                        __/ |
                                       |___/ 
                              ,ddoodd,
                         .cc' ,ooccoo,'cc.
                      .ccldo;....,,...;oxdc.
                   .,,:cc;.''..;lol;;,'..lkl.
                  .dkc';:ccl;..;dl,.''.....oc
                .,lc',cdddddlccld;.,;c::'..,cc:.
                cNo..:ddddddddddd;':clll;,c,';xc
               .lo;,clddddddddddd;':clll;:kc..;'
             .,:;..:ddddddddddddd:';clll;;ll,..
             ;Oc..';:ldddddddddddl,.,c:;';dd;..
           .''',:lc,'cdddddddddddo:,'...'cdd;..
         .cdc';lddd:';lddddddddddddd;.';lddl,..
      .,;::;,cdddddol;;lllllodddddddlcodddd:.'l,
     .dOc..,lddddddddlccc;'';cclddddddddddd;,ll.
   .coc,;::ldddddddddddddl:ccc:ldddddddddlc,ck;
,dl::,..,cccccccccccccccccccccccccccccccc:;':xx,
cNd.........................................;lOc
`),
			},
		},
	}
}
