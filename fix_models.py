import sys

path = r'H:\FlutterProject\pintarlabs_license_platform\backend\internal\database\models.go'
with open(path, 'r', encoding='utf-8') as f:
    content = f.read()

# We need to remove the Licenses field from the User struct.
# The user struct looks like:
# type User struct {
# ...
# 	UpdatedAt    time.Time gorm:\"autoUpdateTime\" json:\"updated_at\"
# 	Licenses     []License gorm:\"foreignKey:CustomerID\" json:\"licenses,omitempty\"
# }

import re
pattern = re.compile(r'(type User struct \{.*?\n\s*UpdatedAt    time\.Time gorm:"autoUpdateTime" json:"updated_at")\n\s*Licenses     \[\]License gorm:"foreignKey:CustomerID" json:"licenses,omitempty"\n\}', re.DOTALL)
new_content = pattern.sub(r'\1\n}', content)

with open(path, 'w', encoding='utf-8') as f:
    f.write(new_content)
