import os
import re

directory = r"H:\FlutterProject\pintarlabs_license_platform\backend\internal"
pattern = re.compile(r'(?<!)\b(json|gorm|binding):"([^"]*)"')

for root, dirs, files in os.walk(directory):
    for file in files:
        if file.endswith(".go"):
            filepath = os.path.join(root, file)
            with open(filepath, "r", encoding="utf-8") as f:
                content = f.read()
            
            # Since multiple tags can be on the same line like gorm:"..." json:"..."
            # and they might have been stripped of all backticks entirely
            # let's just do a simpler search: replace things that look like tag keys
            # Actually, Go struct tags are usually surrounded by backticks if they contain spaces or quotes.
            # If the PowerShell script removed the backticks, the text looks like:
            # Field string json:"field" binding:"required"
            # We want to replace it with:
            # Field string json:"field" binding:"required"
            
            # Let's use a regex to find lines that define a struct field and wrap the trailing tags in backticks.
            # Example: ID string gorm:"primaryKey" json:"id"
            pass

def fix_go_struct_tags(filepath):
    with open(filepath, "r", encoding="utf-8") as f:
        lines = f.readlines()
    
    changed = False
    for i, line in enumerate(lines):
        # check if line has struct tags without backticks
        # e.g., looks for json:"..." or gorm:"..."
        if ('json:"' in line or 'gorm:"' in line or 'binding:"' in line) and '' not in line:
            # find where the tags start
            # usually after the type. We can just find the first occurrence of json:, gorm:, binding:
            m = re.search(r'\b(json|gorm|binding):"', line)
            if m:
                start_idx = m.start()
                # find the end of the line (before newline)
                end_idx = len(line.rstrip())
                
                tags = line[start_idx:end_idx]
                new_line = line[:start_idx] + '' + tags + '' + line[end_idx:]
                lines[i] = new_line
                changed = True
                
    if changed:
        with open(filepath, "w", encoding="utf-8") as f:
            f.writelines(lines)
        print(f"Fixed {filepath}")

for root, dirs, files in os.walk(directory):
    for file in files:
        if file.endswith(".go"):
            filepath = os.path.join(root, file)
            fix_go_struct_tags(filepath)

