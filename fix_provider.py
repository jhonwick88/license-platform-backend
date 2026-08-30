import sys
import re

path = r'H:\FlutterProject\tokopintar\lib\presentation\providers\license_provider.dart'
with open(path, 'r', encoding='utf-8') as f:
    content = f.read()

# Replace the state update with SharedPreferences save and delay
pattern = re.compile(r'state = state\.copyWith\(\s*isActivated: true,\s*token: \'activated\',\s*isLoading: false,\s*\);')
replacement = """
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString('license_token', 'activated');

      state = state.copyWith(isLoading: false);
      
      Future.delayed(const Duration(seconds: 2), () {
        state = state.copyWith(
          token: 'activated',
          isActivated: true,
        );
      });
"""
new_content = pattern.sub(replacement, content)

with open(path, 'w', encoding='utf-8') as f:
    f.write(new_content)
