package permissions

const Private = 0 // 私有
const Inherit = 1 // 继承父目录的权限
const Public = 2  // 公开

func CalculatePublic(parentPublic bool, permission uint) bool {
	if permission == Inherit {
		return parentPublic
	}
	return permission == Public
}
