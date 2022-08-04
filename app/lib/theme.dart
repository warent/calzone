import 'package:yaru/yaru.dart';
import 'package:flutter/material.dart';

class AppTheme {
  static YaruThemeData of(BuildContext context) {
    return SharedAppData.getValue(context, 'theme', () => YaruThemeData());
  }

  static void apply(
    BuildContext context, {
    YaruVariant? variant,
    bool? highContrast,
    ThemeMode? themeMode,
  }) {
    SharedAppData.setValue(
      context,
      'theme',
      AppTheme.of(context).copyWith(
        themeMode: themeMode,
        variant: variant,
        highContrast: highContrast,
      ),
    );
  }
}
