import 'package:flutter/material.dart';
import 'package:yaru/yaru.dart';
import '../calzone_list.dart';
import 'package:yaru_icons/yaru_icons.dart';

class RouteSetup extends StatefulWidget {
  const RouteSetup({Key? key, required this.title}) : super(key: key);

  // This widget is the home page of your application. It is stateful, meaning
  // that it has a State object (defined below) that contains fields that affect
  // how it looks.

  // This class is the configuration for the state. It holds the values (in this
  // case the title) provided by the parent (in this case the App widget) and
  // used by the build method of the State. Fields in a Widget subclass are
  // always marked "final".

  final String title;

  @override
  State<RouteSetup> createState() => _RouteSetupState();
}

class _RouteSetupState extends State<RouteSetup> {
  int _counter = 0;
  bool _connected = false;
  bool _loading = false;

  final ScrollController scrollController = ScrollController();

  void _incrementCounter() {
    setState(() {
      // This call to setState tells the Flutter framework that something has
      // changed in this State, which causes it to rerun the build method below
      // so that the display can reflect the updated values. If we changed
      // _counter without calling setState(), then the build method would not be
      // called again, and so nothing would appear to happen.
      _counter++;
    });
  }

  @override
  Widget build(BuildContext context) {
    // This method is rerun every time setState is called, for instance as done
    // by the _incrementCounter method above.
    //
    // The Flutter framework has been optimized to make rerunning build methods
    // fast, so that you can just rebuild anything that needs updating rather
    // than having to individually change instances of widgets.
    return Scaffold(
      body: Center(
          // Center is a layout widget. It takes a single child and positions it
          // in the middle of the parent.
          child: Column(children: [
        Expanded(flex: 2, child: Container()),
        Expanded(
            flex: 6,
            child: Row(
              children: <Widget>[
                Expanded(
                    child: Container(
                        decoration: BoxDecoration(
                            border:
                                Border(right: BorderSide(color: Colors.black))),
                        child: Column(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              ElevatedButton(
                                  onPressed: () => print("NEW"),
                                  child: Text("Setup Calzone on this system"))
                            ]))),
                Expanded(
                    child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                      Container(
                          padding: EdgeInsets.only(bottom: 16),
                          child: Text("Remote Calzone Connection")),
                      Container(
                        padding: EdgeInsets.only(bottom: 16),
                        width: 200,
                        child: TextField(
                          decoration: InputDecoration(
                              border: OutlineInputBorder(),
                              label: Text("Hostname / IP Address")),
                        ),
                      ),
                      Container(
                        width: 200,
                        padding: EdgeInsets.only(bottom: 16),
                        child: TextField(
                          obscureText: true,
                          decoration: InputDecoration(
                              border: OutlineInputBorder(),
                              label: Text("Password")),
                        ),
                      ),
                      Container(
                          width: 200,
                          child: ElevatedButton(
                              onPressed: () => print("NEW"),
                              child: Text("Connect"))),
                    ])),
              ],
            )),
        Expanded(flex: 2, child: Container()),
      ])),
      // floatingActionButton: FloatingActionButton(
      //   onPressed: _incrementCounter,
      //   tooltip: 'Increment',
      //   child: const Icon(Icons.add),
      // ), // This trailing comma makes auto-formatting nicer for build methods.
    );
  }
}