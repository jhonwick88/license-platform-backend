import sys
import re

path = r'H:\FlutterProject\pintarlabs_license_platform\backend\internal\database\models.go'
with open(path, 'r', encoding='utf-8') as f:
    content = f.read()

# Remove foreignKey tags from Belongs To relationships
content = re.sub(r'Product\s+Product\s+gorm:"foreignKey:ProductID"\s+json:"product,omitempty"', r'Product     Product   json:"product,omitempty"', content)
content = re.sub(r'Plan\s+Plan\s+gorm:"foreignKey:PlanID"\s+json:"plan,omitempty"', r'Plan      Plan      json:"plan,omitempty"', content)
content = re.sub(r'Feature\s+Feature\s+gorm:"foreignKey:FeatureID"\s+json:"feature,omitempty"', r'Feature   Feature   json:"feature,omitempty"', content)
content = re.sub(r'Customer\s+Customer\s+gorm:"foreignKey:CustomerID"\s+json:"customer,omitempty"', r'Customer         Customer   json:"customer,omitempty"', content)
content = re.sub(r'License\s+\*License\s+gorm:"foreignKey:LicenseID"\s+json:"-"', r'License            *License   json:"-"', content)
content = re.sub(r'Installation\s+\*Installation\s+gorm:"foreignKey:InstallationID"\s+json:"installation,omitempty"', r'Installation     *Installation json:"installation,omitempty"', content)

with open(path, 'w', encoding='utf-8') as f:
    f.write(content)
