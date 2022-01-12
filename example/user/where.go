// Code generated by bcurd. DO NOT EDIT.

package user

import (
	"github.com/hongshengjie/crud/xsql"
)

type UserWhere func(s *xsql.Selector)

// IdEQ  =
func IdEQ(arg int64) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.EQ(Id, arg))
	})
}

// IdNEQ <>
func IdNEQ(arg int64) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.NEQ(Id, arg))
	})
}

// IdLT <
func IdLT(arg int64) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.LT(Id, arg))
	})
}

// IdLET <=
func IdLTE(arg int64) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.LTE(Id, arg))
	})
}

// IdGT >
func IdGT(arg int64) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.GT(Id, arg))
	})
}

// IdGTE >=
func IdGTE(arg int64) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.GTE(Id, arg))
	})
}

// IdIn in(...)
func IdIn(args ...int64) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		if len(args) == 0 {
			s.Where(xsql.False())
			return
		}
		v := make([]interface{}, len(args))
		for i := range v {
			v[i] = args[i]
		}
		s.Where(xsql.In(Id, v...))
	})
}

// IdNotIn not in(...)
func IdNotIn(args ...int64) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		if len(args) == 0 {
			s.Where(xsql.False())
			return
		}
		v := make([]interface{}, len(args))
		for i := range v {
			v[i] = args[i]
		}
		s.Where(xsql.NotIn(Id, v...))
	})
}

// NameEQ  =
func NameEQ(arg string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.EQ(Name, arg))
	})
}

// NameNEQ <>
func NameNEQ(arg string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.NEQ(Name, arg))
	})
}

// NameLT <
func NameLT(arg string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.LT(Name, arg))
	})
}

// NameLET <=
func NameLTE(arg string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.LTE(Name, arg))
	})
}

// NameGT >
func NameGT(arg string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.GT(Name, arg))
	})
}

// NameGTE >=
func NameGTE(arg string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.GTE(Name, arg))
	})
}

// NameIn in(...)
func NameIn(args ...string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		if len(args) == 0 {
			s.Where(xsql.False())
			return
		}
		v := make([]interface{}, len(args))
		for i := range v {
			v[i] = args[i]
		}
		s.Where(xsql.In(Name, v...))
	})
}

// NameNotIn not in(...)
func NameNotIn(args ...string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		if len(args) == 0 {
			s.Where(xsql.False())
			return
		}
		v := make([]interface{}, len(args))
		for i := range v {
			v[i] = args[i]
		}
		s.Where(xsql.NotIn(Name, v...))
	})
}

// NameHasPrefix HasPrefix
func NameHasPrefix(arg string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.HasPrefix(Name, arg))
	})
}

// NameHasSuffix HasSuffix
func NameHasSuffix(arg string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.HasSuffix(Name, arg))
	})
}

// NameContains Contains
func NameContains(arg string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.Contains(Name, arg))
	})
}

// AgeEQ  =
func AgeEQ(arg int64) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.EQ(Age, arg))
	})
}

// AgeNEQ <>
func AgeNEQ(arg int64) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.NEQ(Age, arg))
	})
}

// AgeLT <
func AgeLT(arg int64) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.LT(Age, arg))
	})
}

// AgeLET <=
func AgeLTE(arg int64) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.LTE(Age, arg))
	})
}

// AgeGT >
func AgeGT(arg int64) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.GT(Age, arg))
	})
}

// AgeGTE >=
func AgeGTE(arg int64) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.GTE(Age, arg))
	})
}

// AgeIn in(...)
func AgeIn(args ...int64) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		if len(args) == 0 {
			s.Where(xsql.False())
			return
		}
		v := make([]interface{}, len(args))
		for i := range v {
			v[i] = args[i]
		}
		s.Where(xsql.In(Age, v...))
	})
}

// AgeNotIn not in(...)
func AgeNotIn(args ...int64) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		if len(args) == 0 {
			s.Where(xsql.False())
			return
		}
		v := make([]interface{}, len(args))
		for i := range v {
			v[i] = args[i]
		}
		s.Where(xsql.NotIn(Age, v...))
	})
}

// CtimeEQ  =
func CtimeEQ(arg string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.EQ(Ctime, arg))
	})
}

// CtimeNEQ <>
func CtimeNEQ(arg string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.NEQ(Ctime, arg))
	})
}

// CtimeLT <
func CtimeLT(arg string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.LT(Ctime, arg))
	})
}

// CtimeLET <=
func CtimeLTE(arg string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.LTE(Ctime, arg))
	})
}

// CtimeGT >
func CtimeGT(arg string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.GT(Ctime, arg))
	})
}

// CtimeGTE >=
func CtimeGTE(arg string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.GTE(Ctime, arg))
	})
}

// CtimeIn in(...)
func CtimeIn(args ...string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		if len(args) == 0 {
			s.Where(xsql.False())
			return
		}
		v := make([]interface{}, len(args))
		for i := range v {
			v[i] = args[i]
		}
		s.Where(xsql.In(Ctime, v...))
	})
}

// CtimeNotIn not in(...)
func CtimeNotIn(args ...string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		if len(args) == 0 {
			s.Where(xsql.False())
			return
		}
		v := make([]interface{}, len(args))
		for i := range v {
			v[i] = args[i]
		}
		s.Where(xsql.NotIn(Ctime, v...))
	})
}

// MtimeEQ  =
func MtimeEQ(arg string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.EQ(Mtime, arg))
	})
}

// MtimeNEQ <>
func MtimeNEQ(arg string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.NEQ(Mtime, arg))
	})
}

// MtimeLT <
func MtimeLT(arg string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.LT(Mtime, arg))
	})
}

// MtimeLET <=
func MtimeLTE(arg string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.LTE(Mtime, arg))
	})
}

// MtimeGT >
func MtimeGT(arg string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.GT(Mtime, arg))
	})
}

// MtimeGTE >=
func MtimeGTE(arg string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.GTE(Mtime, arg))
	})
}

// MtimeIn in(...)
func MtimeIn(args ...string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		if len(args) == 0 {
			s.Where(xsql.False())
			return
		}
		v := make([]interface{}, len(args))
		for i := range v {
			v[i] = args[i]
		}
		s.Where(xsql.In(Mtime, v...))
	})
}

// MtimeNotIn not in(...)
func MtimeNotIn(args ...string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		if len(args) == 0 {
			s.Where(xsql.False())
			return
		}
		v := make([]interface{}, len(args))
		for i := range v {
			v[i] = args[i]
		}
		s.Where(xsql.NotIn(Mtime, v...))
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...UserWhere) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...UserWhere) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p UserWhere) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		p(s.Not())
	})
}
