// Copyright 2016 NDP Systèmes. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testmodule

import (
	"fmt"

	"github.com/npiganeau/yep/pool"
	"github.com/npiganeau/yep/yep/models"
)

func declareModels() {
	user := models.NewModel("User")
	user.AddCharField("Name", models.StringFieldParams{String: "Name", Help: "The user's username", Unique: true})
	user.AddCharField("DecoratedName", models.StringFieldParams{Compute: "computeDecoratedName"})
	user.AddCharField("Email", models.StringFieldParams{Help: "The user's email address", Size: 100, Index: true})
	user.AddCharField("Password", models.StringFieldParams{})
	user.AddIntegerField("Status", models.SimpleFieldParams{JSON: "status_json", GoType: new(int16)})
	user.AddBooleanField("IsStaff", models.SimpleFieldParams{})
	user.AddBooleanField("IsActive", models.SimpleFieldParams{})
	user.AddMany2OneField("Profile", models.ForeignKeyFieldParams{RelationModel: "Profile"})
	user.AddIntegerField("Age", models.SimpleFieldParams{Compute: "computeAge", Depends: []string{"Profile", "Profile.Age"}, Stored: true, GoType: new(int16)})
	user.AddOne2ManyField("Posts", models.ReverseFieldParams{RelationModel: "Post", ReverseFK: "User"})
	user.AddFloatField("PMoney", models.FloatFieldParams{Related: "Profile.Money"})
	user.AddMany2OneField("LastPost", models.ForeignKeyFieldParams{RelationModel: "Post", Embed: true})
	user.AddCharField("Email2", models.StringFieldParams{})
	user.AddBooleanField("IsPremium", models.SimpleFieldParams{})
	user.AddIntegerField("Nums", models.SimpleFieldParams{GoType: new(int)})

	user.AddMethod("PrefixedUser",
		`PrefixedUser is a sample method layer for testing`,
		func(rs pool.UserSet, prefix string) []string {
			var res []string
			for _, u := range rs.Records() {
				res = append(res, fmt.Sprintf("%s: %s", prefix, u.Name()))
			}
			return res
		})

	user.AddMethod("DecorateEmail",
		`DecorateEmail is a sample method layer for testing`,
		func(rs pool.UserSet, email string) string {
			return fmt.Sprintf("<%s>", email)
		})

	pool.User().Methods().DecorateEmail().Extend(
		`DecorateEmailExtension is a sample method layer for testing`,
		func(rs pool.UserSet, email string) string {
			res := rs.Super().DecorateEmail(email)
			return fmt.Sprintf("[%s]", res)
		})

	user.AddMethod("computeAge",
		`ComputeAge is a sample method layer for testing`,
		func(rs pool.UserSet) (*pool.UserData, []models.FieldNamer) {
			res := pool.UserData{
				Age: rs.Profile().Age(),
			}
			return &res, []models.FieldNamer{pool.User().Age()}
		})

	pool.User().Methods().PrefixedUser().Extend("",
		func(rs pool.UserSet, prefix string) []string {
			res := rs.Super().PrefixedUser(prefix)
			for i, u := range rs.Records() {
				res[i] = fmt.Sprintf("%s %s", res[i], rs.DecorateEmail(u.Email()))
			}
			return res
		})

	user.AddMethod("computeDecoratedName", "",
		func(rs pool.UserSet) (*pool.UserData, []models.FieldNamer) {
			res := pool.UserData{
				DecoratedName: rs.PrefixedUser("User")[0],
			}
			return &res, []models.FieldNamer{pool.User().DecoratedName()}
		})

	user.AddMethod("UpdateCity", "",
		func(rs pool.UserSet, value string) {
			rs.Profile().SetCity(value)
		})

	profile := models.NewModel("Profile")
	profile.AddIntegerField("Age", models.SimpleFieldParams{GoType: new(int16)})
	profile.AddFloatField("Money", models.FloatFieldParams{})
	profile.AddMany2OneField("User", models.ForeignKeyFieldParams{RelationModel: "User"})
	profile.AddOne2OneField("BestPost", models.ForeignKeyFieldParams{RelationModel: "Post"})
	profile.AddCharField("City", models.StringFieldParams{})
	profile.AddCharField("Country", models.StringFieldParams{})

	pool.Profile().AddMethod("PrintAddress",
		`PrintAddress is a sample method layer for testing`,
		func(rs pool.ProfileSet) string {
			res := rs.Super().PrintAddress()
			return fmt.Sprintf("%s, %s", res, rs.Country())
		})

	pool.Profile().Methods().PrintAddress().Extend("",
		func(rs pool.ProfileSet) string {
			res := rs.Super().PrintAddress()
			return fmt.Sprintf("[%s]", res)
		})

	post := models.NewModel("Post")
	post.AddMany2OneField("User", models.ForeignKeyFieldParams{RelationModel: "User"})
	post.AddCharField("Title", models.StringFieldParams{})
	post.AddTextField("Content", models.StringFieldParams{})
	post.AddMany2ManyField("Tags", models.Many2ManyFieldParams{RelationModel: "Tag"})

	pool.Post().Methods().Create().Extend("",
		func(rs pool.PostSet, data models.FieldMapper) pool.PostSet {
			res := rs.Super().Create(data)
			return res
		})

	tag := models.NewModel("Tag")
	tag.AddCharField("Name", models.StringFieldParams{})
	tag.AddMany2OneField("BestPost", models.ForeignKeyFieldParams{RelationModel: "Post"})
	tag.AddMany2ManyField("Posts", models.Many2ManyFieldParams{RelationModel: "Post"})
	tag.AddCharField("Description", models.StringFieldParams{})

	addressMI := models.NewMixinModel("AddressMixIn")
	addressMI.AddCharField("Street", models.StringFieldParams{})
	addressMI.AddCharField("Zip", models.StringFieldParams{})
	addressMI.AddCharField("City", models.StringFieldParams{})
	profile.InheritModel(addressMI)

	addressMI2 := pool.AddressMixIn()
	addressMI2.AddMethod("SayHello",
		`SayHello is a sample method layer for testing`,
		func(rs pool.AddressMixInSet) string {
			return "Hello !"
		})

	addressMI2.AddMethod("PrintAddress",
		`PrintAddressMixIn is a sample method layer for testing`,
		func(rs pool.AddressMixInSet) string {
			return fmt.Sprintf("%s, %s %s", rs.Street(), rs.Zip(), rs.City())
		})

	addressMI2.Methods().PrintAddress().Extend("",
		func(rs pool.AddressMixInSet) string {
			res := rs.Super().PrintAddress()
			return fmt.Sprintf("<%s>", res)
		})

	activeMI := models.NewMixinModel("ActiveMixIn")
	activeMI.AddBooleanField("Active", models.SimpleFieldParams{})
	pool.ModelMixin().InheritModel(activeMI)

	// Chained declaration
	activeMI1 := pool.ActiveMixIn()
	activeMI2 := activeMI1
	activeMI2.AddMethod("IsActivated",
		`IsACtivated is a sample method of ActiveMixIn"`,
		func(rs pool.ActiveMixInSet) bool {
			return rs.Active()
		})

	viewModel := models.NewManualModel("UserView")
	viewModel.AddCharField("Name", models.StringFieldParams{})
	viewModel.AddCharField("City", models.StringFieldParams{})
}
