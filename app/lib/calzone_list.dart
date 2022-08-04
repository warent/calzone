library calzone;

import 'package:flutter/material.dart';

class CalzoneListView extends StatefulWidget {
  CalzoneListView({required this.scrollController});

  final ScrollController scrollController;

  @override
  _CalzoneListView createState() => _CalzoneListView();
}

class _CalzoneListView extends State<CalzoneListView> {
  Widget build(BuildContext context) {
    return ListView.builder(
        scrollDirection: Axis.vertical,
        // shrinkWrap: true,
        // controller: widget.scrollController,
        itemCount: 10,
        itemBuilder: (context, index) => Container(
              height: 50,
              child: TextButton(
                child: Text('Entry $index'),
                onPressed: () => print(index),
              ),
            ));
  }
}
