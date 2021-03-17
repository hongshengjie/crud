package user

import (
	"github.com/hongshengjie/crud/xsql"
	"time"
)

type UserWhere func(s *xsql.Selector)

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

// IDEQ  =
func IDEQ(arg uint32) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.EQ(ID, arg))
	})
}

// IDNEQ <>
func IDNEQ(arg uint32) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.NEQ(ID, arg))
	})
}

// IDLT <
func IDLT(arg uint32) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.LT(ID, arg))
	})
}

// IDLET <=
func IDLTE(arg uint32) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.LTE(ID, arg))
	})
}

// IDGT >
func IDGT(arg uint32) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.GT(ID, arg))
	})
}

// IDGTE >=
func IDGTE(arg uint32) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.GTE(ID, arg))
	})
}

// IDIn in(...)
func IDIn(args ...uint32) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		if len(args) == 0 {
			s.Where(xsql.False())
			return
		}
		v := make([]interface{}, len(args))
		for i := range v {
			v[i] = args[i]
		}
		s.Where(xsql.In(ID, v...))
	})
}

// IDNotIn not in(...)
func IDNotIn(args ...uint32) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		if len(args) == 0 {
			s.Where(xsql.False())
			return
		}
		v := make([]interface{}, len(args))
		for i := range v {
			v[i] = args[i]
		}
		s.Where(xsql.NotIn(ID, v...))
	})
}

// NameEQ =
func NameEQ(v string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.EQ(Name, v))
	})
}

// NameEQ <>
func NameNEQ(v string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.NEQ(Name, v))
	})
}

// NameIn in(...)
func NameIn(vs ...string) UserWhere {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return UserWhere(func(s *xsql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(xsql.False())
			return
		}
		s.Where(xsql.In(Name, v...))
	})
}

// NameNotIn not int(...)
func NameNotIn(vs ...string) UserWhere {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return UserWhere(func(s *xsql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(xsql.False())
			return
		}
		s.Where(xsql.NotIn(Name, v...))
	})
}

// NameHasPrefix HasPrefix
func NameHasPrefix(v string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.HasPrefix(Name, v))
	})
}

// NameHasSuffix HasSuffix
func NameHasSuffix(v string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.HasSuffix(Name, v))
	})
}

// NameContains Contains
func NameContains(v string) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.Contains(Name, v))
	})
}

// MtimeEQ  =
func MtimeEQ(arg time.Time) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.EQ(Mtime, arg))
	})
}

// MtimeNEQ <>
func MtimeNEQ(arg time.Time) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.NEQ(Mtime, arg))
	})
}

// MtimeLT <
func MtimeLT(arg time.Time) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.LT(Mtime, arg))
	})
}

// MtimeLET <=
func MtimeLTE(arg time.Time) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.LTE(Mtime, arg))
	})
}

// MtimeGT >
func MtimeGT(arg time.Time) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.GT(Mtime, arg))
	})
}

// MtimeGTE >=
func MtimeGTE(arg time.Time) UserWhere {
	return UserWhere(func(s *xsql.Selector) {
		s.Where(xsql.GTE(Mtime, arg))
	})
}

// MtimeIn in(...)
func MtimeIn(args ...time.Time) UserWhere {
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
func MtimeNotIn(args ...time.Time) UserWhere {
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
