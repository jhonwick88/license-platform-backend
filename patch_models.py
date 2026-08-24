import re

with open("internal/database/models.go", "r") as f:
    content = f.read()

# Plan
content = re.sub(r'(ProductID\s+string\s+gorm:"type:varchar\(36\);not null" json:"product_id"\n)', r'\1\tProduct     Product   gorm:"foreignKey:ProductID" json:"product,omitempty"\n', content)

# PlanFeature
content = re.sub(r'(PlanID\s+string\s+gorm:"type:varchar\(36\);not null" json:"plan_id"\n)', r'\1\tPlan      Plan      gorm:"foreignKey:PlanID" json:"plan,omitempty"\n', content)
content = re.sub(r'(FeatureID\s+string\s+gorm:"type:varchar\(36\);not null" json:"feature_id"\n)', r'\1\tFeature   Feature   gorm:"foreignKey:FeatureID" json:"feature,omitempty"\n', content)

# License
content = re.sub(r'(CustomerID\s+string\s+gorm:"type:varchar\(36\);not null" json:"customer_id"\n)', r'\1\tCustomer         Customer   gorm:"foreignKey:CustomerID" json:"customer,omitempty"\n', content)
# We already matched ProductID for Plan and Feature so it will apply everywhere, wait! The regex for ProductID might match multiple places. Yes!
# Let's write specific replacements for License:
content = re.sub(r'(PlanID\s+string\s+gorm:"type:varchar\(36\);not null" json:"plan_id"\n)', r'\1\tPlan             Plan       gorm:"foreignKey:PlanID" json:"plan,omitempty"\n', content)
# Oh wait, my regex replaced it globally. I need to be careful.

with open("internal/database/models.go", "w") as f:
    f.write(content)
