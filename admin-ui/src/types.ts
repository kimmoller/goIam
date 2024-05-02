type Identity = {
    id: string;
    firstName: string;
    lastName: string;
    email: string;
};

type Account = {
    username: string;
    systemId: string;
    enabledAt: string;
    disabledAt: string;
    deletedAt: string;
};

type Group = {
    id: string;
    name: string;
}

type Membership = {
    id: string;
    identityId: string;
    group: Group;
    enabledAt: string;
    disabledAt: string;
};

type ExtendedIdentity = {
    id: string;
    firstName: string;
    lastName: string;
    email: string;
    accounts: Account[];
    memberships: Membership[];
};

type CreateIdentity = {
    firstName: string;
    lastName: string;
    email: string;
}

type CreateMembership = {
    identityId: string;
    groupId: string;
    enabledAt: string;
    disabledAt: string | null;
}
